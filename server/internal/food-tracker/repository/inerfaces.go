package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

type SystemRecipeRepository interface {
	GetSystemRecipeByID(ctx context.Context, id int) (*types.SystemRecipeDetail, error)
	GetSystemRecipeAll(ctx context.Context, filters types.RecipeFilters) ([]types.SystemRecipe, error)
	CreateSystemRecipe(ctx context.Context, recipe *types.SystemRecipe) (int, error)
	UpdateSystemRecipe(ctx context.Context, recipe *types.SystemRecipe) error
	DeleteSystemRecipe(ctx context.Context, id int) error
	SetActiveSystemRecipe(ctx context.Context, id int, isActive bool) error

	AddSystemRecipesIngredient(ctx context.Context, ingredient *types.SystemRecipesIngredient) error
	UpdateSystemRecipesIngredient(ctx context.Context, ingredient *types.SystemRecipesIngredient) error
	DeleteSystemRecipesIngredient(ctx context.Context, id int) error
	GetSystemRecipesIngredients(ctx context.Context, recipeID int) ([]types.SystemRecipesIngredient, error)

	AddSystemRecipesInstruction(ctx context.Context, instruction *types.SystemRecipesInstruction) error
	UpdateSystemRecipesInstruction(ctx context.Context, instruction *types.SystemRecipesInstruction) error
	DeleteSystemRecipesInstruction(ctx context.Context, id int) error
	GetSystemRecipesInstructions(ctx context.Context, recipeID int) ([]types.SystemRecipesInstruction, error)

	AddSystemRecipesTag(ctx context.Context, tag *types.SystemRecipesTag) error
	DeleteSystemRecipesTag(ctx context.Context, id int) error
	GetSystemRecipesTags(ctx context.Context, recipeID int) ([]types.SystemRecipesTag, error)
	SearchSystemRecipesByTag(ctx context.Context, tag string) ([]types.SystemRecipe, error)
}

type UserRecipeRepository interface {
	GetUserRecipeByID(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error)
	GetAllUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error)
	CreateUserRecipe(ctx context.Context, recipe *types.UserRecipe) (int, error)
	UpdateUserRecipe(ctx context.Context, recipe *types.UserRecipe) error
	DeleteUserRecipe(ctx context.Context, id int, userID string) error
	SetUserFavorite(ctx context.Context, id int, userID string, isFavorite bool) error

	AddUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error
	UpdateUserRecipeIngredient(ctx context.Context, ingredient *types.UserRecipesIngredient) error
	DeleteUserRecipeIngredient(ctx context.Context, id int) error
	GetUserRecipeIngredients(ctx context.Context, recipeID int) ([]types.UserRecipesIngredient, error)

	AddUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error
	UpdateUserRecipeInstruction(ctx context.Context, instruction *types.UserRecipesInstruction) error
	DeleteUserRecipeInstruction(ctx context.Context, id int) error
	GetUserRecipeInstructions(ctx context.Context, recipeID int) ([]types.UserRecipesInstruction, error)

	AddUserRecipeTag(ctx context.Context, tag *types.UserRecipesTag) error
	DeleteUserRecipeTag(ctx context.Context, id int) error
	GetUserRecipeTags(ctx context.Context, recipeID int) ([]types.UserRecipesTag, error)
	SearchUserRecipesByTag(ctx context.Context, userID string, tag string) ([]types.UserRecipe, error)
}

type UserFavoriteRepository interface {
	AddFavorite(ctx context.Context, userID string, recipeID int) error
	RemoveFavorite(ctx context.Context, userID string, recipeID int) error
	GetFavorites(ctx context.Context, userID string) ([]types.SystemRecipe, error)
	IsFavorite(ctx context.Context, userID string, recipeID int) (bool, error)
}

type FoodLogRepository interface {
	CreateFoodLogEntry(ctx context.Context, entry *types.FoodLogEntry) (int, error)
	GetFoodLogEntryByID(ctx context.Context, id int, userID string) (*types.FoodLogEntryWithRecipe, error)
	UpdateFoodLogEntry(ctx context.Context, entry *types.FoodLogEntry) error
	DeleteFoodLogEntry(ctx context.Context, id int, userID string) error

	GetFoodLogEntriesByDate(ctx context.Context, userID string, date string) ([]types.FoodLogEntryWithRecipe, error)
	GetFoodLogEntriesByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]types.FoodLogEntryWithRecipe, error)
	GetFoodLogEntriesByMealType(ctx context.Context, userID string, date string, mealType types.MealType) ([]types.FoodLogEntryWithRecipe, error)

	GetDailySummary(ctx context.Context, userID string, date string) (*types.DailyNutritionSummary, error)
	GetWeeklySummary(ctx context.Context, userID string, startDate string) ([]types.DailyNutritionSummary, error)
	GetMonthlySummary(ctx context.Context, userID string, year int, month int) ([]types.DailyNutritionSummary, error)
}

type RecipeSearchRepository interface {
	SearchAll(ctx context.Context, userID string, query types.SearchQuery) ([]types.UserAllRecipesView, error)
	GetAllRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserAllRecipesView, error)
}

type NutritionGoalsRepository interface {
	GetUserNutritionGoals(ctx context.Context, userID string) (*types.NutritionGoals, error)
	CreateUserNutritionGoals(ctx context.Context, goals *types.NutritionGoals) error
	UpdateUserNutritionGoals(ctx context.Context, goals *types.NutritionGoals) error
	DeleteUserNutritionGoals(ctx context.Context, userID string) error

}

type FoodTrackerRepo interface {
	SystemRecipes() SystemRecipeRepository
	UserRecipes() UserRecipeRepository
	FoodLogs() FoodLogRepository
	RecipeSearch() RecipeSearchRepository
	UserFavorites() UserFavoriteRepository
	NutritionGoals() NutritionGoalsRepository
}
