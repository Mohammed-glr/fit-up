
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

interface MessageReadStatus {
  read_status_id: number;
  message_id: number;
  user_id: string;
  read_at: string;
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

interface UploadAttachmentRequest {
  message_id: number;
  attachment_type: AttachmentType;
  file_name: string;
  file_url: string;
  file_size?: number;
  mime_type?: string;
  metadata?: Record<string, any>;
}

interface MarkAsReadRequest {
  message_ids: number[];
}

interface PaginationParams {
  limit: number;
  offset: number;
}

interface MessageFilters extends PaginationParams {
  conversation_id: number;
  sender_id?: string;
  start_date?: string;
  end_date?: string;
  include_deleted: boolean;
}

interface ConversationFilters extends PaginationParams {
  user_id: string;
  include_archived: boolean;
}

interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

interface MessageResponse {
  message: MessageWithDetails;
}

interface MessagesResponse {
  messages: MessageWithDetails[];
  total: number;
  has_more: boolean;
}

interface ConversationResponse {
  conversation: ConversationWithDetails;
}

interface ConversationsResponse {
  conversations: ConversationOverview[];
  total: number;
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
  date?: string;
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

interface DeleteMessageResponse {
  message: string;
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

interface Connection {
  user_id: string;
  conversation_ids: string[];
  last_ping: string;
  connected_at: string;
}

export type {
  // Enums
  AttachmentType,
  WebSocketMessageType,
  
  // Conversation Types
  Conversation,
  ConversationWithDetails,
  ConversationOverview,
  
  // Message Types
  Message,
  MessageWithDetails,
  MessageReadStatus,
  MessageAttachment,
  
  // Request Types
  CreateConversationRequest,
  SendMessageRequest,
  UpdateMessageRequest,
  UploadAttachmentRequest,
  MarkAsReadRequest,
  
  // Filter & Pagination
  PaginationParams,
  MessageFilters,
  ConversationFilters,
  
  // Response Types
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
  
  // WebSocket Types
  WebSocketMessage,
  Connection,
};
