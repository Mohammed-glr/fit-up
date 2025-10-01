package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

// =============================================================================
// FITNESS PROFILE REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) CreateFitnessAssessment(ctx context.Context, userID int, assessment *types.FitnessAssessmentRequest) (*types.FitnessAssessment, error) {
	q := `
		INSERT INTO fitness_assessments (user_id, assessment_date, overall_level, strength_level, cardio_level, flexibility_level, assessment_data)
		VALUES ($1, NOW(), $2, $3, $4, $5, $6)
		RETURNING assessment_id, user_id, assessment_date, overall_level, strength_level, cardio_level, flexibility_level, assessment_data
	`

	assessmentDataJSON, err := json.Marshal(assessment.AssessmentData)
	if err != nil {
		return nil, err
	}

	var result types.FitnessAssessment
	err = s.db.QueryRow(ctx, q,
		userID,
		assessment.OverallLevel,
		assessment.StrengthLevel,
		assessment.CardioLevel,
		assessment.FlexibilityLevel,
		assessmentDataJSON,
	).Scan(
		&result.AssessmentID,
		&result.UserID,
		&result.AssessmentDate,
		&result.OverallLevel,
		&result.StrengthLevel,
		&result.CardioLevel,
		&result.FlexibilityLevel,
		&result.AssessmentData,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Store) GetUserFitnessProfile(ctx context.Context, userID int) (*types.FitnessProfile, error) {
	// Get latest assessment
	assessmentQuery := `
		SELECT overall_level, strength_level, cardio_level, flexibility_level, assessment_date
		FROM fitness_assessments
		WHERE user_id = $1
		ORDER BY assessment_date DESC
		LIMIT 1
	`

	var profile types.FitnessProfile
	var lastAssessment time.Time

	err := s.db.QueryRow(ctx, assessmentQuery, userID).Scan(
		&profile.CurrentLevel,
		&profile.StrengthLevel,
		&profile.CardioLevel,
		&profile.FlexibilityLevel,
		&lastAssessment,
	)

	if err != nil {
		// No assessment found, return basic profile from user table
		userQuery := `
			SELECT level, goal, equipment
			FROM users
			WHERE user_id = $1
		`

		var goal types.FitnessGoal
		var equipmentJSON json.RawMessage

		err = s.db.QueryRow(ctx, userQuery, userID).Scan(
			&profile.CurrentLevel,
			&goal,
			&equipmentJSON,
		)

		if err != nil {
			return nil, err
		}

		// Parse equipment
		var equipmentStrings []string
		if err := json.Unmarshal(equipmentJSON, &equipmentStrings); err == nil {
			for _, eq := range equipmentStrings {
				profile.Equipment = append(profile.Equipment, types.EquipmentType(eq))
			}
		}

		profile.UserID = userID
		return &profile, nil
	}

	profile.UserID = userID
	profile.LastAssessment = &lastAssessment

	// Get active goals
	goalsQuery := `
		SELECT goal_id, goal_type, target_value, current_value, target_date, created_at, metadata
		FROM fitness_goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, goalsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []types.FitnessGoalTarget
	for rows.Next() {
		var goal types.FitnessGoalTarget
		err := rows.Scan(
			&goal.GoalID,
			&goal.GoalType,
			&goal.TargetValue,
			&goal.CurrentValue,
			&goal.TargetDate,
			&goal.CreatedAt,
			&goal.Metadata,
		)
		if err != nil {
			return nil, err
		}
		goal.UserID = userID
		goal.IsActive = true
		goals = append(goals, goal)
	}

	profile.Goals = goals

	// Get equipment from user table
	equipmentQuery := `SELECT equipment FROM users WHERE user_id = $1`
	var equipmentJSON json.RawMessage
	err = s.db.QueryRow(ctx, equipmentQuery, userID).Scan(&equipmentJSON)
	if err == nil {
		var equipmentStrings []string
		if err := json.Unmarshal(equipmentJSON, &equipmentStrings); err == nil {
			for _, eq := range equipmentStrings {
				profile.Equipment = append(profile.Equipment, types.EquipmentType(eq))
			}
		}
	}

	return &profile, nil
}

func (s *Store) UpdateFitnessLevel(ctx context.Context, userID int, level types.FitnessLevel) error {
	q := `
		UPDATE users 
		SET level = $1
		WHERE user_id = $2
	`

	_, err := s.db.Exec(ctx, q, level, userID)
	return err
}

func (s *Store) UpdateFitnessGoals(ctx context.Context, userID int, goals []types.FitnessGoalTarget) error {
	// First deactivate all existing goals
	deactivateQuery := `
		UPDATE fitness_goals 
		SET is_active = false 
		WHERE user_id = $1
	`

	_, err := s.db.Exec(ctx, deactivateQuery, userID)
	if err != nil {
		return err
	}

	// Insert new active goals
	insertQuery := `
		INSERT INTO fitness_goals (user_id, goal_type, target_value, current_value, target_date, is_active, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, true, NOW(), $6)
	`

	for _, goal := range goals {
		metadataJSON, err := json.Marshal(goal.Metadata)
		if err != nil {
			metadataJSON = []byte("{}")
		}

		_, err = s.db.Exec(ctx, insertQuery,
			userID,
			goal.GoalType,
			goal.TargetValue,
			goal.CurrentValue,
			goal.TargetDate,
			metadataJSON,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) EstimateOneRepMax(ctx context.Context, userID int, exerciseID int, performance *types.PerformanceData) (*types.OneRepMaxEstimate, error) {
	// Use Epley formula: 1RM = weight * (1 + reps/30)
	estimatedMax := performance.Weight * (1 + float64(performance.Reps)/30.0)

	// Confidence based on rep range (higher reps = lower confidence)
	confidence := 1.0
	if performance.Reps > 10 {
		confidence = 0.7
	} else if performance.Reps > 5 {
		confidence = 0.85
	}

	q := `
		INSERT INTO one_rep_max_estimates (user_id, exercise_id, estimated_max, estimate_date, method, confidence)
		VALUES ($1, $2, $3, NOW(), 'epley', $4)
		RETURNING estimate_id, user_id, exercise_id, estimated_max, estimate_date, method, confidence
	`

	var estimate types.OneRepMaxEstimate
	err := s.db.QueryRow(ctx, q,
		userID,
		exerciseID,
		estimatedMax,
		confidence,
	).Scan(
		&estimate.EstimateID,
		&estimate.UserID,
		&estimate.ExerciseID,
		&estimate.EstimatedMax,
		&estimate.EstimateDate,
		&estimate.Method,
		&estimate.Confidence,
	)

	if err != nil {
		return nil, err
	}

	return &estimate, nil
}

func (s *Store) GetOneRepMaxHistory(ctx context.Context, userID int, exerciseID int) ([]types.OneRepMaxEstimate, error) {
	q := `
		SELECT estimate_id, user_id, exercise_id, estimated_max, estimate_date, method, confidence
		FROM one_rep_max_estimates
		WHERE user_id = $1 AND exercise_id = $2
		ORDER BY estimate_date DESC
		LIMIT 20
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estimates []types.OneRepMaxEstimate
	for rows.Next() {
		var estimate types.OneRepMaxEstimate
		err := rows.Scan(
			&estimate.EstimateID,
			&estimate.UserID,
			&estimate.ExerciseID,
			&estimate.EstimatedMax,
			&estimate.EstimateDate,
			&estimate.Method,
			&estimate.Confidence,
		)
		if err != nil {
			return nil, err
		}
		estimates = append(estimates, estimate)
	}

	return estimates, nil
}

