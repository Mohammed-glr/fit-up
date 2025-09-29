package types


type SchemaError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *SchemaError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound		 = &SchemaError{Code: "USER_NOT_FOUND", Message: "User not found"}
)