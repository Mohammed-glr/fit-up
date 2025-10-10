import { API } from '../endpoints';
import { executeAPI } from '../client';

const recipeService = {
    ListSystemRecipes: () => executeAPI(API.recipe.system.list()),
    RetrieveSystemRecipe: (id: number) => executeAPI(API.recipe.system.retrieve(id)),
    SearchSystemRecipes: () => executeAPI(API.recipe.system.search()),
    CreateSystemRecipe: (data: any) => executeAPI(API.recipe.system.create(), data),
    UpdateSystemRecipe: (id: number, data: any) => executeAPI(API.recipe.system.update(id), data),
    DeleteSystemRecipe: (id: number) => executeAPI(API.recipe.system.delete(id)),  

    ListUserRecipes: () => executeAPI(API.recipe.user.list()),
    CreateUserRecipe: (data: any) => executeAPI(API.recipe.user.create(), data),
    RetrieveUserRecipe: (id: number) => executeAPI(API.recipe.user.retrieve(id)),
    UpdateUserRecipe: (id: number, data: any) => executeAPI(API.recipe.user.update(id), data),
    DeleteUserRecipe: (id: number) => executeAPI(API.recipe.user.delete(id)),

    GetFavoriteRecipes: () => executeAPI(API.recipe.user.getFavorites()),
    ToggleFavoriteRecipe: (recipeId: number) => executeAPI(API.recipe.user.toggleFavorite(recipeId)),

    LogRecipe: (data: any) => executeAPI(API.recipe.logs.logRecipe(), data),
    GetLogsByDate: (date: string) => executeAPI(API.recipe.logs.getlogsByDate(date)),
    GetLogsInRange: () => executeAPI(API.recipe.logs.getLogsInRange()),
    GetFoodLogEntry: (id: number) => executeAPI(API.recipe.logs.getFoodLogEntry(id)),
    UpdateFoodLogEntry: (id: number, data: any) => executeAPI(API.recipe.logs.updateFoodLogEntry(id), data),
    DeleteFoodLogEntry: (id: number) => executeAPI(API.recipe.logs.deleteFoodLogEntry(id)),

    GetDailyNutritionSummary: (date: string) => executeAPI(API.recipe.nutrition.getDailySummary(date)),
    GetWeeklyNutritionSummary: () => executeAPI(API.recipe.nutrition.getWeeklySummary()),
    GetMonthlyNutritionSummary: () => executeAPI(API.recipe.nutrition.getMonthlySummary()),
    GetNutritionGoals: () => executeAPI(API.recipe.nutrition.getGoals()),
    UpdateNutritionGoals: (data: any) => executeAPI(API.recipe.nutrition.updateGoals(), data),
    CompareWithNutritionGoals: (date: string) => executeAPI(API.recipe.nutrition.compareWithGoals(date)),
    GetNutritionInsights: (date: string) => executeAPI(API.recipe.nutrition.getNutritionInsights(date)),
    
};

export default recipeService;