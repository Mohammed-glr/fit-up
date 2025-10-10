import { API } from '../endpoints';
import { executeAPI } from '../client';
import {
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
} from '@/types/food-tracker'

const systemRecipeService = {
    List: async (): Promise<ListSystemRecipesResponse> => {
        const response = await executeAPI(API.recipe.system.list());
        return response.data as ListSystemRecipesResponse;
    },

    Retrieve: async (id: number): Promise<any> => {
        const response = await executeAPI(API.recipe.system.retrieve(id));
        return response.data;
    },

    Search: async (): Promise<SearchRecipesResponse> => {
        const response = await executeAPI(API.recipe.system.search());
        return response.data as SearchRecipesResponse;
    },

    Create: async (data: CreateRecipeRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.system.create(), data);
        return response.data;
    },

    Update: async (id: number, data: CreateRecipeRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.system.update(id), data);
        return response.data;
    },

    Delete: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.system.delete(id));
    },
}

const userRecipeService = {
    List: async (): Promise<ListUserRecipesResponse> => {
        const response = await executeAPI(API.recipe.user.list());
        return response.data as ListUserRecipesResponse;
    },

    Create: async (data: CreateRecipeRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.user.create(), data);
        return response.data;
    },

    Retrieve: async (id: number): Promise<any> => {
        const response = await executeAPI(API.recipe.user.retrieve(id));
        return response.data;
    },

    Update: async (id: number, data: CreateRecipeRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.user.update(id), data);
        return response.data;
    },

    Delete: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.user.delete(id));
    },

    GetFavorites: async (): Promise<GetFavoritesResponse> => {
        const response = await executeAPI(API.recipe.user.getFavorites());
        return response.data as GetFavoritesResponse;
    },

    ToggleFavorite: async (recipeId: number): Promise<void> => {
        await executeAPI(API.recipe.user.toggleFavorite(recipeId));
    },
}

const foodLogService = {
    Log: async (data: CreateFoodLogRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.logs.log(), data);
        return response.data;
    },

    LogRecipe: async (data: LogRecipeRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.logs.logRecipe(), data);
        return response.data;
    },

    GetLogsByDate: async (date: string): Promise<GetLogsByDateResponse> => {
        const response = await executeAPI(API.recipe.logs.getlogsByDate(date));
        return response.data as GetLogsByDateResponse;
    },

    GetLogsInRange: async (): Promise<GetLogsByDateRangeResponse> => {
        const response = await executeAPI(API.recipe.logs.getLogsInRange());
        return response.data as GetLogsByDateRangeResponse;
    },

    GetFoodLogEntry: async (id: number): Promise<any> => {
        const response = await executeAPI(API.recipe.logs.getFoodLogEntry(id));
        return response.data;
    },

    UpdateFoodLogEntry: async (id: number, data: CreateFoodLogRequest): Promise<any> => {
        const response = await executeAPI(API.recipe.logs.updateFoodLogEntry(id), data);
        return response.data;
    },

    DeleteFoodLogEntry: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.logs.deleteFoodLogEntry(id));
    },
}

const nutritionService = {
    GetDailySummary: async (date: string): Promise<any> => {
        const response = await executeAPI(API.recipe.nutrition.getDailySummary(date));
        return response.data;
    },

    GetWeeklySummary: async (): Promise<GetWeeklyNutritionResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getWeeklySummary());
        return response.data as GetWeeklyNutritionResponse;
    }, 

    GetMonthlySummary: async (): Promise<GetMonthlyNutritionResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getMonthlySummary());
        return response.data as GetMonthlyNutritionResponse;
    },

    GetGoals: async (): Promise<any> => {
        const response = await executeAPI(API.recipe.nutrition.getGoals());
        return response.data;
    },

    UpdateGoals: async (data: any): Promise<any> => {
        const response = await executeAPI(API.recipe.nutrition.updateGoals(), data);
        return response.data;
    },

    CompareWithGoals: async (date: string): Promise<GetNutritionComparisonResponse> => {
        const response = await executeAPI(API.recipe.nutrition.compareWithGoals(date));
        return response.data as GetNutritionComparisonResponse;
    },

    GetNutritionInsights: async (date: string): Promise<GetNutritionInsightsResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getNutritionInsights(date));
        return response.data as GetNutritionInsightsResponse;
    },
}

export default { systemRecipeService, userRecipeService, foodLogService, nutritionService };