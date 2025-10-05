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
	ErrInvalidMealType = Error{Code: "invalid_meal_type", Message: "Invalid meal type provided"}
	ErrFailedToCreateLogEntry = Error{Code: "failed_to_create_log_entry", Message: "Failed to create food log entry"}
	ErrFailedToUpdateLogEntry = Error{Code: "failed_to_update_log_entry", Message: "Failed to update food log entry"}
	ErrUnauthorized = Error{Code: "unauthorized", Message: "Unauthorized access"}
	ErrNurtritionValues = Error{Code: "invalid_nutrition_values", Message: "Invalid nutrition values provided"}
)


