package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) GetUserNutritionGoals(ctx context.Context, userID string) (*types.NutritionGoals, error) {
	q := `
		SELECT 
			user_id, 
			calories_goal, 
			protein_goal, 
			carbs_goal, 
			fat_goal, 
			fiber_goal
		FROM nutrition_goals
		WHERE user_id = $1
	`

	var goals types.NutritionGoals
	err := s.db.QueryRow(ctx, q, userID).Scan(
		&goals.UserID,
		&goals.CaloriesGoal,
		&goals.ProteinGoal,
		&goals.CarbsGoal,
		&goals.FatGoal,
		&goals.FiberGoal,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return &types.NutritionGoals{
				UserID:       userID,
				CaloriesGoal: 2000,
				ProteinGoal:  50,
				CarbsGoal:    275,
				FatGoal:      78,
				FiberGoal:    25,
			}, nil
		}
		return nil, fmt.Errorf("failed to get nutrition goals: %w", err)
	}

	return &goals, nil
}

func (s *Store) CreateUserNutritionGoals(ctx context.Context, goals *types.NutritionGoals) error {
	if goals == nil {
		return fmt.Errorf("goals cannot be nil")
	}

	q := `
		INSERT INTO nutrition_goals (
			user_id, 
			calories_goal, 
			protein_goal, 
			carbs_goal, 
			fat_goal, 
			fiber_goal
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			calories_goal = EXCLUDED.calories_goal,
			protein_goal = EXCLUDED.protein_goal,
			carbs_goal = EXCLUDED.carbs_goal,
			fat_goal = EXCLUDED.fat_goal,
			fiber_goal = EXCLUDED.fiber_goal,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := s.db.Exec(ctx, q,
		goals.UserID,
		goals.CaloriesGoal,
		goals.ProteinGoal,
		goals.CarbsGoal,
		goals.FatGoal,
		goals.FiberGoal,
	)

	if err != nil {
		return fmt.Errorf("failed to create nutrition goals: %w", err)
	}

	return nil
}

func (s *Store) UpdateUserNutritionGoals(ctx context.Context, goals *types.NutritionGoals) error {
	if goals == nil {
		return fmt.Errorf("goals cannot be nil")
	}

	q := `
		UPDATE nutrition_goals
		SET 
			calories_goal = $2,
			protein_goal = $3,
			carbs_goal = $4,
			fat_goal = $5,
			fiber_goal = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`

	result, err := s.db.Exec(ctx, q,
		goals.UserID,
		goals.CaloriesGoal,
		goals.ProteinGoal,
		goals.CarbsGoal,
		goals.FatGoal,
		goals.FiberGoal,
	)

	if err != nil {
		return fmt.Errorf("failed to update nutrition goals: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no nutrition goals found for user %s", goals.UserID)
	}

	return nil
}

func (s *Store) DeleteUserNutritionGoals(ctx context.Context, userID string) error {
	q := `
		DELETE FROM nutrition_goals
		WHERE user_id = $1
	`

	result, err := s.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("failed to delete nutrition goals: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no nutrition goals found for user %s", userID)
	}

	return nil
}
