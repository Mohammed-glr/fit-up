package interfaces

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
)

type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUserByUsername(ctx context.Context, username string) (*types.User, error)
	CreateUser(ctx context.Context, user *types.User) error
	UpdateUser(ctx context.Context, id string, updates *types.UpdateUserRequest) error
	UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error
	RefreshTokenStore

	CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error)
	GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error
}

type PasswordResetStore interface {
	CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error)
	GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error
}

type AuthService interface {
	ResetPassword(ctx context.Context, payload types.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	GetUser(ctx context.Context, userID string) (*types.User, error)
	GetUserByUsername(ctx context.Context, username string) (*types.User, error)
	Logout(ctx context.Context, userID string) error
	GenerateTokenPair(ctx context.Context, user *types.User) (*types.TokenPair, error)
	RotateTokens(ctx context.Context, refreshToken string) (*types.TokenPair, error)
}

type RefreshTokenStore interface {
	CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time, accessTokenJTI string) error
	GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	CleanupExpiredRefreshTokens(ctx context.Context) error
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID string) error
	UpdateRefreshTokenLastUsed(ctx context.Context, token string) error
}

type OAuthService interface {
	GetAuthorizationURL(ctx context.Context, provider, redirectURL string) (string, error)
	HandleCallback(ctx context.Context, provider, code, state string) (*types.OAuthUserInfo, error)
	LinkAccount(ctx context.Context, userID, provider string, userInfo *types.OAuthUserInfo) error
	UnlinkAccount(ctx context.Context, userID, provider string) error
	GetLinkedAccounts(ctx context.Context, userID string) ([]*types.Account, error)
}

type OAuthStore interface {
	CreateOAuthState(ctx context.Context, state *types.OAuthState) error
	GetOAuthState(ctx context.Context, state string) (*types.OAuthState, error)
	DeleteOAuthState(ctx context.Context, state string) error
	CleanupExpiredOAuthStates(ctx context.Context) error

	CreateAccount(ctx context.Context, account *types.Account) error
	GetAccountByProvider(ctx context.Context, provider, providerAccountID string) (*types.Account, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]*types.Account, error)
	DeleteAccount(ctx context.Context, userID, provider string) error
	UpdateAccountTokens(ctx context.Context, accountID, accessToken, refreshToken string, expiresAt int) error
}
