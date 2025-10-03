package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (s *Store) SearchAll(ctx context.Context, userID string, query types.SearchQuery) ([]types.UserAllRecipesView, error) {
	q := `
		SELECT 
			source, id, name, category, calories, protein, carbs, fat, fiber, 
			prep_time, servings, image_url, user_id, is_favorite
		FROM (
			SELECT 
				'system' AS source,
				sr.id,
				sr.name,
				sr.category,
				sr.calories,
				sr.protein,
				sr.carbs,
				sr.fat,
				sr.fiber,
				sr.prep_time,
				sr.servings,
				sr.image_url,
				NULL::text AS user_id,
				EXISTS(SELECT 1 FROM user_favorite_recipes ufr WHERE ufr.recipe_id = sr.id AND ufr.user_id = $2) AS is_favorite
			FROM system_recipes sr
			LEFT JOIN system_recipe_tags srt ON sr.id = srt.recipe_id
			WHERE sr.is_active = TRUE AND (sr.name ILIKE $1 OR sr.description ILIKE $1 OR srt.tag ILIKE $1)
			UNION
			SELECT 
				'user' AS source,
				ur.id,
				ur.name,
				ur.category,
				ur.calories,
				ur.protein,
				ur.carbs,
				ur.fat,
				ur.fiber,
				ur.prep_time,
				ur.servings,
				ur.image_url,
				ur.user_id,
				ur.is_favorite
			FROM user_recipes ur
			LEFT JOIN user_recipe_tags urt ON ur.id = urt.recipe_id
			WHERE ur.user_id = $2 AND (ur.name ILIKE $1 OR ur.description ILIKE $1 OR urt.tag ILIKE $1)
		) AS combined
		ORDER BY name
	`
	rows, err := s.db.Query(ctx, q, "%"+query.Term+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []types.UserAllRecipesView
	for rows.Next() {
		var recipe types.UserAllRecipesView
		if err := rows.Scan(
			&recipe.Source,
			&recipe.ID,
			&recipe.Name,
			&recipe.Category,
			&recipe.Calories,
			&recipe.Protein,
			&recipe.Carbs,
			&recipe.Fat,
			&recipe.Fiber,
			&recipe.PrepTime,
			&recipe.Servings,
			&recipe.ImageURL,
			&recipe.UserID,
			&recipe.IsFavorite,
		); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (s *Store) GetAllRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserAllRecipesView, error) {
	q := `
		SELECT 
			source, id, name, category, calories, protein, carbs, fat, fiber,
			prep_time, servings, image_url, user_id, is_favorite
		FROM (
			SELECT 
				'system' AS source,
				sr.id,
				sr.name,
				sr.category,
				sr.calories,
				sr.protein,
				sr.carbs,
				sr.fat,
				sr.fiber,
				sr.prep_time,
				sr.servings,
				sr.image_url,
				NULL::text AS user_id,
				EXISTS(SELECT 1 FROM user_favorite_recipes ufr WHERE ufr.recipe_id = sr.id AND ufr.user_id = $1) AS is_favorite
			FROM system_recipes sr
			WHERE sr.is_active = TRUE
			UNION
			SELECT 
				'user' AS source,
				ur.id,
				ur.name,
				ur.category,
				ur.calories,
				ur.protein,
				ur.carbs,
				ur.fat,
				ur.fiber,
				ur.prep_time,
				ur.servings,
				ur.image_url,
				ur.user_id,
				ur.is_favorite
			FROM user_recipes ur
			WHERE ur.user_id = $1
		) AS combined
		ORDER BY name
	`

	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []types.UserAllRecipesView
	for rows.Next() {
		var recipe types.UserAllRecipesView
		if err := rows.Scan(
			&recipe.Source,
			&recipe.ID,
			&recipe.Name,
			&recipe.Category,
			&recipe.Calories,
			&recipe.Protein,
			&recipe.Carbs,
			&recipe.Fat,
			&recipe.Fiber,
			&recipe.PrepTime,
			&recipe.Servings,
			&recipe.ImageURL,
			&recipe.UserID,
			&recipe.IsFavorite,
		); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
