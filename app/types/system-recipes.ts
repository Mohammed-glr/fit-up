
interface RecipeCategory {
  id: string;
  name: string;
}

interface RecipeDifficulty {
  id: string;
  name: string;
}

interface SystemRecipe {
  id: string;
  name: string;
  description: string;
  category: RecipeCategory;
  difficulty: RecipeDifficulty;
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  caloriesPerServing: number;
  proteinPerServing: number;
  carbsPerServing: number;
  fatPerServing: number;
  imageUrl?: string;
}

interface SystemRecipeIngredient {
  id: string;
  name: string;
  quantity: string;
}

interface SystemRecipeStep { 
  id: string;
  description: string;
  order: number;
  recipeId: string;
}

interface SystemRecipeTags {
  id: string;
  name: string;
}

type ListSystemRecipesResponse = {
  recipes: SystemRecipe[];
};

type GetSystemRecipeResponse = {
  recipe: SystemRecipe;
  ingredients: SystemRecipeIngredient[];
  steps: SystemRecipeStep[];
  tags: SystemRecipeTags[];
};
export type {
  RecipeCategory,
  RecipeDifficulty,
  SystemRecipe,
  ListSystemRecipesResponse,
  GetSystemRecipeResponse,
};


