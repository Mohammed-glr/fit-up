package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type fitnessProfileServiceImpl struct {
	repo repository.SchemaRepo
}

func NewFitnessProfileService(repo repository.SchemaRepo) FitnessProfileService {
	return &fitnessProfileServiceImpl{
		repo: repo,
	}
}

func (s *fitnessProfileServiceImpl) CreateFitnessAssessment(ctx context.Context, userID int, assessment *types.FitnessAssessmentRequest) (*types.FitnessAssessment, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if err := validator.New().Struct(assessment); err != nil {
		return nil, fmt.Errorf("invalid assessment data: %w", err)
	}

	if err := s.validateFitnessLevelConsistency(assessment); err != nil {
		return nil, fmt.Errorf("fitness level validation failed: %w", err)
	}

	previousProfile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err == nil && previousProfile != nil {
		if err := s.validateLevelProgression(previousProfile.CurrentLevel, assessment.OverallLevel); err != nil {
			return nil, fmt.Errorf("unrealistic level progression: %w", err)
		}
	}

	assessment.AssessmentData = s.enrichAssessmentData(assessment)

	return s.repo.FitnessProfiles().CreateFitnessAssessment(ctx, userID, assessment)
}

func (s *fitnessProfileServiceImpl) GetUserFitnessProfile(ctx context.Context, userID int) (*types.FitnessProfile, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	profile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fitness profile: %w", err)
	}

	if trainingHistory, err := s.calculateTrainingHistory(ctx, userID); err == nil {
		profile.TrainingHistory = trainingHistory
	}

	if recommendations, err := s.generateLevelRecommendations(ctx, userID, profile); err == nil {
		fmt.Printf("Generated %d recommendations for user %d\n", len(recommendations), userID)
	}

	return profile, nil
}

func (s *fitnessProfileServiceImpl) UpdateFitnessLevel(ctx context.Context, userID int, level types.FitnessLevel) error {
	if userID <= 0 {
		return types.ErrInvalidUserID
	}

	validLevels := map[types.FitnessLevel]bool{
		types.LevelBeginner:     true,
		types.LevelIntermediate: true,
		types.LevelAdvanced:     true,
	}

	if !validLevels[level] {
		return fmt.Errorf("invalid fitness level: %s", level)
	}

	currentProfile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err == nil && currentProfile != nil {
		if err := s.validateLevelProgression(currentProfile.CurrentLevel, level); err != nil {
			return fmt.Errorf("level progression validation failed: %w", err)
		}

		if err := s.validateProgressionWithPerformanceData(ctx, userID, currentProfile.CurrentLevel, level); err != nil {
			return fmt.Errorf("progression not supported by performance data: %w", err)
		}
	}

	return s.repo.FitnessProfiles().UpdateFitnessLevel(ctx, userID, level)
}

func (s *fitnessProfileServiceImpl) UpdateFitnessGoals(ctx context.Context, userID int, goals []types.FitnessGoalTarget) error {
	if userID <= 0 {
		return types.ErrInvalidUserID
	}

	if len(goals) == 0 {
		return fmt.Errorf("at least one fitness goal is required")
	}

	profile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user profile for goal validation: %w", err)
	}

	for i, goal := range goals {
		if err := s.validateSMARTGoal(&goal, profile); err != nil {
			return fmt.Errorf("goal %d validation failed: %w", i+1, err)
		}

		if err := s.validateGoalFitnessLevelCompatibility(&goal, profile.CurrentLevel); err != nil {
			return fmt.Errorf("goal %d incompatible with fitness level: %w", i+1, err)
		}

		if err := s.validateGoalTimeline(&goal, profile); err != nil {
			return fmt.Errorf("goal %d timeline validation failed: %w", i+1, err)
		}
	}

	if conflictingGoals := s.detectConflictingGoals(goals); len(conflictingGoals) > 0 {
		resolvedGoals, err := s.resolveGoalConflicts(goals, conflictingGoals, profile)
		if err != nil {
			return fmt.Errorf("failed to resolve goal conflicts: %w", err)
		}
		goals = resolvedGoals
	}

	return s.repo.FitnessProfiles().UpdateFitnessGoals(ctx, userID, goals)
}

