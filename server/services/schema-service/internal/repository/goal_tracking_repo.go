package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// =============================================================================
// GOAL TRACKING REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) CreateFitnessGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error) {
	metadataJSON, err := json.Marshal(goal.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	q := `
		INSERT INTO fitness_goals (user_id, goal_type, target_value, current_value, target_date, is_active, created_at, metadata)
		VALUES ($1, $2, $3, 0.0, $4, true, NOW(), $5)
		RETURNING goal_id, user_id, goal_type, target_value, current_value, target_date, is_active, created_at, metadata
	`

	var result types.FitnessGoalTarget
	err = s.db.QueryRow(ctx, q,
		userID,
		goal.GoalType,
		goal.TargetValue,
		goal.TargetDate,
		metadataJSON,
	).Scan(
		&result.GoalID,
		&result.UserID,
		&result.GoalType,
		&result.TargetValue,
		&result.CurrentValue,
		&result.TargetDate,
		&result.IsActive,
		&result.CreatedAt,
		&result.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Store) UpdateGoalProgress(ctx context.Context, goalID int, progress float64) error {
	q := `
		UPDATE fitness_goals 
		SET current_value = $1
		WHERE goal_id = $2
	`

	_, err := s.db.Exec(ctx, q, progress, goalID)
	return err
}

func (s *Store) GetActiveGoals(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error) {
	q := `
		SELECT goal_id, user_id, goal_type, target_value, current_value, target_date, is_active, created_at, metadata
		FROM fitness_goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []types.FitnessGoalTarget
	for rows.Next() {
		var goal types.FitnessGoalTarget
		err := rows.Scan(
			&goal.GoalID,
			&goal.UserID,
			&goal.GoalType,
			&goal.TargetValue,
			&goal.CurrentValue,
			&goal.TargetDate,
			&goal.IsActive,
			&goal.CreatedAt,
			&goal.Metadata,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

func (s *Store) CompleteGoal(ctx context.Context, goalID int) error {
	q := `
		UPDATE fitness_goals 
		SET is_active = false, current_value = target_value
		WHERE goal_id = $1
	`

	_, err := s.db.Exec(ctx, q, goalID)
	return err
}

func (s *Store) CalculateGoalProgress(ctx context.Context, goalID int) (*types.GoalProgress, error) {
	q := `
		SELECT goal_type, target_value, current_value, target_date, created_at
		FROM fitness_goals
		WHERE goal_id = $1
	`

	var goalType types.FitnessGoal
	var targetValue, currentValue float64
	var targetDate, createdAt time.Time

	err := s.db.QueryRow(ctx, q, goalID).Scan(
		&goalType,
		&targetValue,
		&currentValue,
		&targetDate,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Calculate progress percentage
	var progressPercent float64
	if targetValue > 0 {
		progressPercent = (currentValue / targetValue) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}
	}

	// Determine if on track
	timeElapsed := time.Since(createdAt)
	totalTime := targetDate.Sub(createdAt)
	expectedProgress := 0.0

	if totalTime > 0 {
		expectedProgress = (timeElapsed.Seconds() / totalTime.Seconds()) * 100
	}

	onTrack := progressPercent >= expectedProgress*0.8 // Allow 20% tolerance

	// Estimate completion date
	var estimatedCompletion *time.Time
	if progressPercent > 0 && progressPercent < 100 {
		remainingProgress := 100 - progressPercent
		progressRate := progressPercent / (timeElapsed.Hours() / 24) // Progress per day

		if progressRate > 0 {
			daysToCompletion := remainingProgress / progressRate
			completion := time.Now().AddDate(0, 0, int(daysToCompletion))
			estimatedCompletion = &completion
		}
	}

	return &types.GoalProgress{
		GoalID:              goalID,
		ProgressPercent:     progressPercent,
		OnTrack:             onTrack,
		EstimatedCompletion: estimatedCompletion,
	}, nil
}

func (s *Store) EstimateTimeToGoal(ctx context.Context, goalID int) (*types.TimeToGoalEstimate, error) {
	progress, err := s.CalculateGoalProgress(ctx, goalID)
	if err != nil {
		return nil, err
	}

	// Get goal details for more accurate estimation
	q := `
		SELECT goal_type, target_value, current_value, target_date, created_at
		FROM fitness_goals
		WHERE goal_id = $1
	`

	var goalType types.FitnessGoal
	var targetValue, currentValue float64
	var targetDate, createdAt time.Time

	err = s.db.QueryRow(ctx, q, goalID).Scan(
		&goalType,
		&targetValue,
		&currentValue,
		&targetDate,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	timeElapsed := time.Since(createdAt).Hours() / 24 // days
	var estimatedDays int
	var confidence float64
	var assumptions []string

	if progress.ProgressPercent >= 100 {
		estimatedDays = 0
		confidence = 1.0
		assumptions = []string{"Goal already achieved"}
	} else if progress.ProgressPercent > 0 && timeElapsed > 0 {
		progressRate := progress.ProgressPercent / timeElapsed
		remainingProgress := 100 - progress.ProgressPercent
		estimatedDays = int(remainingProgress / progressRate)

		// Adjust confidence based on goal type and current progress
		confidence = 0.7
		if progress.ProgressPercent > 50 {
			confidence = 0.8
		}
		if progress.OnTrack {
			confidence += 0.1
		}

		assumptions = []string{
			"Maintains current progress rate",
			"No major setbacks or plateaus",
			"Consistent training and nutrition",
		}

		// Adjust based on goal type
		switch goalType {
		case types.GoalStrength:
			estimatedDays = int(float64(estimatedDays) * 1.2) // Strength gains slow down
			assumptions = append(assumptions, "Progressive overload continues")
		case types.GoalFatLoss:
			confidence *= 0.9 // Fat loss can be variable
			assumptions = append(assumptions, "Maintains caloric deficit")
		case types.GoalMuscleGain:
			estimatedDays = int(float64(estimatedDays) * 1.1) // Muscle gain is steady but slow
			assumptions = append(assumptions, "Adequate protein intake and recovery")
		}
	} else {
		// No progress yet or insufficient data
		estimatedDays = int(targetDate.Sub(time.Now()).Hours() / 24)
		confidence = 0.5
		assumptions = []string{
			"Based on target date only",
			"No historical progress data available",
		}
	}

	return &types.TimeToGoalEstimate{
		GoalID:        goalID,
		EstimatedDays: estimatedDays,
		Confidence:    confidence,
		Assumptions:   assumptions,
	}, nil
}

func (s *Store) SuggestGoalAdjustments(ctx context.Context, userID int) ([]types.GoalAdjustment, error) {
	// Get all active goals for the user
	goals, err := s.GetActiveGoals(ctx, userID)
	if err != nil {
		return nil, err
	}

	var adjustments []types.GoalAdjustment

	for _, goal := range goals {
		progress, err := s.CalculateGoalProgress(ctx, goal.GoalID)
		if err != nil {
			continue
		}

		// Check if goal needs adjustment
		timeUntilTarget := goal.TargetDate.Sub(time.Now()).Hours() / 24 // days

		// Suggest adjustments based on progress and time remaining
		if timeUntilTarget <= 0 && progress.ProgressPercent < 100 {
			// Past due date
			adjustments = append(adjustments, types.GoalAdjustment{
				GoalID:             goal.GoalID,
				RecommendationType: "extend_deadline",
				Adjustment:         "Extend target date by 30-60 days",
				Reason:             "Goal deadline has passed with incomplete progress",
			})
		} else if !progress.OnTrack && timeUntilTarget > 0 {
			// Behind schedule
			if progress.ProgressPercent < 25 && timeUntilTarget < 30 {
				adjustments = append(adjustments, types.GoalAdjustment{
					GoalID:             goal.GoalID,
					RecommendationType: "reduce_target",
					Adjustment:         "Reduce target value by 20-30%",
					Reason:             "Current progress rate suggests target may be too ambitious",
				})
			} else {
				adjustments = append(adjustments, types.GoalAdjustment{
					GoalID:             goal.GoalID,
					RecommendationType: "increase_intensity",
					Adjustment:         "Increase training frequency or intensity",
					Reason:             "Need to accelerate progress to meet target date",
				})
			}
		} else if progress.ProgressPercent > 80 && timeUntilTarget > 60 {
			// Ahead of schedule
			adjustments = append(adjustments, types.GoalAdjustment{
				GoalID:             goal.GoalID,
				RecommendationType: "increase_target",
				Adjustment:         "Consider increasing target value by 15-25%",
				Reason:             "Current progress suggests target may be too easy",
			})
		}

		// Goal-specific adjustments
		switch goal.GoalType {
		case types.GoalStrength:
			if progress.ProgressPercent > 0 && progress.ProgressPercent < 20 {
				adjustments = append(adjustments, types.GoalAdjustment{
					GoalID:             goal.GoalID,
					RecommendationType: "technique_focus",
					Adjustment:         "Focus on form and technique before increasing weight",
					Reason:             "Strength gains require proper movement patterns",
				})
			}
		case types.GoalFatLoss:
			if progress.ProgressPercent < 30 && timeUntilTarget < 45 {
				adjustments = append(adjustments, types.GoalAdjustment{
					GoalID:             goal.GoalID,
					RecommendationType: "nutrition_review",
					Adjustment:         "Review and adjust nutrition plan",
					Reason:             "Fat loss progress may be limited by dietary factors",
				})
			}
		}
	}

	return adjustments, nil
}
