import { API } from '../endpoints';
import { executeAPI } from '../client';
import { 
    CreateConversationRequest,
    CreateConversationResponse,
    GetConversationResponse,
    ListConversationsParams,
    ListConversationsResponse,
    GetMessagesParams,
    GetMessagesResponse,
    SendMessageRequest,
    SendMessageResponse,
    UpdateMessageRequest,
    UpdateMessageResponse,
    UnreadCountResponse,
 } from '@/types/message'

const conversationService = {
    Create: async (data: CreateConversationRequest): Promise<CreateConversationResponse> => {
        const response = await executeAPI(API.message.conversations.create(), data);
        return response.data as CreateConversationResponse;
    },

    List: async (params?: ListConversationsParams): Promise<ListConversationsResponse> => {
        const response = await executeAPI(API.message.conversations.list(), undefined, { params });
        return response.data as ListConversationsResponse;
    },

    Get: async (conversation_id: number): Promise<GetConversationResponse> => {
        const response = await executeAPI(API.message.conversations.get(conversation_id));
        return response.data as GetConversationResponse;
    },

    GetUnreadCount: async (conversation_id: number): Promise<UnreadCountResponse> => {
        const response = await executeAPI(API.message.conversations.getUnreadCount(conversation_id));
        return response.data as UnreadCountResponse;
    },

    GetMessages: async (conversation_id: number, params?: GetMessagesParams): Promise<GetMessagesResponse> => {
        const response = await executeAPI(API.message.conversations.getMessages(conversation_id), undefined, { params });
        return response.data as GetMessagesResponse;
    },

    MarkAllAsRead: async (conversation_id: number): Promise<void> => {
        await executeAPI(API.message.conversations.markAllAsRead(conversation_id));
    },
}

const messageService = {
    Send: async (data: SendMessageRequest): Promise<SendMessageResponse> => {
        const response = await executeAPI(API.message.messages.send(), data);
        return response.data as SendMessageResponse;
    },

    Update: async (message_id: number, data: UpdateMessageRequest): Promise<UpdateMessageResponse> => {
        const response = await executeAPI(API.message.messages.update(message_id), data);
        return response.data as UpdateMessageResponse;
    },

    Delete: async (message_id: number): Promise<void> => {
        await executeAPI(API.message.messages.delete(message_id));
    },  

    MarkAsRead: async (message_id: number): Promise<void> => {
        await executeAPI(API.message.messages.markAsRead(message_id));
    }
}

export { conversationService, messageService };

