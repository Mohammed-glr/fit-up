import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import recipeAPI from '@/api/services/recipe-service';
import type {
  CreateFoodLogRequest,
  FoodLogEntryWithRecipe,
  GetLogsByDateRangeResponse,
  GetLogsByDateResponse,
  LogRecipeRequest,
} from '@/types/food-tracker';
import { APIError } from '@/api/client';
import { foodLogKeys, nutritionKeys } from './keys';

const { foodLogService } = recipeAPI;

const invalidateLogDependencies = (queryClient: ReturnType<typeof useQueryClient>, logDate?: string) => {
  queryClient.invalidateQueries({ queryKey: foodLogKeys.all });
  if (logDate) {
    queryClient.invalidateQueries({ queryKey: foodLogKeys.byDate(logDate) });
    queryClient.invalidateQueries({ queryKey: nutritionKeys.daily(logDate) });
    queryClient.invalidateQueries({ queryKey: nutritionKeys.comparison(logDate) });
    queryClient.invalidateQueries({ queryKey: nutritionKeys.insights(logDate) });
  }
  queryClient.invalidateQueries({ queryKey: nutritionKeys.all });
};

export const useFoodLogsByDate = (date?: string | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof date === 'string' && date.length > 0;

  return useQuery<GetLogsByDateResponse, APIError>({
    queryKey: foodLogKeys.byDate(date ?? ''),
    queryFn: () => foodLogService.GetLogsByDate(date as string),
    enabled,
  });
};

export const useFoodLogsInRange = (
  startDate?: string | null,
  endDate?: string | null,
  options?: { enabled?: boolean },
) => {
  const enabled =
    (options?.enabled ?? true) &&
    typeof startDate === 'string' &&
    startDate.length > 0 &&
    typeof endDate === 'string' &&
    endDate.length > 0;

  return useQuery<GetLogsByDateRangeResponse, APIError>({
    queryKey: foodLogKeys.range(startDate ?? '', endDate ?? ''),
    queryFn: () => foodLogService.GetLogsInRange(startDate as string, endDate as string),
    enabled,
  });
};

export const useFoodLogEntry = (logId?: number | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof logId === 'number' && logId > 0;

  return useQuery<FoodLogEntryWithRecipe, APIError>({
    queryKey: foodLogKeys.entry(logId ?? null),
    queryFn: () => foodLogService.GetFoodLogEntry(logId as number),
    enabled,
  });
};

export const useLogFood = () => {
  const queryClient = useQueryClient();

  return useMutation<FoodLogEntryWithRecipe, APIError, CreateFoodLogRequest>({
    mutationFn: (payload) => foodLogService.Log(payload),
    onSuccess: (data) => {
      invalidateLogDependencies(queryClient, data.log_date);
    },
  });
};

export const useLogRecipe = () => {
  const queryClient = useQueryClient();

  return useMutation<FoodLogEntryWithRecipe, APIError, LogRecipeRequest>({
    mutationFn: (payload) => foodLogService.LogRecipe(payload),
    onSuccess: (data) => {
      invalidateLogDependencies(queryClient, data.log_date);
    },
  });
};

export const useUpdateFoodLogEntry = () => {
  const queryClient = useQueryClient();

  return useMutation<FoodLogEntryWithRecipe, APIError, { logId: number; payload: CreateFoodLogRequest }>({
    mutationFn: ({ logId, payload }) => foodLogService.UpdateFoodLogEntry(logId, payload),
    onSuccess: (data) => {
      invalidateLogDependencies(queryClient, data.log_date);
      queryClient.invalidateQueries({ queryKey: foodLogKeys.entry(data.id) });
    },
  });
};

export const useDeleteFoodLogEntry = () => {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, { logId: number; logDate?: string }>({
    mutationFn: ({ logId }) => foodLogService.DeleteFoodLogEntry(logId),
    onSuccess: (_, variables) => {
      invalidateLogDependencies(queryClient, variables.logDate);
      queryClient.removeQueries({ queryKey: foodLogKeys.entry(variables.logId) });
    },
  });
};
