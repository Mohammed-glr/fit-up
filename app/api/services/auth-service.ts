import { AuthResponse, LoginRequest, RefreshTokenResponse, RegisterRequest, User } from "@/types/auth";
import { API } from '../endpoints';
import { executeAPI } from '../client';


const authService = {
        VerifyEmail: async (token: string): Promise<{ message: string }> => {
            const response = await executeAPI(API.auth.verifyEmail(), { token });
            return response.data as { message: string };
        },

        ResendVerification: async (email: string): Promise<{ message: string }> => {
            const response = await executeAPI(API.auth.resendVerification(), { email });
            return response.data as { message: string };
        },
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
        let base64Image: string | undefined;
        
        if (data.image) {
            try {
                const response = await fetch(data.image);
                const blob = await response.blob();
                
                const base64 = await new Promise<string>((resolve, reject) => {
                    const reader = new FileReader();
                    reader.onloadend = () => {
                        const base64data = reader.result as string;
                        resolve(base64data);
                    };
                    reader.onerror = reject;
                    reader.readAsDataURL(blob);
                });
                
                base64Image = base64;
            } catch (error) {
                console.error('Failed to convert image to base64:', error);
                throw new Error('Failed to process image');
            }
        }


        const payload: any = {};
        if (data.name) payload.name = data.name;
        if (data.bio) payload.bio = data.bio;
        if (base64Image) payload.image = base64Image;

        const response = await executeAPI(API.auth.updateProfile(), payload);
        return response.data as { message: string; user: User };
    },
}

export { authService };