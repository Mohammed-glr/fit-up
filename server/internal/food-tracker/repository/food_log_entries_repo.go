package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) CreateFoodLogEntry(ctx context.Context, entry *types.FoodLogEntry) (int, error) {
	q := `
		INSERT INTO food_log_entries (
			user_id, log_date, meal_type, system_recipe_id, user_recipe_id,
			calories, protein, carbs, fat, fiber, servings
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	var id int
	err := s.db.QueryRow(ctx, q,
		entry.UserID,
		entry.LogDate,
		entry.MealType,
		entry.SystemRecipeID,
		entry.UserRecipeID,
		entry.Calories,
		entry.Protein,
		entry.Carbs,
		entry.Fat,
		entry.Fiber,
		entry.Servings,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) GetFoodLogEntryByID(ctx context.Context, id int, userID string) (*types.FoodLogEntryWithRecipe, error) {
	q := `
		SELECT 
			fle.id, fle.user_id, fle.log_date, fle.meal_type,
			fle.system_recipe_id, fle.user_recipe_id,
			fle.calories, fle.protein, fle.carbs, fle.fat, fle.fiber, fle.servings,
			fle.created_at, fle.updated_at,
			COALESCE(sr.name, ur.name, '') as recipe_name,
			CASE 
				WHEN fle.system_recipe_id IS NOT NULL THEN 'system'
				WHEN fle.user_recipe_id IS NOT NULL THEN 'user'
				ELSE ''
			END as recipe_source
		FROM food_log_entries fle
		LEFT JOIN system_recipes sr ON fle.system_recipe_id = sr.id
		LEFT JOIN user_recipes ur ON fle.user_recipe_id = ur.id
		WHERE fle.id = $1 AND fle.user_id = $2
	`

	var entry types.FoodLogEntryWithRecipe
	err := s.db.QueryRow(ctx, q, id, userID).Scan(
		&entry.EntryID,
		&entry.UserID,
		&entry.LogDate,
		&entry.MealType,
		&entry.SystemRecipeID,
		&entry.UserRecipeID,
		&entry.Calories,
		&entry.Protein,
		&entry.Carbs,
		&entry.Fat,
		&entry.Fiber,
		&entry.Servings,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&entry.RecipeName,
		&entry.RecipeSource,
	)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Store) UpdateFoodLogEntry(ctx context.Context, entry *types.FoodLogEntry) error {
	q := `
		UPDATE food_log_entries
		SET log_date = $1, meal_type = $2, system_recipe_id = $3, user_recipe_id = $4,
		    calories = $5, protein = $6, carbs = $7, fat = $8, fiber = $9, servings = $10,
		    updated_at = NOW()
		WHERE id = $11 AND user_id = $12
	`
	_, err := s.db.Exec(ctx, q,
		entry.LogDate,
		entry.MealType,
		entry.SystemRecipeID,
		entry.UserRecipeID,
		entry.Calories,
		entry.Protein,
		entry.Carbs,
		entry.Fat,
		entry.Fiber,
		entry.Servings,
		entry.EntryID,
		entry.UserID,
	)
	return err
}

func (s *Store) DeleteFoodLogEntry(ctx context.Context, id int, userID string) error {
	q := `
		DELETE FROM food_log_entries
		WHERE id = $1 AND user_id = $2
	`
	_, err := s.db.Exec(ctx, q, id, userID)
	return err
}

func (s *Store) GetFoodLogEntriesByDate(ctx context.Context, userID string, date string) ([]types.FoodLogEntryWithRecipe, error) {
	q := `
		SELECT 
			fle.id, fle.user_id, fle.log_date, fle.meal_type,
			fle.system_recipe_id, fle.user_recipe_id,
			fle.calories, fle.protein, fle.carbs, fle.fat, fle.fiber, fle.servings,
			fle.created_at, fle.updated_at,
			COALESCE(sr.name, ur.name, '') as recipe_name,
			CASE 
				WHEN fle.system_recipe_id IS NOT NULL THEN 'system'
				WHEN fle.user_recipe_id IS NOT NULL THEN 'user'
				ELSE ''
			END as recipe_source
		FROM food_log_entries fle
		LEFT JOIN system_recipes sr ON fle.system_recipe_id = sr.id
		LEFT JOIN user_recipes ur ON fle.user_recipe_id = ur.id
		WHERE fle.user_id = $1 AND fle.log_date = $2
		ORDER BY fle.meal_type, fle.created_at
	`
	rows, err := s.db.Query(ctx, q, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []types.FoodLogEntryWithRecipe
	for rows.Next() {
		var entry types.FoodLogEntryWithRecipe
		err := rows.Scan(
			&entry.EntryID,
			&entry.UserID,
			&entry.LogDate,
			&entry.MealType,
			&entry.SystemRecipeID,
			&entry.UserRecipeID,
			&entry.Calories,
			&entry.Protein,
			&entry.Carbs,
			&entry.Fat,
			&entry.Fiber,
			&entry.Servings,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.RecipeName,
			&entry.RecipeSource,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Store) GetFoodLogEntriesByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]types.FoodLogEntryWithRecipe, error) {
	q := `
		SELECT 
			fle.id, fle.user_id, fle.log_date, fle.meal_type,
			fle.system_recipe_id, fle.user_recipe_id,
			fle.calories, fle.protein, fle.carbs, fle.fat, fle.fiber, fle.servings,
			fle.created_at, fle.updated_at,
			COALESCE(sr.name, ur.name, '') as recipe_name,
			CASE 
				WHEN fle.system_recipe_id IS NOT NULL THEN 'system'
				WHEN fle.user_recipe_id IS NOT NULL THEN 'user'
				ELSE ''
			END as recipe_source
		FROM food_log_entries fle
		LEFT JOIN system_recipes sr ON fle.system_recipe_id = sr.id
		LEFT JOIN user_recipes ur ON fle.user_recipe_id = ur.id
		WHERE fle.user_id = $1 AND fle.log_date BETWEEN $2 AND $3
		ORDER BY fle.log_date, fle.meal_type, fle.created_at
	`

	rows, err := s.db.Query(ctx, q, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []types.FoodLogEntryWithRecipe
	for rows.Next() {
		var entry types.FoodLogEntryWithRecipe
		err := rows.Scan(
			&entry.EntryID,
			&entry.UserID,
			&entry.LogDate,
			&entry.MealType,
			&entry.SystemRecipeID,
			&entry.UserRecipeID,
			&entry.Calories,
			&entry.Protein,
			&entry.Carbs,
			&entry.Fat,
			&entry.Fiber,
			&entry.Servings,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.RecipeName,
			&entry.RecipeSource,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Store) GetFoodLogEntriesByMealType(ctx context.Context, userID string, date string, mealType types.MealType) ([]types.FoodLogEntryWithRecipe, error) {
	q := `
		SELECT 
			fle.id, fle.user_id, fle.log_date, fle.meal_type,
			fle.system_recipe_id, fle.user_recipe_id,
			fle.calories, fle.protein, fle.carbs, fle.fat, fle.fiber, fle.servings,
			fle.created_at, fle.updated_at,
			COALESCE(sr.name, ur.name, '') as recipe_name,
			CASE 
				WHEN fle.system_recipe_id IS NOT NULL THEN 'system'
				WHEN fle.user_recipe_id IS NOT NULL THEN 'user'
				ELSE ''
			END as recipe_source
		FROM food_log_entries fle
		LEFT JOIN system_recipes sr ON fle.system_recipe_id = sr.id
		LEFT JOIN user_recipes ur ON fle.user_recipe_id = ur.id
		WHERE fle.user_id = $1 AND fle.log_date = $2 AND fle.meal_type = $3
		ORDER BY fle.created_at
	`
	rows, err := s.db.Query(ctx, q, userID, date, mealType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []types.FoodLogEntryWithRecipe
	for rows.Next() {
		var entry types.FoodLogEntryWithRecipe
		err := rows.Scan(
			&entry.EntryID,
			&entry.UserID,
			&entry.LogDate,
			&entry.MealType,
			&entry.SystemRecipeID,
			&entry.UserRecipeID,
			&entry.Calories,
			&entry.Protein,
			&entry.Carbs,
			&entry.Fat,
			&entry.Fiber,
			&entry.Servings,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.RecipeName,
			&entry.RecipeSource,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Store) GetDailySummary(ctx context.Context, userID string, date string) (*types.DailyNutritionSummary, error) {
	q := `
		SELECT 
			$1::text as user_id,
			$2::text as log_date,
			COALESCE(SUM(fle.calories), 0) AS total_calories,
			COALESCE(SUM(fle.protein), 0) AS total_protein,
			COALESCE(SUM(fle.carbs), 0) AS total_carbs,
			COALESCE(SUM(fle.fat), 0) AS total_fat,
			COALESCE(SUM(fle.fiber), 0) AS total_fiber,
			COUNT(fle.id) AS total_entries
		FROM food_log_entries fle
		WHERE fle.user_id = $1 AND fle.log_date = $2
	`
	var summary types.DailyNutritionSummary
	err := s.db.QueryRow(ctx, q, userID, date).Scan(
		&summary.UserID,
		&summary.LogDate,
		&summary.TotalCalories,
		&summary.TotalProtein,
		&summary.TotalCarbs,
		&summary.TotalFat,
		&summary.TotalFiber,
		&summary.TotalEntries,
	)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (s *Store) GetWeeklySummary(ctx context.Context, userID string, startDate string) ([]types.DailyNutritionSummary, error) {
	q := `
		SELECT 
			$1::text as user_id,
			fle.log_date,
			COALESCE(SUM(fle.calories), 0) AS total_calories,
			COALESCE(SUM(fle.protein), 0) AS total_protein,
			COALESCE(SUM(fle.carbs), 0) AS total_carbs,
			COALESCE(SUM(fle.fat), 0) AS total_fat,
			COALESCE(SUM(fle.fiber), 0) AS total_fiber,
			COUNT(fle.id) AS total_entries
		FROM food_log_entries fle
		WHERE fle.user_id = $1 AND fle.log_date BETWEEN $2 AND $2::date + INTERVAL '6 days'
		GROUP BY fle.log_date
		ORDER BY fle.log_date
	`
	rows, err := s.db.Query(ctx, q, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.DailyNutritionSummary
	for rows.Next() {
		var summary types.DailyNutritionSummary
		err := rows.Scan(
			&summary.UserID,
			&summary.LogDate,
			&summary.TotalCalories,
			&summary.TotalProtein,
			&summary.TotalCarbs,
			&summary.TotalFat,
			&summary.TotalFiber,
			&summary.TotalEntries,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (s *Store) GetMonthlySummary(ctx context.Context, userID string, year int, month int) ([]types.DailyNutritionSummary, error) {
	q := `
		SELECT 
			$1::text as user_id,
			fle.log_date,
			COALESCE(SUM(fle.calories), 0) AS total_calories,
			COALESCE(SUM(fle.protein), 0) AS total_protein,
			COALESCE(SUM(fle.carbs), 0) AS total_carbs,
			COALESCE(SUM(fle.fat), 0) AS total_fat,
			COALESCE(SUM(fle.fiber), 0) AS total_fiber,
			COUNT(fle.id) AS total_entries
		FROM food_log_entries fle
		WHERE fle.user_id = $1 AND EXTRACT(YEAR FROM fle.log_date) = $2 AND EXTRACT(MONTH FROM fle.log_date) = $3
		GROUP BY fle.log_date
		ORDER BY fle.log_date
	`
	rows, err := s.db.Query(ctx, q, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.DailyNutritionSummary
	for rows.Next() {
		var summary types.DailyNutritionSummary
		err := rows.Scan(
			&summary.UserID,
			&summary.LogDate,
			&summary.TotalCalories,
			&summary.TotalProtein,
			&summary.TotalCarbs,
			&summary.TotalFat,
			&summary.TotalFiber,
			&summary.TotalEntries,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
