import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface WorkoutHistoryItem {
  date: string;
  plan_id?: number;
  plan_name?: string;
  day_title?: string;
  total_exercises: number;
  completed_sets: number;
  total_volume: number;
  duration_minutes: number;
  exercises: string[];
}

export interface WorkoutHistoryResponse {
  workouts: WorkoutHistoryItem[];
  total_count: number;
  page: number;
  page_size: number;
  has_more: boolean;
}

interface UseWorkoutHistoryOptions {
  page?: number;
  pageSize?: number;
  startDate?: string; // Format: YYYY-MM-DD
  endDate?: string; // Format: YYYY-MM-DD
}

export const useWorkoutHistory = (options: UseWorkoutHistoryOptions = {}) => {
  const { page = 1, pageSize = 20, startDate, endDate } = options;

  return useQuery<WorkoutHistoryResponse>({
    queryKey: ['workout-history', page, pageSize, startDate, endDate],
    queryFn: async () => {
      const params = new URLSearchParams();
      params.append('page', page.toString());
      params.append('page_size', pageSize.toString());
      if (startDate) params.append('start_date', startDate);
      if (endDate) params.append('end_date', endDate);

      const response = await httpClient.get<WorkoutHistoryResponse>(
        `/auth/workout-history?${params.toString()}`
      );
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 15 * 60 * 1000, // 15 minutes
    retry: 2,
  });
};