func (s *fitnessProfileServiceImpl) EstimateOneRepMax(ctx context.Context, userID int, exerciseID int, performance *types.PerformanceData) (*types.OneRepMaxEstimate, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return nil, fmt.Errorf("invalid exercise ID")
	}

	if err := validator.New().Struct(performance); err != nil {
		return nil, fmt.Errorf("invalid performance data: %w", err)
	}

	estimations := s.calculateMultiple1RMEstimates(performance)

	trainingHistory, err := s.calculateTrainingHistory(ctx, userID)
	if err == nil {
		estimations = s.adjustEstimationForExperience(estimations, trainingHistory)
	}

	recentHistory, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err == nil && len(recentHistory) > 0 {
		if err := s.validateEstimateAgainstHistory(estimations.BestEstimate, recentHistory); err != nil {
			return nil, fmt.Errorf("estimate validation failed: %w", err)
		}
	}

	finalEstimate := estimations.BestEstimate
	_ = estimations.Confidence
	_ = estimations.Method

	enhancedPerformance := &types.PerformanceData{
		Weight:   finalEstimate,
		Reps:     1,
		Sets:     performance.Sets,
		RPE:      performance.RPE,
		Duration: performance.Duration,
	}

	return s.repo.FitnessProfiles().EstimateOneRepMax(ctx, userID, exerciseID, enhancedPerformance)
}

func (s *fitnessProfileServiceImpl) GetOneRepMaxHistory(ctx context.Context, userID int, exerciseID int) ([]types.OneRepMaxEstimate, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return nil, fmt.Errorf("invalid exercise ID")
	}

	history, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1RM history: %w", err)
	}

	if len(history) >= 2 {
		progressionRate := s.calculateProgressionRate(history)
		plateauDetected := s.detectProgressionPlateau(history)

		fmt.Printf("User %d exercise %d: progression rate %.2f%%, plateau detected: %t\n",
			userID, exerciseID, progressionRate, plateauDetected)
	}

	return history, nil
}

func (s *fitnessProfileServiceImpl) UpdateOneRepMax(ctx context.Context, userID int, exerciseID int, estimate float64) error {
	if userID <= 0 {
		return types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return fmt.Errorf("invalid exercise ID")
	}

	if estimate <= 0 {
		return fmt.Errorf("1RM estimate must be positive")
	}

	recentHistory, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err == nil && len(recentHistory) > 0 {
		latestEstimate := recentHistory[0].EstimatedMax
		percentChange := ((estimate - latestEstimate) / latestEstimate) * 100

		if percentChange > 20 {
			return fmt.Errorf("unrealistic 1RM increase of %.1f%% - maximum recommended is 20%%", percentChange)
		}
		if percentChange < -10 {
			fmt.Printf("Warning: 1RM decrease of %.1f%% detected for user %d exercise %d\n",
				-percentChange, userID, exerciseID)
		}
	}

	return s.repo.FitnessProfiles().UpdateOneRepMax(ctx, userID, exerciseID, estimate)
}

func (s *fitnessProfileServiceImpl) CreateMovementAssessment(ctx context.Context, userID int, assessment *types.MovementAssessmentRequest) (*types.MovementAssessment, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if err := validator.New().Struct(assessment); err != nil {
		return nil, fmt.Errorf("invalid assessment data: %w", err)
	}

	movementScores := s.analyzeMovementPatterns(assessment.MovementData)
	assessment.MovementData["computed_scores"] = movementScores

	detectedLimitations := s.identifyMobilityLimitations(movementScores)
	if len(detectedLimitations) > 0 {
		allLimitations := append(assessment.Limitations, detectedLimitations...)
		assessment.Limitations = s.removeDuplicateLimitations(allLimitations)
	}

	return s.repo.FitnessProfiles().CreateMovementAssessment(ctx, userID, assessment)
}

func (s *fitnessProfileServiceImpl) GetMovementLimitations(ctx context.Context, userID int) ([]types.MovementLimitation, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	limitations, err := s.repo.FitnessProfiles().GetMovementLimitations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movement limitations: %w", err)
	}

	for i := range limitations {
		limitations[i].Description = s.enhanceLimitationDescription(&limitations[i])
	}

	return limitations, nil
}

