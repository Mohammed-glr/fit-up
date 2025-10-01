package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type performanceAnalyticsServiceImpl struct {
	repo repository.SchemaRepo
}

func NewPerformanceAnalyticsService(repo repository.SchemaRepo) PerformanceAnalyticsService {
	return &performanceAnalyticsServiceImpl{
		repo: repo,
	}
}

func (s *performanceAnalyticsServiceImpl) CalculateStrengthProgression(ctx context.Context, userID int, exerciseID int, timeframe int) (*types.StrengthProgression, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return nil, fmt.Errorf("invalid exercise ID")
	}

	if timeframe <= 0 || timeframe > 365 {
		return nil, fmt.Errorf("timeframe must be between 1 and 365 days")
	}

	progression, err := s.repo.PerformanceAnalytics().CalculateStrengthProgression(ctx, userID, exerciseID, timeframe)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate strength progression: %w", err)
	}

	progressionAnalysis := s.analyzeStrengthProgressionAdvanced(ctx, userID, exerciseID, timeframe, progression)

	if progressionAnalysis.HasSufficientData {
		progression.ProgressionRate = s.calculateStatisticalProgressionRate(progressionAnalysis.DataPoints)
		progression.Trend = s.determineTrendWithConfidence(progressionAnalysis.DataPoints)
	}

	if imbalances := s.identifyStrengthImbalances(ctx, userID, exerciseID); len(imbalances) > 0 {
		fmt.Printf("Strength imbalances detected for user %d exercise %d: %v\n", userID, exerciseID, imbalances)
	}

	if futurePrediction := s.predictFutureStrengthGains(progressionAnalysis); futurePrediction != nil {
		fmt.Printf("Predicted strength gains for user %d exercise %d: %+v\n", userID, exerciseID, futurePrediction)
	}

	trainingImpact := s.analyzeTrainingFrequencyImpact(ctx, userID, exerciseID, timeframe)
	if trainingImpact.OptimizationRecommended {
		fmt.Printf("Training optimization recommended for user %d exercise %d: %s\n",
			userID, exerciseID, trainingImpact.Recommendation)
	}

	return progression, nil
}

func (s *performanceAnalyticsServiceImpl) DetectPerformancePlateau(ctx context.Context, userID int, exerciseID int) (*types.PlateauDetection, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return nil, fmt.Errorf("invalid exercise ID")
	}

	plateau, err := s.repo.PerformanceAnalytics().DetectPerformancePlateau(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to detect performance plateau: %w", err)
	}

	advancedAnalysis := s.applyAdvancedPlateauDetection(ctx, userID, exerciseID)

	if advancedAnalysis.ConsecutiveWeeksNoProgress >= 3 {
		plateau.PlateauDetected = true
		plateau.PlateauDuration = advancedAnalysis.ConsecutiveWeeksNoProgress * 7

		plateau.Recommendation = s.generatePlateauBreakingRecommendation(advancedAnalysis)
	}

	if advancedAnalysis.PerformanceDeclineDetected {
		plateau.PlateauDetected = true
		plateau.Recommendation = "Performance decline detected - consider deload week and form review"
	}

	if advancedAnalysis.VolumeStagnation || advancedAnalysis.IntensityStagnation {
		plateau.PlateauDetected = true
		if plateau.Recommendation == "" {
			plateau.Recommendation = "Training volume/intensity stagnation - increase progressive overload"
		}
	}

	return plateau, nil
}

