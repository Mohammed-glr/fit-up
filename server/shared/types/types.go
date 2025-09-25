package types

// import (
// 	"time"
// )

// type UserRole string

// const (
// 	RoleAdmin UserRole = "admin"
// 	RoleUser  UserRole = "user"
// )

// type User struct {
// 	ID                 string    `json:"id" db:"id"`
// 	Username           string    `json:"username" db:"username"`
// 	Name               string    `json:"name" db:"name"`
// 	Bio                string    `json:"bio" db:"bio"`
// 	Email              string    `json:"email" db:"email"`
// 	EmailVerified      time.Time `json:"email_verified" db:"email_verified"`
// 	Image              string    `json:"image" db:"image"`
// 	Password           string    `json:"-" db:"password"`
// 	PasswordHash       string    `json:"-" db:"password_hash"`
// 	Role               UserRole  `json:"role" db:"role"`
// 	IsTwoFactorEnabled bool      `json:"is_two_factor_enabled" db:"is_two_factor_enabled"`
// 	SubroleID          int       `json:"subrole_id" db:"subrole_id"`
// 	CreatedAt          time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
// }

// type Account struct {
// 	ID                string `json:"id" db:"id"`
// 	UserID            string `json:"user_id" db:"user_id"`
// 	Type              string `json:"type" db:"type"`
// 	Provider          string `json:"provider" db:"provider"`
// 	ProviderAccountID string `json:"provider_account_id" db:"provider_account_id"`
// 	RefreshToken      string `json:"refresh_token" db:"refresh_token"`
// 	AccessToken       string `json:"access_token" db:"access_token"`
// 	ExpiresAt         int    `json:"expires_at" db:"expires_at"`
// 	TokenType         string `json:"token_type" db:"token_type"`
// 	Scope             string `json:"scope" db:"scope"`
// 	IDToken           string `json:"id_token" db:"id_token"`
// 	SessionState      string `json:"session_state" db:"session_state"`
// }

// type Session struct {
// 	ID           string    `json:"id" db:"id"`
// 	UserID       string    `json:"user_id" db:"user_id"`
// 	SessionToken string    `json:"session_token" db:"session_token"`
// 	Expires      time.Time `json:"expires" db:"expires"`
// 	CreatedAt    time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
// }

// type VerificationToken struct {
// 	ID      string    `json:"id" db:"id"`
// 	Email   string    `json:"email" db:"email"`
// 	Token   string    `json:"token" db:"token"`
// 	Expires time.Time `json:"expires" db:"expires"`
// }

// type PasswordResetToken struct {
// 	ID      string    `json:"id" db:"id"`
// 	Email   string    `json:"email" db:"email"`
// 	Token   string    `json:"token" db:"token"`
// 	Expires time.Time `json:"expires" db:"expires"`
// }

// type TwoFactorToken struct {
// 	ID      string    `json:"id" db:"id"`
// 	Email   string    `json:"email" db:"email"`
// 	Token   string    `json:"token" db:"token"`
// 	Expires time.Time `json:"expires" db:"expires"`
// }

// type TwoFactorConfirmation struct {
// 	ID     string `json:"id" db:"id"`
// 	UserID string `json:"user_id" db:"user_id"`
// }

// type LoginActivity struct {
// 	ID        string    `json:"id" db:"id"`
// 	UserID    string    `json:"user_id" db:"user_id"`
// 	IPAddress *string   `json:"ip_address" db:"ip_address"`
// 	UserAgent *string   `json:"user_agent" db:"user_agent"`
// 	Location  *string   `json:"location" db:"location"`
// 	CreatedAt time.Time `json:"created_at" db:"created_at"`
// }

// type Subrole struct {
// 	ID          int    `json:"id" db:"id"`
// 	Name        string `json:"name" db:"name"`
// 	Description string `json:"description" db:"description"`
// }

// type LoginRequest struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8"`
// }

// type RegisterRequest struct {
// 	Username string `json:"username" validate:"required,min=3,max=50"`
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8"`
// 	Name     string `json:"name" validate:"max=100"`
// }

// type LoginResponse struct {
// 	User         *User  `json:"user"`
// 	SessionToken string `json:"session_token"`
// 	ExpiresAt    int64  `json:"expires_at"`
// }

// type UserResponse struct {
// 	ID                 string    `json:"id"`
// 	Username           string    `json:"username"`
// 	Name               string    `json:"name"`
// 	Bio                string    `json:"bio"`
// 	Email              string    `json:"email"`
// 	EmailVerified      time.Time `json:"email_verified"`
// 	Image              *string   `json:"image"`
// 	Role               UserRole  `json:"role"`
// 	IsTwoFactorEnabled bool      `json:"is_two_factor_enabled"`
// 	SubroleID          *int      `json:"subrole_id"`
// 	CreatedAt          time.Time `json:"created_at"`
// 	UpdatedAt          time.Time `json:"updated_at"`
// }

// type UpdateUserRequest struct {
// 	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
// 	Name     *string `json:"name" validate:"omitempty,max=100"`
// 	Bio      *string `json:"bio"`
// 	Image    *string `json:"image"`
// }

// type ChangePasswordRequest struct {
// 	CurrentPassword string `json:"current_password" validate:"required"`
// 	NewPassword     string `json:"new_password" validate:"required,min=8"`
// }

