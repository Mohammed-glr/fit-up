package types

import (
	"time"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleCoach  UserRole = "coach"
	RoleUser   UserRole = "user"
	RoleClient UserRole = "client"
)

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
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID                 string     `json:"id"`
	Username           string     `json:"username"`
	Name               string     `json:"name"`
	Bio                string     `json:"bio"`
	Email              string     `json:"email"`
	EmailVerified      *time.Time `json:"email_verified"`
	Image              *string    `json:"image"`
	Role               UserRole   `json:"role"`
	IsTwoFactorEnabled bool       `json:"is_two_factor_enabled"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
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

type OAuthProvider struct {
	Name         string   `json:"name"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"-"`
	RedirectURI  string   `json:"redirect_uri"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	UserInfoURL  string   `json:"user_info_url"`
	Scopes       []string `json:"scopes"`
	SupportsPKCE bool     `json:"supports_pkce"`
}

type OAuthAuthRequest struct {
	Provider    string `json:"provider" validate:"required,oneof=google github facebook"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type OAuthPKCECallbackRequest struct {
	Code         string `json:"code" validate:"required"`
	CodeVerifier string `json:"code_verifier" validate:"required"`
	State        string `json:"state,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
}

type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Username      string `json:"username,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	EmailVerified bool   `json:"email_verified"`
}

type LinkAccountRequest struct {
	Provider string `json:"provider" validate:"required"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

