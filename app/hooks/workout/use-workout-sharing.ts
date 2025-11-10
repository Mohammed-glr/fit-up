import { useQuery, useMutation } from '@tanstack/react-query';
import { httpClient } from '@/api/client';
import type {
  WorkoutShareSummary,
  ShareWorkoutRequest,
  ShareWorkoutResponse,
} from '@/types/workout-sharing';
import { Alert } from 'react-native';

/**
 * Hook to fetch workout share summary
 * @param sessionId - The workout session ID
 * @param enabled - Whether to enable the query
 */
export const useWorkoutShareSummary = (sessionId?: number, enabled = true) => {
  return useQuery<WorkoutShareSummary>({
    queryKey: ['workout-share-summary', sessionId],
    queryFn: async () => {
      if (!sessionId) throw new Error('Session ID is required');
      const response = await httpClient.get<WorkoutShareSummary>(
        `/workout-sessions/${sessionId}/share-summary`
      );
      return response.data;
    },
    enabled: enabled && !!sessionId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });
};


export const useShareWorkout = () => {
  return useMutation<ShareWorkoutResponse, Error, ShareWorkoutRequest>({
    mutationFn: async (data: ShareWorkoutRequest) => {
      const response = await httpClient.post<ShareWorkoutResponse>(
        '/workout-sessions/share',
        data
      );
      return response.data;
    },
    onSuccess: (data) => {
      if (data.message) {
        Alert.alert('Success', data.message);
      }
    },
    onError: (error) => {
      Alert.alert('Error', error.message || 'Failed to share workout');
    },
  });
};