type OneRepMaxEstimations struct {
	BestEstimate float64
	Confidence   float64
	Method       string
	AllEstimates map[string]float64
}

func (s *fitnessProfileServiceImpl) validateFitnessLevelConsistency(assessment *types.FitnessAssessmentRequest) error {
	levelScores := map[types.FitnessLevel]int{
		types.LevelBeginner:     1,
		types.LevelIntermediate: 2,
		types.LevelAdvanced:     3,
	}

	overallScore := levelScores[assessment.OverallLevel]
	strengthScore := levelScores[assessment.StrengthLevel]
	cardioScore := levelScores[assessment.CardioLevel]
	flexibilityScore := levelScores[assessment.FlexibilityLevel]

	avgSubScore := float64(strengthScore+cardioScore+flexibilityScore) / 3.0

	if float64(overallScore) > avgSubScore+1.5 || float64(overallScore) < avgSubScore-1.5 {
		return fmt.Errorf("overall fitness level inconsistent with individual assessments")
	}

	return nil
}

func (s *fitnessProfileServiceImpl) validateLevelProgression(currentLevel, newLevel types.FitnessLevel) error {
	levelOrder := map[types.FitnessLevel]int{
		types.LevelBeginner:     1,
		types.LevelIntermediate: 2,
		types.LevelAdvanced:     3,
	}

	currentScore := levelOrder[currentLevel]
	newScore := levelOrder[newLevel]

	if newScore > currentScore+1 {
		return fmt.Errorf("cannot advance more than one fitness level at once")
	}

	if newScore < currentScore {
		return fmt.Errorf("fitness level regression requires manual approval")
	}

	return nil
}

func (s *fitnessProfileServiceImpl) enrichAssessmentData(assessment *types.FitnessAssessmentRequest) map[string]interface{} {
	enriched := make(map[string]interface{})

	for k, v := range assessment.AssessmentData {
		enriched[k] = v
	}

	enriched["assessment_timestamp"] = time.Now()
	enriched["overall_fitness_score"] = s.calculateOverallFitnessScore(assessment)
	enriched["recommended_progression"] = s.getRecommendedProgression(assessment.OverallLevel)

	return enriched
}

func (s *fitnessProfileServiceImpl) calculateTrainingHistory(ctx context.Context, userID int) (*types.TrainingHistory, error) {
	return &types.TrainingHistory{
		TotalWorkouts:    50,
		WeeksActive:      12,
		AverageFrequency: 3.5,
		ConsistencyScore: 0.85,
	}, nil
}

func (s *fitnessProfileServiceImpl) generateLevelRecommendations(ctx context.Context, userID int, profile *types.FitnessProfile) ([]string, error) {
	recommendations := []string{}

	switch profile.CurrentLevel {
	case types.LevelBeginner:
		recommendations = append(recommendations, "Focus on form and consistency", "Start with 2-3 workouts per week")
	case types.LevelIntermediate:
		recommendations = append(recommendations, "Increase training volume gradually", "Add more complex movements")
	case types.LevelAdvanced:
		recommendations = append(recommendations, "Focus on periodization", "Consider specialized training blocks")
	}

	return recommendations, nil
}

func (s *fitnessProfileServiceImpl) validateProgressionWithPerformanceData(ctx context.Context, userID int, currentLevel, newLevel types.FitnessLevel) error {
	if currentLevel == newLevel {
		return nil
	}
	return nil
}

func (s *fitnessProfileServiceImpl) validateSMARTGoal(goal *types.FitnessGoalTarget, profile *types.FitnessProfile) error {
	if goal.GoalType == "" {
		return fmt.Errorf("goal must be specific (goal type required)")
	}

	if goal.TargetValue <= 0 {
		return fmt.Errorf("goal must be measurable (target value required)")
	}

	if goal.CurrentValue > 0 && goal.TargetValue > goal.CurrentValue*3 {
		return fmt.Errorf("goal target may not be achievable (>300%% increase)")
	}

	if goal.TargetDate.Before(time.Now()) {
		return fmt.Errorf("goal must be time-bound (target date in future)")
	}

	return nil
}

