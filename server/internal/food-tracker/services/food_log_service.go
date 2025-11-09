package services

import (
	"context"
	"fmt"
	"time"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/repository"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

type foodLogService struct {
	repo repository.FoodTrackerRepo

}


func NewFoodLogService(repo repository.FoodTrackerRepo) FoodLogService {
	return &foodLogService{
		repo: repo,
	}
}

func parseDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return parsedDate
}

func (s *foodLogService) LogFood(ctx context.Context, userID string, req *types.CreateFoodLogRequest) (*types.FoodLogEntryWithRecipe, error) {
	
	if err := s.ValidateFoodLogRequest(req); err != nil {
		return nil, err
	}

	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if req.SystemRecipeID != nil && req.UserRecipeID != nil {
		return nil, types.ErrInvalidRequest
	}

	entry := &types.FoodLogEntry{
		UserID:     userID,
		EntryID:  0,
		LogDate:   parseDate(req.LogDate),
		SystemRecipeID: req.SystemRecipeID,
		UserRecipeID:  req.UserRecipeID,
		Calories:  req.Calories,
		Protein:   req.Protein,
		Carbs:     req.Carbs,
		Fat:       req.Fat,
		Fiber:     req.Fiber,
		Servings:  req.Servings,
		MealType:  req.MealType,
	}

	entryID, err := s.repo.FoodLogs().CreateFoodLogEntry(ctx, entry)
	if err != nil {
		return nil, types.ErrFailedToCreateLogEntry
	}

	logEntry, err := s.repo.FoodLogs().GetFoodLogEntryByID(ctx, entryID, userID)
	if err != nil {
		return nil, err
	}
	
	return logEntry, nil
}

func (s *foodLogService) UpdateLog(ctx context.Context, id int, userID string, req *types.CreateFoodLogRequest) (*types.FoodLogEntryWithRecipe, error) {
	if id <= 0 || userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if err := s.ValidateFoodLogRequest(req); err != nil {
		return nil, err
	}

	existingEntry, err := s.repo.FoodLogs().GetFoodLogEntryByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if existingEntry.UserID != userID {
		return nil, types.ErrUnauthorized
	}
	
	entry := &types.FoodLogEntry{
		EntryID: 	  id,
		UserID:    userID,
		LogDate:  parseDate(req.LogDate),
		SystemRecipeID: req.SystemRecipeID,
		UserRecipeID:  req.UserRecipeID,
		Calories: req.Calories,
		Protein:  req.Protein,
		Carbs:    req.Carbs,
		Fat:      req.Fat,
		Fiber:    req.Fiber,
		Servings: req.Servings,
		MealType: req.MealType,
		UpdatedAt: time.Now(),
  	}

	if err := s.repo.FoodLogs().UpdateFoodLogEntry(ctx, entry); err != nil {
		return nil, types.ErrFailedToUpdateLogEntry
	}

	return s.repo.FoodLogs().GetFoodLogEntryByID(ctx, id, userID)

}

func (s *foodLogService) DeleteLog(ctx context.Context, id int, userID string) error {
	if id <= 0 || userID == "" {
		return types.ErrInvalidRequest
	}

	entry, err := s.repo.FoodLogs().GetFoodLogEntryByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if entry.UserID != userID {
		return types.ErrUnauthorized
	}

	return s.repo.FoodLogs().DeleteFoodLogEntry(ctx, id, userID)
}

func (s *foodLogService) GetLogsByDate(ctx context.Context, userID string, date string) ([]types.FoodLogEntryWithRecipe, error) {
	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid date format: %w", types.ErrInvalidRequest)
	}

	return s.repo.FoodLogs().GetFoodLogEntriesByDate(ctx, userID, date)

}

func (s *foodLogService) GetLogsByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]types.FoodLogEntryWithRecipe, error) {
	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", types.ErrInvalidRequest)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", types.ErrInvalidRequest)
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end date cannot be before start date: %w", types.ErrInvalidRequest)
	}

	return s.repo.FoodLogs().GetFoodLogEntriesByDateRange(ctx, userID, startDate, endDate)

}

func (s *foodLogService) GetDailyNutrition(ctx context.Context, userID string, date string) (*types.DailyNutritionSummary, error) {
	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid date format: %w", types.ErrInvalidRequest)
	}

	return s.repo.FoodLogs().GetDailySummary(ctx, userID, date)
}

