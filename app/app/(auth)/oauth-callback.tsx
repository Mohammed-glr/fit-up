import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ActivityIndicator, Alert } from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { useAuth } from '@/context/auth-context';

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
            
            // TODO: Implement OAuth callback handling
            // This would typically involve:
            // 1. Send code and state to your backend
            // 2. Backend exchanges code for tokens with OAuth provider
            // 3. Backend creates/finds user and returns JWT tokens
            // 4. Frontend stores tokens and redirects to app
            
            // For now, show a message that OAuth is not yet implemented
            Alert.alert(
                'OAuth Not Implemented',
                `OAuth login with ${provider} is not yet implemented. Please use email/password login.`,
                [
                    {
                        text: 'OK',
                        onPress: () => router.replace('/(auth)/login')
                    }
                ]
            );
        } catch (error: any) {
            console.error('OAuth callback error:', error);
            
            Alert.alert(
                'Authentication Failed',
                error.message || 'Failed to complete OAuth login. Please try again.',
                [
                    {
                        text: 'OK',
                        onPress: () => router.replace('/(auth)/login')
                    }
                ]
            );
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