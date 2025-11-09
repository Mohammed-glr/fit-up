import { useMutation, useQueryClient } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface ExerciseSetLog {
  exercise_id?: number;
  exercise_name: string;
  set_number: number;
  reps: number;
  weight: number;
  completed: boolean;
  notes?: string;
}

export interface WorkoutCompletionRequest {
  plan_id: number;
  plan_day_id?: number;
  day_index: number;
  duration_seconds: number;
  completed_at: string;
  exercises: ExerciseSetLog[];
  notes?: string;
}

export interface WorkoutCompletionResponse {
  success: boolean;
  message: string;
  workout_date: string;
  total_sets: number;
  completed_sets: number;
  completion_rate: number;
  total_volume: number;
  duration_minutes: number;
  new_streak: number;
  is_personal_best: boolean;
}

export const useWorkoutCompletion = () => {
  const queryClient = useQueryClient();

  return useMutation<WorkoutCompletionResponse, Error, WorkoutCompletionRequest>({
    mutationFn: async (data: WorkoutCompletionRequest) => {
      const response = await httpClient.post<WorkoutCompletionResponse>(
        '/auth/workout-complete',
        data
      );
      return response.data;
    },
    onSuccess: () => {
      // Invalidate relevant queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['user-stats'] });
      queryClient.invalidateQueries({ queryKey: ['today-workout'] });
      queryClient.invalidateQueries({ queryKey: ['activity-feed'] });
    },
  });
};
