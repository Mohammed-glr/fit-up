package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)


type workoutSessionServiceImpl struct {
	repo repository.SchemaRepo
}

func NewWorkoutSessionService(repo repository.SchemaRepo) WorkoutSessionService {
	return &workoutSessionServiceImpl{
		repo: repo,
	}
}


func (s *workoutSessionServiceImpl) StartWorkoutSession(ctx context.Context, userID int, workoutID int) (*types.WorkoutSession, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if workoutID <= 0 {
		return nil, fmt.Errorf("invalid workout ID")
	}

	activeSession, err := s.repo.WorkoutSessions().GetActiveSession(ctx, userID)
	if err == nil && activeSession != nil {
		return nil, fmt.Errorf("user already has an active workout session (ID: %d)", activeSession.SessionID)
	}

	if err := s.validateWorkoutAccessForUser(ctx, userID, workoutID); err != nil {
		return nil, fmt.Errorf("workout access validation failed: %w", err)
	}

	if recoveryStatus, err := s.assessRecoveryBeforeWorkout(ctx, userID); err == nil {
		if recoveryStatus.RestDayRecommended {
			fmt.Printf("Warning: Rest day recommended for user %d, but starting workout anyway\n", userID)
		}
	}

	session, err := s.repo.WorkoutSessions().StartWorkoutSession(ctx, userID, workoutID)
	if err != nil {
		return nil, fmt.Errorf("failed to start workout session: %w", err)
	}

	fmt.Printf("Started workout session %d for user %d (workout %d)\n",
		session.SessionID, userID, workoutID)

	return session, nil
}

func (s *workoutSessionServiceImpl) CompleteWorkoutSession(ctx context.Context, sessionID int, summary *types.SessionSummary) (*types.WorkoutSession, error) {
	if sessionID <= 0 {
		return nil, fmt.Errorf("invalid session ID")
	}

	if summary == nil {
		return nil, fmt.Errorf("session summary is required")
	}

	session, err := s.validateActiveSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session validation failed: %w", err)
	}

	_ = s.calculateSessionPerformanceMetrics(summary)

	completionAnalysis := s.analyzeWorkoutCompletion(session, summary)

	if err := s.updateUserProgressFromSession(ctx, session.UserID, sessionID, summary); err != nil {
		fmt.Printf("Warning: Failed to update user progress: %v\n", err)
	}

	completedSession, err := s.repo.WorkoutSessions().CompleteWorkoutSession(ctx, sessionID, summary)
	if err != nil {
		return nil, fmt.Errorf("failed to complete session: %w", err)
	}

	if err := s.triggerPlanAdaptationIfNeeded(ctx, session.UserID, completionAnalysis); err != nil {
		fmt.Printf("Warning: Failed to trigger plan adaptation: %v\n", err)
	}

	fmt.Printf("Completed workout session %d: %.1f%% completion, %.1f average RPE\n",
		sessionID, completionAnalysis.CompletionRate*100, summary.AverageRPE)

	return completedSession, nil
}

func (s *workoutSessionServiceImpl) SkipWorkout(ctx context.Context, userID int, workoutID int, reason string) (*types.SkippedWorkout, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if workoutID <= 0 {
		return nil, fmt.Errorf("invalid workout ID")
	}

	if reason == "" {
		reason = "No reason provided"
	}

	skipPattern := s.analyzeSkipPattern(ctx, userID, reason)

	skippedWorkout, err := s.repo.WorkoutSessions().SkipWorkout(ctx, userID, workoutID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to record skipped workout: %w", err)
	}

	if err := s.updatePlanBasedOnSkipPattern(ctx, userID, skipPattern, reason); err != nil {
		fmt.Printf("Warning: Failed to update plan based on skip pattern: %v\n", err)
	}

	if err := s.adjustFuturePlansForSkipReason(ctx, userID, reason); err != nil {
		fmt.Printf("Warning: Failed to adjust future plans: %v\n", err)
	}

	fmt.Printf("Workout skipped for user %d (reason: %s). Skip pattern analysis: %+v\n",
		userID, reason, skipPattern)

	return skippedWorkout, nil
}


func (s *workoutSessionServiceImpl) LogExercisePerformance(ctx context.Context, sessionID int, exerciseID int, performance *types.ExercisePerformance) error {
	if sessionID <= 0 {
		return fmt.Errorf("invalid session ID")
	}

	if exerciseID <= 0 {
		return fmt.Errorf("invalid exercise ID")
	}

	if performance == nil {
		return fmt.Errorf("performance data is required")
	}

	session, err := s.validateActiveSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session validation failed: %w", err)
	}

	analysis := s.analyzeExercisePerformanceRealTime(performance)

	if analysis.FormWarnings != nil {
		fmt.Printf("Form warnings detected for user %d exercise %d: %v\n",
			session.UserID, exerciseID, analysis.FormWarnings)
	}

	if analysis.PlateauIndicator {
		fmt.Printf("Plateau indicator detected for user %d exercise %d\n",
			session.UserID, exerciseID)
	}

	if err := s.repo.WorkoutSessions().LogExercisePerformance(ctx, sessionID, exerciseID, performance); err != nil {
		return fmt.Errorf("failed to log exercise performance: %w", err)
	}

	if performance.BestSet.Weight > 0 && performance.BestSet.Reps > 0 {
		if err := s.updateStrengthProgressions(ctx, session.UserID, exerciseID, &performance.BestSet); err != nil {
			fmt.Printf("Warning: Failed to update strength progressions: %v\n", err)
		}
	}

	return nil
}

