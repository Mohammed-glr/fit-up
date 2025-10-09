import { httpClient } from "@/api/client";

export const oauthService = {
    async getOAuthProviders(): Promise<string[]> {
        const response = await httpClient.get('/auth/oauth/providers');
        return response.data.providers;
    },

    async getOAuthUrl(provider: string, redirectUri: string): Promise<string> {
        const response = await httpClient.get('/auth/oauth/url', {
            params: {
                provider,
                redirect_uri: redirectUri,
            }
        });
        return response.data.url;
    },

    async handleOAuthCallback(provider: string, code: string, redirectUri: string): Promise<void> {
        const response = await httpClient.post('/auth/oauth/callback', {
            provider,
            code,
            redirect_uri: redirectUri,
        });
        const data = response.data;

        if (data.access_token) {
            await import('@/api/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
    },

    async revokeOAuthToken(): Promise<void> {
        await httpClient.post('/auth/oauth/revoke');
    },

    async getUserInfo(): Promise<any> {
        const response = await httpClient.get('/auth/oauth/userinfo');
        return response.data;
    },

    async linkOAuthAccount(provider: string, code: string, redirectUri: string): Promise<void> {
        await httpClient.post('/auth/oauth/link', {
            provider,
            code,
            redirect_uri: redirectUri,
        });
    },

    async unlinkOAuthAccount(provider: string): Promise<void> {
        await httpClient.post('/auth/oauth/unlink', { provider });
    },

    async getLinkedAccounts(): Promise<string[]> {
        const response = await httpClient.get('/auth/oauth/linked-accounts');
        return response.data.linked_accounts;
    },

    async refreshOAuthToken(provider: string): Promise<void> {
        await httpClient.post('/auth/oauth/refresh-token', { provider });
    },

    async getOAuthTokenStatus(provider: string): Promise<{ expires_in: number; is_valid: boolean }> {
        const response = await httpClient.get('/auth/oauth/token-status', {
            params: { provider }
        });
        return response.data;
    },

    async getAvailableProviders(): Promise<string[]> {
        const response = await httpClient.get('/auth/oauth/available-providers');
        return response.data.providers;
    },

    async authenticateWithOAuthToken(provider: string, oauthToken: string): Promise<void> {
        const response = await httpClient.post('/auth/oauth/authenticate', {
            provider,
            oauth_token: oauthToken,
        });
        const data = response.data;
        
        if (data.access_token) {
            await import('@/api/storage/secure-storage').then(({ secureStorage }) => {
                secureStorage.setToken('access_token', data.access_token);
                if (data.refresh_token) {
                    secureStorage.setToken('refresh_token', data.refresh_token);
                }
            });
        }
    },

    async getOAuthToken(provider: string): Promise<string> {
        const response = await httpClient.get('/auth/oauth/get-token', {
            params: { provider }
        });
        return response.data.oauth_token;
    },

    async disconnectOAuthProvider(provider: string): Promise<void> {
        await httpClient.post('/auth/oauth/disconnect', { provider });
    }
}