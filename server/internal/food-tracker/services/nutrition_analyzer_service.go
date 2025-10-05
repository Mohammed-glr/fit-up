package services

import (
	"context"
	"fmt"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/repository"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

type nutritionAnalyzerService struct {
	repo repository.FoodTrackerRepo

	ingredientDB IngredientNutritionDB
}

type IngredientNutritionDB interface {
	GetIngredientNutrition(ingredient string, amount float64, unit string) (*types.IngredientNutrition, error)
}

func NewNutritionAnalyzer(repo repository.FoodTrackerRepo, ingredientDB IngredientNutritionDB) NutritionAnalyzer {
	return &nutritionAnalyzerService{
		repo:         repo,
		ingredientDB: ingredientDB,
	}
}

func (s *nutritionAnalyzerService) CalculateRecipeNutrition(ingredients []types.SystemRecipesIngredient) (calories, protein, carbs, fat, fiber int, err error) {
	if len(ingredients) == 0 {
		return 0, 0, 0, 0, 0, fmt.Errorf("no ingredients provided")
	}

	totalCalories := 0
	totalProtein := 0
	totalCarbs := 0
	totalFat := 0
	totalFiber := 0

	for _, ing := range ingredients {
		nutrition, err := s.ingredientDB.GetIngredientNutrition(ing.IngredientItem, ing.IngredientAmount, ing.IngredientUnit)
		if err != nil {

			return 0, 0, 0, 0, 0, fmt.Errorf("failed to get nutrition for %s: %w", ing.IngredientItem, err)
		}

		totalCalories += nutrition.Calories
		totalProtein += nutrition.Protein
		totalCarbs += nutrition.Carbs
		totalFat += nutrition.Fat
		totalFiber += nutrition.Fiber
	}

	return totalCalories, totalProtein, totalCarbs, totalFat, totalFiber, nil
}

func (s *nutritionAnalyzerService) CalculateMealNutrition(entries []types.FoodLogEntry) (calories, protein, carbs, fat, fiber int) {
	totalCalories := 0
	totalProtein := 0
	totalCarbs := 0
	totalFat := 0
	totalFiber := 0

	for _, entry := range entries {
		totalCalories += entry.Calories
		totalProtein += entry.Protein
		totalCarbs += entry.Carbs
		totalFat += entry.Fat
		totalFiber += entry.Fiber
	}

	return totalCalories, totalProtein, totalCarbs, totalFat, totalFiber
}

func (s *nutritionAnalyzerService) GetNutritionGoals(ctx context.Context, userID string) (*types.NutritionGoals, error) {
	if userID == "" {
		return nil, types.ErrInvalidID
	}


	goals, err := s.repo.NutritionGoals().GetUserNutritionGoals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition goals: %w", err)
	}
	
	return goals, nil
}

func (s *nutritionAnalyzerService) CompareToGoals(summary *types.DailyNutritionSummary, goals *types.NutritionGoals) *types.NutritionComparison {
	if summary == nil || goals == nil {
		return nil
	}

	comparison := &types.NutritionComparison{
		CaloriesPercent: calculatePercent(summary.TotalCalories, goals.CaloriesGoal),
		ProteinPercent:  calculatePercent(summary.TotalProtein, goals.ProteinGoal),
		CarbsPercent:    calculatePercent(summary.TotalCarbs, goals.CarbsGoal),
		FatPercent:      calculatePercent(summary.TotalFat, goals.FatGoal),
		FiberPercent:    calculatePercent(summary.TotalFiber, goals.FiberGoal),
	}

	comparison.IsOverCalories = summary.TotalCalories > goals.CaloriesGoal
	comparison.IsMeetingProtein = summary.TotalProtein >= goals.ProteinGoal

	return comparison
}

func calculatePercent(actual, goal int) float64 {
	if goal == 0 {
		return 0
	}
	return (float64(actual) / float64(goal)) * 100
}



func (s *nutritionAnalyzerService) GetMacroDistribution(calories, protein, carbs, fat int) map[string]float64 {
	proteinCals := protein * 4
	carbsCals := carbs * 4
	fatCals := fat * 9

	totalCals := float64(proteinCals + carbsCals + fatCals)
	if totalCals == 0 {
		return map[string]float64{
			"protein": 0,
			"carbs":   0,
			"fat":     0,
		}
	}

	return map[string]float64{
		"protein": (float64(proteinCals) / totalCals) * 100,
		"carbs":   (float64(carbsCals) / totalCals) * 100,
		"fat":     (float64(fatCals) / totalCals) * 100,
	}
}

