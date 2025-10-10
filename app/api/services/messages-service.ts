import { API } from '../endpoints';
import { executeAPI } from '../client';
import {
    AttachmentType,    
    Conversation,
    ConversationWithDetails,
    ConversationOverview,
    
    Message,
    MessageWithDetails,
    MessageReadStatus,
    MessageAttachment,
    
    CreateConversationRequest,
    SendMessageRequest,
    UpdateMessageRequest,
    UploadAttachmentRequest,
    MarkAsReadRequest,
    
    PaginationParams,
    MessageFilters,
    ConversationFilters,
    
    PaginatedResponse,
    MessageResponse,
    MessagesResponse,
    ConversationResponse,
    ConversationsResponse,
    CreateConversationResponse,
    GetConversationResponse,
    ListConversationsResponse,
    GetMessagesResponse,
    SendMessageResponse,
    UpdateMessageResponse,
    DeleteMessageResponse,
    UnreadCountResponse,

} from '@/types/message'

const messagesService = {
    
}

export default messagesService;


