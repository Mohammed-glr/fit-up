package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) SystemRecipes() SystemRecipeRepository {
	return s
}

func (s *Store) UserRecipes() UserRecipeRepository {
	return s
}

func (s *Store) FoodLogs() FoodLogRepository {
	return s
}

func (s *Store) RecipeSearch() RecipeSearchRepository {
	return s
}

func (s *Store) UserFavorites() UserFavoriteRepository {
	return s
}

func (s *Store) NutritionGoals() NutritionGoalsRepository {
	return s
}
