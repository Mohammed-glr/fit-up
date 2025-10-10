
type RecipeCategory = 'breakfast' | 'lunch' | 'dinner' | 'snack' | 'dessert';
type RecipeDifficulty = 'easy' | 'medium' | 'hard';
type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snack';

interface Recipe {
  id: number;
  name: string;
  description: string;
  category: RecipeCategory;
  difficulty: RecipeDifficulty;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  prep_time: number;
  cook_time: number;
  image_url: string;
  servings: number;
  created_at: string;
  updated_at: string;
}

interface SystemRecipe extends Recipe {
  is_active: boolean;
}

interface UserRecipe extends Recipe {
  user_id: string;
  is_favorite: boolean;
}

interface RecipeIngredient {
  id: number;
  recipe_id: number;
  item: string;
  amount: number;
  unit: string;
  order_index: number;
}

interface RecipeInstruction {
  id: number;
  recipe_id: number;
  step_number: number;
  instruction: string;
}

interface RecipeTag {
  id: number;
  recipe_id: number;
  tag_name: string;
}

interface SystemRecipeDetail extends SystemRecipe {
  ingredients: RecipeIngredient[];
  instructions: RecipeInstruction[];
  tags: RecipeTag[];
}

interface UserRecipeDetail extends UserRecipe {
  ingredients: RecipeIngredient[];
  instructions: RecipeInstruction[];
  tags: RecipeTag[];
}

interface RecipeView {
  source: 'system' | 'user';
  id: number;
  name: string;
  category: RecipeCategory;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  prep_time: number;
  servings: number;
  image_url: string;
  user_id?: string;
  is_favorite: boolean;
}


interface FoodLogEntry {
  id: number;
  user_id: string;
  log_date: string;
  system_recipe_id?: number;
  user_recipe_id?: number;
  meal_type: MealType;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  servings: number;
  created_at: string;
  updated_at: string;
}

interface FoodLogEntryWithRecipe extends FoodLogEntry {
  recipe_name?: string;
  recipe_source?: string;
}


interface DailyNutritionSummary {
  user_id: string;
  log_date: string;
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
  total_fiber: number;
  total_entries: number;
}

interface NutritionGoals {
  user_id?: string;
  calories_goal: number;
  protein_goal: number;
  carbs_goal: number;
  fat_goal: number;
  fiber_goal: number;
}

interface NutritionComparison {
  calories_percent: number;
  protein_percent: number;
  carbs_percent: number;
  fat_percent: number;
  fiber_percent: number;
  is_over_calories: boolean;
  is_meeting_protein: boolean;
}

interface MacroDistribution {
  protein_percent: number;
  carbs_percent: number;
  fat_percent: number;
}

interface CreateRecipeRequest {
  name: string;
  description: string;
  category: RecipeCategory;
  difficulty: RecipeDifficulty;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  prep_time: number;
  cook_time: number;
  image_url: string;
  servings: number;
  ingredients: Array<{
    item: string;
    amount: number;
    unit: string;
    order_index: number;
  }>;
  instructions: Array<{
    step_number: number;
    instruction: string;
  }>;
  tags: string[];
}

interface CreateFoodLogRequest {
  log_date: string;
  meal_type: MealType;
  system_recipe_id?: number;
  user_recipe_id?: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  servings: number;
}

interface LogRecipeRequest {
  recipe_id: number;
  is_system_recipe: boolean;
  date: string;
  meal_type: MealType;
}

interface ListSystemRecipesResponse {
  recipes: SystemRecipe[];
  limit: number;
  offset: number;
}

interface ListUserRecipesResponse {
  recipes: UserRecipe[];
  limit: number;
  offset: number;
}

interface SearchRecipesResponse {
  recipes: RecipeView[];
  limit: number;
  offset: number;
}

interface GetFavoritesResponse {
  favorites: RecipeView[];
}

interface GetLogsByDateResponse {
  date: string;
  logs: FoodLogEntryWithRecipe[];
}

interface GetLogsByDateRangeResponse {
  start_date: string;
  end_date: string;
  logs: FoodLogEntryWithRecipe[];
}

interface GetWeeklyNutritionResponse {
  start_date: string;
  summaries: DailyNutritionSummary[];
}

interface GetMonthlyNutritionResponse {
  year: number;
  month: number;
  summaries: DailyNutritionSummary[];
}

interface GetNutritionComparisonResponse {
  date: string;
  summary: DailyNutritionSummary;
  goals: NutritionGoals;
  comparison: NutritionComparison;
}

interface GetNutritionInsightsResponse {
  date: string;
  summary: DailyNutritionSummary;
  goals: NutritionGoals;
  comparison: NutritionComparison;
  macro_distribution: MacroDistribution;
}

interface ApiError {
  error: string;
}

export type {
  RecipeCategory,
  RecipeDifficulty,
  MealType,
  
  Recipe,
  SystemRecipe,
  UserRecipe,
  RecipeIngredient,
  RecipeInstruction,
  RecipeTag,
  SystemRecipeDetail,
  UserRecipeDetail,
  RecipeView,
  
  FoodLogEntry,
  FoodLogEntryWithRecipe,
  
  DailyNutritionSummary,
  NutritionGoals,
  NutritionComparison,
  MacroDistribution,
  
  CreateRecipeRequest,
  CreateFoodLogRequest,
  LogRecipeRequest,
  
  ListSystemRecipesResponse,
  ListUserRecipesResponse,
  SearchRecipesResponse,
  GetFavoritesResponse,
  GetLogsByDateResponse,
  GetLogsByDateRangeResponse,
  GetWeeklyNutritionResponse,
  GetMonthlyNutritionResponse,
  GetNutritionComparisonResponse,
  GetNutritionInsightsResponse,
  ApiError,
};