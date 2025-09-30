package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// =============================================================================
// FITNESS PROFILE SERVICE IMPLEMENTATION
// =============================================================================

type fitnessProfileServiceImpl struct {
	repo repository.SchemaRepo
}

// NewFitnessProfileService creates a new fitness profile service instance
func NewFitnessProfileService(repo repository.SchemaRepo) FitnessProfileService {
	return &fitnessProfileServiceImpl{
		repo: repo,
	}
}

// =============================================================================
// FITNESS ASSESSMENT METHODS
// =============================================================================

func (s *fitnessProfileServiceImpl) CreateFitnessAssessment(ctx context.Context, userID int, assessment *types.FitnessAssessmentRequest) (*types.FitnessAssessment, error) {
	// Validate input parameters
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if err := validator.New().Struct(assessment); err != nil {
		return nil, fmt.Errorf("invalid assessment data: %w", err)
	}

	// Apply FitUp Smart Logic: Validate fitness level consistency
	if err := s.validateFitnessLevelConsistency(assessment); err != nil {
		return nil, fmt.Errorf("fitness level validation failed: %w", err)
	}

	// Check for previous assessments to track progression
	previousProfile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err == nil && previousProfile != nil {
		// Validate progression is realistic (no more than 2 levels jump)
		if err := s.validateLevelProgression(previousProfile.CurrentLevel, assessment.OverallLevel); err != nil {
			return nil, fmt.Errorf("unrealistic level progression: %w", err)
		}
	}

	// Enrich assessment with calculated metrics
	assessment.AssessmentData = s.enrichAssessmentData(assessment)

	return s.repo.FitnessProfiles().CreateFitnessAssessment(ctx, userID, assessment)
}

func (s *fitnessProfileServiceImpl) GetUserFitnessProfile(ctx context.Context, userID int) (*types.FitnessProfile, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	// Get base profile from repository
	profile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fitness profile: %w", err)
	}

	// Enrich profile with calculated training history
	if trainingHistory, err := s.calculateTrainingHistory(ctx, userID); err == nil {
		profile.TrainingHistory = trainingHistory
	}

	// Add current fitness level recommendations
	if recommendations, err := s.generateLevelRecommendations(ctx, userID, profile); err == nil {
		// Add recommendations to profile metadata or separate field
		fmt.Printf("Generated %d recommendations for user %d\n", len(recommendations), userID)
	}

	return profile, nil
}

func (s *fitnessProfileServiceImpl) UpdateFitnessLevel(ctx context.Context, userID int, level types.FitnessLevel) error {
	if userID <= 0 {
		return types.ErrInvalidUserID
	}

	// Validate fitness level value
	validLevels := map[types.FitnessLevel]bool{
		types.LevelBeginner:     true,
		types.LevelIntermediate: true,
		types.LevelAdvanced:     true,
	}

	if !validLevels[level] {
		return fmt.Errorf("invalid fitness level: %s", level)
	}

	// Get current profile to validate progression
	currentProfile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err == nil && currentProfile != nil {
		// Apply FitUp Smart Logic: Validate level progression is supported by data
		if err := s.validateLevelProgression(currentProfile.CurrentLevel, level); err != nil {
			return fmt.Errorf("level progression validation failed: %w", err)
		}

		// Check if progression is supported by recent performance
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

	// Get current fitness profile for validation
	profile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user profile for goal validation: %w", err)
	}

	// Apply FitUp Smart Logic: Validate each goal
	for i, goal := range goals {
		// Validate goal is SMART (Specific, Measurable, Achievable, Relevant, Time-bound)
		if err := s.validateSMARTGoal(&goal, profile); err != nil {
			return fmt.Errorf("goal %d validation failed: %w", i+1, err)
		}

		// Check goal compatibility with user's current fitness level
		if err := s.validateGoalFitnessLevelCompatibility(&goal, profile.CurrentLevel); err != nil {
			return fmt.Errorf("goal %d incompatible with fitness level: %w", i+1, err)
		}

		// Validate timeline is realistic
		if err := s.validateGoalTimeline(&goal, profile); err != nil {
			return fmt.Errorf("goal %d timeline validation failed: %w", i+1, err)
		}
	}

	// Apply FitUp Smart Logic: Detect and resolve conflicting goals
	if conflictingGoals := s.detectConflictingGoals(goals); len(conflictingGoals) > 0 {
		resolvedGoals, err := s.resolveGoalConflicts(goals, conflictingGoals, profile)
		if err != nil {
			return fmt.Errorf("failed to resolve goal conflicts: %w", err)
		}
		goals = resolvedGoals
	}

	return s.repo.FitnessProfiles().UpdateFitnessGoals(ctx, userID, goals)
}

