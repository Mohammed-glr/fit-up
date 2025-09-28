import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { useAuth } from '@/context/auth-context';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function OAuthCallback() {
    const { code, state, provider, error } = useLocalSearchParams<{
        code?: string;
        state?: string;
        provider?: string;
        error?: string;
    }>();
    const { login } = useAuth();
    const [isProcessing, setIsProcessing] = useState(true);
    const [message, setMessage] = useState('Processing OAuth login...');
    const { showError, showInfo, showWarning } = useToastMethods();

    useEffect(() => {
        handleOAuthCallback();
    }, []);

    const handleOAuthCallback = async () => {
        try {
            if (error) {
                throw new Error(`OAuth error: ${error}`);
            }

            if (!code || !state || !provider) {
                throw new Error('Missing OAuth parameters');
            }

            setMessage(`Authenticating with ${provider}...`);
            
            showInfo(`Processing ${provider} authentication...`, {
                position: 'top',
                duration: 3000,
            });
            
            // TODO: Implement OAuth callback handling
            // This would typically involve:
            // 1. Send code and state to your backend
            // 2. Backend exchanges code for tokens with OAuth provider
            // 3. Backend creates/finds user and returns JWT tokens
            // 4. Frontend stores tokens and redirects to app

            
            
            showWarning(
                `OAuth login with ${provider} is not yet implemented. Please use email/password login.`,
                {
                    position: 'center',
                    duration: 8000,
                    actionButton: {
                        text: 'Login',
                        onPress: () => router.replace('/(auth)/login')
                    }
                }
            );
            
            setTimeout(() => {
                router.replace('/(auth)/login');
            }, 3000);
        } catch (error: any) {
            console.error('OAuth callback error:', error);
            
            showError(
                error.message || 'Failed to complete OAuth login. Please try again.',
                {
                    position: 'center',
                    duration: 6000,
                    actionButton: {
                        text: 'Try Again',
                        onPress: () => router.replace('/(auth)/login')
                    }
                }
            );
            
            setTimeout(() => {
                router.replace('/(auth)/login');
            }, 3000);
        } finally {
            setIsProcessing(false);
        }
    };

    return (
        <View style={styles.container}>
            <ActivityIndicator size="large" color="#007AFF" style={styles.spinner} />
            <Text style={styles.message}>{message}</Text>
            <Text style={styles.subMessage}>
                {provider && `Completing ${provider} authentication...`}
            </Text>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#FFFFFF',
        padding: 24,
    },
    spinner: {
        marginBottom: 24,
    },
    message: {
        fontSize: 18,
        fontWeight: '600',
        textAlign: 'center',
        marginBottom: 8,
        color: '#333',
    },
    subMessage: {
        fontSize: 14,
        textAlign: 'center',
        color: '#666',
    },
});