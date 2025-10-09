type RecipeCategory = 'breakfast' | 'lunch' | 'dinner' | 'snack' | 'dessert';
type RecipeDifficulty = 'easy' | 'medium' | 'hard';
type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snack';

interface SystemRecipe {
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
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface SystemRecipeIngredient {
  id: number;
  recipe_id: number;
  item: string;
  amount: number;
  unit: string;
  order_index: number;
}

interface SystemRecipeInstruction { 
  id: number;
  recipe_id: number;
  step_number: number;
  instruction: string;
}

interface SystemRecipeTag {
  id: number;
  recipe_id: number;
  tag_name: string;
}

interface SystemRecipeDetail extends SystemRecipe {
  ingredients: SystemRecipeIngredient[];
  instructions: SystemRecipeInstruction[];
  tags: SystemRecipeTag[];
}

interface UserRecipe {
  id: number;
  user_id: string;
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
  is_favorite: boolean;
  created_at: string;
  updated_at: string;
}

interface UserRecipeIngredient {
  id: number;
  recipe_id: number;
  item: string;
  amount: number;
  unit: string;
  order_index: number;
}

interface UserRecipeInstruction { 
  id: number;
  recipe_id: number;
  step_number: number;
  instruction: string;
}

interface UserRecipeTag {
  id: number;
  recipe_id: number;
  tag_name: string;
}

interface UserRecipeDetail extends UserRecipe {
  ingredients: UserRecipeIngredient[];
  instructions: UserRecipeInstruction[];
  tags: UserRecipeTag[];
}

interface UserFavoriteRecipe {
  id: number;
  user_id: string;
  recipe_id: number;
  created_at: string;
}

interface UserAllRecipesView {
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
  recipe_source?: 'system' | 'user';
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
  user_id: string;
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

interface IngredientNutrition {
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
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
    ingredient_id?: number;
    item: string;
    amount: number;
    unit: string;
    order_index: number;
  }>;
  instructions: Array<{
    instruction_id?: number;
    step_number: number;
    instruction?: string;
    text?: string;
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

interface RecipeFilters {
  category?: RecipeCategory;
  difficulty?: RecipeDifficulty;
  max_calories?: number;
  min_protein?: number;
  max_prep_time?: number;
  tags?: string[];
  is_favorite?: boolean;
  search_term?: string;
  limit?: number;
  offset?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

interface SearchQuery {
  term: string;
  category?: RecipeCategory;
  difficulty?: RecipeDifficulty;
  max_calories?: number;
  min_protein?: number;
  max_prep_time?: number;
  tags?: string[];
  include_system?: boolean;
  include_user?: boolean;
  favorites_only?: boolean;
  limit?: number;
  offset?: number;
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
  recipes: UserAllRecipesView[];
  limit: number;
  offset: number;
}

interface GetFavoritesResponse {
  favorites: UserAllRecipesView[];
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


interface ApiError {
  error: string;
  status?: number;
}

export type {
  RecipeCategory,
  RecipeDifficulty,
  MealType,
  
  SystemRecipe,
  SystemRecipeIngredient,
  SystemRecipeInstruction,
  SystemRecipeTag,
  SystemRecipeDetail,
  
  UserRecipe,
  UserRecipeIngredient,
  UserRecipeInstruction,
  UserRecipeTag,
  UserRecipeDetail,
  UserFavoriteRecipe,
  UserAllRecipesView,
  
  FoodLogEntry,
  FoodLogEntryWithRecipe,
  DailyNutritionSummary,
  
  NutritionGoals,
  NutritionComparison,
  IngredientNutrition,
  
  CreateRecipeRequest,
  CreateFoodLogRequest,
  LogRecipeRequest,
  RecipeFilters,
  SearchQuery,
  
  ListSystemRecipesResponse,
  ListUserRecipesResponse,
  SearchRecipesResponse,
  GetFavoritesResponse,
  GetLogsByDateResponse,
  GetLogsByDateRangeResponse,
  GetWeeklyNutritionResponse,
  GetMonthlyNutritionResponse,
  GetNutritionComparisonResponse,
  ApiError,
};