import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Button } from '@/components/forms';
import { SPACING } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';

interface OAuthButtonsProps {
    disabled?: boolean;
}

export default function OAuthButtons({ disabled = false }: OAuthButtonsProps) {
    const { showInfo, showError } = useToastMethods();

    const handleGoogleLogin = async () => {
        try {
            // TODO: Implement Google OAuth login
            showInfo('Google login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('Google login failed:', error);
            showError('Google login failed. Please try again.');
        }
    };

    const handleGitHubLogin = async () => {
        try {
            // TODO: Implement GitHub OAuth login
            showInfo('GitHub login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('GitHub login failed:', error);
            showError('GitHub login failed. Please try again.');
        }
    };

    const handleFacebookLogin = async () => {
        try {
            // TODO: Implement Facebook OAuth login
            showInfo('Facebook login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('Facebook login failed:', error);
            showError('Facebook login failed. Please try again.');
        }
    };

    return (
        <View style={styles.container}>
            <Button
                title="Continue with Google"
                onPress={handleGoogleLogin}
                variant="outline"
                disabled={disabled}
                style={styles.button}
            />
            <Button
                title="Continue with GitHub"
                onPress={handleGitHubLogin}
                variant="outline"
                disabled={disabled}
                style={styles.button}
            />
            <Button
                title="Continue with Facebook"
                onPress={handleFacebookLogin}
                variant="outline"
                disabled={disabled}
                style={styles.button}
            />
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        width: '100%',
        gap: SPACING.md,
    },
    button: {
        width: '100%',
    },
});