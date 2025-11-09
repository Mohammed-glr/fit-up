import { executeAPI } from '../client';
import { API } from '../endpoints';
import type {
  CreateFoodLogRequest,
  LogRecipeRequest,
  FoodLogEntryWithRecipe,
  GetLogsByDateResponse,
  GetLogsByDateRangeResponse,
} from '@/types/food-tracker';

export const foodLogService = {
  async logFood(data: CreateFoodLogRequest): Promise<FoodLogEntryWithRecipe> {
    const response = await executeAPI<FoodLogEntryWithRecipe>(
      API.recipe.logs.log(),
      data
    );
    return response.data;
  },

  async logRecipe(data: LogRecipeRequest): Promise<FoodLogEntryWithRecipe> {
    const response = await executeAPI<FoodLogEntryWithRecipe>(
      API.recipe.logs.logRecipe(),
      data
    );
    return response.data;
  },

  async getLogsByDate(date: string): Promise<GetLogsByDateResponse> {
    const response = await executeAPI<GetLogsByDateResponse>(
      API.recipe.logs.getByDate(date)
    );
    return response.data;
  },

  async getLogsByDateRange(
    startDate: string,
    endDate: string
  ): Promise<GetLogsByDateRangeResponse> {
    const response = await executeAPI<GetLogsByDateRangeResponse>(
      API.recipe.logs.getInRange(startDate, endDate)
    );
    return response.data;
  },

  async getFoodLogEntry(id: number): Promise<FoodLogEntryWithRecipe> {
    const response = await executeAPI<FoodLogEntryWithRecipe>(
      API.recipe.logs.getFoodLogEntry(id)
    );
    return response.data;
  },

  async updateFoodLog(id: number, data: CreateFoodLogRequest): Promise<FoodLogEntryWithRecipe> {
    const response = await executeAPI<FoodLogEntryWithRecipe>(
      API.recipe.logs.updateFoodLogEntry(id),
      data
    );
    return response.data;
  },

  async deleteFoodLog(id: number): Promise<{ message: string }> {
    const response = await executeAPI<{ message: string }>(
      API.recipe.logs.deleteFoodLogEntry(id)
    );
    return response.data;
  },
};
