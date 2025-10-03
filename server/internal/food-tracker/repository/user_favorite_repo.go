package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) AddFavorite(ctx context.Context, userID string, recipeID int) error {
	q := `
	INSERT INTO user_favorites (user_id, recipe_id)
	VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, q, userID, recipeID)
	return err
}

func (s *Store)	RemoveFavorite(ctx context.Context, userID string, recipeID int) error {
	q := `
	DELETE FROM user_favorites
	WHERE user_id = $1 AND recipe_id = $2;
	`
	_, err := s.db.Exec(ctx, q, userID, recipeID)
	return err
}

func (s *Store)	GetFavorites(ctx context.Context, userID string) ([]types.SystemRecipe, error) {
	q := `
	SELECT sr.*
	FROM system_recipes sr
	JOIN user_favorites uf ON sr.id = uf.recipe_id
	WHERE uf.user_id = $1;
	`
	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []types.SystemRecipe
	for rows.Next() {
		var recipe types.SystemRecipe
		if err := rows.Scan(
			&recipe.RecipeID,
			&recipe.CookTime,
			&recipe.CreatedAt,
			&recipe.IsActive,
			&recipe.PrepTime,
			&recipe.RecipeDesc,
			&recipe.RecipeName,
			&recipe.RecipesCalories,
			&recipe.RecipesDifficulty,
			&recipe.RecipesFat,
			&recipe.RecipesFiber,
			&recipe.RecipesImageURL,
			&recipe.RecipesProtein,
			&recipe.RecipesCategory,
			&recipe.RecipesCarbs,
			&recipe.UpdatedAt,
			); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (s *Store)	IsFavorite(ctx context.Context, userID string, recipeID int) (bool, error) {
	q := `
	SELECT EXISTS (
		SELECT 1
		FROM user_favorite_recipes
		WHERE user_id = $1 AND recipe_id = $2
	);
	`

	var exists bool
	err := s.db.QueryRow(ctx, q, userID, recipeID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

