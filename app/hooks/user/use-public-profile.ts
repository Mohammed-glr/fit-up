import { useQuery } from '@tanstack/react-query';

import { authService } from '@/api/services/auth-service';
import type { PublicUserResponse } from '@/types/auth';

const FIVE_MINUTES = 5 * 60 * 1000;

export const usePublicProfile = (username?: string) => {
  return useQuery<PublicUserResponse>({
    queryKey: ['public-profile', username],
    queryFn: () => authService.GetPublicProfile(username as string),
    enabled: !!username,
    staleTime: FIVE_MINUTES,
  });
};