func (s *performanceAnalyticsServiceImpl) PredictGoalAchievement(ctx context.Context, userID int, goalID int) (*types.GoalPrediction, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if goalID <= 0 {
		return nil, fmt.Errorf("invalid goal ID")
	}

	prediction, err := s.repo.PerformanceAnalytics().PredictGoalAchievement(ctx, userID, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to predict goal achievement: %w", err)
	}

	predictionFactors := s.gatherGoalPredictionFactors(ctx, userID, goalID)

	progressionAnalysis := s.analyzeCurrentProgressionRate(predictionFactors)
	prediction.ProbabilityOfSuccess *= progressionAnalysis.ProgressionMultiplier

	historicalPatterns := s.analyzeHistoricalPerformancePatterns(ctx, userID)
	prediction.Confidence *= historicalPatterns.ConsistencyFactor

	consistencyFactor := s.calculateTrainingConsistency(ctx, userID)
	prediction.ProbabilityOfSuccess *= consistencyFactor.AdjustmentMultiplier

	timelineAdjustment := s.calculateRealisticTimeline(predictionFactors, historicalPatterns)
	prediction.EstimatedTime = int(float64(prediction.EstimatedTime) * timelineAdjustment.TimeMultiplier)

	confidenceInterval := s.calculateConfidenceInterval(predictionFactors, prediction.ProbabilityOfSuccess)
	prediction.Confidence = math.Min(prediction.Confidence*confidenceInterval.ConfidenceFactor, 1.0)

	fmt.Printf("Goal prediction for user %d goal %d: %.1f%% success probability, %d days estimated\n",
		userID, goalID, prediction.ProbabilityOfSuccess*100, prediction.EstimatedTime)

	return prediction, nil
}

func (s *performanceAnalyticsServiceImpl) CalculateTrainingVolume(ctx context.Context, userID int, weekStart time.Time) (*types.TrainingVolume, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	volume, err := s.repo.PerformanceAnalytics().CalculateTrainingVolume(ctx, userID, weekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate training volume: %w", err)
	}

	volumeAnalysis := s.calculateVolumeAnalysisAdvanced(ctx, userID, weekStart, volume)

	volume.VolumeLoad = s.calculateAdvancedVolumeLoad(volumeAnalysis)

	muscleGroupDistribution := s.analyzeMuscleGroupVolumeDistribution(ctx, userID, weekStart)
	if muscleGroupDistribution.ImbalanceDetected {
		fmt.Printf("Volume imbalance detected for user %d week %v: %s\n",
			userID, weekStart, muscleGroupDistribution.RecommendedAdjustment)
	}

	progressionCheck := s.checkVolumeProgressionSafety(ctx, userID, weekStart, volume)
	if progressionCheck.ExceedsSafetyLimits {
		fmt.Printf("Warning: Volume progression exceeds safety limits for user %d: %s\n",
			userID, progressionCheck.WarningMessage)
	}

	if progressionCheck.RecommendedVolumeAdjustment != 0 {
		adjustedVolume := float64(volume.TotalSets) * (1.0 + progressionCheck.RecommendedVolumeAdjustment)
		volume.TotalSets = int(math.Min(adjustedVolume, float64(volume.TotalSets)*1.10))
	}

	return volume, nil
}

func (s *performanceAnalyticsServiceImpl) TrackIntensityProgression(ctx context.Context, userID int, exerciseID int) (*types.IntensityProgression, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	if exerciseID <= 0 {
		return nil, fmt.Errorf("invalid exercise ID")
	}

	intensity, err := s.repo.PerformanceAnalytics().TrackIntensityProgression(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to track intensity progression: %w", err)
	}

	intensityAnalysis := s.analyzeIntensityProgressionAdvanced(ctx, userID, exerciseID, intensity)

	intensityZones := s.calculateIntensityZoneDistribution(intensityAnalysis)
	if intensityZones.RecommendedAdjustment != "" {
		fmt.Printf("Intensity zone adjustment recommended for user %d exercise %d: %s\n",
			userID, exerciseID, intensityZones.RecommendedAdjustment)
	}

	optimalAdjustment := s.calculateOptimalIntensityAdjustment(intensityAnalysis)

	if optimalAdjustment.RecommendedIncrease > 0 {
		safeIncrease := math.Min(optimalAdjustment.RecommendedIncrease, 0.05)
		intensity.RecommendedNext = intensity.CurrentIntensity * (1.0 + safeIncrease)
	}

	return intensity, nil
}

