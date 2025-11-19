import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

interface CoachInfo {
  assignment_id: number;
  coach_id: string;
  user_id: number;
  is_active: boolean;
  assigned_at: string;
  coach_name?: string;
  coach_email?: string;
}

export const useAssignedCoach = (workoutProfileId?: number) => {
  return useQuery({
    queryKey: ['assigned-coach', workoutProfileId],
    queryFn: async () => {
      if (!workoutProfileId) throw new Error('Workout profile ID required');
      const response = await httpClient.get(`coach/assigned/${workoutProfileId}`);
      return response.data as CoachInfo;
    },
    enabled: !!workoutProfileId,
  });
};
