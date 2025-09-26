
type User = {
    identifier: string;
    password?: string;
    email?: string;
    name?: string;
    bio?: string;
    role?: string;
    image?: string;
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
    accessToken: string;
    refreshToken: string;
}

type RefreshTokenResponse = {
    accessToken: string;
    refreshToken: string;
}

export { User, LoginRequest, RegisterRequest, AuthResponse, RefreshTokenResponse };