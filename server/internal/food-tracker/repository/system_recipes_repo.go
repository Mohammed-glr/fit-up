package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) GetSystemRecipeByID(ctx context.Context, id int) (*types.SystemRecipeDetail, error) {
	q := `
	SELECT
		id, name, description, category, difficulty, calories, protein, carbs, fat, fiber,
		prep_time, cook_time, servings, image_url, is_active, created_at, updated_at
	FROM system_recipes
	WHERE id = $1;
	`

	var recipe types.SystemRecipeDetail
	var description, imageURL, difficulty sql.NullString
	var cookTime, servings sql.NullInt32

	err := s.db.QueryRow(ctx, q, id).Scan(
		&recipe.RecipeID,
		&recipe.RecipeName,
		&description,
		&recipe.RecipesCategory,
		&difficulty,
		&recipe.RecipesCalories,
		&recipe.RecipesProtein,
		&recipe.RecipesCarbs,
		&recipe.RecipesFat,
		&recipe.RecipesFiber,
		&recipe.PrepTime,
		&cookTime,
		&servings,
		&imageURL,
		&recipe.IsActive,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if description.Valid {
		recipe.RecipeDesc = description.String
	}
	if difficulty.Valid {
		recipe.RecipesDifficulty = types.RecipeDifficulty(difficulty.String)
	}
	if imageURL.Valid {
		recipe.RecipesImageURL = imageURL.String
	}
	if cookTime.Valid {
		recipe.CookTime = int(cookTime.Int32)
	}
	if servings.Valid {
		recipe.Servings = int(servings.Int32)
	}
	if recipe.Servings == 0 {
		recipe.Servings = 1
	}

	return &recipe, nil
}

func (s *Store) GetSystemRecipeAll(ctx context.Context, filters types.RecipeFilters) ([]types.SystemRecipe, error) {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	base := strings.Builder{}
	base.WriteString(`
	SELECT
		id, name, description, category, difficulty, calories, protein, carbs, fat, fiber,
		prep_time, cook_time, servings, image_url, is_active, created_at, updated_at
	FROM system_recipes
	WHERE is_active = true`)

	args := []interface{}{}
	argPos := 1

	if filters.Category != nil {
		base.WriteString(fmt.Sprintf(" AND category = $%d", argPos))
		args = append(args, string(*filters.Category))
		argPos++
	}

	if filters.Difficulty != nil {
		base.WriteString(fmt.Sprintf(" AND difficulty = $%d", argPos))
		args = append(args, string(*filters.Difficulty))
		argPos++
	}

	if filters.MaxCalories != nil {
		base.WriteString(fmt.Sprintf(" AND calories <= $%d", argPos))
		args = append(args, *filters.MaxCalories)
		argPos++
	}

	if filters.MinProtein != nil {
		base.WriteString(fmt.Sprintf(" AND protein >= $%d", argPos))
		args = append(args, *filters.MinProtein)
		argPos++
	}

	if filters.MaxPrepTime != nil {
		base.WriteString(fmt.Sprintf(" AND prep_time <= $%d", argPos))
		args = append(args, *filters.MaxPrepTime)
		argPos++
	}

	if filters.SearchTerm != "" {
		base.WriteString(fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filters.SearchTerm+"%")
		argPos++
	}

	sortColumn := "created_at"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "created_at", "updated_at", "name", "calories", "prep_time":
			sortColumn = filters.SortBy
		}
	}

	sortDirection := "DESC"
	if strings.EqualFold(filters.SortOrder, "asc") {
		sortDirection = "ASC"
	}

	base.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortDirection))
	base.WriteString(fmt.Sprintf(" OFFSET $%d LIMIT $%d", argPos, argPos+1))

	args = append(args, filters.Offset, filters.Limit)

	rows, err := s.db.Query(ctx, base.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := make([]types.SystemRecipe, 0)
	for rows.Next() {
		var recipe types.SystemRecipe
		var description, imageURL, difficulty sql.NullString
		var cookTime, servings sql.NullInt32

		err := rows.Scan(
			&recipe.RecipeID,
			&recipe.RecipeName,
			&description,
			&recipe.RecipesCategory,
			&difficulty,
			&recipe.RecipesCalories,
			&recipe.RecipesProtein,
			&recipe.RecipesCarbs,
			&recipe.RecipesFat,
			&recipe.RecipesFiber,
			&recipe.PrepTime,
			&cookTime,
			&servings,
			&imageURL,
			&recipe.IsActive,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if description.Valid {
			recipe.RecipeDesc = description.String
		}
		if difficulty.Valid {
			recipe.RecipesDifficulty = types.RecipeDifficulty(difficulty.String)
		}
		if imageURL.Valid {
			recipe.RecipesImageURL = imageURL.String
		}
		if cookTime.Valid {
			recipe.CookTime = int(cookTime.Int32)
		}
		if servings.Valid {
			recipe.Servings = int(servings.Int32)
		}
		if recipe.Servings == 0 {
			recipe.Servings = 1
		}

		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (s *Store) CreateSystemRecipe(ctx context.Context, recipe *types.SystemRecipe) (int, error) {
	q := `
	INSERT INTO system_recipes
	(name, description, category, difficulty, calories, protein, carbs, fat, fiber,
	prep_time, cook_time, image_url, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id;
	`
	var id int
	err := s.db.QueryRow(ctx, q,
		recipe.RecipeName,
		recipe.RecipeDesc,
		recipe.RecipesCategory,
		recipe.RecipesDifficulty,
		recipe.RecipesCalories,
		recipe.RecipesProtein,
		recipe.RecipesCarbs,
		recipe.RecipesFat,
		recipe.RecipesFiber,
		recipe.PrepTime,
		recipe.CookTime,
		recipe.RecipesImageURL,
		recipe.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateSystemRecipe(ctx context.Context, recipe *types.SystemRecipe) error {
	q := `
	UPDATE system_recipes
	SET name = $1,
		description = $2,
		category = $3,
		difficulty = $4,
		calories = $5,
		protein = $6,
		carbs = $7,
		fat = $8,
		fiber = $9,
		prep_time = $10,
		cook_time = $11,
		image_url = $12,
		is_active = $13,
		updated_at = NOW()
	WHERE id = $14;
	`
	_, err := s.db.Exec(ctx, q,
		recipe.RecipeName,
		recipe.RecipeDesc,
		recipe.RecipesCategory,
		recipe.RecipesDifficulty,
		recipe.RecipesCalories,
		recipe.RecipesProtein,
		recipe.RecipesCarbs,
		recipe.RecipesFat,
		recipe.RecipesFiber,
		recipe.PrepTime,
		recipe.CookTime,
		recipe.RecipesImageURL,
		recipe.IsActive,
		recipe.RecipeID,
	)
	return err
}

func (s *Store) DeleteSystemRecipe(ctx context.Context, id int) error {
	q := `
	DELETE FROM system_recipes
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) SetActiveSystemRecipe(ctx context.Context, id int, isActive bool) error {
	q := `
	UPDATE system_recipes
	SET is_active = $1,
		updated_at = NOW()
	WHERE id = $2;
	`
	_, err := s.db.Exec(ctx, q, isActive, id)
	return err
}

func (s *Store) AddSystemRecipesIngredient(ctx context.Context, ingredient *types.SystemRecipesIngredient) error {
	q := `
	INSERT INTO system_recipes_ingredients
	(recipe_id, item, amount, unit, order_index)
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err := s.db.Exec(ctx, q,
		ingredient.IngredientAmount,
		ingredient.IngredientItem,
		ingredient.IngredientAmount,
		ingredient.IngredientUnit,
		ingredient.OrderIndex,
		ingredient.RecipeID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateSystemRecipesIngredient(ctx context.Context, ingredient *types.SystemRecipesIngredient) error {
	q := `
	UPDATE system_recipes_ingredients
	SET item = $1,
		amount = $2,
		unit = $3,
		order_index = $4
	WHERE id = $5;
	`
	_, err := s.db.Exec(ctx, q,
		ingredient.IngredientAmount,
		ingredient.IngredientItem,
		ingredient.IngredientID,
		ingredient.IngredientUnit,
		ingredient.OrderIndex,
		ingredient.RecipeID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteSystemRecipesIngredient(ctx context.Context, id int) error {
	q := `
	DELETE FROM system_recipes_ingredients
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetSystemRecipesIngredients(ctx context.Context, recipeID int) ([]types.SystemRecipesIngredient, error) {
	q := `
	SELECT
		id, recipe_id, item, amount, unit, order_index
	FROM system_recipes_ingredients
	WHERE recipe_id = $1
	ORDER BY order_index;
	`

	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []types.SystemRecipesIngredient
	for rows.Next() {
		var ingredient types.SystemRecipesIngredient
		err := rows.Scan(
			&ingredient.IngredientID,
			&ingredient.RecipeID,
			&ingredient.IngredientItem,
			&ingredient.IngredientAmount,
			&ingredient.IngredientUnit,
			&ingredient.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (s *Store) AddSystemRecipesInstruction(ctx context.Context, instruction *types.SystemRecipesInstruction) error {
	q := `
	INSERT INTO system_recipes_instructions
	(recipe_id, step_number, instruction)
	VALUES ($1, $2, $3);
	`
	_, err := s.db.Exec(ctx, q,
		instruction.RecipeID,
		instruction.InstructionStep,
		instruction.InstructionText,
		instruction.InstructionID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateSystemRecipesInstruction(ctx context.Context, instruction *types.SystemRecipesInstruction) error {
	q := `
	UPDATE system_recipes_instructions
	SET step_number = $1,
		instruction = $2
	WHERE id = $3;
	`
	_, err := s.db.Exec(ctx, q,
		instruction.InstructionID,
		instruction.InstructionStep,
		instruction.InstructionText,
		instruction.RecipeID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteSystemRecipesInstruction(ctx context.Context, id int) error {
	q := `
	DELETE FROM system_recipes_instructions
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetSystemRecipesInstructions(ctx context.Context, recipeID int) ([]types.SystemRecipesInstruction, error) {
	q := `
	SELECT
		id, recipe_id, step_number, instruction
	FROM system_recipes_instructions
	WHERE recipe_id = $1
	ORDER BY step_number;
	`
	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []types.SystemRecipesInstruction
	for rows.Next() {
		var instruction types.SystemRecipesInstruction
		err := rows.Scan(
			&instruction.InstructionID,
			&instruction.InstructionStep,
			&instruction.InstructionText,
			&instruction.RecipeID,
		)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func (s *Store) AddSystemRecipesTag(ctx context.Context, tag *types.SystemRecipesTag) error {
	q := `
	INSERT INTO system_recipes_tags
	(recipe_id, tag_name)
	VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, q,
		tag.RecipeID,
		tag.TagName,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteSystemRecipesTag(ctx context.Context, id int) error {
	q := `
	DELETE FROM system_recipes_tags
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetSystemRecipesTags(ctx context.Context, recipeID int) ([]types.SystemRecipesTag, error) {
	q := `
	SELECT
		id, recipe_id, tag_name
	FROM system_recipes_tags
	WHERE recipe_id = $1;
	`
	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []types.SystemRecipesTag
	for rows.Next() {
		var tag types.SystemRecipesTag
		err := rows.Scan(
			&tag.TagID,
			&tag.RecipeID,
			&tag.TagName,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Store) SearchSystemRecipesByTag(ctx context.Context, tag string) ([]types.SystemRecipe, error) {
	q := `
	SELECT
		sr.id, sr.name, sr.description, sr.category, sr.difficulty, sr.calories, sr.protein, sr.carbs, sr.fat, sr.fiber,
		sr.prep_time, sr.cook_time, sr.image_url, sr.is_active, sr.created_at, sr.updated_at
	FROM system_recipes sr
	JOIN system_recipes_tags srt ON sr.id = srt.recipe_id
	WHERE srt.tag_name = $1;
	`
	rows, err := s.db.Query(ctx, q, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []types.SystemRecipe
	for rows.Next() {
		var recipe types.SystemRecipe
		err := rows.Scan(
			&recipe.RecipeID,
			&recipe.RecipeName,
			&recipe.RecipesCalories,
			&recipe.RecipesCarbs,
			&recipe.RecipesDifficulty,
			&recipe.RecipesCategory,
			&recipe.RecipesFat,
			&recipe.RecipesFiber,
			&recipe.RecipesImageURL,
			&recipe.RecipesProtein,
			&recipe.CookTime,
			&recipe.IsActive,
			&recipe.PrepTime,
			&recipe.RecipeDesc,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
