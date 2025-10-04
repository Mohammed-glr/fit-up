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
		return nil, types.ErrInvalidID
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
		Tags:         tags,
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
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	recipe := &types.SystemRecipe{
		RecipeID:          0,
		RecipeName:        req.Name,
		RecipeDesc:        req.Description,
		RecipesCategory:   req.Category,
		RecipesDifficulty: req.Difficulty,
		RecipesCalories:   req.Calories,
		RecipesProtein:    req.Protein,
		RecipesCarbs:      req.Carbs,
		RecipesFat:        req.Fat,
		RecipesFiber:      req.Fiber,
		RecipesImageURL:   req.ImageURL,
		PrepTime:          req.PrepTime,
		CookTime:          req.CookTime,
		IsActive:          true,
		Servings:          req.Servings,
		CreatedAt:         "now()",
		UpdatedAt:         "now()",
	}
	recipeID, err := s.repo.SystemRecipes().CreateSystemRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	var ingredients []types.SystemRecipesIngredient
	var instructions []types.SystemRecipesInstruction
	var tags []types.SystemRecipesTag

	for _, ing := range req.Ingredients {
		ingredient := &types.SystemRecipesIngredient{
			RecipeID:         recipeID,
			IngredientID:     0,
			IngredientItem:   ing.Item,
			IngredientAmount: ing.Amount,
			IngredientUnit:   ing.Unit,
			OrderIndex:       ing.OrderIndex,
		}
		if err := s.repo.SystemRecipes().AddSystemRecipesIngredient(ctx, ingredient); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, *ingredient)
	}

	for _, instr := range req.Instructions {
		instruction := &types.SystemRecipesInstruction{
			RecipeID:        recipeID,
			InstructionID:   0,
			InstructionText: instr.Instruction,
			InstructionStep: instr.StepNumber,
		}
		if err := s.repo.SystemRecipes().AddSystemRecipesInstruction(ctx, instruction); err != nil {
			return nil, err
		}
		instructions = append(instructions, *instruction)
	}

	for _, tagName := range req.Tags {
		tag := &types.SystemRecipesTag{
			RecipeID: recipeID,
			TagID:    0,
			TagName:  tagName,
		}
		if err := s.repo.SystemRecipes().AddSystemRecipesTag(ctx, tag); err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	recipe.RecipeID = recipeID

	return &types.SystemRecipeDetail{
		SystemRecipe: *recipe,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:         tags,
	}, nil
}

