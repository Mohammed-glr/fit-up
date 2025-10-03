package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) GetUserRecipeByID(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error) {
	q := `
	SELECT
		ur.id, ur.name, ur.description, ur.category, ur.difficulty, ur.calories,
		ur.protein, ur.carbs, ur.fat, ur.fiber, ur.prep_time, ur.cook_time,
		ur.image_url, ur.is_favorite, ur.is_active, ur.created_at, ur.updated_at
	FROM user_recipes ur
	WHERE ur.id = $1 AND ur.user_id = $2 AND ur.is_active = TRUE;
	`

	var recipe types.UserRecipeDetail
	err := s.db.QueryRow(ctx, q, id, userID).Scan(
		&recipe.RecipeID, &recipe.RecipeName, &recipe.RecipeDesc, &recipe.RecipesCategory,
		&recipe.RecipesDifficulty, &recipe.RecipesCalories, &recipe.RecipesProtein,
		&recipe.RecipesCarbs, &recipe.RecipesFat, &recipe.RecipesFiber,
		&recipe.PrepTime, &recipe.CookTime, &recipe.RecipesImageURL,
		&recipe.IsFavorite, &recipe.IsFavorite, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (s *Store)	GetAllUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error) {
	q := `
	SELECT
		ur.id, ur.name, ur.description, ur.category, ur.difficulty, ur.calories,
		ur.protein, ur.carbs, ur.fat, ur.fiber, ur.prep_time, ur.cook_time,
		ur.image_url, ur.is_favorite, ur.is_active, ur.created_at, ur.updated_at
	FROM user_recipes ur
	WHERE ur.user_id = $1 AND ur.is_active = TRUE
	ORDER BY ur.created_at DESC
	LIMIT $2 OFFSET $3;
	`
	rows, err := s.db.Query(ctx, q, userID, filters.Limit, filters.Offset)
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
			&recipe.IsFavorite, &recipe.IsFavorite, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (s *Store)	CreateUserRecipe(ctx context.Context, recipe *types.UserRecipe) (int, error) {
	q := `
	INSERT INTO user_recipes
		(user_id, name, description, category, difficulty, calories,
		 protein, carbs, fat, fiber, prep_time, cook_time, image_url, is_favorite, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id;
	`

	var id int
	err := s.db.QueryRow(ctx, q,
		recipe.UserID, recipe.RecipeName, recipe.RecipeDesc, recipe.RecipesCategory,
		recipe.RecipesDifficulty, recipe.RecipesCalories, recipe.RecipesProtein,
		recipe.RecipesCarbs, recipe.RecipesFat, recipe.RecipesFiber,
		recipe.PrepTime, recipe.CookTime, recipe.RecipesImageURL,
		recipe.IsFavorite, recipe.IsFavorite, recipe.CreatedAt, recipe.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store)	UpdateUserRecipe(ctx context.Context, recipe *types.UserRecipe) error {
	q := `
	UPDATE user_recipes
	SET name = $1, description = $2, category = $3, difficulty = $4, calories = $5,
		protein = $6, carbs = $7, fat = $8, fiber = $9, prep_time = $10,
		cook_time = $11, image_url = $12, is_favorite = $13, updated_at = NOW()
	WHERE id = $14 AND user_id = $15 AND is_active = TRUE;
	`

	_, err := s.db.Exec(ctx, q,
		recipe.RecipeName, recipe.RecipeDesc, recipe.RecipesCategory,
		recipe.RecipesDifficulty, recipe.RecipesCalories, recipe.RecipesProtein,
		recipe.RecipesCarbs, recipe.RecipesFat, recipe.RecipesFiber,
		recipe.PrepTime, recipe.CookTime, recipe.RecipesImageURL,
		recipe.IsFavorite, recipe.CreatedAt, recipe.UpdatedAt,
		recipe.RecipeID, recipe.UserID,
	)
	return err
}

func (s *Store)	DeleteUserRecipe(ctx context.Context, id int, userID string) error {
	q := `
	DELETE FROM user_recipes
	WHERE id = $1 AND user_id = $2;
	`
	_, err := s.db.Exec(ctx, q, id, userID)
	return err
}

func (s *Store)	SetUserFavorite(ctx context.Context, id int, userID string, isFavorite bool) error {
	q := `
	UPDATE user_recipes
	SET is_favorite = $1, updated_at = NOW()
	WHERE id = $2 AND user_id = $3 AND is_active = TRUE;
	`
	_, err := s.db.Exec(ctx, q, isFavorite, id, userID)
	return err
}

func (s *Store)	AddUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error {
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

func (s *Store)	UpdateUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error {
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

func (s *Store)	DeleteUserRecipeIngredient(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_ingredients
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store)	GetUserRecipeIngredients(ctx context.Context, recipeID int) ([]types.UserRecipesIngredient, error) {
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

func (s *Store)	AddUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error {
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

func (s *Store)	UpdateUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error {
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

func (s *Store)	DeleteUserRecipeInstruction(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_instructions
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store)	GetUserRecipeInstructions(ctx context.Context, recipeID int) ([]types.UserRecipesInstruction, error) {
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

func (s *Store)	AddUserRecipeTag(ctx context.Context, tag *types.UserRecipesTag) error {
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

func (s *Store)	DeleteUserRecipeTag(ctx context.Context, id int) error {
	q := `
	DELETE FROM user_recipes_tags
	WHERE id = $1;
	`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *Store)	GetUserRecipeTags(ctx context.Context, recipeID int) ([]types.UserRecipesTag, error) {
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

func (s *Store)	SearchUserRecipesByTag(ctx context.Context, userID string, tag string) ([]types.UserRecipe, error) {
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

