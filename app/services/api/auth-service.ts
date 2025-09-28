import { httpClient } from "@/services/api/http-client";
import { AuthResponse, LoginRequest, RefreshTokenResponse, RegisterRequest, User } from "@/types/auth";


export const authService = {
    async login(credentials: LoginRequest): Promise<AuthResponse> {
        const response = await httpClient.post('/auth/login', credentials);
        const data = response.data;
        
        if (data.access_token) {
            await import('@/services/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
        
        return data;
    },

    async logout(): Promise<void> {
        await httpClient.post('/auth/logout');
    },

    async register(userData: RegisterRequest): Promise<AuthResponse> {
        const response = await httpClient.post('/auth/register', userData);
        const data = response.data;
        
        if (data.access_token) {
            await import('@/services/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
        
        return data;
    },

    async refreshToken(): Promise<RefreshTokenResponse> {
        const response = await httpClient.post('/auth/refresh-token');
        return response.data;
    },
    async validateToken(): Promise<{ user: User }> {
        const response = await httpClient.get('/auth/validate-token');
        return response.data;
    },

    async forgetPassword(email: string): Promise<void> {
        return httpClient.post('/auth/forget-password', { email });
    },

    async resetPassword(token: string, newPassword: string): Promise<void> {
        await httpClient.post('/auth/reset-password', {
            token,
            new_password: newPassword
        });
    },

    async changePassword(currentPassword: string, newPassword: string): Promise<void> {
        await httpClient.post('/auth/change-password', {
            current_password: currentPassword,
            new_password: newPassword
        });
    }
};