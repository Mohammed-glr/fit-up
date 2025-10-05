package services

import (
	"context"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/repository"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

type RecipeService interface {
	GetSystemRecipe(ctx context.Context, id int) (*types.SystemRecipeDetail, error)
	ListSystemRecipes(ctx context.Context, filters types.RecipeFilters) ([]types.SystemRecipe, error)
	CreateSystemRecipe(ctx context.Context, req *types.CreateRecipeRequest) (*types.SystemRecipeDetail, error)
	UpdateSystemRecipe(ctx context.Context, id int, req *types.CreateRecipeRequest) (*types.SystemRecipeDetail, error)
	DeleteSystemRecipe(ctx context.Context, id int) error

	GetUserRecipe(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error)
	ListUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error)
	CreateUserRecipe(ctx context.Context, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error)
	UpdateUserRecipe(ctx context.Context, id int, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error)
	DeleteUserRecipe(ctx context.Context, id int, userID string) error

	ToggleFavorite(ctx context.Context, userID string, recipeID int) error
	GetFavorites(ctx context.Context, userID string) ([]types.SystemRecipe, error)

	SearchRecipes(ctx context.Context, userID string, query types.SearchQuery) ([]types.UserAllRecipesView, error)
}
type FoodLogService interface {
	LogFood(ctx context.Context, userID string, req *types.CreateFoodLogRequest) (*types.FoodLogEntryWithRecipe, error)
	UpdateLog(ctx context.Context, id int, userID string, req *types.CreateFoodLogRequest) (*types.FoodLogEntryWithRecipe, error)
	DeleteLog(ctx context.Context, id int, userID string) error

	GetLogsByDate(ctx context.Context, userID string, date string) ([]types.FoodLogEntryWithRecipe, error)
	GetLogsByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]types.FoodLogEntryWithRecipe, error)

	GetDailyNutrition(ctx context.Context, userID string, date string) (*types.DailyNutritionSummary, error)
	GetWeeklyNutrition(ctx context.Context, userID string, startDate string) ([]types.DailyNutritionSummary, error)
	GetMonthlyNutrition(ctx context.Context, userID string, year int, month int) ([]types.DailyNutritionSummary, error)

	LogRecipe(ctx context.Context, userID string, recipeID int, isSystemRecipe bool, date string, mealType types.MealType) (*types.FoodLogEntryWithRecipe, error)
}

type NutritionAnalyzer interface {
	CalculateRecipeNutrition(ingredients []types.SystemRecipesIngredient) (calories, protein, carbs, fat, fiber int, err error)
	CalculateMealNutrition(entries []types.FoodLogEntry) (calories, protein, carbs, fat, fiber int)
	GetNutritionGoals(userID string) (*types.NutritionGoals, error)
	CompareToGoals(summary *types.DailyNutritionSummary, goals *types.NutritionGoals) *types.NutritionComparison
}



type FoodTrackerService interface {
	Recipes()  RecipeService
	FoodLogs() FoodLogService
	Nutrition() NutritionAnalyzer
}

type Service struct {
	repo             repository.FoodTrackerRepo
	recipeService    RecipeService
	foodLogService   FoodLogService
	nutritionAnalyzer NutritionAnalyzer
}

func NewService(repo repository.FoodTrackerRepo) FoodTrackerService {
	return &Service{
		repo:             repo,
		recipeService:    NewRecipeService(repo),
		foodLogService:   NewFoodLogService(repo),
		// nutritionAnalyzer: NewNutritionAnalyzer(repo),
	}
}

func (s *Service) Recipes() RecipeService {
	return s.recipeService
}

func (s *Service) FoodLogs() FoodLogService {
	return s.foodLogService
}

func (s *Service) Nutrition() NutritionAnalyzer {
	return s.nutritionAnalyzer
}