func (s *workoutSessionServiceImpl) GetActiveSession(ctx context.Context, userID int) (*types.WorkoutSession, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	session, err := s.repo.WorkoutSessions().GetActiveSession(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	if session != nil {
		if err := s.validateSessionNotStale(session); err != nil {
			fmt.Printf("Session %d for user %d marked as stale: %v\n", session.SessionID, userID, err)
			return nil, nil
		}
	}

	return session, nil
}



func (s *workoutSessionServiceImpl) GetSessionHistory(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutSession], error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	// Apply FitUp Smart Logic: Add filtering and sorting options
	sessions, err := s.repo.WorkoutSessions().GetSessionHistory(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get session history: %w", err)
	}

	if len(sessions.Data) > 0 {
		trends := s.calculateSessionTrends(sessions.Data)
		fmt.Printf("Session trends for user %d: %+v\n", userID, trends)
	}

	return sessions, nil
}

func (s *workoutSessionServiceImpl) GetSessionMetrics(ctx context.Context, sessionID int) (*types.SessionMetrics, error) {
	if sessionID <= 0 {
		return nil, fmt.Errorf("invalid session ID")
	}

	metrics, err := s.repo.WorkoutSessions().GetSessionMetrics(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session metrics: %w", err)
	}

	enhancedMetrics := s.enhanceSessionMetrics(ctx, metrics)

	return enhancedMetrics, nil
}

func (s *workoutSessionServiceImpl) GetWeeklySessionStats(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySessionStats, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	stats, err := s.repo.WorkoutSessions().GetWeeklySessionStats(ctx, userID, weekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly session stats: %w", err)
	}

	trends := s.calculateWeeklyTrends(ctx, userID, weekStart, stats)
	fmt.Printf("Weekly trends for user %d (week %v): %+v\n", userID, weekStart, trends)

	return stats, nil
}


type SessionCompletionAnalysis struct {
	CompletionRate     float64
	QualityScore       float64
	PerformanceLevel   string
	AdaptationRequired bool
	RecommendedChanges []string
}

type PerformanceAnalysis struct {
	FormWarnings     []string
	PlateauIndicator bool
	ProgressionReady bool
	RPEConsistency   float64
}

type SkipPattern struct {
	RecentSkips      int
	SkipFrequency    float64
	CommonReasons    []string
	AdaptationNeeded bool
}

func (s *workoutSessionServiceImpl) validateWorkoutAccessForUser(ctx context.Context, userID int, workoutID int) error {

	_, err := s.repo.Workouts().GetWorkoutByID(ctx, workoutID)
	return err
}

func (s *workoutSessionServiceImpl) assessRecoveryBeforeWorkout(ctx context.Context, userID int) (*types.RecoveryStatus, error) {
	return s.repo.RecoveryMetrics().GetRecoveryStatus(ctx, userID)
}

func (s *workoutSessionServiceImpl) validateActiveSession(ctx context.Context, sessionID int) (*types.WorkoutSession, error) {

	return &types.WorkoutSession{SessionID: sessionID}, nil 
}

func (s *workoutSessionServiceImpl) calculateSessionPerformanceMetrics(summary *types.SessionSummary) map[string]interface{} {
	metrics := make(map[string]interface{})

	metrics["completion_rate"] = float64(summary.ExercisesCompleted) / float64(len(summary.Exercises))
	metrics["average_rpe"] = summary.AverageRPE
	metrics["total_volume"] = summary.TotalVolume
	metrics["duration_minutes"] = summary.TotalDuration / 60

	return metrics
}

func (s *workoutSessionServiceImpl) analyzeWorkoutCompletion(session *types.WorkoutSession, summary *types.SessionSummary) *SessionCompletionAnalysis {
	completionRate := float64(summary.ExercisesCompleted) / float64(len(summary.Exercises))

	analysis := &SessionCompletionAnalysis{
		CompletionRate: completionRate,
		QualityScore:   s.calculateQualityScore(summary),
	}

	if completionRate >= 0.90 {
		analysis.PerformanceLevel = "excellent"
		analysis.AdaptationRequired = false
	} else if completionRate >= 0.70 {
		analysis.PerformanceLevel = "good"
		analysis.AdaptationRequired = false
	} else {
		analysis.PerformanceLevel = "poor"
		analysis.AdaptationRequired = true
		analysis.RecommendedChanges = []string{"reduce_intensity", "extend_timeline"}
	}

	return analysis
}

