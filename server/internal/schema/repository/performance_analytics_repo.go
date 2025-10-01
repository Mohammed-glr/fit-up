package repository

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

// =============================================================================
// PERFORMANCE ANALYTICS REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) CalculateStrengthProgression(ctx context.Context, userID int, exerciseID int, timeframe int) (*types.StrengthProgression, error) {
	// Get 1RM estimates over the timeframe
	q := `
		SELECT estimated_max, estimate_date
		FROM one_rep_max_estimates
		WHERE user_id = $1 AND exercise_id = $2
		AND estimate_date >= CURRENT_DATE - INTERVAL '%d days'
		ORDER BY estimate_date ASC
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID, timeframe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estimates []float64
	var dates []time.Time

	for rows.Next() {
		var estimate float64
		var date time.Time
		err := rows.Scan(&estimate, &date)
		if err != nil {
			return nil, err
		}
		estimates = append(estimates, estimate)
		dates = append(dates, date)
	}

	if len(estimates) < 2 {
		// Not enough data for progression analysis
		return &types.StrengthProgression{
			ExerciseID:      exerciseID,
			StartingMax:     0,
			CurrentMax:      0,
			ProgressionRate: 0,
			Trend:           "insufficient_data",
		}, nil
	}

	startingMax := estimates[0]
	currentMax := estimates[len(estimates)-1]

	// Calculate progression rate (% change per week)
	totalDays := dates[len(dates)-1].Sub(dates[0]).Hours() / 24
	weeks := totalDays / 7

	var progressionRate float64
	if weeks > 0 && startingMax > 0 {
		totalChange := ((currentMax - startingMax) / startingMax) * 100
		progressionRate = totalChange / weeks
	}

	// Determine trend
	trend := "stable"
	if progressionRate > 1.0 {
		trend = "increasing"
	} else if progressionRate < -1.0 {
		trend = "decreasing"
	}

	return &types.StrengthProgression{
		ExerciseID:      exerciseID,
		StartingMax:     startingMax,
		CurrentMax:      currentMax,
		ProgressionRate: progressionRate,
		Trend:           trend,
	}, nil
}

func (s *Store) DetectPerformancePlateau(ctx context.Context, userID int, exerciseID int) (*types.PlateauDetection, error) {
	// Get recent performance data (last 4 weeks)
	q := `
		SELECT pl.weight_used, pl.reps_completed, pl.date
		FROM progress_logs pl
		WHERE pl.user_id = $1 AND pl.exercise_id = $2
		AND pl.date >= CURRENT_DATE - INTERVAL '28 days'
		AND pl.weight_used IS NOT NULL AND pl.reps_completed IS NOT NULL
		ORDER BY pl.date DESC
		LIMIT 20
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []struct {
		weight float64
		reps   int
		date   time.Time
	}

	for rows.Next() {
		var perf struct {
			weight float64
			reps   int
			date   time.Time
		}
		err := rows.Scan(&perf.weight, &perf.reps, &perf.date)
		if err != nil {
			return nil, err
		}
		performances = append(performances, perf)
	}

	if len(performances) < 4 {
		return &types.PlateauDetection{
			ExerciseID:      exerciseID,
			PlateauDetected: false,
			PlateauDuration: 0,
			Recommendation:  "Insufficient data for plateau analysis",
		}, nil
	}

	// Calculate volume (weight * reps) for each performance
	volumes := make([]float64, len(performances))
	dates := make([]time.Time, len(performances))

	for i, perf := range performances {
		volumes[i] = perf.weight * float64(perf.reps)
		dates[i] = perf.date
	}

	// Check for plateau - no improvement in volume over multiple sessions
	plateauDetected := true
	bestVolume := volumes[0]

	// Check if any recent performance exceeds the best from 2+ weeks ago
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	recentImprovement := false

	for i, vol := range volumes {
		if dates[i].After(twoWeeksAgo) && vol > bestVolume*1.05 { // 5% improvement threshold
			recentImprovement = true
			break
		}
	}

	plateauDetected = !recentImprovement

	var plateauDuration int
	var recommendation string

	if plateauDetected {
		// Calculate plateau duration
		plateauDuration = int(time.Since(dates[len(dates)-1]).Hours() / 24)

		if plateauDuration >= 21 {
			recommendation = "Extended plateau detected. Consider deload week or exercise variation."
		} else if plateauDuration >= 14 {
			recommendation = "Plateau detected. Try increasing rest time or adjusting rep ranges."
		} else {
			recommendation = "Minor plateau detected. Monitor for another week."
		}
	} else {
		recommendation = "Performance is progressing normally."
	}

	return &types.PlateauDetection{
		ExerciseID:      exerciseID,
		PlateauDetected: plateauDetected,
		PlateauDuration: plateauDuration,
		Recommendation:  recommendation,
	}, nil
}