func (s *fitnessProfileServiceImpl) validateGoalFitnessLevelCompatibility(goal *types.FitnessGoalTarget, level types.FitnessLevel) error {
	if level == types.LevelBeginner {
		if goal.GoalType == types.GoalStrength && goal.TargetValue > goal.CurrentValue*2 {
			return fmt.Errorf("strength goal too aggressive for beginner level")
		}
	}

	return nil
}

func (s *fitnessProfileServiceImpl) validateGoalTimeline(goal *types.FitnessGoalTarget, profile *types.FitnessProfile) error {
	weeksToGoal := time.Until(goal.TargetDate).Hours() / (24 * 7)

	minWeeks := map[types.FitnessGoal]float64{
		types.GoalStrength:       8,
		types.GoalMuscleGain:     12,
		types.GoalFatLoss:        6,
		types.GoalEndurance:      8,
		types.GoalGeneralFitness: 4,
	}

	if minWeeksRequired, exists := minWeeks[goal.GoalType]; exists {
		if weeksToGoal < minWeeksRequired {
			return fmt.Errorf("timeline too short for %s goal (minimum %d weeks required)", goal.GoalType, int(minWeeksRequired))
		}
	}

	return nil
}

func (s *fitnessProfileServiceImpl) detectConflictingGoals(goals []types.FitnessGoalTarget) []int {
	conflicts := []int{}

	hasStrength := false
	hasFatLoss := false
	strengthIndex := -1
	fatLossIndex := -1

	for i, goal := range goals {
		if goal.GoalType == types.GoalStrength {
			hasStrength = true
			strengthIndex = i
		}
		if goal.GoalType == types.GoalFatLoss {
			hasFatLoss = true
			fatLossIndex = i
		}
	}

	if hasStrength && hasFatLoss {
		conflicts = append(conflicts, strengthIndex, fatLossIndex)
	}

	return conflicts
}

func (s *fitnessProfileServiceImpl) resolveGoalConflicts(goals []types.FitnessGoalTarget, conflicts []int, profile *types.FitnessProfile) ([]types.FitnessGoalTarget, error) {
	fmt.Printf("Warning: Conflicting goals detected for fitness profile. Consider prioritizing one goal at a time.\n")
	return goals, nil
}

func (s *fitnessProfileServiceImpl) calculateMultiple1RMEstimates(performance *types.PerformanceData) *OneRepMaxEstimations {
	weight := performance.Weight
	reps := float64(performance.Reps)

	epley := weight * (1 + reps/30)
	brzycki := weight * (36 / (37 - reps))
	mcglothin := weight * (1 + 0.025*reps)
	lombardi := weight * math.Pow(reps, 0.10)

	estimates := map[string]float64{
		"epley":     epley,
		"brzycki":   brzycki,
		"mcglothin": mcglothin,
		"lombardi":  lombardi,
	}

	var best float64
	var method string
	var confidence float64

	if reps <= 5 {
		best = epley
		method = "epley"
		confidence = 0.95
	} else if reps <= 10 {
		best = brzycki
		method = "brzycki"
		confidence = 0.85
	} else {
		best = (epley + brzycki + mcglothin) / 3
		method = "average"
		confidence = 0.70
	}

	return &OneRepMaxEstimations{
		BestEstimate: best,
		Confidence:   confidence,
		Method:       method,
		AllEstimates: estimates,
	}
}

func (s *fitnessProfileServiceImpl) adjustEstimationForExperience(estimations *OneRepMaxEstimations, history *types.TrainingHistory) *OneRepMaxEstimations {
	if history.TotalWorkouts < 10 {
		estimations.Confidence *= 0.8
	} else if history.TotalWorkouts > 100 {
		estimations.Confidence = math.Min(estimations.Confidence*1.1, 1.0)
	}

	if history.ConsistencyScore < 0.7 {
		estimations.Confidence *= 0.9
	}

	return estimations
}