func (s *performanceAnalyticsServiceImpl) GetOptimalTrainingLoad(ctx context.Context, userID int) (*types.OptimalLoad, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	load, err := s.repo.PerformanceAnalytics().GetOptimalTrainingLoad(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get optimal training load: %w", err)
	}

	loadOptimization := s.calculateOptimalLoadAdvanced(ctx, userID)

	recoveryCapacity := s.assessUserRecoveryCapacity(ctx, userID)
	if recoveryCapacity.AdjustmentNeeded {
		load.RecommendedSets = int(float64(load.RecommendedSets) * recoveryCapacity.VolumeMultiplier)
		load.VolumeTarget = load.VolumeTarget * recoveryCapacity.VolumeMultiplier
	}

	performanceResponse := s.analyzeHistoricalPerformanceResponse(ctx, userID)
	if performanceResponse.OptimizationRecommended {
		load.RecommendedReps = performanceResponse.OptimalRepRange
		load.IntensityRange = performanceResponse.OptimalIntensityRange
	}

	personalization := s.calculatePersonalizedRecommendations(loadOptimization, recoveryCapacity, performanceResponse)
	load.VolumeTarget = personalization.OptimalVolumeTarget

	stressBalance := s.calculateTrainingStressBalance(ctx, userID)
	if stressBalance.AdjustmentRequired {
		fmt.Printf("Training stress adjustment required for user %d: %s\n", userID, stressBalance.Recommendation)
	}

	if conflicts := s.detectTrainingLoadConflicts(load, recoveryCapacity, stressBalance); len(conflicts) > 0 {
		resolvedLoad := s.resolveTrainingLoadConflicts(load, conflicts)
		load = resolvedLoad
	}

	return load, nil
}

type StrengthProgressionAnalysis struct {
	HasSufficientData bool
	DataPoints        []ProgressionDataPoint
	TrendConfidence   float64
	ProjectedGains    float64
}

type ProgressionDataPoint struct {
	Date       time.Time
	Value      float64
	Confidence float64
}

type PlateauAnalysis struct {
	ConsecutiveWeeksNoProgress int
	PerformanceDeclineDetected bool
	VolumeStagnation           bool
	IntensityStagnation        bool
}

type GoalPredictionFactors struct {
	CurrentProgress float64
	TargetValue     float64
	TimeRemaining   int
	UserLevel       types.FitnessLevel
	GoalType        types.FitnessGoal
}

type ProgressionRateAnalysis struct {
	ProgressionMultiplier float64
	TrendDirection        string
}

type HistoricalPatterns struct {
	ConsistencyFactor  float64
	PerformancePattern string
}

type ConsistencyFactor struct {
	AdjustmentMultiplier float64
	ConsistencyScore     float64
}

type TimelineAdjustment struct {
	TimeMultiplier float64
	Confidence     float64
}

type ConfidenceInterval struct {
	ConfidenceFactor float64
	LowerBound       float64
	UpperBound       float64
}

func (s *performanceAnalyticsServiceImpl) analyzeStrengthProgressionAdvanced(ctx context.Context, userID int, exerciseID int, timeframe int, progression *types.StrengthProgression) *StrengthProgressionAnalysis {
	analysis := &StrengthProgressionAnalysis{
		HasSufficientData: progression.StartingMax > 0 && progression.CurrentMax > 0,
		DataPoints:        []ProgressionDataPoint{},
		TrendConfidence:   0.8,
	}

	if analysis.HasSufficientData {
		analysis.DataPoints = append(analysis.DataPoints, ProgressionDataPoint{
			Date:       time.Now().AddDate(0, 0, -timeframe),
			Value:      progression.StartingMax,
			Confidence: 0.9,
		})
		analysis.DataPoints = append(analysis.DataPoints, ProgressionDataPoint{
			Date:       time.Now(),
			Value:      progression.CurrentMax,
			Confidence: 0.95,
		})
	}

	return analysis
}

