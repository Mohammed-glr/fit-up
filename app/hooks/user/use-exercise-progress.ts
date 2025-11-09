import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface ExerciseProgressDataPoint {
  date: string;
  weight: number;
  reps: number;
  sets: number;
  volume: number;
  is_personal_record: boolean;
}

export interface ExerciseProgressData {
  exercise_id?: number;
  exercise_name: string;
  data_points: ExerciseProgressDataPoint[];
  max_weight: number;
  max_volume: number;
  total_sets: number;
}

interface UseExerciseProgressOptions {
  exerciseName: string;
  startDate?: string; // Format: YYYY-MM-DD
  endDate?: string; // Format: YYYY-MM-DD
  enabled?: boolean;
}

export const useExerciseProgress = (options: UseExerciseProgressOptions) => {
  const { exerciseName, startDate, endDate, enabled = true } = options;

  return useQuery<ExerciseProgressData>({
    queryKey: ['exercise-progress', exerciseName, startDate, endDate],
    queryFn: async () => {
      const params = new URLSearchParams();
      params.append('exercise', exerciseName);
      if (startDate) params.append('start_date', startDate);
      if (endDate) params.append('end_date', endDate);

      const response = await httpClient.get<ExerciseProgressData>(
        `/auth/exercise-progress?${params.toString()}`
      );
      return response.data;
    },
    enabled: enabled && !!exerciseName,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 15 * 60 * 1000, // 15 minutes
    retry: 2,
  });
};
