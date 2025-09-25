package pkg

// TODO: Step 1 - Create authentication client for other services:
//   - AuthClient interface (VerifyToken, GetUserInfo, RefreshToken)
//   - HTTPAuthClient implementation for inter-service communication
//   - MockAuthClient for testing
// TODO: Step 2 - Implement OAuth provider clients:
//   - OAuthProvider interface (GetAuthURL, ExchangeCode, GetUserInfo)
//   - GoogleOAuthClient implementation
//   - GitHubOAuthClient implementation
//   - FacebookOAuthClient implementation
// TODO: Step 3 - Create token management utilities:
//   - TokenManager interface (Generate, Validate, Refresh, Revoke)
//   - JWTTokenManager implementation
//   - RedisTokenManager for token blacklisting
// TODO: Step 4 - Implement security middleware for reuse:
//   - RequireAuth middleware for protecting endpoints
//   - RequireRole middleware for role-based access
//   - RateLimitMiddleware for authentication endpoints
// TODO: Step 5 - Create validation utilities:
//   - InputValidator for auth requests
//   - PasswordValidator with configurable rules
//   - EmailValidator with domain checking
// TODO: Step 6 - Add monitoring and metrics:
//   - AuthMetrics for tracking login success/failure rates
//   - SecurityEventLogger for audit trails

// Flow: External services -> pkg interfaces -> auth service business logic
// Exports: Client interfaces, middleware, validators, metrics collectors