func (s *performanceAnalyticsServiceImpl) calculateStatisticalProgressionRate(dataPoints []ProgressionDataPoint) float64 {
	if len(dataPoints) < 2 {
		return 0
	}

	first := dataPoints[0]
	last := dataPoints[len(dataPoints)-1]

	if first.Value == 0 {
		return 0
	}

	return ((last.Value - first.Value) / first.Value) * 100
}

func (s *performanceAnalyticsServiceImpl) determineTrendWithConfidence(dataPoints []ProgressionDataPoint) string {
	if len(dataPoints) < 2 {
		return "insufficient_data"
	}

	rate := s.calculateStatisticalProgressionRate(dataPoints)

	if rate > 5 {
		return "strong_increasing"
	} else if rate > 1 {
		return "increasing"
	} else if rate < -5 {
		return "declining"
	} else if rate < -1 {
		return "decreasing"
	}

	return "stable"
}

func (s *performanceAnalyticsServiceImpl) identifyStrengthImbalances(ctx context.Context, userID int, exerciseID int) []string {
	imbalances := []string{}
	return imbalances
}

func (s *performanceAnalyticsServiceImpl) predictFutureStrengthGains(analysis *StrengthProgressionAnalysis) map[string]interface{} {
	if !analysis.HasSufficientData {
		return nil
	}

	prediction := make(map[string]interface{})

	currentRate := s.calculateStatisticalProgressionRate(analysis.DataPoints)
	prediction["projected_monthly_gain"] = currentRate * 4
	prediction["confidence"] = analysis.TrendConfidence

	return prediction
}

func (s *performanceAnalyticsServiceImpl) analyzeTrainingFrequencyImpact(ctx context.Context, userID int, exerciseID int, timeframe int) struct {
	OptimizationRecommended bool
	Recommendation          string
} {
	return struct {
		OptimizationRecommended bool
		Recommendation          string
	}{
		OptimizationRecommended: false,
		Recommendation:          "Current frequency is appropriate",
	}
}

func (s *performanceAnalyticsServiceImpl) applyAdvancedPlateauDetection(ctx context.Context, userID int, exerciseID int) *PlateauAnalysis {
	analysis := &PlateauAnalysis{}

	analysis.ConsecutiveWeeksNoProgress = 2
	analysis.PerformanceDeclineDetected = false
	analysis.VolumeStagnation = false
	analysis.IntensityStagnation = false

	return analysis
}

func (s *performanceAnalyticsServiceImpl) generatePlateauBreakingRecommendation(analysis *PlateauAnalysis) string {
	if analysis.ConsecutiveWeeksNoProgress >= 4 {
		return "Extended plateau detected - implement deload week and exercise variation"
	} else if analysis.ConsecutiveWeeksNoProgress >= 3 {
		return "Plateau detected - try increasing rest time, adjusting rep ranges, or adding volume"
	}

	if analysis.VolumeStagnation {
		return "Volume stagnation detected - gradually increase training volume"
	}

	if analysis.IntensityStagnation {
		return "Intensity stagnation detected - implement progressive overload"
	}

	return "Continue monitoring performance trends"
}

func (s *performanceAnalyticsServiceImpl) gatherGoalPredictionFactors(ctx context.Context, userID int, goalID int) *GoalPredictionFactors {
	return &GoalPredictionFactors{
		CurrentProgress: 50.0,
		TargetValue:     100.0,
		TimeRemaining:   30,
		UserLevel:       types.LevelIntermediate,
		GoalType:        types.GoalStrength,
	}
}

func (s *performanceAnalyticsServiceImpl) analyzeCurrentProgressionRate(factors *GoalPredictionFactors) *ProgressionRateAnalysis {
	weeklyProgressRate := factors.CurrentProgress / 12

	multiplier := 1.0
	if weeklyProgressRate > 10 {
		multiplier = 1.2
	} else if weeklyProgressRate < 5 {
		multiplier = 0.8
	}

	return &ProgressionRateAnalysis{
		ProgressionMultiplier: multiplier,
		TrendDirection:        "stable",
	}
}

