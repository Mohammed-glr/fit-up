package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) GetUserRecipeByID(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error) {
	q := `
	SELECT
		ur.id,
		ur.user_id,
		ur.name,
		ur.description,
		ur.category,
		ur.difficulty,
		ur.calories,
		ur.protein,
		ur.carbs,
		ur.fat,
		ur.fiber,
		ur.prep_time,
		ur.cook_time,
		ur.image_url,
		ur.is_favorite,
		ur.created_at,
		ur.updated_at
	FROM user_recipes ur
	WHERE ur.id = $1 AND ur.user_id = $2;
	`

	var recipe types.UserRecipeDetail
	err := s.db.QueryRow(ctx, q, id, userID).Scan(
		&recipe.RecipeID,
		&recipe.UserID,
		&recipe.RecipeName,
		&recipe.RecipeDesc,
		&recipe.RecipesCategory,
		&recipe.RecipesDifficulty,
		&recipe.RecipesCalories,
		&recipe.RecipesProtein,
		&recipe.RecipesCarbs,
		&recipe.RecipesFat,
		&recipe.RecipesFiber,
		&recipe.PrepTime,
		&recipe.CookTime,
		&recipe.RecipesImageURL,
		&recipe.IsFavorite,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if recipe.Servings == 0 {
		recipe.Servings = 1
	}
	return &recipe, nil
}

func (s *Store) GetAllUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error) {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	base := strings.Builder{}
	base.WriteString(`
SELECT
	id,
	user_id,
	name,
	description,
	category,
	difficulty,
	calories,
	protein,
	carbs,
	fat,
	fiber,
	prep_time,
	cook_time,
	image_url,
	is_favorite,
	created_at,
	updated_at
FROM user_recipes
WHERE user_id = $1`)

	args := []interface{}{userID}
	argPos := 2

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

	if filters.IsFavorite != nil {
		base.WriteString(fmt.Sprintf(" AND is_favorite = $%d", argPos))
		args = append(args, *filters.IsFavorite)
		argPos++
	}

	if filters.SearchTerm != "" {
		base.WriteString(fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filters.SearchTerm+"%")
		argPos++
	}

	sortColumn := "updated_at"
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

	var recipes []types.UserRecipe
	for rows.Next() {
		var recipe types.UserRecipe
		if err := rows.Scan(
			&recipe.RecipeID,
			&recipe.UserID,
			&recipe.RecipeName,
			&recipe.RecipeDesc,
			&recipe.RecipesCategory,
			&recipe.RecipesDifficulty,
			&recipe.RecipesCalories,
			&recipe.RecipesProtein,
			&recipe.RecipesCarbs,
			&recipe.RecipesFat,
			&recipe.RecipesFiber,
			&recipe.PrepTime,
			&recipe.CookTime,
			&recipe.RecipesImageURL,
			&recipe.IsFavorite,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		); err != nil {
			return nil, err
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

func (s *Store) CreateUserRecipe(ctx context.Context, recipe *types.UserRecipe) (int, error) {
	q := `
	INSERT INTO user_recipes
		(user_id, name, description, category, difficulty, calories,
		 protein, carbs, fat, fiber, prep_time, cook_time, image_url, is_favorite)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id;
	`

	var id int
	err := s.db.QueryRow(ctx, q,
		recipe.UserID,
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
		recipe.IsFavorite,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateUserRecipe(ctx context.Context, recipe *types.UserRecipe) error {
	q := `
	UPDATE user_recipes
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
		is_favorite = $13,
		updated_at = NOW()
	WHERE id = $14 AND user_id = $15;
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
		recipe.IsFavorite,
		recipe.RecipeID,
		recipe.UserID,
	)
	return err
}

func (s *Store) DeleteUserRecipe(ctx context.Context, id int, userID string) error {
	q := `
	DELETE FROM user_recipes
	WHERE id = $1 AND user_id = $2;
	`
	_, err := s.db.Exec(ctx, q, id, userID)
	return err
}

func (s *Store) SetUserFavorite(ctx context.Context, id int, userID string, isFavorite bool) error {
	q := `
	UPDATE user_recipes
	SET is_favorite = $1,
		updated_at = NOW()
	WHERE id = $2 AND user_id = $3;
	`
	_, err := s.db.Exec(ctx, q, isFavorite, id, userID)
	return err
}

func (s *Store) AddUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error {
	q := `
	INSERT INTO user_recipes_ingredients
		(recipe_id, item, amount, unit, order_index)
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err := s.db.Exec(ctx, q,
		ingredient.RecipeID, ingredient.IngredientItem, ingredient.IngredientAmount, ingredient.IngredientUnit, ingredient.OrderIndex,
	)
	return err
}

func (s *Store) UpdateUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error {
	q := `
	UPDATE user_recipes_ingredients
	SET item = $1, amount = $2, unit = $3, order_index = $4
	WHERE id = $5 AND recipe_id = $6;
	`
	_, err := s.db.Exec(ctx, q,
		ingredient.IngredientItem, ingredient.IngredientAmount, ingredient.IngredientUnit, ingredient.OrderIndex,
		ingredient.IngredientID, ingredient.RecipeID,
	)
	return err
}

func (s *Store) DeleteUserRecipeIngredient(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_ingredients
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetUserRecipeIngredients(ctx context.Context, recipeID int) ([]types.UserRecipesIngredient, error) {
	q := `
	SELECT
		id, recipe_id, item, amount, unit, order_index
	FROM user_recipes_ingredients
	WHERE recipe_id = $1
	ORDER BY order_index;
	`
	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []types.UserRecipesIngredient
	for rows.Next() {
		var ingredient types.UserRecipesIngredient
		if err := rows.Scan(&ingredient.IngredientID, &ingredient.RecipeID, &ingredient.IngredientItem, &ingredient.IngredientAmount, &ingredient.IngredientUnit, &ingredient.OrderIndex); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (s *Store) AddUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error {
	q := `
	INSERT INTO user_recipes_instructions
		(recipe_id, step_number, instruction)
	VALUES ($1, $2, $3);
	`
	_, err := s.db.Exec(ctx, q,
		instruction.RecipeID, instruction.InstructionStep, instruction.InstructionText,
	)
	return err
}

func (s *Store) UpdateUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error {
	q := `
	UPDATE user_recipes_instructions
	SET step_number = $1, instruction = $2
	WHERE id = $3 AND recipe_id = $4;
	`
	_, err := s.db.Exec(ctx, q,
		instruction.InstructionStep, instruction.InstructionText,
		instruction.InstructionID, instruction.RecipeID,
	)
	return err
}

func (s *Store) DeleteUserRecipeInstruction(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_instructions
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetUserRecipeInstructions(ctx context.Context, recipeID int) ([]types.UserRecipesInstruction, error) {
	q := `
	SELECT
		id, recipe_id, step_number, instruction
	FROM user_recipes_instructions
	WHERE recipe_id = $1
	ORDER BY step_number;
	`
	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []types.UserRecipesInstruction
	for rows.Next() {
		var instruction types.UserRecipesInstruction
		if err := rows.Scan(
			&instruction.InstructionID, &instruction.RecipeID,
			&instruction.InstructionStep, &instruction.InstructionText,
		); err != nil {
			return nil, err
		}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func (s *Store) AddUserRecipeTag(ctx context.Context, tag *types.UserRecipesTag) error {
	q := `
	INSERT INTO user_recipes_tags
		(recipe_id, tag_name)
	VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, q,
		tag.RecipeID, tag.TagName,
	)
	return err
}

func (s *Store) DeleteUserRecipeTag(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_tags
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store) GetUserRecipeTags(ctx context.Context, recipeID int) ([]types.UserRecipesTag, error) {
	q := `
	SELECT
		id, recipe_id, tag_name
	FROM user_recipes_tags
	WHERE recipe_id = $1;
	`
	rows, err := s.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []types.UserRecipesTag
	for rows.Next() {
		var tag types.UserRecipesTag
		if err := rows.Scan(&tag.TagID, &tag.RecipeID, &tag.TagName); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Store) SearchUserRecipesByTag(ctx context.Context, userID string, tag string) ([]types.UserRecipe, error) {
	q := `
	SELECT
		ur.id, ur.name, ur.description, ur.category, ur.difficulty, ur.calories,
		ur.protein, ur.carbs, ur.fat, ur.fiber, ur.prep_time, ur.cook_time,
		ur.image_url, ur.is_favorite, ur.created_at, ur.updated_at
	FROM user_recipes ur
	JOIN user_recipes_tags urt ON ur.id = urt.recipe_id
	WHERE ur.user_id = $1 AND urt.tag_name ILIKE $2
	ORDER BY ur.created_at DESC;
	`
	rows, err := s.db.Query(ctx, q, userID, "%"+tag+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []types.UserRecipe
	for rows.Next() {
		var recipe types.UserRecipe
		err := rows.Scan(
			&recipe.RecipeID, &recipe.RecipeName, &recipe.RecipeDesc, &recipe.RecipesCategory,
			&recipe.RecipesDifficulty, &recipe.RecipesCalories, &recipe.RecipesProtein,
			&recipe.RecipesCarbs, &recipe.RecipesFat, &recipe.RecipesFiber,
			&recipe.PrepTime, &recipe.CookTime, &recipe.RecipesImageURL,
			&recipe.IsFavorite, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
