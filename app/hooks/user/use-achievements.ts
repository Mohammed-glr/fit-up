import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface Achievement {
  achievement_id: number;
  name: string;
  description: string;
  badge_icon: string;
  badge_color: string;
  category: 'streak' | 'volume' | 'pr' | 'milestone' | 'consistency';
  requirement_type: string;
  requirement_value: number;
  points: number;
  progress: number;
  earned_at?: string;
  is_completed: boolean;
  completion_rate: number;
}

export interface AchievementStats {
  total_achievements: number;
  earned_achievements: number;
  total_points: number;
  earned_points: number;
  completion_rate: number;
}

export const useAchievements = () => {
  return useQuery<Achievement[]>({
    queryKey: ['achievements'],
    queryFn: async () => {
      const response = await httpClient.get<Achievement[]>('/auth/achievements');
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    retry: 2,
  });
};

export const useAchievementStats = () => {
  return useQuery<AchievementStats>({
    queryKey: ['achievement-stats'],
    queryFn: async () => {
      const response = await httpClient.get<AchievementStats>('/auth/achievement-stats');
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    retry: 2,
  });
};
