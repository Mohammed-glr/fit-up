package repository

import (
	"context"
	"encoding/json"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

// =============================================================================
// PLAN GENERATION REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) GetPlanID(ctx context.Context, planID int) (*types.GeneratedPlan, error) {
	q := `
		SELECT plan_id, user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
		FROM generated_plans
		WHERE plan_id = $1
	`

	var plan types.GeneratedPlan
	err := s.db.QueryRow(ctx, q, planID).Scan(
		&plan.PlanID,
		&plan.UserID,
		&plan.WeekStart,
		&plan.GeneratedAt,
		&plan.Algorithm,
		&plan.Effectiveness,
		&plan.IsActive,
		&plan.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &plan, nil
}


func (s *Store) CreatePlanGeneration(ctx context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error) {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	q := `
		INSERT INTO generated_plans (user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata)
		VALUES ($1, $2, NOW(), $3, 0.0, true, $4)
		RETURNING plan_id, user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
	`

	var plan types.GeneratedPlan
	err = s.db.QueryRow(ctx, q,
		userID,
		metadata.Parameters["week_start"], // Assuming week_start is in parameters
		metadata.Algorithm,
		metadataJSON,
	).Scan(
		&plan.PlanID,
		&plan.UserID,
		&plan.WeekStart,
		&plan.GeneratedAt,
		&plan.Algorithm,
		&plan.Effectiveness,
		&plan.IsActive,
		&plan.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *Store) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	q := `
		SELECT plan_id, user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
		FROM generated_plans
		WHERE user_id = $1 AND is_active = true
		ORDER BY generated_at DESC
		LIMIT 1
	`

	var plan types.GeneratedPlan
	err := s.db.QueryRow(ctx, q, userID).Scan(
		&plan.PlanID,
		&plan.UserID,
		&plan.WeekStart,
		&plan.GeneratedAt,
		&plan.Algorithm,
		&plan.Effectiveness,
		&plan.IsActive,
		&plan.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *Store) GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error) {
	q := `
		SELECT plan_id, user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
		FROM generated_plans
		WHERE user_id = $1
		ORDER BY generated_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []types.GeneratedPlan
	for rows.Next() {
		var plan types.GeneratedPlan
		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.WeekStart,
			&plan.GeneratedAt,
			&plan.Algorithm,
			&plan.Effectiveness,
			&plan.IsActive,
			&plan.Metadata,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func (s *Store) TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error {
	effectivenessScore := s.calculateEffectivenessScore(performance)

	q := `
		UPDATE generated_plans 
		SET effectiveness = $1
		WHERE plan_id = $2
	`

	_, err := s.db.Exec(ctx, q, effectivenessScore, planID)
	if err != nil {
		return err
	}

	// Store detailed performance data
	performanceJSON, err := json.Marshal(performance)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO plan_performance_data (plan_id, recorded_at, completion_rate, average_rpe, progress_rate, user_satisfaction, injury_rate, performance_data)
		VALUES ($1, NOW(), $2, $3, $4, $5, $6, $7)
	`

	_, err = s.db.Exec(ctx, insertQuery,
		planID,
		performance.CompletionRate,
		performance.AverageRPE,
		performance.ProgressRate,
		performance.UserSatisfaction,
		performance.InjuryRate,
		performanceJSON,
	)

	return err
}

func (s *Store) calculateEffectivenessScore(performance *types.PlanPerformanceData) float64 {
	// Simple effectiveness calculation - can be made more sophisticated
	score := 0.0

	// Completion rate (40% weight)
	score += performance.CompletionRate * 0.4

	// Progress rate (30% weight)
	score += performance.ProgressRate * 0.3

	// User satisfaction (20% weight)
	score += performance.UserSatisfaction * 0.2

	// Injury penalty (10% weight) - lower injury rate is better
	score += (1.0 - performance.InjuryRate) * 0.1

	// Ensure score is between 0 and 1
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score
}

func (s *Store) GetPlanEffectivenessScore(ctx context.Context, planID int) (float64, error) {
	q := `
		SELECT effectiveness 
		FROM generated_plans 
		WHERE plan_id = $1
	`

	var effectiveness float64
	err := s.db.QueryRow(ctx, q, planID).Scan(&effectiveness)
	if err != nil {
		return 0.0, err
	}

	return effectiveness, nil
}

func (s *Store) MarkPlanForRegeneration(ctx context.Context, planID int, reason string) error {
	// Deactivate current plan
	updateQuery := `
		UPDATE generated_plans 
		SET is_active = false 
		WHERE plan_id = $1
	`

	_, err := s.db.Exec(ctx, updateQuery, planID)
	if err != nil {
		return err
	}

	// Log the regeneration trigger
	logQuery := `
		INSERT INTO plan_adaptations (plan_id, adaptation_date, reason, trigger, changes)
		VALUES ($1, NOW(), $2, 'regeneration_required', '{"action": "plan_marked_for_regeneration"}')
	`

	_, err = s.db.Exec(ctx, logQuery, planID, reason)
	return err
}

func (s *Store) LogPlanAdaptation(ctx context.Context, planID int, adaptation *types.PlanAdaptation) error {
	changesJSON, err := json.Marshal(adaptation.Changes)
	if err != nil {
		return err
	}

	q := `
		INSERT INTO plan_adaptations (plan_id, adaptation_date, reason, trigger, changes)
		VALUES ($1, NOW(), $2, $3, $4)
		RETURNING adaptation_id
	`

	var adaptationID int
	err = s.db.QueryRow(ctx, q,
		planID,
		adaptation.Reason,
		adaptation.Trigger,
		changesJSON,
	).Scan(&adaptationID)

	return err
}

func (s *Store) GetAdaptationHistory(ctx context.Context, userID int) ([]types.PlanAdaptation, error) {
	q := `
		SELECT pa.adaptation_id, pa.plan_id, pa.adaptation_date, pa.reason, pa.trigger, pa.changes
		FROM plan_adaptations pa
		JOIN generated_plans gp ON pa.plan_id = gp.plan_id
		WHERE gp.user_id = $1
		ORDER BY pa.adaptation_date DESC
		LIMIT 50
	`

	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adaptations []types.PlanAdaptation
	for rows.Next() {
		var adaptation types.PlanAdaptation
		err := rows.Scan(
			&adaptation.AdaptationID,
			&adaptation.PlanID,
			&adaptation.AdaptationDate,
			&adaptation.Reason,
			&adaptation.Trigger,
			&adaptation.Changes,
		)
		if err != nil {
			return nil, err
		}
		adaptations = append(adaptations, adaptation)
	}

	return adaptations, nil
}