type OAuthState struct {
	ID          string    `json:"id" db:"id"`
	State       string    `json:"state" db:"state"`
	Provider    string    `json:"provider" db:"provider"`
	RedirectURL string    `json:"redirect_url" db:"redirect_url"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Name     string   `json:"name" validate:"max=100"`
	Role     UserRole `json:"role" validate:"omitempty,oneof=user coach client"`
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

type UpdateRoleRequest struct {
	Role UserRole `json:"role" validate:"required,oneof=user coach client"`
}

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

type VerificationToken struct {
	ID      string    `json:"id" db:"id"`
	Email   string    `json:"email" db:"email"`
	Token   string    `json:"token" db:"token"`
	Expires time.Time `json:"expires" db:"expires_at"`
}

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

type RefreshToken struct {
	ID             string     `json:"id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	TokenHash      string     `json:"-" db:"token_hash"`
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RevokeTokenRequest struct {
	Token     string `json:"token" validate:"required"`
	TokenType string `json:"token_type,omitempty"`
}

type TokenInfoResponse struct {
	Active    bool                   `json:"active"`
	Claims    *TokenClaims           `json:"claims,omitempty"`
	ExpiresIn int64                  `json:"expires_in,omitempty"`
	IssuedAt  int64                  `json:"issued_at,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e AuthError) Error() string {
	return e.Message
}

var (
	ErrInvalidCredentials        = AuthError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrUserNotFound              = AuthError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUserAlreadyExists         = AuthError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
	ErrUsernameAlreadyExists     = AuthError{Code: "USERNAME_ALREADY_EXISTS", Message: "Username already exists"}
	ErrInvalidToken              = AuthError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrSessionExpired            = AuthError{Code: "SESSION_EXPIRED", Message: "Session expired"}
	ErrUnauthorized              = AuthError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrEmailNotVerified          = AuthError{Code: "EMAIL_NOT_VERIFIED", Message: "Email not verified"}
	ErrEmailAlreadyVerified      = AuthError{Code: "EMAIL_ALREADY_VERIFIED", Message: "Email is already verified"}
	ErrVerificationTokenNotFound = AuthError{Code: "VERIFICATION_TOKEN_NOT_FOUND", Message: "Verification token not found"}
	ErrVerificationTokenExpired  = AuthError{Code: "VERIFICATION_TOKEN_EXPIRED", Message: "Verification token has expired"}

	ErrAccountLocked        = AuthError{Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked"}
	ErrAccountDisabled      = AuthError{Code: "ACCOUNT_DISABLED", Message: "Account has been disabled"}
	ErrPasswordTooWeak      = AuthError{Code: "PASSWORD_TOO_WEAK", Message: "Password does not meet security requirements"}
	ErrPasswordRecentlyUsed = AuthError{Code: "PASSWORD_RECENTLY_USED", Message: "Password was recently used"}

	ErrTooManyAttempts    = AuthError{Code: "TOO_MANY_ATTEMPTS", Message: "Too many failed attempts, please try again later"}
	ErrSuspiciousActivity = AuthError{Code: "SUSPICIOUS_ACTIVITY", Message: "Suspicious activity detected"}

	ErrEmailNotFound      = AuthError{Code: "EMAIL_NOT_FOUND", Message: "No account found with this email"}
	ErrInvalidEmailFormat = AuthError{Code: "INVALID_EMAIL_FORMAT", Message: "Invalid email format"}

	ErrInvalidUserID    = AuthError{Code: "INVALID_USER_ID", Message: "Invalid user ID provided"}
	ErrSessionNotFound  = AuthError{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrTokenExpired     = AuthError{Code: "TOKEN_EXPIRED", Message: "Token has expired"}
	ErrTokenAlreadyUsed = AuthError{Code: "TOKEN_ALREADY_USED", Message: "Token has already been used"}

	ErrRefreshTokenNotFound = AuthError{Code: "REFRESH_TOKEN_NOT_FOUND", Message: "Refresh token not found"}
	ErrRefreshTokenExpired  = AuthError{Code: "REFRESH_TOKEN_EXPIRED", Message: "Refresh token has expired"}
	ErrInvalidRefreshToken  = AuthError{Code: "INVALID_REFRESH_TOKEN", Message: "Invalid or expired refresh token"}

	ErrPasswordResetTokenNotFound = AuthError{Code: "PASSWORD_RESET_TOKEN_NOT_FOUND", Message: "Password reset token not found"}
	ErrPasswordResetTokenExpired  = AuthError{Code: "PASSWORD_RESET_TOKEN_EXPIRED", Message: "Password reset token has expired"}
	ErrPasswordResetTokenUsed     = AuthError{Code: "PASSWORD_RESET_TOKEN_USED", Message: "Password reset token has already been used"}
	ErrSamePassword               = AuthError{Code: "SAME_PASSWORD", Message: "New password cannot be the same as current password"}
	ErrIncorrectCurrentPassword   = AuthError{Code: "INCORRECT_CURRENT_PASSWORD", Message: "Current password is incorrect"}

	ErrInvalidInput         = AuthError{Code: "INVALID_INPUT", Message: "Invalid input provided"}
	ErrMissingRequiredField = AuthError{Code: "MISSING_REQUIRED_FIELD", Message: "Required field is missing"}
	ErrUsernameTaken        = AuthError{Code: "USERNAME_TAKEN", Message: "Username is already taken"}

	ErrInsufficientPermissions = AuthError{Code: "INSUFFICIENT_PERMISSIONS", Message: "Insufficient permissions to perform this action"}
	ErrInvalidRole             = AuthError{Code: "INVALID_ROLE", Message: "Invalid role specified"}

	ErrProviderError        = AuthError{Code: "PROVIDER_ERROR", Message: "Authentication provider error"}
	ErrProviderNotSupported = AuthError{Code: "PROVIDER_NOT_SUPPORTED", Message: "Authentication provider not supported"}
	ErrAccountNotLinked     = AuthError{Code: "ACCOUNT_NOT_LINKED", Message: "Account is not linked to this provider"}

	ErrInternalServerError = AuthError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error"}
	ErrServiceUnavailable  = AuthError{Code: "SERVICE_UNAVAILABLE", Message: "Service is currently unavailable"}
	ErrMaintenanceMode     = AuthError{Code: "MAINTENANCE_MODE", Message: "Service is in maintenance mode"}

	ErrJWTSecretNotSet = AuthError{Code: "JWT_SECRET_NOT_SET", Message: "JWT secret is not set in the configuration"}
)
