import { API } from '../endpoints';
import { executeAPI } from '../client';
import {
    CreateRecipeRequest,
    CreateFoodLogRequest,
    LogRecipeRequest,
    RecipeListParams,
    RecipeSearchParams,
    UpsertNutritionGoalsRequest,
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
    UserRecipeDetail,
    SystemRecipeDetail,
    FoodLogEntryWithRecipe,
    DailyNutritionSummary,
    NutritionGoals,
} from '@/types/food-tracker';

const sanitizeParams = <T extends Record<string, unknown>>(params?: T): Partial<T> | undefined => {
    if (!params) {
        return undefined;
    }

    return Object.fromEntries(
        Object.entries(params).filter(([, value]) => value !== undefined && value !== null)
    ) as Partial<T>;
};

const systemRecipeService = {
    List: async (params?: RecipeListParams): Promise<ListSystemRecipesResponse> => {
        const response = await executeAPI(
            API.recipe.system.list(),
            undefined,
            { params: sanitizeParams(params) }
        );
        return response.data as ListSystemRecipesResponse;
    },

    Retrieve: async (id: number): Promise<SystemRecipeDetail> => {
        const response = await executeAPI(API.recipe.system.retrieve(id));
        return response.data as SystemRecipeDetail;
    },

    Search: async (params?: RecipeSearchParams): Promise<SearchRecipesResponse> => {
        const response = await executeAPI(
            API.recipe.search(),
            undefined,
            { params: sanitizeParams(params) }
        );
        return response.data as SearchRecipesResponse;
    },

    Create: async (data: CreateRecipeRequest): Promise<SystemRecipeDetail> => {
        const response = await executeAPI(API.recipe.system.create(), data);
        return response.data as SystemRecipeDetail;
    },

    Update: async (id: number, data: CreateRecipeRequest): Promise<SystemRecipeDetail> => {
        const response = await executeAPI(API.recipe.system.update(id), data);
        return response.data as SystemRecipeDetail;
    },

    Delete: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.system.delete(id));
    },
};

const userRecipeService = {
    List: async (params?: RecipeListParams): Promise<ListUserRecipesResponse> => {
        const response = await executeAPI(
            API.recipe.user.list(),
            undefined,
            { params: sanitizeParams(params) }
        );
        return response.data as ListUserRecipesResponse;
    },

    Create: async (data: CreateRecipeRequest): Promise<UserRecipeDetail> => {
        const response = await executeAPI(API.recipe.user.create(), data);
        return response.data as UserRecipeDetail;
    },

    Retrieve: async (id: number): Promise<UserRecipeDetail> => {
        const response = await executeAPI(API.recipe.user.retrieve(id));
        return response.data as UserRecipeDetail;
    },

    Update: async (id: number, data: CreateRecipeRequest): Promise<UserRecipeDetail> => {
        const response = await executeAPI(API.recipe.user.update(id), data);
        return response.data as UserRecipeDetail;
    },

    Delete: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.user.delete(id));
    },

    GetFavorites: async (): Promise<GetFavoritesResponse> => {
        const response = await executeAPI(API.recipe.user.getFavorites());
        return response.data as GetFavoritesResponse;
    },

    ToggleFavorite: async (recipeId: number): Promise<string | undefined> => {
        const response = await executeAPI(API.recipe.user.toggleFavorite(recipeId));
        return (response.data as { message?: string } | undefined)?.message;
    },

    Search: async (params?: RecipeSearchParams): Promise<SearchRecipesResponse> => {
        const mergedParams: RecipeSearchParams = {
            ...params,
            include_user: params?.include_user ?? true,
            include_system: params?.include_system ?? false,
        };

        const response = await executeAPI(
            API.recipe.search(),
            undefined,
            { params: sanitizeParams(mergedParams) }
        );
        return response.data as SearchRecipesResponse;
    },
};

const foodLogService = {
    Log: async (data: CreateFoodLogRequest): Promise<FoodLogEntryWithRecipe> => {
        const response = await executeAPI(API.recipe.logs.log(), data);
        return response.data as FoodLogEntryWithRecipe;
    },

    LogRecipe: async (data: LogRecipeRequest): Promise<FoodLogEntryWithRecipe> => {
        const response = await executeAPI(API.recipe.logs.logRecipe(), data);
        return response.data as FoodLogEntryWithRecipe;
    },

    GetLogsByDate: async (date: string): Promise<GetLogsByDateResponse> => {
        const response = await executeAPI(API.recipe.logs.getByDate(date));
        return response.data as GetLogsByDateResponse;
    },

    GetLogsInRange: async (startDate: string, endDate: string): Promise<GetLogsByDateRangeResponse> => {
        const response = await executeAPI(API.recipe.logs.getInRange(startDate, endDate));
        return response.data as GetLogsByDateRangeResponse;
    },

    GetFoodLogEntry: async (id: number): Promise<FoodLogEntryWithRecipe> => {
        const response = await executeAPI(API.recipe.logs.getFoodLogEntry(id));
        return response.data as FoodLogEntryWithRecipe;
    },

    UpdateFoodLogEntry: async (id: number, data: CreateFoodLogRequest): Promise<FoodLogEntryWithRecipe> => {
        const response = await executeAPI(API.recipe.logs.updateFoodLogEntry(id), data);
        return response.data as FoodLogEntryWithRecipe;
    },

    DeleteFoodLogEntry: async (id: number): Promise<void> => {
        await executeAPI(API.recipe.logs.deleteFoodLogEntry(id));
    },
};

const nutritionService = {
    GetDailySummary: async (date: string): Promise<DailyNutritionSummary> => {
        const response = await executeAPI(API.recipe.nutrition.getDailySummary(date));
        return response.data as DailyNutritionSummary;
    },

    GetWeeklySummary: async (startDate: string): Promise<GetWeeklyNutritionResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getWeeklySummary(startDate));
        return response.data as GetWeeklyNutritionResponse;
    },

    GetMonthlySummary: async (year: number, month: number): Promise<GetMonthlyNutritionResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getMonthlySummary(year, month));
        return response.data as GetMonthlyNutritionResponse;
    },

    GetGoals: async (): Promise<NutritionGoals> => {
        const response = await executeAPI(API.recipe.nutrition.getGoals());
        return response.data as NutritionGoals;
    },

    UpdateGoals: async (data: UpsertNutritionGoalsRequest): Promise<NutritionGoals> => {
        const response = await executeAPI(API.recipe.nutrition.updateGoals(), data);
        return response.data as NutritionGoals;
    },

    CompareWithGoals: async (date: string): Promise<GetNutritionComparisonResponse> => {
        const response = await executeAPI(API.recipe.nutrition.compareWithGoals(date));
        return response.data as GetNutritionComparisonResponse;
    },

    GetNutritionInsights: async (date: string): Promise<GetNutritionInsightsResponse> => {
        const response = await executeAPI(API.recipe.nutrition.getNutritionInsights(date));
        return response.data as GetNutritionInsightsResponse;
    },
};

export default { systemRecipeService, userRecipeService, foodLogService, nutritionService };