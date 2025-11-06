

type UserRole = 'admin' | 'coach' | 'user' | 'client';


interface User {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image: string;
  role: UserRole;
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
  email_verified?: boolean;
}

interface UserResponse {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image?: string | null;
  role: UserRole;
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
}

interface PublicUserResponse {
  id: string;
  username: string;
  name: string;
  bio: string;
  image?: string | null;
  role: UserRole;
  created_at: string;
}

interface Account {
  id: string;
  user_id: string;
  type: string;
  provider: string;
  provider_account_id: string;
  refresh_token: string;
  access_token: string;
  expires_at: number;
  token_type: string;
  scope: string;
  id_token: string;
  session_state: string;
}

interface OAuthProvider {
  name: string;
  client_id: string;
  redirect_uri: string;
  auth_url: string;
  token_url: string;
  user_info_url: string;
  scopes: string[];
}

interface OAuthUserInfo {
  id: string;
  email: string;
  name: string;
  username?: string;
  avatar_url?: string;
  email_verified: boolean;
}

interface OAuthState {
  id: string;
  state: string;
  provider: string;
  redirect_url: string;
  expires_at: string;
  created_at: string;
}

interface TokenClaims {
  user_id: string;
  email: string;
  role: UserRole;
  jti: string;
  iss: string;
  sub: string;
  aud: string;
  exp: number;
  iat: number;
  nbf: number;
}

interface RefreshToken {
  id: string;
  user_id: string;
  access_token_jti?: string;
  expires_at: string;
  created_at: string;
  last_used_at: string;
  is_revoked: boolean;
  revoked_at?: string | null;
  user_agent?: string;
  ip_address?: string;
}

interface TokenPair {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
  token_type: string;
}

interface PasswordResetToken {
  id: string;
  email: string;
  token: string;
  expires: string;
}

interface LoginRequest {
  identifier: string;
  password: string;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  name?: string;
  role?: UserRole;
}

interface UpdateRoleRequest {
  role: UserRole;
}

interface UpdateRoleResponse {
  message: string;
  user: UserResponse;
}

interface OAuthAuthRequest {
  provider: 'google' | 'github' | 'facebook';
  redirect_url?: string;
}

interface OAuthCallbackRequest {
  code: string;
  state: string;
}

interface LinkAccountRequest {
  provider: string;
  code: string;
  state: string;
}

interface UpdateUserRequest {
  username?: string;
  name?: string;
  bio?: string;
  image?: string;
}

interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

interface ForgotPasswordRequest {
  email: string;
}

interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

interface RefreshTokenRequest {
  refresh_token: string;
}

interface RevokeTokenRequest {
  token: string;
  token_type?: string;
}

interface LoginResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  user: User;
  expires_at: number;
}

interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
  user: User;
}

interface TokenResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
}

interface RefreshTokenResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
}

interface TokenInfoResponse {
  active: boolean;
  claims?: TokenClaims;
  expires_in?: number;
  issued_at?: number;
  extra?: Record<string, any>;
}

interface AuthError {
  code: string;
  message: string;
}

export type {
  UserRole,
  User,
  UserResponse,
  PublicUserResponse,
  
  Account,
  OAuthProvider,
  OAuthUserInfo,
  OAuthState,
  
  TokenClaims,
  RefreshToken,
  TokenPair,
  PasswordResetToken,
  
  LoginRequest,
  RegisterRequest,
  UpdateRoleRequest,
  UpdateRoleResponse,
  OAuthAuthRequest,
  OAuthCallbackRequest,
  LinkAccountRequest,
  UpdateUserRequest,
  ChangePasswordRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  RefreshTokenRequest,
  RevokeTokenRequest,
  
  LoginResponse,
  AuthResponse,
  TokenResponse,
  RefreshTokenResponse,
  TokenInfoResponse,

  AuthError,
};