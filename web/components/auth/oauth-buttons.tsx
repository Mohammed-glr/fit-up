import React from 'react';
import { View, StyleSheet, Alert } from 'react-native';
import { Button } from '@/components/forms';

interface OAuthButtonsProps {
    disabled?: boolean;
}

export default function OAuthButtons({ disabled = false }: OAuthButtonsProps) {
    const handleGoogleLogin = async () => {
        try {
            // TODO: Implement Google OAuth login
            Alert.alert('Coming Soon', 'Google login will be implemented soon!');
        } catch (error) {
            console.error('Google login failed:', error);
            Alert.alert('Error', 'Google login failed. Please try again.');
        }
    };

    const handleGitHubLogin = async () => {
        try {
            // TODO: Implement GitHub OAuth login
            Alert.alert('Coming Soon', 'GitHub login will be implemented soon!');
        } catch (error) {
            console.error('GitHub login failed:', error);
            Alert.alert('Error', 'GitHub login failed. Please try again.');
        }
    };

    const handleFacebookLogin = async () => {
        try {
            // TODO: Implement Facebook OAuth login
            Alert.alert('Coming Soon', 'Facebook login will be implemented soon!');
        } catch (error) {
            console.error('Facebook login failed:', error);
            Alert.alert('Error', 'Facebook login failed. Please try again.');
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
        gap: 12,
    },
    button: {
        width: '100%',
    },
});