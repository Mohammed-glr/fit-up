package services

import (
	"fmt"
	"strings"

	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)


type SimpleIngredientDB struct {
	nutritionData map[string]*types.IngredientNutrition
}

func NewSimpleIngredientDB() IngredientNutritionDB {
	db := &SimpleIngredientDB{
		nutritionData: make(map[string]*types.IngredientNutrition),
	}

	db.addIngredient("chicken breast", 165, 31, 0, 4, 0)
	db.addIngredient("brown rice", 370, 8, 77, 3, 4)
	db.addIngredient("broccoli", 34, 3, 7, 0, 3)
	db.addIngredient("egg", 155, 13, 1, 11, 0)
	db.addIngredient("salmon", 208, 20, 0, 13, 0)
	db.addIngredient("sweet potato", 86, 2, 20, 0, 3)
	db.addIngredient("oats", 389, 17, 66, 7, 11)
	db.addIngredient("banana", 89, 1, 23, 0, 3)
	db.addIngredient("apple", 52, 0, 14, 0, 2)
	db.addIngredient("almond", 579, 21, 22, 50, 13)
	db.addIngredient("milk", 61, 3, 5, 3, 0)
	db.addIngredient("yogurt", 59, 10, 4, 0, 0)
	db.addIngredient("beef", 250, 26, 0, 15, 0)
	db.addIngredient("pasta", 371, 13, 75, 2, 3)
	db.addIngredient("bread", 265, 9, 49, 3, 3)
	db.addIngredient("cheese", 402, 25, 1, 33, 0)
	db.addIngredient("olive oil", 884, 0, 0, 100, 0)
	db.addIngredient("avocado", 160, 2, 9, 15, 7)
	db.addIngredient("spinach", 23, 3, 4, 0, 2)
	db.addIngredient("tomato", 18, 1, 4, 0, 1)

	return db
}

func (db *SimpleIngredientDB) addIngredient(name string, calories, protein, carbs, fat, fiber int) {
	db.nutritionData[strings.ToLower(name)] = &types.IngredientNutrition{
		Calories: calories,
		Protein:  protein,
		Carbs:    carbs,
		Fat:      fat,
		Fiber:    fiber,
	}
}

// GetIngredientNutrition returns nutrition information for an ingredient
// Amount is assumed to be in grams or ml, and nutrition data is scaled accordingly
func (db *SimpleIngredientDB) GetIngredientNutrition(ingredient string, amount float64, unit string) (*types.IngredientNutrition, error) {
	// Normalize ingredient name
	ingredientKey := strings.ToLower(strings.TrimSpace(ingredient))

	// Look up nutrition data (per 100g/ml)
	baseNutrition, exists := db.nutritionData[ingredientKey]
	if !exists {
		// If not found, return a default approximation
		// In a real implementation, this would call an external API or return an error
		return &types.IngredientNutrition{
			Calories: 100,
			Protein:  5,
			Carbs:    15,
			Fat:      3,
			Fiber:    2,
		}, nil
	}

	// Scale nutrition based on amount (assuming amount is in grams or ml)
	// The base nutrition is per 100g/ml
	scale := amount / 100.0

	// Handle different units
	switch strings.ToLower(unit) {
	case "kg":
		scale = amount * 10 // 1kg = 1000g, so scale * 10
	case "oz":
		scale = (amount * 28.35) / 100.0 // 1oz = 28.35g
	case "lb":
		scale = (amount * 453.592) / 100.0 // 1lb = 453.592g
	case "cup":
		// Approximate: 1 cup ≈ 240g for liquids, varies for solids
		scale = (amount * 240.0) / 100.0
	case "tbsp", "tablespoon":
		scale = (amount * 15.0) / 100.0 // 1 tbsp ≈ 15g
	case "tsp", "teaspoon":
		scale = (amount * 5.0) / 100.0 // 1 tsp ≈ 5g
	case "g", "gram", "ml":
		// Already in grams or ml
		scale = amount / 100.0
	default:
		// If unit is unknown, assume grams
		scale = amount / 100.0
	}

	return &types.IngredientNutrition{
		Calories: int(float64(baseNutrition.Calories) * scale),
		Protein:  int(float64(baseNutrition.Protein) * scale),
		Carbs:    int(float64(baseNutrition.Carbs) * scale),
		Fat:      int(float64(baseNutrition.Fat) * scale),
		Fiber:    int(float64(baseNutrition.Fiber) * scale),
	}, nil
}

// AddIngredient allows adding new ingredients to the database (useful for testing or admin features)
func (db *SimpleIngredientDB) AddIngredient(name string, calories, protein, carbs, fat, fiber int) error {
	if name == "" {
		return fmt.Errorf("ingredient name cannot be empty")
	}

	db.addIngredient(name, calories, protein, carbs, fat, fiber)
	return nil
}