func (s *workoutSessionServiceImpl) calculateQualityScore(summary *types.SessionSummary) float64 {
	rpeScore := 1.0
	if summary.AverageRPE > 9 {
		rpeScore = 0.7 
	} else if summary.AverageRPE < 5 {
		rpeScore = 0.8 
	}

	return rpeScore
}

func (s *workoutSessionServiceImpl) updateUserProgressFromSession(ctx context.Context, userID int, sessionID int, summary *types.SessionSummary) error {
	for _, exercise := range summary.Exercises {
		progressLog := &types.ProgressLogRequest{
			UserID:        userID,
			ExerciseID:    exercise.ExerciseID,
			Date:          time.Now(),
			SetsCompleted: &exercise.SetsCompleted,
			RepsCompleted: &exercise.BestSet.Reps,
			WeightUsed:    &exercise.BestSet.Weight,
		}

		if _, err := s.repo.Progress().CreateProgressLog(ctx, progressLog); err != nil {
			fmt.Printf("Warning: Failed to create progress log for exercise %d: %v\n", exercise.ExerciseID, err)
		}
	}

	return nil
}

func (s *workoutSessionServiceImpl) triggerPlanAdaptationIfNeeded(ctx context.Context, userID int, analysis *SessionCompletionAnalysis) error {
	if !analysis.AdaptationRequired {
		return nil
	}

	fmt.Printf("Plan adaptation triggered for user %d: %s performance\n", userID, analysis.PerformanceLevel)
	return nil
}

func (s *workoutSessionServiceImpl) analyzeSkipPattern(ctx context.Context, userID int, reason string) *SkipPattern {
	return &SkipPattern{
		RecentSkips:      1,
		SkipFrequency:    0.1,
		CommonReasons:    []string{reason},
		AdaptationNeeded: false,
	}
}

func (s *workoutSessionServiceImpl) updatePlanBasedOnSkipPattern(ctx context.Context, userID int, pattern *SkipPattern, reason string) error {
	if pattern.AdaptationNeeded {
		fmt.Printf("Updating plan for user %d based on skip pattern\n", userID)
	}
	return nil
}

func (s *workoutSessionServiceImpl) adjustFuturePlansForSkipReason(ctx context.Context, userID int, reason string) error {
	switch reason {
	case "lack_of_time":
		fmt.Printf("Considering shorter workouts for user %d\n", userID)
	case "fatigue":
		fmt.Printf("Considering lower intensity for user %d\n", userID)
	case "injury":
		fmt.Printf("Considering exercise modifications for user %d\n", userID)
	}
	return nil
}

func (s *workoutSessionServiceImpl) analyzeExercisePerformanceRealTime(performance *types.ExercisePerformance) *PerformanceAnalysis {
	analysis := &PerformanceAnalysis{
		FormWarnings:     []string{},
		PlateauIndicator: false,
		ProgressionReady: false,
		RPEConsistency:   performance.RPE,
	}

	if performance.RPE > 9 {
		analysis.FormWarnings = append(analysis.FormWarnings, "Very high RPE - check form")
	}

	if performance.SetsCompleted < performance.SetsCompleted {
		analysis.PlateauIndicator = true
	}

	return analysis
}

func (s *workoutSessionServiceImpl) updateStrengthProgressions(ctx context.Context, userID int, exerciseID int, bestSet *types.SetPerformance) error {
	performanceData := &types.PerformanceData{
		Weight:   bestSet.Weight,
		Reps:     bestSet.Reps,
		Sets:     1,
		RPE:      bestSet.RPE,
		Duration: 0,
	}

	_, err := s.repo.FitnessProfiles().EstimateOneRepMax(ctx, userID, exerciseID, performanceData)
	return err
}

func (s *workoutSessionServiceImpl) validateSessionNotStale(session *types.WorkoutSession) error {
	if time.Since(session.StartTime) > 24*time.Hour {
		return fmt.Errorf("session started more than 24 hours ago")
	}
	return nil
}

func (s *workoutSessionServiceImpl) calculateSessionTrends(sessions []types.WorkoutSession) map[string]interface{} {
	if len(sessions) == 0 {
		return nil
	}

	trends := make(map[string]interface{})
	trends["total_sessions"] = len(sessions)
	trends["completion_trend"] = "stable"

	return trends
}

func (s *workoutSessionServiceImpl) enhanceSessionMetrics(ctx context.Context, metrics *types.SessionMetrics) *types.SessionMetrics {

	return metrics
}

func (s *workoutSessionServiceImpl) calculateWeeklyTrends(ctx context.Context, userID int, weekStart time.Time, stats *types.WeeklySessionStats) map[string]interface{} {
	trends := make(map[string]interface{})

	trends["volume_change"] = "increasing"
	trends["consistency_score"] = stats.CompletionRate

	if stats.CompletionRate < 0.7 {
		trends["recommendation"] = "consider_plan_adjustment"
	}

	return trends
}