func (s *Store) UpdateOneRepMax(ctx context.Context, userID int, exerciseID int, estimate float64) error {
	q := `
		INSERT INTO one_rep_max_estimates (user_id, exercise_id, estimated_max, estimate_date, method, confidence)
		VALUES ($1, $2, $3, NOW(), 'manual', 1.0)
	`

	_, err := s.db.Exec(ctx, q, userID, exerciseID, estimate)
	return err
}

func (s *Store) CreateMovementAssessment(ctx context.Context, userID int, assessment *types.MovementAssessmentRequest) (*types.MovementAssessment, error) {
	movementDataJSON, err := json.Marshal(assessment.MovementData)
	if err != nil {
		return nil, err
	}

	q := `
		INSERT INTO movement_assessments (user_id, assessment_date, movement_data)
		VALUES ($1, NOW(), $2)
		RETURNING assessment_id, user_id, assessment_date, movement_data
	`

	var result types.MovementAssessment
	err = s.db.QueryRow(ctx, q,
		userID,
		movementDataJSON,
	).Scan(
		&result.AssessmentID,
		&result.UserID,
		&result.AssessmentDate,
		&result.MovementData,
	)

	if err != nil {
		return nil, err
	}

	// Insert limitations if any
	if len(assessment.Limitations) > 0 {
		limitationQuery := `
			INSERT INTO movement_limitations (user_id, movement_type, severity, description)
			VALUES ($1, $2, 'moderate', $3)
		`

		for _, limitation := range assessment.Limitations {
			_, err = s.db.Exec(ctx, limitationQuery, userID, limitation, limitation)
			if err != nil {
				// Continue with other limitations even if one fails
				continue
			}
		}
	}

	return &result, nil
}

func (s *Store) GetMovementLimitations(ctx context.Context, userID int) ([]types.MovementLimitation, error) {
	q := `
		SELECT limitation_id, user_id, movement_type, severity, description
		FROM movement_limitations
		WHERE user_id = $1
		ORDER BY limitation_id DESC
	`

	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limitations []types.MovementLimitation
	for rows.Next() {
		var limitation types.MovementLimitation
		err := rows.Scan(
			&limitation.LimitationID,
			&limitation.UserID,
			&limitation.MovementType,
			&limitation.Severity,
			&limitation.Description,
		)
		if err != nil {
			return nil, err
		}
		limitations = append(limitations, limitation)
	}

	return limitations, nil
}
