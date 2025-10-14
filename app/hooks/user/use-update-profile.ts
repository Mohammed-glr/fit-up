import { useMutation, useQueryClient } from '@tanstack/react-query';
import { authService } from '@/api/services/auth-service';

export const useUpdateProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: { name?: string; bio?: string; image?: string }) => {
            return await authService.UpdateProfile(data);
        },
        onSuccess: (response) => {
            queryClient.setQueryData(['currentUser'], response.user);
            queryClient.invalidateQueries({ queryKey: ['currentUser'] });
        },
    });
};
