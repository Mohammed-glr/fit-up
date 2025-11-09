import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export type ActivityType =
  | 'workout_completed'
  | 'pr_achieved'
  | 'streak_milestone'
  | 'coach_message'
  | 'new_plan'
  | 'goal_achieved'
  | 'plan_completed';

export interface ActivityFeedItem {
  id: string;
  type: ActivityType;
  title: string;
  description: string;
  timestamp: string;
  icon: string;
  metadata?: Record<string, any>;
}

interface ActivityFeedResponse {
  activities: ActivityFeedItem[];
  count: number;
}

export const useActivityFeed = (limit: number = 10) => {
  return useQuery<ActivityFeedResponse>({
    queryKey: ['activity-feed', limit],
    queryFn: async () => {
      const response = await httpClient.get<ActivityFeedResponse>(
        `/auth/activity-feed?limit=${limit}`
      );
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    retry: 2,
  });
};