func (s *nutritionAnalyzerService) CalculateCaloriesFromMacros(protein, carbs, fat int) int {
	return (protein * 4) + (carbs * 4) + (fat * 9)
}

func (s *nutritionAnalyzerService) ValidateNutritionGoals(goals *types.NutritionGoals) error {
	if goals == nil {
		return fmt.Errorf("goals cannot be nil")
	}

	if goals.CaloriesGoal < 1000 || goals.CaloriesGoal > 10000 {
		return fmt.Errorf("calories goal must be between 1000 and 10000")
	}

	if goals.ProteinGoal < 0 || goals.ProteinGoal > 500 {
		return fmt.Errorf("protein goal must be between 0 and 500g")
	}

	if goals.CarbsGoal < 0 || goals.CarbsGoal > 1000 {
		return fmt.Errorf("carbs goal must be between 0 and 1000g")
	}

	if goals.FatGoal < 0 || goals.FatGoal > 300 {
		return fmt.Errorf("fat goal must be between 0 and 300g")
	}

	if goals.FiberGoal < 0 || goals.FiberGoal > 100 {
		return fmt.Errorf("fiber goal must be between 0 and 100g")
	}

	calculatedCals := s.CalculateCaloriesFromMacros(goals.ProteinGoal, goals.CarbsGoal, goals.FatGoal)
	difference := abs(calculatedCals - goals.CaloriesGoal)
	
	tolerance := goals.CaloriesGoal / 10
	if difference > tolerance {
		return fmt.Errorf("macro calories (%d) don't match calorie goal (%d)", calculatedCals, goals.CaloriesGoal)
	}

	return nil
}

func (s *nutritionAnalyzerService) GetNutritionInsights(summary *types.DailyNutritionSummary, goals *types.NutritionGoals) []string {
	insights := []string{}

	if summary == nil || goals == nil {
		return insights
	}

	comparison := s.CompareToGoals(summary, goals)
	if comparison == nil {
		return insights
	}

	if comparison.IsOverCalories {
		over := summary.TotalCalories - goals.CaloriesGoal
		insights = append(insights, fmt.Sprintf("You exceeded your calorie goal by %d calories", over))
	} else if comparison.CaloriesPercent < 80 {
		insights = append(insights, "You're significantly under your calorie goal")
	} else if comparison.CaloriesPercent >= 95 && comparison.CaloriesPercent <= 105 {
		insights = append(insights, "Great job hitting your calorie goal!")
	}

	if comparison.IsMeetingProtein {
		insights = append(insights, "Excellent protein intake!")
	} else if comparison.ProteinPercent < 80 {
		needed := goals.ProteinGoal - summary.TotalProtein
		insights = append(insights, fmt.Sprintf("Try to add %dg more protein", needed))
	}

	if comparison.FiberPercent < 50 {
		insights = append(insights, "Consider adding more fiber-rich foods")
	} else if comparison.FiberPercent >= 100 {
		insights = append(insights, "Great fiber intake!")
	}

	distribution := s.GetMacroDistribution(summary.TotalCalories, summary.TotalProtein, summary.TotalCarbs, summary.TotalFat)
	
	if distribution["protein"] < 15 {
		insights = append(insights, "Your protein intake is low relative to total calories")
	}
	if distribution["fat"] > 40 {
		insights = append(insights, "Your fat intake is high relative to total calories")
	}

	return insights
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}


func (s *nutritionAnalyzerService) SetNutritionGoals(ctx context.Context, goals *types.NutritionGoals) error {
	existing, err := s.repo.NutritionGoals().GetUserNutritionGoals(ctx, goals.UserID)
	if err != nil && err != types.ErrNotFound {
		return fmt.Errorf("failed to check existing goals: %w", err)
	}
	if existing != nil {
		return s.repo.NutritionGoals().UpdateUserNutritionGoals(ctx, goals)
	}
	return s.repo.NutritionGoals().CreateUserNutritionGoals(ctx, goals)
}

