import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import recipeAPI from '@/api/services/recipe-service';
import type {
  DailyNutritionSummary,
  GetMonthlyNutritionResponse,
  GetNutritionComparisonResponse,
  GetNutritionInsightsResponse,
  GetWeeklyNutritionResponse,
  NutritionGoals,
  UpsertNutritionGoalsRequest,
} from '@/types/food-tracker';
import { APIError } from '@/api/client';
import { nutritionKeys } from './keys';

const { nutritionService } = recipeAPI;

export const useDailyNutritionSummary = (date?: string | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof date === 'string' && date.length > 0;

  return useQuery<DailyNutritionSummary, APIError>({
    queryKey: nutritionKeys.daily(date ?? ''),
    queryFn: () => nutritionService.GetDailySummary(date as string),
    enabled,
  });
};

export const useWeeklyNutritionSummary = (startDate?: string | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof startDate === 'string' && startDate.length > 0;

  return useQuery<GetWeeklyNutritionResponse, APIError>({
    queryKey: nutritionKeys.weekly(startDate ?? ''),
    queryFn: () => nutritionService.GetWeeklySummary(startDate as string),
    enabled,
  });
};

export const useMonthlyNutritionSummary = (
  year?: number | null,
  month?: number | null,
  options?: { enabled?: boolean },
) => {
  const enabled =
    (options?.enabled ?? true) &&
    typeof year === 'number' &&
    typeof month === 'number';

  return useQuery<GetMonthlyNutritionResponse, APIError>({
    queryKey: nutritionKeys.monthly(year ?? 0, month ?? 0),
    queryFn: () => nutritionService.GetMonthlySummary(year as number, month as number),
    enabled,
  });
};

export const useNutritionGoals = (options?: { enabled?: boolean }) => {
  return useQuery<NutritionGoals, APIError>({
    queryKey: nutritionKeys.goals,
    queryFn: () => nutritionService.GetGoals(),
    enabled: options?.enabled ?? true,
  });
};

export const useUpdateNutritionGoals = () => {
  const queryClient = useQueryClient();

  return useMutation<NutritionGoals, APIError, UpsertNutritionGoalsRequest>({
    mutationFn: (payload) => nutritionService.UpdateGoals(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: nutritionKeys.goals });
      queryClient.invalidateQueries({ queryKey: nutritionKeys.all });
    },
  });
};

export const useNutritionComparison = (date?: string | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof date === 'string' && date.length > 0;

  return useQuery<GetNutritionComparisonResponse, APIError>({
    queryKey: nutritionKeys.comparison(date ?? ''),
    queryFn: () => nutritionService.CompareWithGoals(date as string),
    enabled,
  });
};

export const useNutritionInsights = (date?: string | null, options?: { enabled?: boolean }) => {
  const enabled = (options?.enabled ?? true) && typeof date === 'string' && date.length > 0;

  return useQuery<GetNutritionInsightsResponse, APIError>({
    queryKey: nutritionKeys.insights(date ?? ''),
    queryFn: () => nutritionService.GetNutritionInsights(date as string),
    enabled,
  });
};
