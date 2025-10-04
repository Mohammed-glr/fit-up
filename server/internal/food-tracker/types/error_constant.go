package types

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrNotFound      = Error{Code: "not_found", Message: "Resource not found"}
	ErrInvalidID     = Error{Code: "invalid_id", Message: "Invalid ID provided"}
	ErrInvalidRequest = Error{Code: "invalid_request", Message: "Invalid request data"}
)