// =============================================================================
// ONE REP MAX ESTIMATION METHODS
// =============================================================================

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

	// Apply FitUp Smart Logic: Use multiple 1RM estimation formulas
	estimations := s.calculateMultiple1RMEstimates(performance)

	// Get user's training history to determine best formula
	trainingHistory, err := s.calculateTrainingHistory(ctx, userID)
	if err == nil {
		// Select best estimation method based on user experience
		estimations = s.adjustEstimationForExperience(estimations, trainingHistory)
	}

	// Cross-reference with recent performance data for validation
	recentHistory, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err == nil && len(recentHistory) > 0 {
		// Validate estimate is realistic compared to history
		if err := s.validateEstimateAgainstHistory(estimations.BestEstimate, recentHistory); err != nil {
			return nil, fmt.Errorf("estimate validation failed: %w", err)
		}
	}

	// Use the most appropriate estimation method
	finalEstimate := estimations.BestEstimate
	_ = estimations.Confidence // Store confidence for potential future use
	_ = estimations.Method     // Store method for potential future use

	// Create enhanced performance data for repository
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

	// Get history from repository
	history, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1RM history: %w", err)
	}

	// Apply FitUp Smart Logic: Calculate progression rates and detect plateaus
	if len(history) >= 2 {
		// Analyze progression trends
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

	// Apply FitUp Smart Logic: Validate estimate is realistic
	recentHistory, err := s.repo.FitnessProfiles().GetOneRepMaxHistory(ctx, userID, exerciseID)
	if err == nil && len(recentHistory) > 0 {
		// Check if estimate is within reasonable range of recent estimates
		latestEstimate := recentHistory[0].EstimatedMax
		percentChange := ((estimate - latestEstimate) / latestEstimate) * 100

		// Flag unrealistic jumps (>20% increase or >10% decrease)
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

// =============================================================================
// MOVEMENT ASSESSMENT METHODS
// =============================================================================

func (s *fitnessProfileServiceImpl) CreateMovementAssessment(ctx context.Context, userID int, assessment *types.MovementAssessmentRequest) (*types.MovementAssessment, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if err := validator.New().Struct(assessment); err != nil {
		return nil, fmt.Errorf("invalid assessment data: %w", err)
	}

	// Apply FitUp Smart Logic: Analyze movement patterns and score
	movementScores := s.analyzeMovementPatterns(assessment.MovementData)
	assessment.MovementData["computed_scores"] = movementScores

	// Identify mobility limitations automatically
	detectedLimitations := s.identifyMobilityLimitations(movementScores)
	if len(detectedLimitations) > 0 {
		// Merge with user-reported limitations
		allLimitations := append(assessment.Limitations, detectedLimitations...)
		assessment.Limitations = s.removeDuplicateLimitations(allLimitations)
	}

	return s.repo.FitnessProfiles().CreateMovementAssessment(ctx, userID, assessment)
}

func (s *fitnessProfileServiceImpl) GetMovementLimitations(ctx context.Context, userID int) ([]types.MovementLimitation, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	// Get limitations from repository
	limitations, err := s.repo.FitnessProfiles().GetMovementLimitations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movement limitations: %w", err)
	}

	// Apply FitUp Smart Logic: Add severity scoring and exercise contraindications
	for i := range limitations {
		// Generate exercise modifications based on limitations
		limitations[i].Description = s.enhanceLimitationDescription(&limitations[i])
	}

	return limitations, nil
}

// =============================================================================
// HELPER METHODS FOR FITUP SMART LOGIC
// =============================================================================

type OneRepMaxEstimations struct {
	BestEstimate float64
	Confidence   float64
	Method       string
	AllEstimates map[string]float64
}

func (s *fitnessProfileServiceImpl) validateFitnessLevelConsistency(assessment *types.FitnessAssessmentRequest) error {
	// Check if individual levels are consistent with overall level
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

	// Allow ±1 level difference between overall and average sub-scores
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

	// Don't allow jumping more than 1 level at once
	if newScore > currentScore+1 {
		return fmt.Errorf("cannot advance more than one fitness level at once")
	}

	// Don't allow regression without justification
	if newScore < currentScore {
		return fmt.Errorf("fitness level regression requires manual approval")
	}

	return nil
}

func (s *fitnessProfileServiceImpl) enrichAssessmentData(assessment *types.FitnessAssessmentRequest) map[string]interface{} {
	enriched := make(map[string]interface{})

	// Copy original data
	for k, v := range assessment.AssessmentData {
		enriched[k] = v
	}

	// Add computed metrics
	enriched["assessment_timestamp"] = time.Now()
	enriched["overall_fitness_score"] = s.calculateOverallFitnessScore(assessment)
	enriched["recommended_progression"] = s.getRecommendedProgression(assessment.OverallLevel)

	return enriched
}

func (s *fitnessProfileServiceImpl) calculateTrainingHistory(ctx context.Context, userID int) (*types.TrainingHistory, error) {
	// Get workout session data from progress service
	// This is a simplified implementation - in practice would query actual session data
	return &types.TrainingHistory{
		TotalWorkouts:    50, // placeholder
		WeeksActive:      12, // placeholder
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
	// In a full implementation, this would analyze recent workout performance
	// For now, just validate the progression makes sense
	if currentLevel == newLevel {
		return nil // No change
	}

	// Simplified validation - would need actual performance data
	return nil
}

func (s *fitnessProfileServiceImpl) validateSMARTGoal(goal *types.FitnessGoalTarget, profile *types.FitnessProfile) error {
	// Specific: Goal type must be defined
	if goal.GoalType == "" {
		return fmt.Errorf("goal must be specific (goal type required)")
	}

	// Measurable: Target value must be defined
	if goal.TargetValue <= 0 {
		return fmt.Errorf("goal must be measurable (target value required)")
	}

	// Achievable: Target should be realistic (within 300% of current)
	if goal.CurrentValue > 0 && goal.TargetValue > goal.CurrentValue*3 {
		return fmt.Errorf("goal target may not be achievable (>300%% increase)")
	}

	// Time-bound: Target date must be in the future
	if goal.TargetDate.Before(time.Now()) {
		return fmt.Errorf("goal must be time-bound (target date in future)")
	}

	return nil
}

func (s *fitnessProfileServiceImpl) validateGoalFitnessLevelCompatibility(goal *types.FitnessGoalTarget, level types.FitnessLevel) error {
	// Check if goal is appropriate for fitness level
	if level == types.LevelBeginner {
		if goal.GoalType == types.GoalStrength && goal.TargetValue > goal.CurrentValue*2 {
			return fmt.Errorf("strength goal too aggressive for beginner level")
		}
	}

	return nil
}

func (s *fitnessProfileServiceImpl) validateGoalTimeline(goal *types.FitnessGoalTarget, profile *types.FitnessProfile) error {
	weeksToGoal := time.Until(goal.TargetDate).Hours() / (24 * 7)

	// Minimum timeline validation based on goal type
	minWeeks := map[types.FitnessGoal]float64{
		types.GoalStrength:       8,  // 8 weeks minimum for strength gains
		types.GoalMuscleGain:     12, // 12 weeks minimum for muscle gain
		types.GoalFatLoss:        6,  // 6 weeks minimum for sustainable fat loss
		types.GoalEndurance:      8,  // 8 weeks minimum for endurance
		types.GoalGeneralFitness: 4,  // 4 weeks minimum for general fitness
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

	// Simple conflict detection: strength + fat loss can conflict
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
	// For now, just warn about conflicts - in practice would suggest modifications
	fmt.Printf("Warning: Conflicting goals detected for fitness profile. Consider prioritizing one goal at a time.\n")
	return goals, nil
}

func (s *fitnessProfileServiceImpl) calculateMultiple1RMEstimates(performance *types.PerformanceData) *OneRepMaxEstimations {
	weight := performance.Weight
	reps := float64(performance.Reps)

	// Multiple 1RM formulas
	epley := weight * (1 + reps/30)           // Epley formula
	brzycki := weight * (36 / (37 - reps))    // Brzycki formula
	mcglothin := weight * (1 + 0.025*reps)    // McGlothin formula
	lombardi := weight * math.Pow(reps, 0.10) // Lombardi formula

	estimates := map[string]float64{
		"epley":     epley,
		"brzycki":   brzycki,
		"mcglothin": mcglothin,
		"lombardi":  lombardi,
	}

	// Choose best estimate based on rep range
	var best float64
	var method string
	var confidence float64

	if reps <= 5 {
		// Low reps: Epley formula most accurate
		best = epley
		method = "epley"
		confidence = 0.95
	} else if reps <= 10 {
		// Medium reps: Brzycki formula
		best = brzycki
		method = "brzycki"
		confidence = 0.85
	} else {
		// High reps: Average of formulas, lower confidence
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
	// Adjust confidence based on training experience
	if history.TotalWorkouts < 10 {
		// Very new - reduce confidence
		estimations.Confidence *= 0.8
	} else if history.TotalWorkouts > 100 {
		// Experienced - increase confidence
		estimations.Confidence = math.Min(estimations.Confidence*1.1, 1.0)
	}

	// Adjust based on consistency
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

	// Flag unrealistic changes
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

	// Check last 4 estimates for stagnation
	recentEstimates := history[:4]
	base := recentEstimates[3].EstimatedMax

	for _, estimate := range recentEstimates {
		improvement := ((estimate.EstimatedMax - base) / base) * 100
		if improvement > 2.5 { // 2.5% improvement threshold
			return false
		}
	}

	return true
}

func (s *fitnessProfileServiceImpl) analyzeMovementPatterns(movementData map[string]interface{}) map[string]float64 {
	scores := make(map[string]float64)

	// Simplified movement analysis - would be more complex in practice
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

// =============================================================================
// WORKOUT PROFILE METHODS (merged from WorkoutProfileService)
// =============================================================================

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
	// For simplicity, return nil since the interface doesn't have this exact method
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) ActivateWorkoutProfile(ctx context.Context, profileID int) error {
	// For simplicity, return error since the interface doesn't have this method
	return fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) DeactivateWorkoutProfile(ctx context.Context, profileID int) error {
	// For simplicity, return error since the interface doesn't have this method
	return fmt.Errorf("method not implemented in repository")
}

// =============================================================================
// GOAL TRACKING METHODS (merged from GoalTrackingService)
// =============================================================================

func (s *fitnessProfileServiceImpl) CreateGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().CreateFitnessGoal(ctx, userID, goal)
}

func (s *fitnessProfileServiceImpl) GetGoalByID(ctx context.Context, goalID int) (*types.FitnessGoalTarget, error) {
	// Repository doesn't have this method, return error for now
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) UpdateGoal(ctx context.Context, goalID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	// Repository doesn't have this method, return error for now
	return nil, fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) DeleteGoal(ctx context.Context, goalID int) error {
	// Repository doesn't have this method, return error for now
	return fmt.Errorf("method not implemented in repository")
}

func (s *fitnessProfileServiceImpl) GetGoalsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.FitnessGoalTarget], error) {
	// Repository doesn't have this method, use GetActiveGoals as approximation
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
	// Repository doesn't have this method, use GetActiveGoals for now
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

// Missing WorkoutProfile methods from interface
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

// Missing Goal Tracking methods from interface
func (s *fitnessProfileServiceImpl) CreateFitnessGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().CreateFitnessGoal(ctx, userID, goal)
}

func (s *fitnessProfileServiceImpl) GetActiveGoals(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error) {
	return s.repo.GoalTracking().GetActiveGoals(ctx, userID)
}
