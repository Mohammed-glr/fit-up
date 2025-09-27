
interface User {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image: string;
  role: 'admin' | 'user';
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
}

type LoginRequest = {
    identifier: string;
    password: string;
}

type RegisterRequest = {
    name: string;
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
}

type AuthResponse = {
    user: User;
    access_token: string; 
    refresh_token: string;
    token_type: string;
    expires_at: number;
}

type RefreshTokenResponse = {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
}

export { User, LoginRequest, RegisterRequest, AuthResponse, RefreshTokenResponse };