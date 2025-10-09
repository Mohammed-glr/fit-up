
interface RecipeCategory {
  breakfast: string;
  lunch: string;
  dinner: string;
  snack: string;
  dessert: string;
}

interface RecipeDifficulty {
  easy: string;
  medium: string;
  hard: string;
}

interface SystemRecipe {
  id: string;
  name: string;
  description: string;
  category: keyof RecipeCategory;
  difficulty: keyof RecipeDifficulty;
  calories: number;
  protein: number;
  carbs: number;
  fats: number;
  fiber: number;
  prep_time: number; 
  cook_time: number;
  image_url: string;
  servings: number;
  created_at: string;
  updated_at: string;
  is_active: boolean;
}

interface SystemRecipeIngredient {
  id: string;
  recipe_id: string;
  item: string;
  amount: string;
  unit: string;
  order_index: number;
}

interface SystemRecipeStep { 
  id: string;
  recipe_id: string;
  step_number: number;
  instruction: string;
}

interface SystemRecipeTags {
  id: string;
  tag_name: string;
  recipe_id: string;
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

interface RecipeCategory {
  breakfast: string;
  lunch: string;
  dinner: string;
  snack: string;
  dessert: string;
}

interface RecipeDifficulty {
  easy: string;
  medium: string;
  hard: string;
}

interface UserRecipe {
    id: string;
    user_id: string;
    name: string;
    description: string;
    category: keyof RecipeCategory;
    difficulty: keyof RecipeDifficulty;
    calories: number;
    protein: number;
    carbs: number;
    fats: number;
    fiber: number;
    prep_time: number; 
    cook_time: number;
    image_url: string;
    servings: number;
    created_at: string;
    updated_at: string;
    is_favorite: boolean;
}

interface UserRecipeIngredient {
    id: string;
    recipe_id: string;
    item: string;
    amount: string;
    unit: string;
    order_index: number;
}

interface UserRecipeStep { 
    id: string;
    instruction: string;
    recipe_id: string;
    step_number: number;
}

interface UserRecipeTags {
    id: string;
    tag_name: string;
    recipe_id: string;
}

interface UserFavoriteRecipe {
    id: string;
    user_id: string;
    recipe_id: string;
    created_at: string;
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

interface MealType {
  breakfast: string;
  lunch: string;
  dinner: string;
  snack: string;
}

interface FoodLogEntry {
  id: string;
  user_id: string;
  log_date: string;
  system_recipe_id?: string;
  user_recipe_id?: string;
  meal_type: keyof MealType;
  calories: number;
  protein: number;
  carbs: number;
  fats: number;
  fiber: number;
  servings: number;
  created_at: string;
  updated_at: string;
}

interface DailyNutritionSummary {
    user_id: string;
    log_date: string;
    total_calories: number;
    total_protein: number;
    total_carbs: number;
    total_fats: number;
    total_fiber: number;
    total_entries: number;
}

type ListFoodLogEntriesResponse = {
  entries: FoodLogEntry[];
};

type GetDailyNutritionSummaryResponse = {
  summary: DailyNutritionSummary;
};

export type {
  RecipeCategory,
  RecipeDifficulty,
  SystemRecipe,
  SystemRecipeIngredient,
  SystemRecipeStep,
  SystemRecipeTags,
  ListSystemRecipesResponse,
  GetSystemRecipeResponse,
  UserRecipe,
  UserRecipeIngredient,
  UserRecipeStep,
  UserRecipeTags,
  UserFavoriteRecipe,
  ListUserRecipesResponse,
  GetUserRecipeResponse,
  MealType,
  FoodLogEntry,
  DailyNutritionSummary,
  ListFoodLogEntriesResponse,
  GetDailyNutritionSummaryResponse
};