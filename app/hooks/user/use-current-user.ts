import {
    useQuery,
    useMutation,
    useQueryClient
} from '@tanstack/react-query';
import { useAuth } from '@/context/auth-context'; 
import { APIError } from '@/api/client';

export const useCurrentUser = () => {
    const { getCurrentUser } = useAuth();
    const queryClient = useQueryClient();

    const { data, error, isLoading, refetch } = useQuery({
        queryKey: ['currentUser'],
        queryFn: async () => {
            const user = await getCurrentUser();
            if (!user) throw new APIError('User not found', 404);
            return user;
        }
    });

    return { data, error, isLoading, refetch };
};
