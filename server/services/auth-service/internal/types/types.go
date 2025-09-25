package types

import (
	"time"
)

// =============================================================================
// ENUMS AND CONSTANTS
// =============================================================================

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// =============================================================================
// USER RELATED TYPES
// =============================================================================

type User struct {
	ID                 string     `json:"id" db:"id"`
	Username           string     `json:"username" db:"username"`
	Name               string     `json:"name" db:"name"`
	Bio                string     `json:"bio" db:"bio"`
	Email              string     `json:"email" db:"email"`
	EmailVerified      *time.Time `json:"email_verified" db:"email_verified"`
	Image              string     `json:"image" db:"image"`
	Password           string     `json:"-" db:"password"`
	PasswordHash       string     `json:"-" db:"password_hash"`
	Role               UserRole   `json:"role" db:"role"`
	IsTwoFactorEnabled bool       `json:"is_two_factor_enabled" db:"is_two_factor_enabled"`
	SubroleID          int        `json:"subrole_id" db:"subrole_id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID                 string    `json:"id"`
	Username           string    `json:"username"`
	Name               string    `json:"name"`
	Bio                string    `json:"bio"`
	Email              string    `json:"email"`
	EmailVerified      time.Time `json:"email_verified"`
	Image              *string   `json:"image"`
	Role               UserRole  `json:"role"`
	IsTwoFactorEnabled bool      `json:"is_two_factor_enabled"`
	SubroleID          int       `json:"subrole_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PublicUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	Image     *string   `json:"image"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Subrole struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// =============================================================================
// 0 AUTHENTICATION TYPES
// =============================================================================

type Account struct {
	ID                string `json:"id" db:"id"`
	UserID            string `json:"user_id" db:"user_id"`
	Type              string `json:"type" db:"type"`
	Provider          string `json:"provider" db:"provider"`
	ProviderAccountID string `json:"provider_account_id" db:"provider_account_id"`
	RefreshToken      string `json:"refresh_token" db:"refresh_token"`
	AccessToken       string `json:"access_token" db:"access_token"`
	ExpiresAt         int    `json:"expires_at" db:"expires_at"`
	TokenType         string `json:"token_type" db:"token_type"`
	Scope             string `json:"scope" db:"scope"`
	IDToken           string `json:"id_token" db:"id_token"`
	SessionState      string `json:"session_state" db:"session_state"`
}

// =============================================================================
// OAUTH TYPES
// =============================================================================

// OAuth Provider Configuration
type OAuthProvider struct {
	Name         string   `json:"name"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"-"` // Never expose
	RedirectURI  string   `json:"redirect_uri"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	UserInfoURL  string   `json:"user_info_url"`
	Scopes       []string `json:"scopes"`
}

// OAuth Authorization Request
type OAuthAuthRequest struct {
	Provider    string `json:"provider" validate:"required,oneof=google github facebook"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// OAuth Callback Request
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// OAuth User Info (from provider)
type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Username      string `json:"username,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	EmailVerified bool   `json:"email_verified"`
}

// OAuth Link Account Request
type LinkAccountRequest struct {
	Provider string `json:"provider" validate:"required"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

// OAuth State (for CSRF protection)
type OAuthState struct {
	ID          string    `json:"id" db:"id"`
	State       string    `json:"state" db:"state"`
	Provider    string    `json:"provider" db:"provider"`
	RedirectURL string    `json:"redirect_url" db:"redirect_url"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// =============================================================================
// AUTHENTICATION REQUEST/RESPONSE TYPES
// =============================================================================

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Can be email or username
	Password   string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"max=100"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	User         *User  `json:"user"`
	ExpiresAt    int64  `json:"expires_at"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// =============================================================================
// USER MANAGEMENT REQUEST TYPES
// =============================================================================

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Name     *string `json:"name" validate:"omitempty,max=100"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// =============================================================================
// PASSWORD RESET TYPES
// =============================================================================

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type PasswordResetToken struct {
	ID      string    `json:"id" db:"id"`
	Email   string    `json:"email" db:"email"`
	Token   string    `json:"token" db:"token"`
	Expires time.Time `json:"expires" db:"expires"`
}

// =============================================================================
// EMAIL VERIFICATION TYPES
// =============================================================================

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerificationToken struct {
	ID      string    `json:"id" db:"id"`
	Email   string    `json:"email" db:"email"`
	Token   string    `json:"token" db:"token"`
	Expires time.Time `json:"expires" db:"expires"`
}

// =============================================================================
// TWO-FACTOR AUTHENTICATION TYPES
// =============================================================================

type TwoFactorSetupRequest struct {
	Token string `json:"token" validate:"required"`
}

type TwoFactorVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type TwoFactorSetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type TwoFactorToken struct {
	ID      string    `json:"id" db:"id"`
	Email   string    `json:"email" db:"email"`
	Token   string    `json:"token" db:"token"`
	Expires time.Time `json:"expires" db:"expires"`
}

type TwoFactorConfirmation struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
}

// =============================================================================
// SESSION AND ACCOUNT TYPES
// =============================================================================

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	SessionToken string    `json:"session_token" db:"session_token"`
	Expires      time.Time `json:"expires" db:"expires"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type EnhancedSession struct {
	Session                  // Embed existing Session struct
	AccessTokenJTI *string   `json:"access_token_jti,omitempty" db:"access_token_jti"`
	RefreshTokenID *string   `json:"refresh_token_id,omitempty" db:"refresh_token_id"`
	IPAddress      *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent      *string   `json:"user_agent,omitempty" db:"user_agent"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	LastActivityAt time.Time `json:"last_activity_at" db:"last_activity_at"`
}

type LoginActivity struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	IPAddress *string   `json:"ip_address" db:"ip_address"`
	UserAgent *string   `json:"user_agent" db:"user_agent"`
	Location  *string   `json:"location" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// =============================================================================
// JWT AND TOKEN MANAGEMENT TYPES
// =============================================================================

type TokenClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
	JTI       string   `json:"jti"`
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	Audience  string   `json:"aud"`
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
	NotBefore int64    `json:"nbf"`
}

type JWTRefreshToken struct {
	ID             string     `json:"id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	TokenHash      string     `json:"-" db:"token_hash"` // Never expose the hash
	AccessTokenJTI *string    `json:"access_token_jti,omitempty" db:"access_token_jti"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt     time.Time  `json:"last_used_at" db:"last_used_at"`
	IsRevoked      bool       `json:"is_revoked" db:"is_revoked"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	UserAgent      *string    `json:"user_agent,omitempty" db:"user_agent"`
	IPAddress      *string    `json:"ip_address,omitempty" db:"ip_address"`
}

type JWTBlacklist struct {
	ID            string    `json:"id" db:"id"`
	JTI           string    `json:"jti" db:"jti"`      // JWT ID claim
	TokenHash     string    `json:"-" db:"token_hash"` // Never expose the hash
	UserID        string    `json:"user_id" db:"user_id"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	BlacklistedAt time.Time `json:"blacklisted_at" db:"blacklisted_at"`
	Reason        string    `json:"reason" db:"reason"`
}

type RefreshToken struct {
	ID             string     `json:"id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	TokenHash      string     `json:"-" db:"token_hash"` // Never expose the hash
	AccessTokenJTI string     `json:"access_token_jti,omitempty" db:"access_token_jti"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt     time.Time  `json:"last_used_at" db:"last_used_at"`
	IsRevoked      bool       `json:"is_revoked" db:"is_revoked"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	UserAgent      string     `json:"user_agent,omitempty" db:"user_agent"`
	IPAddress      string     `json:"ip_address,omitempty" db:"ip_address"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// =============================================================================
// JWT REQUEST/RESPONSE TYPES
// =============================================================================

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RevokeTokenRequest struct {
	Token     string `json:"token" validate:"required"`
	TokenType string `json:"token_type,omitempty"` // access_token, refresh_token
}

type TokenInfoResponse struct {
	Active    bool                   `json:"active"`
	Claims    *TokenClaims           `json:"claims,omitempty"`
	ExpiresIn int64                  `json:"expires_in,omitempty"`
	IssuedAt  int64                  `json:"issued_at,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// =============================================================================
// AUDIT AND MONITORING TYPES
// =============================================================================

type AuthAuditLog struct {
	ID        string                 `json:"id" db:"id"`
	UserID    *string                `json:"user_id,omitempty" db:"user_id"`
	Action    string                 `json:"action" db:"action"`
	Success   bool                   `json:"success" db:"success"`
	IPAddress *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string                `json:"user_agent,omitempty" db:"user_agent"`
	Details   map[string]interface{} `json:"details,omitempty" db:"details"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

type RateLimit struct {
	ID           string     `json:"id" db:"id"`
	Identifier   string     `json:"identifier" db:"identifier"`
	Endpoint     string     `json:"endpoint" db:"endpoint"`
	Attempts     int        `json:"attempts" db:"attempts"`
	WindowStart  time.Time  `json:"window_start" db:"window_start"`
	WindowEnd    time.Time  `json:"window_end" db:"window_end"`
	IsBlocked    bool       `json:"is_blocked" db:"is_blocked"`
	BlockedUntil *time.Time `json:"blocked_until,omitempty" db:"blocked_until"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type TokenUsageStats struct {
	ID        string    `json:"id" db:"id"`
	UserID    *string   `json:"user_id,omitempty" db:"user_id"`
	TokenType string    `json:"token_type" db:"token_type"` // access, refresh
	Action    string    `json:"action" db:"action"`         // generate, validate, refresh, revoke
	Success   bool      `json:"success" db:"success"`
	IPAddress *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string   `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AuditLogEntry struct {
	Action    string                 `json:"action"`
	Success   bool                   `json:"success"`
	UserID    *string                `json:"user_id,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// =============================================================================
// AUDIT AND SECURITY TYPES
// =============================================================================

// AuthAuditEvent represents an authentication audit log entry
type AuthAuditEvent struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Action    string                 `json:"action" db:"action"`
	Success   bool                   `json:"success" db:"success"`
	IPAddress string                 `json:"ip_address" db:"ip_address"`
	UserAgent string                 `json:"user_agent" db:"user_agent"`
	Details   map[string]interface{} `json:"details" db:"details"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

// SuspiciousActivityReport represents analysis of potentially suspicious activity
type SuspiciousActivityReport struct {
	IPAddress    string         `json:"ip_address"`
	TimeWindow   time.Duration  `json:"time_window"`
	TotalEvents  int            `json:"total_events"`
	FailedLogins int            `json:"failed_logins"`
	UserAgents   map[string]int `json:"user_agents"`
	Actions      map[string]int `json:"actions"`
	IsSuspicious bool           `json:"is_suspicious"`
}

// =============================================================================
// ERROR TYPES
// =============================================================================

type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e AuthError) Error() string {
	return e.Message
}

// =============================================================================
// ERROR CONSTANTS
// =============================================================================

var (
	// Auth Errors
	ErrInvalidCredentials      = AuthError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrUserNotFound            = AuthError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUserAlreadyExists       = AuthError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
	ErrUsernameAlreadyExists   = AuthError{Code: "USERNAME_ALREADY_EXISTS", Message: "Username already exists"}
	ErrInvalidToken            = AuthError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrSessionExpired          = AuthError{Code: "SESSION_EXPIRED", Message: "Session expired"}
	ErrUnauthorized            = AuthError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrEmailNotVerified        = AuthError{Code: "EMAIL_NOT_VERIFIED", Message: "Email not verified"}
	ErrTwoFactorRequired       = AuthError{Code: "TWO_FACTOR_REQUIRED", Message: "Two-factor authentication required"}
	ErrInvalidTwoFactorToken   = AuthError{Code: "INVALID_TWO_FACTOR_TOKEN", Message: "Invalid two-factor authentication token"}

	// Account Management Errors
	ErrAccountLocked        = AuthError{Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked"}
	ErrAccountDisabled      = AuthError{Code: "ACCOUNT_DISABLED", Message: "Account has been disabled"}
	ErrPasswordTooWeak      = AuthError{Code: "PASSWORD_TOO_WEAK", Message: "Password does not meet security requirements"}
	ErrPasswordRecentlyUsed = AuthError{Code: "PASSWORD_RECENTLY_USED", Message: "Password was recently used"}

	// Rate Limiting & Security
	ErrTooManyAttempts    = AuthError{Code: "TOO_MANY_ATTEMPTS", Message: "Too many failed attempts, please try again later"}
	ErrSuspiciousActivity = AuthError{Code: "SUSPICIOUS_ACTIVITY", Message: "Suspicious activity detected"}

	// Email & Verification
	ErrFailedToDeleteVerificationT	 = AuthError{Code: "FAILED_TO_DELETE_VERIFICATION_TOKEN", Message: "Failed to delete verification token"}
	ErrEmailAlreadyVerified          = AuthError{Code: "EMAIL_ALREADY_VERIFIED", Message: "Email is already verified"}
	ErrEmailNotFound                 = AuthError{Code: "EMAIL_NOT_FOUND", Message: "No account found with this email"}
	ErrInvalidEmailFormat            = AuthError{Code: "INVALID_EMAIL_FORMAT", Message: "Invalid email format"}
	ErrVerificationTokenExpired      = AuthError{Code: "VERIFICATION_TOKEN_EXPIRED", Message: "Verification token has expired"}
	ErrVerificationTokenAlreadyUsed  = AuthError{Code: "VERIFICATION_TOKEN_ALREADY_USED", Message: "Verification token has already been used"}
	ErrVerificationTokenNotFound     = AuthError{Code: "VERIFICATION_TOKEN_NOT_FOUND", Message: "Verification token not found"}
	ErrFailedToVerifyEmail           = AuthError{Code: "FAILED_TO_VERIFY_EMAIL", Message: "Failed to verify email address"}
	ErrFailedToSendVerificationEmail = AuthError{Code: "FAILED_TO_SEND_VERIFICATION_EMAIL", Message: "Failed to send verification email"}
	ErrFailedToResendVerification    = AuthError{Code: "FAILED_TO_RESEND_VERIFICATION", Message: "Failed to resend verification email"}

	// Two-Factor Authentication
	ErrTwoFactorNotEnabled     = AuthError{Code: "TWO_FACTOR_NOT_ENABLED", Message: "Two-factor authentication is not enabled"}
	ErrTwoFactorAlreadyEnabled = AuthError{Code: "TWO_FACTOR_ALREADY_ENABLED", Message: "Two-factor authentication is already enabled"}
	ErrInvalidBackupCode       = AuthError{Code: "INVALID_BACKUP_CODE", Message: "Invalid backup code"}

	// Session & Token Management
	ErrInvalidUserID        = AuthError{Code: "INVALID_USER_ID", Message: "Invalid user ID provided"}
	ErrSessionNotFound      = AuthError{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrTokenExpired         = AuthError{Code: "TOKEN_EXPIRED", Message: "Token has expired"}
	ErrTokenAlreadyUsed     = AuthError{Code: "TOKEN_ALREADY_USED", Message: "Token has already been used"}
	ErrTokenBlacklisted     = AuthError{Code: "TOKEN_BLACKLISTED", Message: "Token has been blacklisted"}
	ErrRefreshTokenNotFound = AuthError{Code: "REFRESH_TOKEN_NOT_FOUND", Message: "Refresh token not found"}
	ErrRefreshTokenExpired  = AuthError{Code: "REFRESH_TOKEN_EXPIRED", Message: "Refresh token has expired"}
	ErrInvalidRefreshToken  = AuthError{Code: "INVALID_REFRESH_TOKEN", Message: "Invalid or expired refresh token"}

	// Password Reset Specific Errors
	ErrPasswordResetTokenNotFound = AuthError{Code: "PASSWORD_RESET_TOKEN_NOT_FOUND", Message: "Password reset token not found"}
	ErrPasswordResetTokenExpired  = AuthError{Code: "PASSWORD_RESET_TOKEN_EXPIRED", Message: "Password reset token has expired"}
	ErrPasswordResetTokenUsed     = AuthError{Code: "PASSWORD_RESET_TOKEN_USED", Message: "Password reset token has already been used"}
	ErrSamePassword               = AuthError{Code: "SAME_PASSWORD", Message: "New password cannot be the same as current password"}
	ErrIncorrectCurrentPassword   = AuthError{Code: "INCORRECT_CURRENT_PASSWORD", Message: "Current password is incorrect"}

	// Validation & Input
	ErrInvalidInput         = AuthError{Code: "INVALID_INPUT", Message: "Invalid input provided"}
	ErrMissingRequiredField = AuthError{Code: "MISSING_REQUIRED_FIELD", Message: "Required field is missing"}
	ErrUsernameTaken        = AuthError{Code: "USERNAME_TAKEN", Message: "Username is already taken"}

	// Permission & Role Errors
	ErrInsufficientPermissions = AuthError{Code: "INSUFFICIENT_PERMISSIONS", Message: "Insufficient permissions to perform this action"}
	ErrInvalidRole             = AuthError{Code: "INVALID_ROLE", Message: "Invalid role specified"}

	// Provider & OAuth Errors
	ErrProviderError        = AuthError{Code: "PROVIDER_ERROR", Message: "Authentication provider error"}
	ErrProviderNotSupported = AuthError{Code: "PROVIDER_NOT_SUPPORTED", Message: "Authentication provider not supported"}
	ErrAccountNotLinked     = AuthError{Code: "ACCOUNT_NOT_LINKED", Message: "Account is not linked to this provider"}

	// Miscellaneous Errors
	ErrInternalServerError = AuthError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error"}
	ErrServiceUnavailable  = AuthError{Code: "SERVICE_UNAVAILABLE", Message: "Service is currently unavailable"}
	ErrMaintenanceMode     = AuthError{Code: "MAINTENANCE_MODE", Message: "Service is in maintenance mode"}

	// Secret Management Errors
	ErrJWTSecretNotSet = AuthError{Code: "JWT_SECRET_NOT_SET", Message: "JWT secret is not set in the configuration"}
)