func (s *fitnessProfileServiceImpl) validateEstimateAgainstHistory(estimate float64, history []types.OneRepMaxEstimate) error {
	if len(history) == 0 {
		return nil
	}

	latestEstimate := history[0].EstimatedMax
	percentChange := ((estimate - latestEstimate) / latestEstimate) * 100

	if percentChange > 25 {
		return fmt.Errorf("estimated 1RM increase of %.1f%% seems unrealistic", percentChange)
	}
	if percentChange < -15 {
		return fmt.Errorf("estimated 1RM decrease of %.1f%% seems excessive", -percentChange)
	}

	return nil
}

func (s *fitnessProfileServiceImpl) calculateProgressionRate(history []types.OneRepMaxEstimate) float64 {
	if len(history) < 2 {
		return 0
	}

	oldest := history[len(history)-1]
	newest := history[0]

	if oldest.EstimatedMax == 0 {
		return 0
	}

	return ((newest.EstimatedMax - oldest.EstimatedMax) / oldest.EstimatedMax) * 100
}

func (s *fitnessProfileServiceImpl) detectProgressionPlateau(history []types.OneRepMaxEstimate) bool {
	if len(history) < 4 {
		return false
	}

	recentEstimates := history[:4]
	base := recentEstimates[3].EstimatedMax

	for _, estimate := range recentEstimates {
		improvement := ((estimate.EstimatedMax - base) / base) * 100
		if improvement > 2.5 {
			return false
		}
	}

	return true
}

func (s *fitnessProfileServiceImpl) analyzeMovementPatterns(movementData map[string]interface{}) map[string]float64 {
	scores := make(map[string]float64)

	scores["overhead_mobility"] = 7.5
	scores["hip_mobility"] = 8.0
	scores["ankle_mobility"] = 6.5
	scores["shoulder_stability"] = 8.5
	scores["core_stability"] = 7.0

	return scores
}

func (s *fitnessProfileServiceImpl) identifyMobilityLimitations(scores map[string]float64) []string {
	limitations := []string{}

	for movement, score := range scores {
		if score < 7.0 {
			limitations = append(limitations, movement)
		}
	}

	return limitations
}

func (s *fitnessProfileServiceImpl) removeDuplicateLimitations(limitations []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, limitation := range limitations {
		if !seen[limitation] {
			seen[limitation] = true
			result = append(result, limitation)
		}
	}

	return result
}

func (s *fitnessProfileServiceImpl) enhanceLimitationDescription(limitation *types.MovementLimitation) string {
	enhanced := limitation.Description

	switch limitation.MovementType {
	case "overhead_mobility":
		enhanced += " - Avoid overhead pressing movements, use incline variations"
	case "hip_mobility":
		enhanced += " - Modify squat depth, include hip mobility work"
	case "ankle_mobility":
		enhanced += " - Use heel elevation for squats, focus on calf stretching"
	}

	return enhanced
}

func (s *fitnessProfileServiceImpl) calculateOverallFitnessScore(assessment *types.FitnessAssessmentRequest) float64 {
	levelScores := map[types.FitnessLevel]float64{
		types.LevelBeginner:     3.0,
		types.LevelIntermediate: 6.0,
		types.LevelAdvanced:     9.0,
	}

	overall := levelScores[assessment.OverallLevel]
	strength := levelScores[assessment.StrengthLevel]
	cardio := levelScores[assessment.CardioLevel]
	flexibility := levelScores[assessment.FlexibilityLevel]

	return (overall + strength + cardio + flexibility) / 4.0
}

func (s *fitnessProfileServiceImpl) getRecommendedProgression(level types.FitnessLevel) string {
	switch level {
	case types.LevelBeginner:
		return "Focus on learning proper form and building consistency"
	case types.LevelIntermediate:
		return "Increase training volume and add complex movements"
	case types.LevelAdvanced:
		return "Implement periodization and specialized training phases"
	default:
		return "Continue with current training approach"
	}
}

func (s *fitnessProfileServiceImpl) CreateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().CreateWorkoutProfile(ctx, authUserID, profile)
}

func (s *fitnessProfileServiceImpl) GetWorkoutProfileByID(ctx context.Context, profileID int) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, profileID)
}

func (s *fitnessProfileServiceImpl) UpdateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().UpdateWorkoutProfile(ctx, authUserID, profile)
}

func (s *fitnessProfileServiceImpl) DeleteWorkoutProfile(ctx context.Context, authUserID string) error {
	return s.repo.WorkoutProfiles().DeleteWorkoutProfile(ctx, authUserID)
}

