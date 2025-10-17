
type AttachmentType = 'image' | 'document' | 'workout_plan';
type WebSocketMessageType = 'new_message' | 'message_edited' | 'message_read' | 'message_deleted' | 'error';


interface Conversation {
  conversation_id: number;
  coach_id: string;
  client_id: string;
  created_at: string;
  updated_at: string;
  last_message_at: string;
  is_archived: boolean;
}

interface ConversationWithDetails extends Conversation {
  coach_name: string;
  coach_image?: string | null;
  client_name: string;
  client_image?: string | null;
  last_message?: MessageWithDetails;
  unread_count: number;
  total_messages: number;
}

interface ConversationOverview {
  conversation_id: number;
  coach_id: string;
  client_id: string;
  coach_name: string;
  coach_image?: string | null;
  client_name: string;
  client_image?: string | null;
  created_at: string;
  last_message_at: string;
  is_archived: boolean;
  last_message_text?: string | null;
  last_message_sender_id?: string | null;
  last_message_sent_at?: string | null;
  total_messages: number;
}


interface Message {
  message_id: number;
  conversation_id: number;
  sender_id: string;
  message_text: string;
  sent_at: string;
  edited_at?: string | null;
  is_deleted: boolean;
  deleted_at?: string | null;
  reply_to_message_id?: number | null;
}

interface MessageWithDetails extends Message {
  sender_name: string;
  sender_image?: string | null;
  attachments?: MessageAttachment[];
  is_read: boolean;
  reply_to_message?: MessageWithDetails;
}

interface MessageAttachment {
  attachment_id: number;
  message_id: number;
  attachment_type: AttachmentType;
  file_name: string;
  file_url: string;
  file_size?: number;
  mime_type?: string;
  uploaded_at: string;
  metadata?: any;
}


interface CreateConversationRequest {
  coach_id: string;
  client_id: string;
}

interface SendMessageRequest {
  conversation_id: number;
  message_text: string;
  reply_to_message_id?: number;
}

interface UpdateMessageRequest {
  message_text: string;
}

interface MarkAsReadRequest {
  message_ids: number[];
}

interface ListConversationsParams {
  is_archived?: boolean;
  limit?: number;
  offset?: number;
}

interface GetMessagesParams {
  limit?: number;
  offset?: number;
}


interface CreateConversationResponse {
  conversation: ConversationWithDetails;
  message?: string;
}

interface GetConversationResponse {
  conversation: ConversationWithDetails;
}

interface ListConversationsResponse {
  conversations: ConversationOverview[];
  total: number;
}

interface GetMessagesResponse {
  messages: MessageWithDetails[];
  total: number;
  has_more: boolean;
}

interface SendMessageResponse {
  message: MessageWithDetails;
}

interface UpdateMessageResponse {
  message: MessageWithDetails;
}

interface UnreadCountResponse {
  unread_count: number;
}

interface WebSocketMessage {
  type: WebSocketMessageType;
  conversation_id: number;
  message?: MessageWithDetails;
  message_id?: number;
  read_by?: string;
  error?: string;
  timestamp: string;
}

export type {
  AttachmentType,
  WebSocketMessageType,
  
  Conversation,
  ConversationWithDetails,
  ConversationOverview,
  
  Message,
  MessageWithDetails,
  MessageAttachment,
  
  CreateConversationRequest,
  SendMessageRequest,
  UpdateMessageRequest,
  MarkAsReadRequest,
  ListConversationsParams,
  GetMessagesParams,
  
  CreateConversationResponse,
  GetConversationResponse,
  ListConversationsResponse,
  GetMessagesResponse,
  SendMessageResponse,
  UpdateMessageResponse,
  UnreadCountResponse,
  
  WebSocketMessage,
};
