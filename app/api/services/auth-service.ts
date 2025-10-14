import { AuthResponse, LoginRequest, RefreshTokenResponse, RegisterRequest, User } from "@/types/auth";
import { API } from '../endpoints';
import { executeAPI } from '../client';


const authService = {
    Login: async (credentials: LoginRequest): Promise<AuthResponse> => {
        const response = await executeAPI(API.auth.login(), credentials);
        const data = response.data;

         if (data.access_token) {
            await import('@/api/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
        return data as AuthResponse;
    },

    Logout: async (): Promise<void> => {
        await executeAPI(API.auth.logout());
    },

    Register: async (userData: RegisterRequest): Promise<AuthResponse> => {
        const response = await executeAPI(API.auth.register(), userData);
        const data = response.data;
             if (data.access_token) {
            await import('@/api/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
        return data as AuthResponse;
    },

    RefreshToken: async (): Promise<RefreshTokenResponse> => {
        const response = await executeAPI(API.auth.refreshToken());
        return response.data as RefreshTokenResponse;
    },

    ValidateToken: async (): Promise<{ user: User }> => {
        const response = await executeAPI(API.auth.validateToken());
        return response.data as { user: User };
    },

    ForgetPassword: async (email: string): Promise<void> => {
        await executeAPI(API.auth.forgetPassword(), { email });
    },

    ResetPassword: async (token: string, newPassword: string): Promise<void> => {
        await executeAPI(API.auth.resetPassword(), {
            token,
            new_password: newPassword
        });
    },

    ChangePassword: async (currentPassword: string, newPassword: string): Promise<void> => {
        await executeAPI(API.auth.changePassword(), {
            current_password: currentPassword,
            new_password: newPassword
        });
    },

    UpdateRole: async (role: 'user' | 'coach'): Promise<{ message: string; user: User }> => {
        const response = await executeAPI(API.auth.updateRole(), { role });
        return response.data as { message: string; user: User };
    },

    UpdateProfile: async (data: { name?: string; bio?: string; image?: string }): Promise<{ message: string; user: User }> => {
        const formData = new FormData();
        
        if (data.name) {
            formData.append('name', data.name);
        }
        if (data.bio) {
            formData.append('bio', data.bio);
        }
        if (data.image) {
            const uriParts = data.image.split('.');
            const fileType = uriParts[uriParts.length - 1];
            
            formData.append('image', {
                uri: data.image,
                name: `profile.${fileType}`,
                type: `image/${fileType}`,
            } as any);
        }

        const response = await executeAPI(API.auth.updateProfile(), formData);
        return response.data as { message: string; user: User };
    },
}

export { authService };