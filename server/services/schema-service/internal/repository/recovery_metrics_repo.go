package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// =============================================================================
// RECOVERY METRICS REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) LogRecoveryMetrics(ctx context.Context, userID int, metrics *types.RecoveryMetrics) error {
	q := `
		INSERT INTO recovery_metrics (user_id, date, sleep_hours, sleep_quality, stress_level, energy_level, soreness)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, date) 
		DO UPDATE SET 
			sleep_hours = EXCLUDED.sleep_hours,
			sleep_quality = EXCLUDED.sleep_quality,
			stress_level = EXCLUDED.stress_level,
			energy_level = EXCLUDED.energy_level,
			soreness = EXCLUDED.soreness
	`

	_, err := s.db.Exec(ctx, q,
		userID,
		metrics.Date,
		metrics.SleepHours,
		metrics.SleepQuality,
		metrics.StressLevel,
		metrics.EnergyLevel,
		metrics.Soreness,
	)

	return err
}

func (s *Store) GetRecoveryStatus(ctx context.Context, userID int) (*types.RecoveryStatus, error) {
	// Get the most recent recovery metrics (last 3 days)
	q := `
		SELECT 
			AVG(sleep_hours) as avg_sleep,
			AVG(sleep_quality) as avg_sleep_quality,
			AVG(stress_level) as avg_stress,
			AVG(energy_level) as avg_energy,
			AVG(soreness) as avg_soreness
		FROM recovery_metrics
		WHERE user_id = $1 
		AND date >= CURRENT_DATE - INTERVAL '3 days'
	`

	var avgSleep, avgSleepQuality, avgStress, avgEnergy, avgSoreness float64
	err := s.db.QueryRow(ctx, q, userID).Scan(
		&avgSleep,
		&avgSleepQuality,
		&avgStress,
		&avgEnergy,
		&avgSoreness,
	)

	if err != nil {
		// No recent data, return default moderate recovery status
		return &types.RecoveryStatus{
			UserID:               userID,
			RecoveryScore:        0.7,
			Recommendation:       "Moderate recovery status - no recent data available",
			RecommendedIntensity: 0.8,
			RestDayRecommended:   false,
		}, nil
	}

	// Calculate recovery score using rule-based logic
	recoveryScore := s.calculateRecoveryScore(avgSleep, avgSleepQuality, avgStress, avgEnergy, avgSoreness)

	// Generate recommendations based on score
	recommendation, recommendedIntensity, restDayRecommended := s.generateRecoveryRecommendations(recoveryScore)

	return &types.RecoveryStatus{
		UserID:               userID,
		RecoveryScore:        recoveryScore,
		Recommendation:       recommendation,
		RecommendedIntensity: recommendedIntensity,
		RestDayRecommended:   restDayRecommended,
	}, nil
}

func (s *Store) calculateRecoveryScore(sleep, sleepQuality, stress, energy, soreness float64) float64 {
	// Rule-based recovery score calculation (0-1 scale)
	score := 0.0

	// Sleep component (40% weight)
	sleepScore := 0.0
	if sleep >= 7.5 {
		sleepScore = 1.0
	} else if sleep >= 6.5 {
		sleepScore = 0.8
	} else if sleep >= 5.5 {
		sleepScore = 0.5
	} else {
		sleepScore = 0.2
	}

	// Sleep quality component (20% weight)
	sleepQualityScore := sleepQuality / 10.0 // Assuming 1-10 scale

	// Stress component (15% weight) - lower is better
	stressScore := (10.0 - stress) / 10.0 // Invert stress (assuming 1-10 scale)

	// Energy component (15% weight)
	energyScore := energy / 10.0 // Assuming 1-10 scale

	// Soreness component (10% weight) - lower is better
	sorenessScore := (10.0 - soreness) / 10.0 // Invert soreness

	score = (sleepScore * 0.4) + (sleepQualityScore * 0.2) + (stressScore * 0.15) + (energyScore * 0.15) + (sorenessScore * 0.1)

	// Ensure score is between 0 and 1
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score
}

func (s *Store) generateRecoveryRecommendations(recoveryScore float64) (string, float64, bool) {
	if recoveryScore >= 0.8 {
		return "Excellent recovery! You're ready for high-intensity training.", 1.0, false
	} else if recoveryScore >= 0.65 {
		return "Good recovery status. Proceed with normal training intensity.", 0.9, false
	} else if recoveryScore >= 0.5 {
		return "Moderate recovery. Consider reducing training intensity slightly.", 0.75, false
	} else if recoveryScore >= 0.35 {
		return "Poor recovery status. Significantly reduce training intensity or take a rest day.", 0.5, true
	} else {
		return "Very poor recovery. Rest day strongly recommended.", 0.0, true
	}
}

