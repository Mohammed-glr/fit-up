import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';
import { UserWorkoutTemplate, UserTemplateListResponse, UserTemplateListParams } from '@/types/workout-template';

/**
 * Hook to fetch user's workout templates with pagination
 * @param params - Pagination parameters (page, page_size)
 */
export const useUserTemplates = (params?: UserTemplateListParams) => {
  const { page = 1, page_size = 20 } = params || {};

  return useQuery<UserTemplateListResponse>({
    queryKey: ['templates', 'user', page, page_size],
    queryFn: async () => {
      const response = await httpClient.get<UserTemplateListResponse>('/templates', {
        params: { page, page_size },
      });
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // Consider data fresh for 2 minutes
    gcTime: 5 * 60 * 1000, // Keep in cache for 5 minutes
  });
};

/**
 * Hook to fetch public workout templates with pagination
 * These are community templates that users can browse and clone
 * @param params - Pagination parameters (page, page_size)
 */
export const usePublicTemplates = (params?: UserTemplateListParams) => {
  const { page = 1, page_size = 20 } = params || {};

  return useQuery<UserTemplateListResponse>({
    queryKey: ['templates', 'public', page, page_size],
    queryFn: async () => {
      const response = await httpClient.get<UserTemplateListResponse>('/templates/public', {
        params: { page, page_size },
      });
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // Public templates change less frequently
    gcTime: 10 * 60 * 1000,
  });
};

/**
 * Hook to fetch a single workout template by ID
 * @param templateId - The ID of the template to fetch
 * @param enabled - Whether to enable the query (default: true if templateId exists)
 */
export const useTemplate = (templateId?: string, enabled = true) => {
  return useQuery<UserWorkoutTemplate>({
    queryKey: ['templates', templateId],
    queryFn: async () => {
      if (!templateId) throw new Error('Template ID is required');
      const response = await httpClient.get<UserWorkoutTemplate>(`/templates/${templateId}`);
      return response.data;
    },
    enabled: enabled && !!templateId,
    staleTime: 2 * 60 * 1000,
    gcTime: 5 * 60 * 1000,
  });
};
