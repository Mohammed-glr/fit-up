package types

import (
	"time"

	"golang.org/x/net/websocket"
)

type AttachmentType string

const (
	AttachmentTypeImage       AttachmentType = "image"
	AttachmentTypeDocument    AttachmentType = "document"
	AttachmentTypeWorkoutPlan AttachmentType = "workout_plan"
)

type Conversation struct {
	ConversationID int       `json:"conversation_id" db:"conversation_id"`
	CoachID        string    `json:"coach_id" db:"coach_id"`
	ClientID       string    `json:"client_id" db:"client_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	LastMessageAt  time.Time `json:"last_message_at" db:"last_message_at"`
	IsArchived     bool      `json:"is_archived" db:"is_archived"`
}

type Message struct {
	MessageID        int64      `json:"message_id" db:"message_id"`
	ConversationID   int        `json:"conversation_id" db:"conversation_id"`
	SenderID         string     `json:"sender_id" db:"sender_id"`
	MessageText      string     `json:"message_text" db:"message_text"`
	SentAt           time.Time  `json:"sent_at" db:"sent_at"`
	EditedAt         *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	IsDeleted        bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	ReplyToMessageID *int64     `json:"reply_to_message_id,omitempty" db:"reply_to_message_id"`
}

type MessageReadStatus struct {
	ReadStatusID int64     `json:"read_status_id" db:"read_status_id"`
	MessageID    int64     `json:"message_id" db:"message_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ReadAt       time.Time `json:"read_at" db:"read_at"`
}

type MessageAttachment struct {
	AttachmentID   int64          `json:"attachment_id" db:"attachment_id"`
	MessageID      int64          `json:"message_id" db:"message_id"`
	AttachmentType AttachmentType `json:"attachment_type" db:"attachment_type"`
	FileName       string         `json:"file_name" db:"file_name"`
	FileURL        string         `json:"file_url" db:"file_url"`
	FileSize       *int           `json:"file_size,omitempty" db:"file_size"`
	MimeType       *string        `json:"mime_type,omitempty" db:"mime_type"`
	UploadedAt     time.Time      `json:"uploaded_at" db:"uploaded_at"`
	Metadata       []byte         `json:"metadata,omitempty" db:"metadata"`
}

type CreateConversationRequest struct {
	CoachID  string `json:"coach_id" validate:"required"`
	ClientID string `json:"client_id" validate:"required"`
}

type SendMessageRequest struct {
	ConversationID   int    `json:"conversation_id" validate:"required"`
	MessageText      string `json:"message_text" validate:"required,min=1,max=5000"`
	ReplyToMessageID *int64 `json:"reply_to_message_id,omitempty"`
}

type UpdateMessageRequest struct {
	MessageText string `json:"message_text" validate:"required,min=1,max=5000"`
}

type UploadAttachmentRequest struct {
	MessageID      int64          `json:"message_id" validate:"required"`
	AttachmentType AttachmentType `json:"attachment_type" validate:"required"`
	FileName       string         `json:"file_name" validate:"required"`
	FileURL        string         `json:"file_url" validate:"required,url"`
	FileSize       *int           `json:"file_size,omitempty"`
	MimeType       *string        `json:"mime_type,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

type MarkAsReadRequest struct {
	MessageIDs []int64 `json:"message_ids" validate:"required,min=1"`
}

type MessageWithDetails struct {
	Message
	SenderName     string              `json:"sender_name"`
	SenderImage    *string             `json:"sender_image,omitempty"`
	Attachments    []MessageAttachment `json:"attachments,omitempty"`
	IsRead         bool                `json:"is_read"`
	ReplyToMessage *MessageWithDetails `json:"reply_to_message,omitempty"`
}

type ConversationWithDetails struct {
	Conversation
	CoachName     string              `json:"coach_name"`
	CoachImage    *string             `json:"coach_image,omitempty"`
	ClientName    string              `json:"client_name"`
	ClientImage   *string             `json:"client_image,omitempty"`
	LastMessage   *MessageWithDetails `json:"last_message,omitempty"`
	UnreadCount   int                 `json:"unread_count"`
	TotalMessages int                 `json:"total_messages"`
}

type ConversationOverview struct {
	ConversationID      int        `json:"conversation_id" db:"conversation_id"`
	CoachID             string     `json:"coach_id" db:"coach_id"`
	ClientID            string     `json:"client_id" db:"client_id"`
	CoachName           string     `json:"coach_name" db:"coach_name"`
	CoachImage          *string    `json:"coach_image,omitempty" db:"coach_image"`
	ClientName          string     `json:"client_name" db:"client_name"`
	ClientImage         *string    `json:"client_image,omitempty" db:"client_image"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	LastMessageAt       time.Time  `json:"last_message_at" db:"last_message_at"`
	IsArchived          bool       `json:"is_archived" db:"is_archived"`
	LastMessageText     *string    `json:"last_message_text,omitempty" db:"last_message_text"`
	LastMessageSenderID *string    `json:"last_message_sender_id,omitempty" db:"last_message_sender_id"`
	LastMessageSentAt   *time.Time `json:"last_message_sent_at,omitempty" db:"last_message_sent_at"`
	TotalMessages       int        `json:"total_messages" db:"total_messages"`
}

type PaginationParams struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

type MessageFilters struct {
	ConversationID int        `json:"conversation_id"`
	SenderID       *string    `json:"sender_id,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	IncludeDeleted bool       `json:"include_deleted"`
	PaginationParams
}

type ConversationFilters struct {
	UserID          string `json:"user_id"`
	IncludeArchived bool   `json:"include_archived"`
	PaginationParams
}

type PaginatedResponse[T any] struct {
	Data    []T  `json:"data"`
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type MessageResponse struct {
	Message MessageWithDetails `json:"message"`
}
type MessagesResponse struct {
	Messages []MessageWithDetails `json:"messages"`
	Total    int                  `json:"total"`
	HasMore  bool                 `json:"has_more"`
}

type ConversationResponse struct {
	Conversation ConversationWithDetails `json:"conversation"`
}

type ConversationsResponse struct {
	Conversations []ConversationOverview `json:"conversations"`
	Total         int                    `json:"total"`
	HasMore       bool                   `json:"has_more"`
}

type WebSocketMessageType string

const (
	WSTypeNewMessage     WebSocketMessageType = "new_message"
	WSTypeMessageEdited  WebSocketMessageType = "message_edited"
	WSTypeMessageRead    WebSocketMessageType = "message_read"
	WSTypeMessageDeleted WebSocketMessageType = "message_deleted"
	WSTypeError          WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type           WebSocketMessageType `json:"type"`
	ConversationID int                  `json:"conversation_id"`
	Message        *MessageWithDetails  `json:"message,omitempty"`
	MessageID      *int64               `json:"message_id,omitempty"`
	ReadBy         *string              `json:"read_by,omitempty"`
	Error          *string              `json:"error,omitempty"`
	Timestamp      time.Time            `json:"timestamp"`
}

type Connection struct {
	Conn            *websocket.Conn `json:"-"`
	UserID          string          `json:"user_id"`
	ConversationIDs []string        `json:"conversation_ids"`
	LastPing        time.Time       `json:"last_ping"`
	ConnectedAt     time.Time       `json:"connected_at"`
}
