package types

import "errors"

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrConversationExists   = errors.New("conversation already exists")
	ErrInvalidConversation  = errors.New("invalid conversation participants")
	ErrConversationArchived = errors.New("conversation is archived")

	ErrMessageNotFound     = errors.New("message not found")
	ErrMessageEmpty        = errors.New("message text cannot be empty")
	ErrMessageTooLong      = errors.New("message exceeds maximum length")
	ErrMessageDeleted      = errors.New("message has been deleted")
	ErrCannotEditMessage   = errors.New("cannot edit this message")
	ErrCannotDeleteMessage = errors.New("cannot delete this message")

	ErrAttachmentNotFound  = errors.New("attachment not found")
	ErrInvalidAttachment   = errors.New("invalid attachment")
	ErrAttachmentTooLarge  = errors.New("attachment exceeds size limit")
	ErrUnsupportedFileType = errors.New("unsupported file type")

	ErrUnauthorized     = errors.New("not authorized to access this resource")
	ErrNotParticipant   = errors.New("user is not a participant in this conversation")
	ErrNotMessageSender = errors.New("user is not the sender of this message")

	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrInvalidConversationID = errors.New("invalid conversation ID")
	ErrInvalidMessageID      = errors.New("invalid message ID")
	ErrInvalidPagination     = errors.New("invalid pagination parameters")

	ErrInternalServer = errors.New("internal server error")
	ErrDatabaseError  = errors.New("database error")
)

const (
	StatusConversationNotFound = 404
	StatusMessageNotFound      = 404
	StatusUnauthorized         = 403
	StatusInvalidRequest       = 400
	StatusConversationExists   = 409
	StatusInternalError        = 500
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(err error, details ...string) *ErrorResponse {
	response := &ErrorResponse{
		Error: err.Error(),
		Code:  getErrorCode(err),
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	return response
}

func getErrorCode(err error) string {
	switch err {
	case ErrConversationNotFound:
		return "CONVERSATION_NOT_FOUND"
	case ErrMessageNotFound:
		return "MESSAGE_NOT_FOUND"
	case ErrUnauthorized:
		return "UNAUTHORIZED"
	case ErrInvalidConversation:
		return "INVALID_CONVERSATION"
	case ErrMessageTooLong:
		return "MESSAGE_TOO_LONG"
	case ErrMessageEmpty:
		return "MESSAGE_EMPTY"
	case ErrInvalidAttachment:
		return "INVALID_ATTACHMENT"
	case ErrConversationExists:
		return "CONVERSATION_EXISTS"
	case ErrNotParticipant:
		return "NOT_PARTICIPANT"
	case ErrMessageDeleted:
		return "MESSAGE_DELETED"
	default:
		return "INTERNAL_ERROR"
	}
}

func GetHTTPStatus(err error) int {
	switch err {
	case ErrConversationNotFound, ErrMessageNotFound, ErrAttachmentNotFound:
		return StatusConversationNotFound
	case ErrUnauthorized, ErrNotParticipant, ErrNotMessageSender:
		return StatusUnauthorized
	case ErrConversationExists:
		return StatusConversationExists
	case ErrInvalidConversation, ErrMessageEmpty, ErrMessageTooLong,
		ErrInvalidAttachment, ErrInvalidUserID, ErrInvalidConversationID,
		ErrInvalidMessageID, ErrInvalidPagination:
		return StatusInvalidRequest
	default:
		return StatusInternalError
	}
}