func (s *performanceAnalyticsServiceImpl) analyzeHistoricalPerformancePatterns(ctx context.Context, userID int) *HistoricalPatterns {
	return &HistoricalPatterns{
		ConsistencyFactor:  0.9,
		PerformancePattern: "steady_progression",
	}
}

func (s *performanceAnalyticsServiceImpl) calculateTrainingConsistency(ctx context.Context, userID int) *ConsistencyFactor {
	return &ConsistencyFactor{
		AdjustmentMultiplier: 1.0,
		ConsistencyScore:     0.85,
	}
}

func (s *performanceAnalyticsServiceImpl) calculateRealisticTimeline(factors *GoalPredictionFactors, patterns *HistoricalPatterns) *TimelineAdjustment {
	multiplier := 1.0

	switch factors.GoalType {
	case types.GoalStrength:
		multiplier = 1.2
	case types.GoalFatLoss:
		multiplier = 0.9
	case types.GoalMuscleGain:
		multiplier = 1.3
	}

	switch factors.UserLevel {
	case types.LevelBeginner:
		multiplier *= 0.8
	case types.LevelAdvanced:
		multiplier *= 1.4
	}

	return &TimelineAdjustment{
		TimeMultiplier: multiplier,
		Confidence:     patterns.ConsistencyFactor,
	}
}

func (s *performanceAnalyticsServiceImpl) calculateConfidenceInterval(factors *GoalPredictionFactors, probability float64) *ConfidenceInterval {
	confidenceFactor := 0.9

	if factors.TimeRemaining < 7 {
		confidenceFactor = 0.95
	} else if factors.TimeRemaining > 90 {
		confidenceFactor = 0.7
	}

	return &ConfidenceInterval{
		ConfidenceFactor: confidenceFactor,
		LowerBound:       probability * 0.8,
		UpperBound:       probability * 1.2,
	}
}

func (s *performanceAnalyticsServiceImpl) calculateVolumeAnalysisAdvanced(ctx context.Context, userID int, weekStart time.Time, volume *types.TrainingVolume) map[string]interface{} {
	analysis := make(map[string]interface{})

	analysis["total_volume"] = volume.TotalSets
	analysis["volume_trend"] = "increasing"
	analysis["muscle_group_balance"] = "balanced"

	return analysis
}

func (s *performanceAnalyticsServiceImpl) calculateAdvancedVolumeLoad(analysis map[string]interface{}) float64 {
	if totalVolume, ok := analysis["total_volume"].(int); ok {
		return float64(totalVolume) * 100
	}
	return 0
}

func (s *performanceAnalyticsServiceImpl) analyzeMuscleGroupVolumeDistribution(ctx context.Context, userID int, weekStart time.Time) struct {
	ImbalanceDetected     bool
	RecommendedAdjustment string
} {
	return struct {
		ImbalanceDetected     bool
		RecommendedAdjustment string
	}{
		ImbalanceDetected:     false,
		RecommendedAdjustment: "Volume distribution is balanced",
	}
}

func (s *performanceAnalyticsServiceImpl) checkVolumeProgressionSafety(ctx context.Context, userID int, weekStart time.Time, volume *types.TrainingVolume) struct {
	ExceedsSafetyLimits         bool
	WarningMessage              string
	RecommendedVolumeAdjustment float64
} {
	return struct {
		ExceedsSafetyLimits         bool
		WarningMessage              string
		RecommendedVolumeAdjustment float64
	}{
		ExceedsSafetyLimits:         false,
		WarningMessage:              "",
		RecommendedVolumeAdjustment: 0.0,
	}
}

