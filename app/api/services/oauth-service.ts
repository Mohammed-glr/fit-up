import { API } from "../endpoints";
import { secureStorage } from "../storage/secure-storage";
import { executeAPI } from '../client';

const oauthService = {
    GetOAuthURL: async (provider: string): Promise<{ url: string }> => {
        const response = await executeAPI(API.auth.oauthLogin(provider));
        return response.data as { url: string };
    },

    HandleOAuthCallback: async (provider: string, code: string, redirectUri: string): Promise<void> => {
        const response = await executeAPI(API.auth.callbackOAuth(provider), { code, redirect_uri: redirectUri });
        const data = response.data;
        if (data.access_token) {
            await import('@/api/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
        return data as void;
    },

}