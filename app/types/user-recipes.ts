interface RecipeCategory {
  id: string;
  name: string;
}

interface RecipeDifficulty {
  id: string;
  name: string;
}

interface UserRecipe {
    id: string;
    name: string;
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

interface UserRecipeIngredient {
    id: string;
    name: string;
    quantity: string;
}

interface UserRecipeStep { 
    id: string;
    description: string;
    order: number;
    recipeId: string;
}

interface UserRecipeTags {
    id: string;
    name: string;
}

interface UserFavoriteRecipe {
    id: string;
    userId: string;
    recipeId: string;
}

type ListUserRecipesResponse = {
  recipes: UserRecipe[];
};


type GetUserRecipeResponse = {
  recipe: UserRecipe;
  ingredients: UserRecipeIngredient[];
  steps: UserRecipeStep[];
  tags: UserRecipeTags[];
  isFavorite: boolean;
  favoriteEntry?: UserFavoriteRecipe;
};

export type {
  RecipeCategory,
  RecipeDifficulty,
  UserRecipe,
  ListUserRecipesResponse,
  GetUserRecipeResponse,
};

