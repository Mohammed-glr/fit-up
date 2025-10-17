import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import {
    conversationService,
    messageService,
} from '@/api/services/messages-service';
import { authService } from '@/api/services/auth-service';
import { APIError } from '@/api/client';
import type { GetMessagesResponse, ListConversationsParams, ListConversationsResponse } from '@/types/message';
import type { PublicUserResponse } from '@/types/auth';

const DEFAULT_PAGE_SIZE = 20;

export const conversationKeys = {
    all: ['conversations'] as const,
    lists: () => [...conversationKeys.all, 'list'] as const,
    list: (filter: string) => [...conversationKeys.lists(), { filter }] as const,
    details: () => [...conversationKeys.all, 'detail'] as const,
    detail: (id: number) => [...conversationKeys.details(), id] as const,
    messages: (conversation_id: number) => [...conversationKeys.detail(conversation_id), 'messages'] as const,
    unreadCount: (conversation_id: number) => [...conversationKeys.detail(conversation_id), 'unreadCount'] as const,
}

export const useConversations = (params?: ListConversationsParams, pageSize: number = DEFAULT_PAGE_SIZE) => {
    const filterParams: Partial<ListConversationsParams> = params ? { ...params } : {};
    delete filterParams.limit;
    delete filterParams.offset;

    const filter = JSON.stringify({ ...filterParams, pageSize });

    return useInfiniteQuery<ListConversationsResponse, APIError>({
        queryKey: conversationKeys.list(filter),
        queryFn: ({ pageParam }) => {
            const offset = typeof pageParam === 'number' ? pageParam : 0;
            return conversationService.List({
                ...filterParams,
                limit: pageSize,
                offset,
            });
        },
        initialPageParam: 0,
        getNextPageParam: (lastPage, allPages) => {
            if (lastPage?.has_more) {
                return allPages.length * pageSize;
            }
            return undefined;
        },
        staleTime: 5 * 60 * 1000,
    });
}

export const useConversation = (conversation_id: number) => {
    return useQuery({
        queryKey: conversationKeys.detail(conversation_id),
        queryFn: () => conversationService.Get(conversation_id),
        staleTime: 5 * 60 * 1000, 
    })
}

export const useConversationMessages = (conversation_id: number, pageSize: number = DEFAULT_PAGE_SIZE) => {
    return useInfiniteQuery<GetMessagesResponse, APIError>({
        queryKey: conversationKeys.messages(conversation_id),
        queryFn: ({ pageParam }) => {
            const page = typeof pageParam === 'number' ? pageParam : 0;
            const offset = page * pageSize;
            return conversationService.GetMessages(conversation_id, {
                limit: pageSize,
                offset,
            });
        },
        initialPageParam: 0,
        getNextPageParam: (lastPage, allPages) => {
            if (lastPage?.has_more) {
                return allPages.length;
            }
            return undefined;
        },
        enabled: !!conversation_id,
    });
}

export const useUnreadCount = (conversation_id: number) => {
    return useQuery({
        queryKey: conversationKeys.unreadCount(conversation_id),
        queryFn: () => conversationService.GetUnreadCount(conversation_id),
        staleTime: 1 * 60 * 1000, 
        enabled: !!conversation_id,
    });
}


export const useCreateConversation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: conversationService.Create,
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: conversationKeys.lists()});
        }
    })
}

export const useSendMessage = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: messageService.Send,
        onSuccess: (data) => {
            queryClient.invalidateQueries({
                queryKey: conversationKeys.messages(data.message.conversation_id),
            });
            queryClient.invalidateQueries({
                queryKey: conversationKeys.detail(data.message.conversation_id),
            });
            queryClient.invalidateQueries({
                queryKey: conversationKeys.lists(),
            })
        }
    })
}

export const useUserLookup = () => {
    return useMutation<PublicUserResponse, APIError, string>({
        mutationFn: (username: string) => authService.GetPublicProfile(username),
    });
}

export const useUpdateMessage = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ message_id, data }: { message_id: number; data: any }) =>
            messageService.Update(message_id, data),
        onSuccess: (data) => {
            queryClient.invalidateQueries({
                queryKey: conversationKeys.messages(data.message.conversation_id),
            });
        },
    })
}

export const useDeleteMessage = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (message_id: number) => messageService.Delete(message_id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: conversationKeys.lists() });
        }
    })
}

export const useMarkAsRead = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (message_id: number) => messageService.MarkAsRead(message_id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: conversationKeys.lists() });
        }
    })
}

export const useMarkAllAsRead = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: conversationService.MarkAllAsRead,
    onSuccess: (_, conversationId) => {
      queryClient.invalidateQueries({
        queryKey: conversationKeys.unreadCount(conversationId),
      });
      queryClient.invalidateQueries({
        queryKey: conversationKeys.messages(conversationId),
      });
    },
  });
};