func (s *Store) PredictGoalAchievement(ctx context.Context, userID int, goalID int) (*types.GoalPrediction, error) {
	// Get goal details
	goalQuery := `
		SELECT goal_type, target_value, current_value, target_date, created_at
		FROM fitness_goals
		WHERE goal_id = $1 AND user_id = $2
	`

	var goalType types.FitnessGoal
	var targetValue, currentValue float64
	var targetDate, createdAt time.Time

	err := s.db.QueryRow(ctx, goalQuery, goalID, userID).Scan(
		&goalType,
		&targetValue,
		&currentValue,
		&targetDate,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	timeElapsed := time.Since(createdAt).Hours() / 24        // days
	timeRemaining := targetDate.Sub(time.Now()).Hours() / 24 // days

	if timeRemaining <= 0 {
		return &types.GoalPrediction{
			GoalID:               goalID,
			ProbabilityOfSuccess: 0.0,
			EstimatedTime:        0,
			Confidence:           1.0,
		}, nil
	}

	if timeElapsed > 0 {
		_ = (currentValue - 0) / timeElapsed 
	}

	var currentProgress float64
	if targetValue != 0 {
		currentProgress = (currentValue / targetValue) * 100
	}

	var probabilityOfSuccess float64
	var estimatedDays int
	var confidence float64

	if currentProgress >= 90 {
		probabilityOfSuccess = 0.95
		estimatedDays = int(timeRemaining * 0.5) // Should finish early
		confidence = 0.9
	} else if currentProgress >= 70 {
		probabilityOfSuccess = 0.8
		estimatedDays = int(timeRemaining * 0.8)
		confidence = 0.8
	} else if currentProgress >= 50 {
		probabilityOfSuccess = 0.6
		estimatedDays = int(timeRemaining * 1.2)
		confidence = 0.7
	} else if currentProgress >= 30 {
		probabilityOfSuccess = 0.4
		estimatedDays = int(timeRemaining * 1.5)
		confidence = 0.6
	} else {
		probabilityOfSuccess = 0.2
		estimatedDays = int(timeRemaining * 2.0)
		confidence = 0.5
	}

	switch goalType {
	case types.GoalStrength:
		probabilityOfSuccess *= 0.9
		estimatedDays = int(float64(estimatedDays) * 1.1)
	case types.GoalFatLoss:
		confidence *= 0.8
	case types.GoalMuscleGain:
		probabilityOfSuccess *= 0.95
	}

	return &types.GoalPrediction{
		GoalID:               goalID,
		ProbabilityOfSuccess: probabilityOfSuccess,
		EstimatedTime:        estimatedDays,
		Confidence:           confidence,
	}, nil
}

func (s *Store) CalculateTrainingVolume(ctx context.Context, userID int, weekStart time.Time) (*types.TrainingVolume, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)

	q := `
		SELECT 
			COUNT(DISTINCT ws.session_id) as sessions,
			COALESCE(SUM(sep.sets_completed), 0) as total_sets,
			COALESCE(SUM(sep.best_reps * sep.sets_completed), 0) as total_reps,
			COALESCE(SUM(sep.best_weight * sep.best_reps * sep.sets_completed), 0) as total_weight,
			COALESCE(AVG(sep.rpe), 0) as avg_intensity
		FROM workout_sessions ws
		LEFT JOIN session_exercise_performances sep ON ws.session_id = sep.session_id
		WHERE ws.user_id = $1 
		AND ws.start_time >= $2 
		AND ws.start_time < $3
		AND ws.status = 'completed'
	`

	var sessions int
	var totalSets, totalReps int
	var totalWeight, avgIntensity float64

	err := s.db.QueryRow(ctx, q, userID, weekStart, weekEnd).Scan(
		&sessions,
		&totalSets,
		&totalReps,
		&totalWeight,
		&avgIntensity,
	)
	if err != nil {
		return nil, err
	}

	volumeLoad := totalWeight // This could be more sophisticated

	return &types.TrainingVolume{
		WeekStart:    weekStart,
		TotalSets:    totalSets,
		TotalReps:    totalReps,
		TotalWeight:  totalWeight,
		VolumeLoad:   volumeLoad,
		IntensityAvg: avgIntensity,
	}, nil
}