func (s *Store) GetRecoveryTrend(ctx context.Context, userID int, days int) ([]types.RecoveryMetrics, error) {
	q := `
		SELECT metric_id, user_id, date, sleep_hours, sleep_quality, stress_level, energy_level, soreness
		FROM recovery_metrics
		WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '%d days'
		ORDER BY date DESC
	`

	rows, err := s.db.Query(ctx, q, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.RecoveryMetrics
	for rows.Next() {
		var metric types.RecoveryMetrics
		err := rows.Scan(
			&metric.MetricID,
			&metric.UserID,
			&metric.Date,
			&metric.SleepHours,
			&metric.SleepQuality,
			&metric.StressLevel,
			&metric.EnergyLevel,
			&metric.Soreness,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (s *Store) CalculateFatigueScore(ctx context.Context, userID int) (float64, error) {
	// Get recent training volume and recovery metrics
	volumeQuery := `
		SELECT COALESCE(SUM(total_volume), 0) as weekly_volume
		FROM workout_sessions
		WHERE user_id = $1 
		AND start_time >= CURRENT_DATE - INTERVAL '7 days'
		AND status = 'completed'
	`

	var weeklyVolume float64
	err := s.db.QueryRow(ctx, volumeQuery, userID).Scan(&weeklyVolume)
	if err != nil {
		weeklyVolume = 0
	}

	// Get recent recovery metrics
	recoveryQuery := `
		SELECT 
			AVG(sleep_hours) as avg_sleep,
			AVG(soreness) as avg_soreness,
			AVG(energy_level) as avg_energy
		FROM recovery_metrics
		WHERE user_id = $1 
		AND date >= CURRENT_DATE - INTERVAL '3 days'
	`

	var avgSleep, avgSoreness, avgEnergy float64
	err = s.db.QueryRow(ctx, recoveryQuery, userID).Scan(&avgSleep, &avgSoreness, &avgEnergy)
	if err != nil {
		// No recovery data, estimate from volume alone
		fatigueScore := weeklyVolume / 10000.0 // Rough estimate
		if fatigueScore > 1.0 {
			fatigueScore = 1.0
		}
		return fatigueScore, nil
	}

	// Calculate fatigue score using rule-based logic
	volumeComponent := weeklyVolume / 15000.0 // Normalize based on expected max volume
	if volumeComponent > 1.0 {
		volumeComponent = 1.0
	}

	sleepComponent := (8.0 - avgSleep) / 8.0 // Sleep debt component
	if sleepComponent < 0 {
		sleepComponent = 0
	}

	sorenessComponent := avgSoreness / 10.0      // Direct soreness component
	energyComponent := (10.0 - avgEnergy) / 10.0 // Inverse energy component

	fatigueScore := (volumeComponent * 0.4) + (sleepComponent * 0.3) + (sorenessComponent * 0.2) + (energyComponent * 0.1)

	if fatigueScore > 1.0 {
		fatigueScore = 1.0
	}
	if fatigueScore < 0.0 {
		fatigueScore = 0.0
	}

	return fatigueScore, nil
}

func (s *Store) RecommendRestDay(ctx context.Context, userID int) (*types.RestDayRecommendation, error) {
	fatigueScore, err := s.CalculateFatigueScore(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get recent rest days
	restDaysQuery := `
		SELECT COUNT(*)
		FROM (
			SELECT date
			FROM recovery_metrics rm
			WHERE rm.user_id = $1 
			AND rm.date >= CURRENT_DATE - INTERVAL '7 days'
			AND NOT EXISTS (
				SELECT 1 
				FROM workout_sessions ws 
				WHERE ws.user_id = rm.user_id 
				AND DATE(ws.start_time) = rm.date
				AND ws.status = 'completed'
			)
		) rest_days
	`

	var recentRestDays int
	err = s.db.QueryRow(ctx, restDaysQuery, userID).Scan(&recentRestDays)
	if err != nil {
		recentRestDays = 0
	}

	// Rule-based rest day recommendation
	recommendation := &types.RestDayRecommendation{
		Recommended: false,
		Reason:      "Current fatigue levels are manageable",
		Duration:    0,
		Activities:  []string{"light stretching", "walk"},
	}

	if fatigueScore >= 0.8 {
		recommendation.Recommended = true
		recommendation.Reason = "High fatigue detected - rest day strongly recommended"
		recommendation.Duration = 2
		recommendation.Activities = []string{"complete rest", "light stretching", "meditation"}
	} else if fatigueScore >= 0.6 {
		recommendation.Recommended = true
		recommendation.Reason = "Moderate fatigue detected - rest day recommended"
		recommendation.Duration = 1
		recommendation.Activities = []string{"active recovery", "yoga", "light walk"}
	} else if recentRestDays == 0 && fatigueScore >= 0.4 {
		recommendation.Recommended = true
		recommendation.Reason = "No recent rest days - preventive rest recommended"
		recommendation.Duration = 1
		recommendation.Activities = []string{"active recovery", "stretching", "mobility work"}
	}

	return recommendation, nil
}

func (s *Store) TrackSleepQuality(ctx context.Context, userID int, quality *types.SleepQuality) error {
	// Insert or update today's sleep quality
	q := `
		INSERT INTO recovery_metrics (user_id, date, sleep_hours, sleep_quality, stress_level, energy_level, soreness)
		VALUES ($1, CURRENT_DATE, $2, $3, 5.0, 5.0, 5.0)
		ON CONFLICT (user_id, date)
		DO UPDATE SET 
			sleep_hours = EXCLUDED.sleep_hours,
			sleep_quality = EXCLUDED.sleep_quality
	`

	_, err := s.db.Exec(ctx, q,
		userID,
		quality.Hours,
		quality.Quality,
	)

	return err
}
