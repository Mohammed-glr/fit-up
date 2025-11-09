import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface Exercise {
  exercise_id: number;
  name: string;
  category?: string;
  muscle_group?: string;
  equipment?: string;
  difficulty_level?: string;
  instructions?: string[];
  tips?: string[];
  image_url?: string;
}

interface UseMostUsedExercisesOptions {
  limit?: number;
  enabled?: boolean;
}

export const useMostUsedExercises = (options: UseMostUsedExercisesOptions = {}) => {
  const { limit = 20, enabled = true } = options;

  return useQuery<Exercise[]>({
    queryKey: ['most-used-exercises', limit],
    queryFn: async () => {
      const params = new URLSearchParams();
      params.append('limit', limit.toString());

      const response = await httpClient.get<Exercise[]>(
        `/schema/exercises/most-used?${params.toString()}`
      );
      return response.data;
    },
    enabled,
    staleTime: 10 * 60 * 1000, // 10 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
    retry: 2,
  });
};