func (s *performanceAnalyticsServiceImpl) analyzeIntensityProgressionAdvanced(ctx context.Context, userID int, exerciseID int, intensity *types.IntensityProgression) map[string]interface{} {
	analysis := make(map[string]interface{})

	analysis["current_intensity"] = intensity.CurrentIntensity
	analysis["progression_rate"] = intensity.ProgressionRate
	analysis["trend"] = "stable"

	return analysis
}

func (s *performanceAnalyticsServiceImpl) calculateIntensityZoneDistribution(analysis map[string]interface{}) struct {
	RecommendedAdjustment string
} {
	return struct {
		RecommendedAdjustment string
	}{
		RecommendedAdjustment: "",
	}
}

func (s *performanceAnalyticsServiceImpl) calculateOptimalIntensityAdjustment(analysis map[string]interface{}) struct {
	RecommendedIncrease float64
} {
	return struct {
		RecommendedIncrease float64
	}{
		RecommendedIncrease: 0.025,
	}
}

func (s *performanceAnalyticsServiceImpl) calculateOptimalLoadAdvanced(ctx context.Context, userID int) map[string]interface{} {
	optimization := make(map[string]interface{})

	optimization["recommended_volume"] = 15
	optimization["recommended_intensity"] = 0.75
	optimization["recovery_factor"] = 1.0

	return optimization
}

func (s *performanceAnalyticsServiceImpl) assessUserRecoveryCapacity(ctx context.Context, userID int) struct {
	AdjustmentNeeded bool
	VolumeMultiplier float64
} {
	return struct {
		AdjustmentNeeded bool
		VolumeMultiplier float64
	}{
		AdjustmentNeeded: false,
		VolumeMultiplier: 1.0,
	}
}

func (s *performanceAnalyticsServiceImpl) analyzeHistoricalPerformanceResponse(ctx context.Context, userID int) struct {
	OptimizationRecommended bool
	OptimalRepRange         int
	OptimalIntensityRange   string
} {
	return struct {
		OptimizationRecommended bool
		OptimalRepRange         int
		OptimalIntensityRange   string
	}{
		OptimizationRecommended: false,
		OptimalRepRange:         8,
		OptimalIntensityRange:   "70-85%",
	}
}

func (s *performanceAnalyticsServiceImpl) calculatePersonalizedRecommendations(optimization map[string]interface{}, recovery struct {
	AdjustmentNeeded bool
	VolumeMultiplier float64
}, performance struct {
	OptimizationRecommended bool
	OptimalRepRange         int
	OptimalIntensityRange   string
}) struct {
	OptimalVolumeTarget float64
} {
	return struct {
		OptimalVolumeTarget float64
	}{
		OptimalVolumeTarget: 4000.0,
	}
}

func (s *performanceAnalyticsServiceImpl) calculateTrainingStressBalance(ctx context.Context, userID int) struct {
	AdjustmentRequired bool
	Recommendation     string
} {
	return struct {
		AdjustmentRequired bool
		Recommendation     string
	}{
		AdjustmentRequired: false,
		Recommendation:     "Training stress is balanced",
	}
}

func (s *performanceAnalyticsServiceImpl) detectTrainingLoadConflicts(load *types.OptimalLoad, recovery struct {
	AdjustmentNeeded bool
	VolumeMultiplier float64
}, stress struct {
	AdjustmentRequired bool
	Recommendation     string
}) []string {
	conflicts := []string{}

	if recovery.AdjustmentNeeded && stress.AdjustmentRequired {
		conflicts = append(conflicts, "recovery_stress_conflict")
	}

	return conflicts
}

func (s *performanceAnalyticsServiceImpl) resolveTrainingLoadConflicts(load *types.OptimalLoad, conflicts []string) *types.OptimalLoad {
	resolvedLoad := *load

	for _, conflict := range conflicts {
		switch conflict {
		case "recovery_stress_conflict":
			resolvedLoad.RecommendedSets = int(float64(resolvedLoad.RecommendedSets) * 0.9)
			resolvedLoad.VolumeTarget = resolvedLoad.VolumeTarget * 0.9
		}
	}

	return &resolvedLoad
}
