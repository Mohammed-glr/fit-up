package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) GetPlanID(ctx context.Context, planID int) (*types.GeneratedPlan, error) {
	q := `
		SELECT gp.plan_id, wp.workout_profile_id, gp.week_start, gp.generated_at, gp.algorithm, gp.effectiveness, gp.is_active, gp.metadata
		FROM generated_plans gp
		JOIN workout_profiles wp ON gp.user_id = wp.auth_user_id
		WHERE gp.plan_id = $1
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.ErrPlanNotFound
		}
		return nil, err
	}

	return &plan, nil
}

func (s *Store) CreatePlanGeneration(ctx context.Context, userID int, authUserID string, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error) {
	if metadata == nil {
		return nil, fmt.Errorf("plan generation metadata cannot be nil")
	}

	weekStart := startOfWeek(time.Now().UTC())
	if metadata.Parameters == nil {
		metadata.Parameters = make(map[string]any)
	}

	if raw, ok := metadata.Parameters["week_start"]; ok {
		if parsed, ok := parseWeekStart(raw); ok {
			weekStart = parsed
		}
	}

	metadata.Parameters["week_start"] = weekStart.Format("2006-01-02")

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	q := `
		INSERT INTO generated_plans (user_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata)
		VALUES ($1, $2, NOW(), $3, 0.0, true, $4)
		RETURNING plan_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
	`

	var plan types.GeneratedPlan
	err = s.db.QueryRow(ctx, q,
		authUserID,
		weekStart,
		metadata.Algorithm,
		metadataJSON,
	).Scan(
		&plan.PlanID,
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

	plan.UserID = userID

	return &plan, nil
}

func (s *Store) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	q := `
		SELECT plan_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
		FROM generated_plans
		WHERE user_id = $1 AND is_active = true
		ORDER BY generated_at DESC
		LIMIT 1
	`

	authUserID, err := s.lookupAuthUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, types.ErrInvalidUserID) {
			return nil, nil
		}
		return nil, err
	}

	var plan types.GeneratedPlan
	err = s.db.QueryRow(ctx, q, authUserID).Scan(
		&plan.PlanID,
		&plan.WeekStart,
		&plan.GeneratedAt,
		&plan.Algorithm,
		&plan.Effectiveness,
		&plan.IsActive,
		&plan.Metadata,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	plan.UserID = userID

	return &plan, nil
}

func (s *Store) GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error) {
	q := `
		SELECT plan_id, week_start, generated_at, algorithm, effectiveness, is_active, metadata
		FROM generated_plans
		WHERE user_id = $1
		ORDER BY generated_at DESC
		LIMIT $2
	`

	authUserID, err := s.lookupAuthUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, q, authUserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []types.GeneratedPlan
	for rows.Next() {
		var plan types.GeneratedPlan
		err := rows.Scan(
			&plan.PlanID,
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
		plan.UserID = userID
		plans = append(plans, plan)
	}

	return plans, nil
}

func (s *Store) TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error {
	return s.db.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		execErr := s.trackPlanPerformanceWithConn(ctx, conn, planID, performance)
		if execErr == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if errors.As(execErr, &pgErr) && pgErr.Code == "08P01" && strings.Contains(pgErr.Message, "prepared statement name is already in use") {
			// Reset prepared statements on the connection before retrying.
			if _, deallocErr := conn.Exec(ctx, "DEALLOCATE ALL"); deallocErr != nil {
				return execErr
			}
			return s.trackPlanPerformanceWithConn(ctx, conn, planID, performance)
		}

		return execErr
	})
}

func (s *Store) trackPlanPerformanceWithConn(ctx context.Context, conn *pgxpool.Conn, planID int, performance *types.PlanPerformanceData) error {
	effectivenessScore := s.calculateEffectivenessScore(performance)

	const updateQuery = `
		UPDATE generated_plans 
		SET effectiveness = $1
		WHERE plan_id = $2
	`

	if _, err := conn.Exec(ctx, updateQuery, effectivenessScore, planID); err != nil {
		return err
	}

	const insertQuery = `
		INSERT INTO plan_performance_data (
			plan_id,
			completion_rate,
			average_rpe,
			progress_rate,
			user_satisfaction,
			injury_rate
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := conn.Exec(ctx, insertQuery,
		planID,
		performance.CompletionRate,
		performance.AverageRPE,
		performance.ProgressRate,
		performance.UserSatisfaction,
		performance.InjuryRate,
	)

	return err
}

func (s *Store) calculateEffectivenessScore(performance *types.PlanPerformanceData) float64 {
	score := 0.0

	score += performance.CompletionRate * 0.4

	score += performance.ProgressRate * 0.3

	score += performance.UserSatisfaction * 0.2

	score += (1.0 - performance.InjuryRate) * 0.1

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
	updateQuery := `
		UPDATE generated_plans 
		SET is_active = false 
		WHERE plan_id = $1
	`

	_, err := s.db.Exec(ctx, updateQuery, planID)
	if err != nil {
		return err
	}

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

	authUserID, err := s.lookupAuthUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, q, authUserID)
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

func (s *Store) CountActivePlans(ctx context.Context, userID int) (int, error) {
	authUserID, err := s.lookupAuthUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, types.ErrInvalidUserID) {
			return 0, nil
		}
		return 0, err
	}

	const q = `SELECT COUNT(*) FROM generated_plans WHERE user_id = $1 AND is_active = true`

	var count int
	if err := s.db.QueryRow(ctx, q, authUserID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Store) SaveGeneratedPlanStructure(ctx context.Context, planID int, structure []types.PlanStructureDayInput) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM generated_plan_exercises WHERE plan_day_id IN (SELECT plan_day_id FROM generated_plan_days WHERE plan_id = $1)`, planID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM generated_plan_days WHERE plan_id = $1`, planID); err != nil {
		return err
	}

	for _, day := range structure {
		var planDayID int
		err := tx.QueryRow(ctx,
			`INSERT INTO generated_plan_days (plan_id, day_index, day_title, focus, is_rest)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING plan_day_id`,
			planID,
			day.DayIndex,
			day.DayTitle,
			day.Focus,
			day.IsRest,
		).Scan(&planDayID)
		if err != nil {
			return err
		}

		for idx, exercise := range day.Exercises {
			var exerciseID interface{}
			if exercise.ExerciseID != nil {
				exerciseID = *exercise.ExerciseID
			}

			if _, err := tx.Exec(ctx,
				`INSERT INTO generated_plan_exercises (plan_day_id, exercise_order, exercise_id, name, sets, reps, rest_seconds, notes)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				planDayID,
				idx+1,
				exerciseID,
				exercise.Name,
				exercise.Sets,
				exercise.Reps,
				exercise.RestSeconds,
				exercise.Notes,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) GetGeneratedPlanStructure(ctx context.Context, planID int) ([]types.GeneratedPlanDay, error) {
	const dayQuery = `
		SELECT plan_day_id, plan_id, day_index, day_title, focus, is_rest
		FROM generated_plan_days
		WHERE plan_id = $1
		ORDER BY day_index
	`

	rows, err := s.db.Query(ctx, dayQuery, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []types.GeneratedPlanDay
	for rows.Next() {
		var day types.GeneratedPlanDay
		if err := rows.Scan(
			&day.PlanDayID,
			&day.PlanID,
			&day.DayIndex,
			&day.DayTitle,
			&day.Focus,
			&day.IsRest,
		); err != nil {
			return nil, err
		}

		exQuery := `
			SELECT plan_exercise_id, plan_day_id, exercise_order, exercise_id, name, sets, reps, rest_seconds, notes
			FROM generated_plan_exercises
			WHERE plan_day_id = $1
			ORDER BY exercise_order
		`

		exRows, err := s.db.Query(ctx, exQuery, day.PlanDayID)
		if err != nil {
			return nil, err
		}

		var exercises []types.GeneratedPlanExercise
		for exRows.Next() {
			var ex types.GeneratedPlanExercise
			var exerciseID pgtype.Int4
			if err := exRows.Scan(
				&ex.PlanExerciseID,
				&ex.PlanDayID,
				&ex.ExerciseOrder,
				&exerciseID,
				&ex.Name,
				&ex.Sets,
				&ex.Reps,
				&ex.RestSeconds,
				&ex.Notes,
			); err != nil {
				exRows.Close()
				return nil, err
			}
			if exerciseID.Valid {
				val := int(exerciseID.Int32)
				ex.ExerciseID = &val
			}
			exercises = append(exercises, ex)
		}
		exRows.Close()

		day.Exercises = exercises
		days = append(days, day)
	}

	return days, nil
}

func (s *Store) DeletePlanForUser(ctx context.Context, planID int, authUserID string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM generated_plans WHERE plan_id = $1 AND user_id = $2`, planID, authUserID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return types.ErrPlanNotFound
	}

	return nil
}

func parseWeekStart(value any) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return startOfWeek(v), true
	case string:
		if v == "" {
			return time.Time{}, false
		}
		for _, layout := range []string{time.RFC3339, "2006-01-02"} {
			if ts, err := time.Parse(layout, v); err == nil {
				return startOfWeek(ts), true
			}
		}
	}
	return time.Time{}, false
}

func (s *Store) lookupAuthUserID(ctx context.Context, userID int) (string, error) {
	const q = `SELECT auth_user_id FROM workout_profiles WHERE workout_profile_id = $1`

	var authUserID string
	err := s.db.QueryRow(ctx, q, userID).Scan(&authUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", types.ErrInvalidUserID
		}
		return "", err
	}

	return authUserID, nil
}

func startOfWeek(t time.Time) time.Time {
	loc := t.Location()
	base := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	daysSinceMonday := (int(base.Weekday()) + 6) % 7
	return base.AddDate(0, 0, -daysSinceMonday)
}
