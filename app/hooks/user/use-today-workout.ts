import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

// Type definitions for today's workout
export interface TodayExercise {
  exercise_id?: number;
  name: string;
  sets: number;
  reps: string;
  rest_seconds: number;
  notes?: string;
}

export interface TodayWorkout {
  plan_id: number;
  plan_name: string;
  day_index: number;
  day_title: string;
  focus: string;
  is_rest: boolean;
  total_exercises: number;
  estimated_minutes: number;
  is_completed: boolean;
  completed_at?: string;
  exercises: TodayExercise[];
}

/**
 * Hook to fetch today's workout from the active plan
 * Returns workout details including exercises, or null if no active plan
 */
export const useTodayWorkout = () => {
  return useQuery<TodayWorkout | null>({
    queryKey: ['workout', 'today'],
    queryFn: async () => {
      const response = await httpClient.get<TodayWorkout>('/auth/today-workout');
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
    gcTime: 15 * 60 * 1000, // Keep in cache for 15 minutes
    retry: 2,
  });
};