func (s *foodLogService) GetWeeklyNutrition(ctx context.Context, userID string, startDate string) ([]types.DailyNutritionSummary, error) {
	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", types.ErrInvalidRequest)
	}

	return s.repo.FoodLogs().GetWeeklySummary(ctx, userID, startDate)
}

func (s *foodLogService) GetMonthlyNutrition(ctx context.Context, userID string, year int, month int) ([]types.DailyNutritionSummary, error) {
	if userID == "" {
		return nil, types.ErrInvalidRequest
	}

	if year < 2000 || year > time.Now().Year() {
		return nil, fmt.Errorf("invalid year: %w", types.ErrInvalidRequest)
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %w", types.ErrInvalidRequest)
	}

	return s.repo.FoodLogs().GetMonthlySummary(ctx, userID, year, month)
}

func (s *foodLogService) LogRecipe(ctx context.Context, userID string, recipeID int, isSystemRecipe bool, date string, mealType types.MealType) (*types.FoodLogEntryWithRecipe, error) {
	if userID == "" || recipeID <= 0 {
		return nil, types.ErrInvalidID
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	var calories, protein, carbs, fat, fiber int
	var recipeName string
	var systemRecipeID, userRecipeID *int

	if isSystemRecipe {
		recipe, err := s.repo.SystemRecipes().GetSystemRecipeByID(ctx, recipeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch system recipe: %w", err)
		}
		if !recipe.IsActive {
			return nil, fmt.Errorf("recipe is not active")
		}

		systemRecipeID = &recipeID
		recipeName = recipe.RecipeName
		calories = recipe.RecipesCalories
		protein = recipe.RecipesProtein
		carbs = recipe.RecipesCarbs
		fat = recipe.RecipesFat
		fiber = recipe.RecipesFiber
	} else {
		recipe, err := s.repo.UserRecipes().GetUserRecipeByID(ctx, recipeID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user recipe: %w", err)
		}

		userRecipeID = &recipeID
		recipeName = recipe.RecipeName
		calories = recipe.RecipesCalories
		protein = recipe.RecipesProtein
		carbs = recipe.RecipesCarbs
		fat = recipe.RecipesFat
		fiber = recipe.RecipesFiber
	}

	entry := &types.FoodLogEntry{
		UserID:         userID,
		LogDate:        parseDate(date),
		MealType:       mealType,
		SystemRecipeID: systemRecipeID,
		UserRecipeID:   userRecipeID,
		Calories:       calories,
		Protein:        protein,
		Carbs:          carbs,
		Fat:            fat,
		Fiber:          fiber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	entryID, err := s.repo.FoodLogs().CreateFoodLogEntry(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create food log entry: %w", err)
	}

	return &types.FoodLogEntryWithRecipe{
		FoodLogEntry: types.FoodLogEntry{
			EntryID:        entryID,
			UserID:         userID,
			LogDate:        parseDate(date),
			MealType:       mealType,
			SystemRecipeID: systemRecipeID,
			UserRecipeID:   userRecipeID,
			Calories:       calories,
			Protein:        protein,
			Carbs:          carbs,
			Fat:            fat,
			Fiber:          fiber,
			CreatedAt:      entry.CreatedAt,
			UpdatedAt:      entry.UpdatedAt,
		},
		RecipeName:   recipeName,
		RecipeSource: getRecipeSource(isSystemRecipe),
	}, nil
}

func (s *foodLogService) ValidateFoodLogRequest(req *types.CreateFoodLogRequest) error {
	if req == nil {
		return types.ErrInvalidRequest
	}

	if _, err := time.Parse("2006-01-02", req.LogDate); err != nil {
		return fmt.Errorf("invalid log date format: %w", types.ErrInvalidRequest)
	}

	validMealTypes := map[types.MealType]bool{
		types.MealTypeBreakfast: true,
		types.MealTypeLunch:     true,
		types.MealTypeDinner:    true,
		types.MealTypeSnack:     true,
	}
	if !validMealTypes[req.MealType] {
		return types.ErrInvalidMealType
	}

	if req.Calories < 0 || req.Protein < 0 || req.Carbs < 0 || req.Fat < 0 || req.Fiber < 0 || req.Servings <= 0 {
		return types.ErrNurtritionValues
	}

	return nil
}

func getRecipeSource(isSystemRecipe bool) string {
	if isSystemRecipe {
		return "system"
	}
	return "user"
}
 
