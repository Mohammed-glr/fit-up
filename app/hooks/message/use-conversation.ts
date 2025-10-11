import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import {
    conversationService,
    messageService,
} from '@/api/services/messages-service';
import { APIError } from '@/api/client';
import { Toast } from '@/components/ui';

export const conversationKeys = {
    all: ['conversations'] as const,
    lists: () => [...conversationKeys.all, 'list'] as const,
    list: (filter: string) => [...conversationKeys.lists(), { filter }] as const,
    details: () => [...conversationKeys.all, 'detail'] as const,
    detail: (id: number) => [...conversationKeys.details(), id] as const,
    messages: (conversation_id: number) => [...conversationKeys.detail(conversation_id), 'messages'] as const,
    unreadCount: (conversation_id: number) => [...conversationKeys.detail(conversation_id), 'unreadCount'] as const,
}

export const useConversations = (params?: { is_archived?: boolean }) => {
    const filter = params ? JSON.stringify(params) : '';
    return useQuery({
        queryKey: conversationKeys.list(filter),
        queryFn: () => conversationService.List(),
        staleTime: 5 * 60 * 1000, 
    })
}

export const useConversation = (conversation_id: number) => {
    return useQuery({
        queryKey: conversationKeys.detail(conversation_id),
        queryFn: () => conversationService.Get(conversation_id),
        staleTime: 5 * 60 * 1000, 
    })
}

export const useConversationMessages = (conversation_id: number) => {
    return useInfiniteQuery({
        queryKey: conversationKeys.messages(conversation_id),
        queryFn: ({ pageParam = 0 }) =>
            conversationService.GetMessages(conversation_id),
        initialPageParam: 0,
        getNextPageParam: (lastPage: any, allPages) => {
            if (lastPage && lastPage.has_more) {
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
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: conversationKeys.lists()});
            Toast({
                message: 'Conversation created successfully',
                type: 'success',
                isVisible: true,
            });
        },
        onError: (error: APIError) => {
            Toast({
                message: error.message || 'Failed to create conversation',
                type: 'error',
                isVisible: true,
            });
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
        },
        onError: (error: APIError) => {
            Toast({
                message: error.message || 'Failed to send message',
                type: 'error',
                isVisible: true,
            });
        },
    })
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
        onError: (error: APIError) => {
            Toast({
                message: error.message || 'Failed to update message',
                type: 'error',
                isVisible: true,
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
        },
        onError: (error: APIError) => {
            Toast({
                message: error.message || 'Failed to delete message',
                type: 'error',
                isVisible: true,
            });
        }
    })
}

export const useMarkAsRead = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (message_id: number) => messageService.MarkAsRead(message_id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: conversationKeys.lists() });
        },
        onError: (error: APIError) => {
            Toast({
                message: error.message || 'Failed to mark message as read',
                type: 'error',
                isVisible: true,
            });
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