// type ForgotPasswordRequest struct {
// 	Email string `json:"email" validate:"required,email"`
// }

// type ResetPasswordRequest struct {
// 	Token       string `json:"token" validate:"required"`
// 	NewPassword string `json:"new_password" validate:"required,min=8"`
// }

// type VerifyEmailRequest struct {
// 	Token string `json:"token" validate:"required"`
// }

// type TwoFactorSetupRequest struct {
// 	Token string `json:"token" validate:"required"`
// }

// type TwoFactorVerifyRequest struct {
// 	Email string `json:"email" validate:"required,email"`
// 	Token string `json:"token" validate:"required"`
// }

// type TwoFactorSetupResponse struct {
// 	Secret      string   `json:"secret"`
// 	QRCodeURL   string   `json:"qr_code_url"`
// 	BackupCodes []string `json:"backup_codes"`
// }

// type AuthError struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }

// func (e AuthError) Error() string {
// 	return e.Message
// }
// var (
// 	// Auth Errors
// 	ErrInvalidCredentials    = AuthError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
// 	ErrUserNotFound          = AuthError{Code: "USER_NOT_FOUND", Message: "User not found"}
// 	ErrUserAlreadyExists     = AuthError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
// 	ErrInvalidToken          = AuthError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
// 	ErrSessionExpired        = AuthError{Code: "SESSION_EXPIRED", Message: "Session expired"}
// 	ErrUnauthorized          = AuthError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
// 	ErrEmailNotVerified      = AuthError{Code: "EMAIL_NOT_VERIFIED", Message: "Email not verified"}
// 	ErrTwoFactorRequired     = AuthError{Code: "TWO_FACTOR_REQUIRED", Message: "Two-factor authentication required"}
// 	ErrInvalidTwoFactorToken = AuthError{Code: "INVALID_TWO_FACTOR_TOKEN", Message: "Invalid two-factor authentication token"}

// 	// Account Management Errors
// 	ErrAccountLocked         = AuthError{Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked"}
// 	ErrAccountDisabled       = AuthError{Code: "ACCOUNT_DISABLED", Message: "Account has been disabled"}
// 	ErrPasswordTooWeak       = AuthError{Code: "PASSWORD_TOO_WEAK", Message: "Password does not meet security requirements"}
// 	ErrPasswordRecentlyUsed  = AuthError{Code: "PASSWORD_RECENTLY_USED", Message: "Password was recently used"}

// 	// Rate Limiting & Security
// 	ErrTooManyAttempts       = AuthError{Code: "TOO_MANY_ATTEMPTS", Message: "Too many failed attempts, please try again later"}
// 	ErrSuspiciousActivity    = AuthError{Code: "SUSPICIOUS_ACTIVITY", Message: "Suspicious activity detected"}

// 	// Email & Verification
// 	ErrEmailAlreadyVerified  = AuthError{Code: "EMAIL_ALREADY_VERIFIED", Message: "Email is already verified"}
// 	ErrEmailNotFound         = AuthError{Code: "EMAIL_NOT_FOUND", Message: "No account found with this email"}
// 	ErrInvalidEmailFormat    = AuthError{Code: "INVALID_EMAIL_FORMAT", Message: "Invalid email format"}

// 	// Two-Factor Authentication
// 	ErrTwoFactorNotEnabled   = AuthError{Code: "TWO_FACTOR_NOT_ENABLED", Message: "Two-factor authentication is not enabled"}
// 	ErrTwoFactorAlreadyEnabled = AuthError{Code: "TWO_FACTOR_ALREADY_ENABLED", Message: "Two-factor authentication is already enabled"}
// 	ErrInvalidBackupCode     = AuthError{Code: "INVALID_BACKUP_CODE", Message: "Invalid backup code"}

// 	// Session & Token Management
// 	ErrSessionNotFound       = AuthError{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
// 	ErrTokenExpired          = AuthError{Code: "TOKEN_EXPIRED", Message: "Token has expired"}
// 	ErrTokenAlreadyUsed      = AuthError{Code: "TOKEN_ALREADY_USED", Message: "Token has already been used"}

// 	// Validation & Input
// 	ErrInvalidInput          = AuthError{Code: "INVALID_INPUT", Message: "Invalid input provided"}
// 	ErrMissingRequiredField  = AuthError{Code: "MISSING_REQUIRED_FIELD", Message: "Required field is missing"}
// 	ErrUsernameTaken         = AuthError{Code: "USERNAME_TAKEN", Message: "Username is already taken"}

// 	// Permission & Role Errors
// 	ErrInsufficientPermissions = AuthError{Code: "INSUFFICIENT_PERMISSIONS", Message: "Insufficient permissions to perform this action"}
// 	ErrInvalidRole           = AuthError{Code: "INVALID_ROLE", Message: "Invalid role specified"}

// 	// Provider & OAuth Errors
// 	ErrProviderError         = AuthError{Code: "PROVIDER_ERROR", Message: "Authentication provider error"}
// 	ErrProviderNotSupported  = AuthError{Code: "PROVIDER_NOT_SUPPORTED", Message: "Authentication provider not supported"}
// 	ErrAccountNotLinked      = AuthError{Code: "ACCOUNT_NOT_LINKED", Message: "Account is not linked to this provider"}
// )