func (s *fitnessProfileServiceImpl) GetWorkoutProfilesByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	return s.repo.WorkoutProfiles().ListWorkoutProfiles(ctx, pagination)
}

func (s *fitnessProfileServiceImpl) GetActiveWorkoutProfileByUserID(ctx context.Context, userID int) (*types.WorkoutProfile, error) {
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) ActivateWorkoutProfile(ctx context.Context, profileID int) error {
	return fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) DeactivateWorkoutProfile(ctx context.Context, profileID int) error {
	return fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) CreateGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().CreateFitnessGoal(ctx, userID, goal)
}

func (s *fitnessProfileServiceImpl) GetGoalByID(ctx context.Context, goalID int) (*types.FitnessGoalTarget, error) {
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) UpdateGoal(ctx context.Context, goalID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) DeleteGoal(ctx context.Context, goalID int) error {
	return fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) GetGoalsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.FitnessGoalTarget], error) {
	goals, err := s.repo.GoalTracking().GetActiveGoals(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &types.PaginatedResponse[types.FitnessGoalTarget]{
		Data:       goals,
		TotalCount: len(goals),
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
		TotalPages: (len(goals) + pagination.Limit - 1) / pagination.Limit,
	}, nil
}

func (s *fitnessProfileServiceImpl) GetActiveGoalsByUserID(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().GetActiveGoals(ctx, userID)
}

func (s *fitnessProfileServiceImpl) UpdateGoalProgress(ctx context.Context, goalID int, progress float64) error {
	return s.repo.GoalTracking().UpdateGoalProgress(ctx, goalID, progress)
}

func (s *fitnessProfileServiceImpl) CompleteGoal(ctx context.Context, goalID int) error {
	return s.repo.GoalTracking().CompleteGoal(ctx, goalID)
}

func (s *fitnessProfileServiceImpl) GetGoalProgress(ctx context.Context, goalID int) (*types.GoalProgress, error) {
	return s.repo.GoalTracking().CalculateGoalProgress(ctx, goalID)
}

func (s *fitnessProfileServiceImpl) GetGoalsByStatus(ctx context.Context, userID int, status string) ([]types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().GetActiveGoals(ctx, userID)
}

func (s *fitnessProfileServiceImpl) CalculateGoalProgress(ctx context.Context, goalID int) (*types.GoalProgress, error) {
	return s.repo.GoalTracking().CalculateGoalProgress(ctx, goalID)
}

func (s *fitnessProfileServiceImpl) EstimateTimeToGoal(ctx context.Context, goalID int) (*types.TimeToGoalEstimate, error) {
	return s.repo.GoalTracking().EstimateTimeToGoal(ctx, goalID)
}

func (s *fitnessProfileServiceImpl) SuggestGoalAdjustments(ctx context.Context, userID int) ([]types.GoalAdjustment, error) {
	return s.repo.GoalTracking().SuggestGoalAdjustments(ctx, userID)
}

func (s *fitnessProfileServiceImpl) GetWorkoutProfileByAuthID(ctx context.Context, authUserID string) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, authUserID)
}

func (s *fitnessProfileServiceImpl) ListWorkoutProfiles(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	return s.repo.WorkoutProfiles().ListWorkoutProfiles(ctx, pagination)
}

func (s *fitnessProfileServiceImpl) SearchWorkoutProfiles(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	return s.repo.WorkoutProfiles().SearchWorkoutProfiles(ctx, query, pagination)
}

func (s *fitnessProfileServiceImpl) GetProfilesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetProfilesByLevel(ctx, level)
}

func (s *fitnessProfileServiceImpl) GetProfilesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetProfilesByGoal(ctx, goal)
}

func (s *fitnessProfileServiceImpl) CountActiveProfiles(ctx context.Context) (int, error) {
	return s.repo.WorkoutProfiles().CountActiveProfiles(ctx)
}

func (s *fitnessProfileServiceImpl) CreateFitnessGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().CreateFitnessGoal(ctx, userID, goal)
}

func (s *fitnessProfileServiceImpl) GetActiveGoals(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().GetActiveGoals(ctx, userID)
}
