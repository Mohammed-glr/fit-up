import type {
  RecipeListParams,
  RecipeSearchParams,
} from '@/types/food-tracker';

const baseKey = ['food-tracker'] as const;

const stringifyParams = (params?: Record<string, unknown>) =>
  JSON.stringify(params ?? {});

export const recipeKeys = {
  all: [...baseKey, 'recipes'] as const,
  list: (params?: RecipeListParams) =>
    [...baseKey, 'recipes', 'list', stringifyParams(params as Record<string, unknown> | undefined)] as const,
  detail: (id: number | null | undefined) =>
    [...baseKey, 'recipes', 'detail', id ?? 'unknown'] as const,
  favorites: [...baseKey, 'recipes', 'favorites'] as const,
  search: (params?: RecipeSearchParams) =>
    [...baseKey, 'recipes', 'search', stringifyParams(params as Record<string, unknown> | undefined)] as const,
};

export const foodLogKeys = {
  all: [...baseKey, 'logs'] as const,
  byDate: (date: string) => [...baseKey, 'logs', 'date', date] as const,
  range: (startDate: string, endDate: string) =>
    [...baseKey, 'logs', 'range', startDate, endDate] as const,
  entry: (id: number | null | undefined) =>
    [...baseKey, 'logs', 'entry', id ?? 'unknown'] as const,
};

export const nutritionKeys = {
  all: [...baseKey, 'nutrition'] as const,
  daily: (date: string) => [...baseKey, 'nutrition', 'daily', date] as const,
  weekly: (startDate: string) => [...baseKey, 'nutrition', 'weekly', startDate] as const,
  monthly: (year: number, month: number) => [...baseKey, 'nutrition', 'monthly', year, month] as const,
  goals: [...baseKey, 'nutrition', 'goals'] as const,
  comparison: (date: string) => [...baseKey, 'nutrition', 'comparison', date] as const,
  insights: (date: string) => [...baseKey, 'nutrition', 'insights', date] as const,
};
