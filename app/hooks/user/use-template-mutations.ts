import { useMutation, useQueryClient } from '@tanstack/react-query';
import { httpClient } from '@/api/client';
import { CreateUserTemplateRequest, UpdateUserTemplateRequest, UserWorkoutTemplate } from '@/types/workout-template';

/**
 * Hook to create a new workout template
 * Invalidates the user templates list after successful creation
 */
export const useCreateTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation<UserWorkoutTemplate, Error, CreateUserTemplateRequest>({
    mutationFn: async (data: CreateUserTemplateRequest) => {
      const response = await httpClient.post<UserWorkoutTemplate>('/auth/templates', data);
      return response.data;
    },
    onSuccess: () => {
      // Invalidate user templates list to refetch
      queryClient.invalidateQueries({ queryKey: ['templates', 'user'] });
    },
  });
};

/**
 * Hook to update an existing workout template
 * Invalidates both the specific template and the user templates list
 */
export const useUpdateTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation<
    UserWorkoutTemplate,
    Error,
    { templateId: string; data: UpdateUserTemplateRequest }
  >({
    mutationFn: async ({ templateId, data }) => {
      const response = await httpClient.put<UserWorkoutTemplate>(
        `/auth/templates/${templateId}`,
        data
      );
      return response.data;
    },
    onSuccess: (_, variables) => {
      // Invalidate the specific template
      queryClient.invalidateQueries({ queryKey: ['templates', variables.templateId] });
      // Invalidate user templates list
      queryClient.invalidateQueries({ queryKey: ['templates', 'user'] });
      // If it's now public, invalidate public templates too
      queryClient.invalidateQueries({ queryKey: ['templates', 'public'] });
    },
  });
};

/**
 * Hook to delete a workout template
 * Invalidates the user templates list after successful deletion
 */
export const useDeleteTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (templateId: string) => {
      await httpClient.delete(`/auth/templates/${templateId}`);
    },
    onSuccess: (_, templateId) => {
      // Remove from cache
      queryClient.removeQueries({ queryKey: ['templates', templateId] });
      // Invalidate user templates list
      queryClient.invalidateQueries({ queryKey: ['templates', 'user'] });
      // Invalidate public templates in case it was public
      queryClient.invalidateQueries({ queryKey: ['templates', 'public'] });
    },
  });
};
