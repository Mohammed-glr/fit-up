import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

// Type definitions for user stats
export interface CoachInfo {
  coach_id: string;
  name: string;
  image: string | null;
  specialty?: string;
  assigned_at: string;
  total_messages: number;
  last_message_at?: string;
}

export interface UserStats {
  user_id: string;
  total_workouts: number;
  active_programs: number;
  days_active: number;
  current_streak: number;
  longest_streak: number;
  total_weeks: number;
  completion_rate: number;
  last_workout_date?: string;
  first_workout_date?: string;
  assigned_coach?: CoachInfo;
}

/**
 * Hook to fetch comprehensive user statistics
 * Returns total workouts, active programs, streaks, completion rate, and coach info
 */
export const useUserStats = () => {
  return useQuery<UserStats>({
    queryKey: ['user', 'stats'],
    queryFn: async () => {
      const response = await httpClient.get<UserStats>('/auth/stats');
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
    gcTime: 10 * 60 * 1000, // Keep in cache for 10 minutes
    retry: 2,
  });
};