func (s *Store) TrackIntensityProgression(ctx context.Context, userID int, exerciseID int) (*types.IntensityProgression, error) {
	q := `
		SELECT pl.weight_used, pl.reps_completed, pl.date
		FROM progress_logs pl
		WHERE pl.user_id = $1 AND pl.exercise_id = $2
		AND pl.weight_used IS NOT NULL AND pl.reps_completed IS NOT NULL
		AND pl.date >= CURRENT_DATE - INTERVAL '60 days'
		ORDER BY pl.date ASC
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []struct {
		weight float64
		reps   int
		date   time.Time
	}

	for rows.Next() {
		var perf struct {
			weight float64
			reps   int
			date   time.Time
		}
		err := rows.Scan(&perf.weight, &perf.reps, &perf.date)
		if err != nil {
			return nil, err
		}
		performances = append(performances, perf)
	}

	if len(performances) < 2 {
		return &types.IntensityProgression{
			ExerciseID:        exerciseID,
			BaselineIntensity: 0,
			CurrentIntensity:  0,
			ProgressionRate:   0,
			RecommendedNext:   0,
		}, nil
	}

	baseline := performances[0].weight
	current := performances[len(performances)-1].weight

	progressionRate := 0.0
	if baseline > 0 {
		progressionRate = ((current - baseline) / baseline) * 100
	}

	recommendedNext := current * 1.025
	if progressionRate > 10 { 
		recommendedNext = current * 1.05
	}

	return &types.IntensityProgression{
		ExerciseID:        exerciseID,
		BaselineIntensity: baseline,
		CurrentIntensity:  current,
		ProgressionRate:   progressionRate,
		RecommendedNext:   recommendedNext,
	}, nil
}

func (s *Store) GetOptimalTrainingLoad(ctx context.Context, userID int) (*types.OptimalLoad, error) {
	userQuery := `
		SELECT level, goal
		FROM users
		WHERE user_id = $1
	`

	var level types.FitnessLevel
	var goal types.FitnessGoal

	err := s.db.QueryRow(ctx, userQuery, userID).Scan(&level, &goal)
	if err != nil {
		return nil, err
	}

	recentVolumeQuery := `
		SELECT COALESCE(AVG(total_volume), 0)
		FROM workout_sessions
		WHERE user_id = $1 
		AND start_time >= CURRENT_DATE - INTERVAL '14 days'
		AND status = 'completed'
	`

	var recentAvgVolume float64
	err = s.db.QueryRow(ctx, recentVolumeQuery, userID).Scan(&recentAvgVolume)
	if err != nil {
		recentAvgVolume = 0
	}

	var recommendedSets, recommendedReps int
	var intensityRange string
	var volumeTarget float64

	switch level {
	case types.LevelBeginner:
		recommendedSets = 8
		recommendedReps = 10
		intensityRange = "60-75%"
		volumeTarget = 2000
	case types.LevelIntermediate:
		recommendedSets = 12
		recommendedReps = 8
		intensityRange = "70-85%"
		volumeTarget = 4000
	case types.LevelAdvanced:
		recommendedSets = 16
		recommendedReps = 6
		intensityRange = "80-95%"
		volumeTarget = 6000
	}

	switch goal {
	case types.GoalStrength:
		recommendedReps = int(float64(recommendedReps) * 0.75)
		intensityRange = "85-95%"
	case types.GoalMuscleGain:
		recommendedSets = int(float64(recommendedSets) * 1.2)
		intensityRange = "65-80%"
	case types.GoalFatLoss:
		recommendedSets = int(float64(recommendedSets) * 1.1)
		recommendedReps = int(float64(recommendedReps) * 1.2)
		intensityRange = "60-75%"
	}

	return &types.OptimalLoad{
		UserID:          userID,
		RecommendedSets: recommendedSets,
		RecommendedReps: recommendedReps,
		IntensityRange:  intensityRange,
		VolumeTarget:    volumeTarget,
	}, nil
}
