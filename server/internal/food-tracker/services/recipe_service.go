package services

import (
	"context"
	"fmt"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/repository"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

type recipeService struct {
	repo repository.FoodTrackerRepo

}

func NewRecipeService(repo repository.FoodTrackerRepo) RecipeService {
	return &recipeService{
		repo: repo,
	}
}

func (s *recipeService) GetSystemRecipe(ctx context.Context, id int) (*types.SystemRecipeDetail, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid recipe ID")
	}
	recipe, err := s.repo.SystemRecipes().GetSystemRecipeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ingredients, err := s.repo.SystemRecipes().GetSystemRecipesIngredients(ctx, id)
	if err != nil {
		return nil, err
	}
	instructions, err := s.repo.SystemRecipes().GetSystemRecipesInstructions(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := s.repo.SystemRecipes().GetSystemRecipesTags(ctx, id)
	if err != nil {
		return nil, err
	}

	return &types.SystemRecipeDetail{
		SystemRecipe: recipe.SystemRecipe,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:        tags,
	}, nil
}

func (s *recipeService) ListSystemRecipes(ctx context.Context, filters types.RecipeFilters) ([]types.SystemRecipe, error) {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	recipes, err := s.repo.SystemRecipes().GetSystemRecipeAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (s *recipeService) CreateSystemRecipe(ctx context.Context, req *types.CreateRecipeRequest) (*types.SystemRecipeDetail, error) {
	return nil, nil
}

func (s *recipeService) UpdateSystemRecipe(ctx context.Context, id int, req *types.CreateRecipeRequest) (*types.SystemRecipeDetail, error) {
	return nil, nil
}

func (s *recipeService) DeleteSystemRecipe(ctx context.Context, id int) error {
	return nil
}

func (s *recipeService)	GetUserRecipe(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error) {
	return nil, nil
}

func (s *recipeService)	ListUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error) {
	return nil, nil
}

func (s *recipeService)	CreateUserRecipe(ctx context.Context, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error) {
	return nil, nil
}

func (s *recipeService)	UpdateUserRecipe(ctx context.Context, id int, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error) {
	return nil, nil
}

func (s *recipeService)	DeleteUserRecipe(ctx context.Context, id int, userID string) error {
	return nil
}

func (s *recipeService)	ToggleFavorite(ctx context.Context, userID string, recipeID int) error {
	return nil
}

func (s *recipeService)	GetFavorites(ctx context.Context, userID string) ([]types.SystemRecipe, error) {
	return nil, nil
}

func (s *recipeService)	SearchRecipes(ctx context.Context, userID string, query types.SearchQuery) ([]types.UserAllRecipesView, error) {
	return nil, nil
}