func (s *recipeService) UpdateSystemRecipe(ctx context.Context, id int, req *types.CreateRecipeRequest) (*types.SystemRecipeDetail, error) {
	if id <= 0 {
		return nil, types.ErrInvalidID
	}

	if req == nil {
		return nil, types.ErrInvalidRequest
	}


	recipe := &types.SystemRecipe{
		RecipeID:          id,
		RecipeName:        req.Name,
		RecipeDesc:        req.Description,
		RecipesCategory:   req.Category,
		RecipesDifficulty: req.Difficulty,
		RecipesCalories:   req.Calories,
		RecipesProtein:    req.Protein,
		RecipesCarbs:      req.Carbs,
		RecipesFat:        req.Fat,
		RecipesFiber:      req.Fiber,
		RecipesImageURL:   req.ImageURL,
		PrepTime:          req.PrepTime,
		CookTime:          req.CookTime,
		IsActive:          true,
		Servings:          req.Servings,
		UpdatedAt:         "now()",
	}
	if err := s.repo.SystemRecipes().UpdateSystemRecipe(ctx, recipe); err != nil {
		return nil, err
	}

	existingIngredients, err := s.repo.SystemRecipes().GetSystemRecipesIngredients(ctx, id)
	if err != nil {
		return nil, err
	}

	existingIngMap := make(map[int]bool)
	for _, ing := range existingIngredients {
		existingIngMap[ing.IngredientID] = true
	}

	var ingredients []types.SystemRecipesIngredient
	requestedIngMap := make(map[int]bool)

	for _, ing := range req.Ingredients {
		ingredient := &types.SystemRecipesIngredient{
			RecipeID:         id,
			IngredientID:     ing.IngredientID,
			IngredientItem:   ing.Item,
			IngredientAmount: ing.Amount,
			IngredientUnit:   ing.Unit,
			OrderIndex:       ing.OrderIndex,
		}
		if ingredient.IngredientID == 0 {
			if err := s.repo.SystemRecipes().AddSystemRecipesIngredient(ctx, ingredient); err != nil {
				return nil, err
			}
		} else {
			if err := s.repo.SystemRecipes().UpdateSystemRecipesIngredient(ctx, ingredient); err != nil {
				return nil, err
			}
			requestedIngMap[ingredient.IngredientID] = true
		}
		ingredients = append(ingredients, *ingredient)
	}

	for ingID := range existingIngMap {
		if !requestedIngMap[ingID] {
			if err := s.repo.SystemRecipes().DeleteSystemRecipesIngredient(ctx, ingID); err != nil {
				return nil, err
			}
		}
	}

	existingInstructions, err := s.repo.SystemRecipes().GetSystemRecipesInstructions(ctx, id)
	if err != nil {
		return nil, err
	}

	existingintrgMap := make(map[int]bool)
	for _, instr := range existingInstructions {
		existingintrgMap[instr.InstructionID] = true
	}

	var instructions []types.SystemRecipesInstruction
	requestedInstrMap := make(map[int]bool)

	for _, instr := range req.Instructions {
		instruction := &types.SystemRecipesInstruction{
			RecipeID:        id,
			InstructionID:   instr.InstructionID,
			InstructionText: instr.Instruction,
			InstructionStep: instr.StepNumber,
		}

		if instruction.InstructionID == 0 {
			if err := s.repo.SystemRecipes().AddSystemRecipesInstruction(ctx, instruction); err != nil {
				return nil, err
			}
		} else {
			if err := s.repo.SystemRecipes().UpdateSystemRecipesInstruction(ctx, instruction); err != nil {
				return nil, err
			}
			requestedInstrMap[instruction.InstructionID] = true
		}
		instructions = append(instructions, *instruction)
	}

	for instrID := range existingintrgMap {
		if !requestedInstrMap[instrID] {
			if err := s.repo.SystemRecipes().DeleteSystemRecipesInstruction(ctx, instrID); err != nil {
				return nil, err
			}
		}
	}

	existingTags, err := s.repo.SystemRecipes().GetSystemRecipesTags(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, tag := range existingTags {
		if err := s.repo.SystemRecipes().DeleteSystemRecipesTag(ctx, tag.TagID); err != nil {
			return nil, err
		}
	}

	var tags []types.SystemRecipesTag
	for _, tagName := range req.Tags {
		tag := &types.SystemRecipesTag{
			RecipeID: id,
			TagID:    0,
			TagName:  tagName,
		}
		if err := s.repo.SystemRecipes().AddSystemRecipesTag(ctx, tag); err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	return &types.SystemRecipeDetail{
		SystemRecipe: *recipe,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:         tags,
	}, nil
}

func (s *recipeService) DeleteSystemRecipe(ctx context.Context, id int) error {
	if id <= 0 {
		return types.ErrInvalidID
	}
	return s.repo.SystemRecipes().DeleteSystemRecipe(ctx, id)
}

func (s *recipeService) GetUserRecipe(ctx context.Context, id int, userID string) (*types.UserRecipeDetail, error) {
	if id <= 0 {
		return nil, types.ErrInvalidID
	}

	userRecipe, err := s.repo.UserRecipes().GetUserRecipeByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	ingredient, err := s.repo.UserRecipes().GetUserRecipeIngredients(ctx, id)
	if err != nil {
		return nil, err
	}

	intruction, err := s.repo.UserRecipes().GetUserRecipeInstructions(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := s.repo.UserRecipes().GetUserRecipeTags(ctx, id)
	if err != nil {
		return nil, err
	}

	return &types.UserRecipeDetail{
		UserRecipe:   userRecipe.UserRecipe,
		Ingredients:  ingredient,
		Instructions: intruction,
		Tags:         tags,
	}, nil
}

func (s *recipeService) ListUserRecipes(ctx context.Context, userID string, filters types.RecipeFilters) ([]types.UserRecipe, error) {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	recipes, err := s.repo.UserRecipes().GetAllUserRecipes(ctx, userID, filters)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (s *recipeService) CreateUserRecipe(ctx context.Context, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}
	if userID == "" {
		return nil, types.ErrInvalidID
	}

	recipe := &types.UserRecipe{
		RecipeID:          0,
		RecipeName:        req.Name,
		RecipeDesc:        req.Description,
		RecipesCategory:   req.Category,
		RecipesDifficulty: req.Difficulty,
		RecipesCalories:   req.Calories,
		RecipesProtein:    req.Protein,
		RecipesCarbs:      req.Carbs,
		RecipesFat:        req.Fat,
		RecipesFiber:      req.Fiber,
		RecipesImageURL:   req.ImageURL,
		PrepTime:          req.PrepTime,
		CookTime:          req.CookTime,
		IsFavorite:        false,
		Servings:          req.Servings,
		UserID:            userID,
		CreatedAt:         "now()",
		UpdatedAt:         "now()",
	}
	recipeID, err := s.repo.UserRecipes().CreateUserRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	var ingredients []types.UserRecipesIngredient
	var instructions []types.UserRecipesInstruction
	var tags []types.UserRecipesTag

	for _, ing := range req.Ingredients {
		ingredient := &types.UserRecipesIngredient{
			RecipeID:         recipeID,
			IngredientID:     0,
			IngredientItem:   ing.Item,
			IngredientAmount: ing.Amount,
			IngredientUnit:   ing.Unit,
			OrderIndex:       ing.OrderIndex,
		}
		if err := s.repo.UserRecipes().AddUserRecipeIngredient(ctx, ingredient); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, *ingredient)
	}
	for _, instr := range req.Instructions {
		instruction := &types.UserRecipesInstruction{
			RecipeID:        recipeID,
			InstructionID:   0,
			InstructionText: instr.Instruction,
			InstructionStep: instr.StepNumber,
		}
		if err := s.repo.UserRecipes().AddUserRecipeInstruction(ctx, instruction); err != nil {
			return nil, err
		}
		instructions = append(instructions, *instruction)
	}
	for _, tagName := range req.Tags {
		tag := &types.UserRecipesTag{
			RecipeID: recipeID,
			TagID:    0,
			TagName:  tagName,
		}
		if err := s.repo.UserRecipes().AddUserRecipeTag(ctx, tag); err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	return &types.UserRecipeDetail{
		UserRecipe:   *recipe,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:         tags,
	}, nil
}

func (s *recipeService) UpdateUserRecipe(ctx context.Context, id int, userID string, req *types.CreateRecipeRequest) (*types.UserRecipeDetail, error) {
	if id <= 0 {
		return nil, types.ErrInvalidID
	}
	if req == nil {
		return nil, types.ErrInvalidRequest
	}
	if userID == "" {
		return nil, types.ErrInvalidID
	}

	recipe := &types.UserRecipe{
		RecipeID:          id,
		RecipeName:        req.Name,
		RecipeDesc:        req.Description,
		RecipesCategory:   req.Category,
		RecipesDifficulty: req.Difficulty,
		RecipesCalories:   req.Calories,
		RecipesProtein:    req.Protein,
		RecipesCarbs:      req.Carbs,
		RecipesFat:        req.Fat,
		RecipesFiber:      req.Fiber,
		RecipesImageURL:   req.ImageURL,
		PrepTime:          req.PrepTime,
		CookTime:          req.CookTime,
		IsFavorite:        false,
		Servings:          req.Servings,
		UserID:            userID,
		UpdatedAt:         "now()",
	}
	if err := s.repo.UserRecipes().UpdateUserRecipe(ctx, recipe); err != nil {
		return nil, err
	}

	existingIngredients, err := s.repo.UserRecipes().GetUserRecipeIngredients(ctx, id)
	if err != nil {
		return nil, err
	}

	existingIngMap := make(map[int]bool)
	for _, ing := range existingIngredients {
		existingIngMap[ing.IngredientID] = true
	}

	var ingredients []types.UserRecipesIngredient
	requestedIngMap := make(map[int]bool)

	for _, ing := range req.Ingredients {
		ingredient := &types.UserRecipesIngredient{
			RecipeID:         id,
			IngredientID:     ing.IngredientID,
			IngredientItem:   ing.Item,
			IngredientAmount: ing.Amount,
			IngredientUnit:   ing.Unit,
			OrderIndex:       ing.OrderIndex,
		}
		if ingredient.IngredientID == 0 {
			if err := s.repo.UserRecipes().AddUserRecipeIngredient(ctx, ingredient); err != nil {
				return nil, err
			}
		} else {
			if err := s.repo.UserRecipes().UpdateUserRecipeIngredient(ctx, ingredient); err != nil {
				return nil, err
			}
			requestedIngMap[ingredient.IngredientID] = true
		}
		ingredients = append(ingredients, *ingredient)
	}

	for ingID := range existingIngMap {
		if !requestedIngMap[ingID] {
			if err := s.repo.UserRecipes().DeleteUserRecipeIngredient(ctx, ingID); err != nil {
				return nil, err
			}
		}
	}

	existingInstructions, err := s.repo.UserRecipes().GetUserRecipeInstructions(ctx, id)
	if err != nil {
		return nil, err
	}

	existingintrgMap := make(map[int]bool)
	for _, instr := range existingInstructions {
		existingintrgMap[instr.InstructionID] = true
	}

	var instructions []types.UserRecipesInstruction
	requestedInstrMap := make(map[int]bool)

	for _, intr := range req.Instructions {
		instruction := &types.UserRecipesInstruction{
			RecipeID:        id,
			InstructionID:   intr.InstructionID,
			InstructionStep: intr.StepNumber,
			InstructionText: intr.Instruction,
		}
		if instruction.InstructionID == 0 {
			if err := s.repo.UserRecipes().AddUserRecipeInstruction(ctx, instruction); err != nil {
				return nil, err
			}
		} else {
			if err := s.repo.UserRecipes().UpdateUserRecipeInstruction(ctx, instruction); err != nil {
				return nil, err
			}
			requestedInstrMap[instruction.InstructionID] = true
		}
		instructions = append(instructions, *instruction)
	}

	for instrID := range existingintrgMap {
		if !requestedInstrMap[instrID] {
			if err := s.repo.UserRecipes().DeleteUserRecipeInstruction(ctx, instrID); err != nil {
				return nil, err
			}
		}
	}

	existingTags, err := s.repo.UserRecipes().GetUserRecipeTags(ctx, id)
	if err != nil {
		return nil, err
	}
	
	for _, tag := range existingTags {
		if err := s.repo.UserRecipes().DeleteUserRecipeInstruction(ctx, tag.TagID); err != nil {
			return nil, err
		}
	}

	var tags []types.UserRecipesTag
	for _, tagName := range req.Tags {
		tag := &types.UserRecipesTag{
			RecipeID: id,
			TagName:  tagName,
			TagID:    0,
		}
		if err := s.repo.UserRecipes().AddUserRecipeTag(ctx, tag); err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	return &types.UserRecipeDetail{
		UserRecipe:   *recipe,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:         tags,
	}, nil
}

func (s *recipeService) DeleteUserRecipe(ctx context.Context, id int, userID string) error {
	if id <= 0 {
		return types.ErrInvalidID
	}
	return s.repo.UserRecipes().DeleteUserRecipe(ctx, id, userID)
}

func (s *recipeService) ToggleFavorite(ctx context.Context, userID string, recipeID int) error {
	if recipeID <= 0 {
		return types.ErrInvalidID
	}
	isFav, err := s.repo.UserFavorites().IsFavorite(ctx, userID, recipeID)
	if err != nil {
		return err
	}
	if isFav {
		return s.repo.UserFavorites().RemoveFavorite(ctx, userID, recipeID)
	}
	return s.repo.UserFavorites().AddFavorite(ctx, userID, recipeID)
}

func (s *recipeService) GetFavorites(ctx context.Context, userID string) ([]types.SystemRecipe, error) {
	if userID == "" {
		return nil, types.ErrInvalidID
	}

	return s.repo.UserFavorites().GetFavorites(ctx, userID)
}

func (s *recipeService) SearchRecipes(ctx context.Context, userID string, query types.SearchQuery) ([]types.UserAllRecipesView, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	return s.repo.RecipeSearch().SearchAll(ctx, userID, query)
}


func (s *recipeService) ValidateRecipeRequest(req *types.CreateRecipeRequest) error {
	if req == nil {
		return types.ErrInvalidRequest
	}
	if req.Name == "" {
		return fmt.Errorf("recipe name is required")
	}
	if req.Calories < 0 || req.Protein < 0 || req.Carbs < 0 || req.Fat < 0 {
		return fmt.Errorf("nutrition values cannot be negative")
	}
	if req.PrepTime < 0 || req.CookTime < 0 {
		return fmt.Errorf("time values cannot be negative")
	}
	if req.Servings <= 0 {
		return fmt.Errorf("servings must be greater than 0")
	}
	if len(req.Ingredients) == 0 {
		return fmt.Errorf("recipe must have at least one ingredient")
	}
	if len(req.Instructions) == 0 {
		return fmt.Errorf("recipe must have at least one instruction")
	}
	return nil
}


func (s *recipeService) GetRecipeNutritionPerServing(ctx context.Context, recipe *types.SystemRecipe) map[string]float64 {
	if recipe.Servings <= 0 {
		return map[string]float64{}
	}

	return map[string]float64{
		"calories": float64(recipe.RecipesCalories) / float64(recipe.Servings),
		"protein":  float64(recipe.RecipesProtein) / float64(recipe.Servings),
		"carbs":    float64(recipe.RecipesCarbs) / float64(recipe.Servings),
		"fat":      float64(recipe.RecipesFat) / float64(recipe.Servings),
		"fiber":    float64(recipe.RecipesFiber) / float64(recipe.Servings),
	}
}