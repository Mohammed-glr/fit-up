import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Link, useLocalSearchParams, router } from 'expo-router';
import { FormContainer, Button } from '@/components/forms';
import PasswordInput from '@/components/auth/password-input';
import { authService } from '@/api/services/auth-service'; 
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function ResetPassword() {
    const { token } = useLocalSearchParams<{ token: string }>();
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [errors, setErrors] = useState<{ password?: string; confirmPassword?: string; general?: string }>({});
    const [success, setSuccess] = useState(false);
    const { showError, showSuccess, showWarning } = useToastMethods();

    useEffect(() => {
        if (!token) {
            showError(
                'This password reset link is invalid or has expired.',
                {
                    position: 'top',
                    duration: 6000,
                    actionButton: {
                        text: 'New Link',
                        onPress: () => router.replace('/(auth)/forgot-password')
                    }
                }
            );
        }
    }, [token, showError]);

    const validate = (): boolean => {
        const newErrors: typeof errors = {};

        if (!password) {
            newErrors.password = 'New password is required.';
        } else if (password.length < 8) {
            newErrors.password = 'Password must be at least 8 characters long.';
        }

        if (!confirmPassword) {
            newErrors.confirmPassword = 'Please confirm your password.';
        } else if (password !== confirmPassword) {
            newErrors.confirmPassword = 'Passwords do not match.';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async () => {
        if (!validate() || !token) {
            return;
        }

        try {
            setIsSubmitting(true);
            setErrors({});
            
            await authService.ResetPassword(token, password);
            setSuccess(true);
            
            showSuccess(
                'Your password has been reset successfully. You can now log in with your new password.',
                {
                    position: 'top',
                    duration: 5000,
                    actionButton: {
                        text: 'Login',
                        onPress: () => router.replace('/(auth)/login')
                    }
                }
            );
        } catch (error: any) {
            let errorMessage = 'Failed to reset password. Please try again.';
            
            if (error.response?.status === 400) {
                errorMessage = 'Invalid or expired reset token.';
            } else if (error.response?.status === 422) {
                errorMessage = 'Password does not meet requirements.';
            } else if (error.response?.data?.message) {
                errorMessage = error.response.data.message;
            }
            
            showError(errorMessage, {
                position: 'top',
                duration: error.response?.status === 400 ? 8000 : 5000,
                actionButton: error.response?.status === 400 ? {
                    text: 'New Link',
                    onPress: () => router.replace('/(auth)/forgot-password')
                } : undefined
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    if (!token) {
        return (
            <FormContainer>
                <View style={styles.errorContainer}>
                    <Text style={styles.title}>
                        Invalid Reset Link
                    </Text>
                    <Text style={styles.description}>
                        This password reset link is invalid or has expired.
                    </Text>
                    <Link href="/(auth)/forgot-password" style={styles.link}>
                        <Text style={styles.linkText}>Request New Reset Link</Text>
                    </Link>
                </View>
            </FormContainer>
        );
    }

    if (success) {
        return (
            <FormContainer>
                <View style={styles.successContainer}>
                    <Text style={styles.title}>
                        Password Reset Successful
                    </Text>
                    <Text style={styles.description}>
                        Your password has been reset successfully.
                    </Text>
                    <Link href="/(auth)/login" style={styles.link}>
                        <Text style={styles.linkText}>Go to Login</Text>
                    </Link>
                </View>
            </FormContainer>
        );
    }

    return (
        <FormContainer>
            <View style={styles.container}>
                <Text style={styles.title}>
                    Reset Your Password
                </Text>
                <Text style={styles.description}>
                    Enter your new password below.
                </Text>
                
                <PasswordInput
                    label="New Password"
                    value={password}
                    onChangeText={setPassword}
                    error={errors.password}
                    placeholder="Enter your new password"
                    disabled={isSubmitting}
                    showStrengthIndicator
                    style={styles.input}
                />
                
                <PasswordInput
                    label="Confirm New Password"
                    value={confirmPassword}
                    onChangeText={setConfirmPassword}
                    error={errors.confirmPassword}
                    placeholder="Confirm your new password"
                    disabled={isSubmitting}
                    style={styles.input}
                />
                
                <Button
                    title="Reset Password"
                    onPress={handleSubmit}
                    loading={isSubmitting}
                    disabled={!password || !confirmPassword || isSubmitting}
                    style={styles.button}
                />
                
                <Link href="/(auth)/login" style={styles.link}>
                    <Text style={styles.linkText}>Back to Login</Text>
                </Link>
            </View>
        </FormContainer>
    );
}

const styles = StyleSheet.create({
    container: {
        width: '100%',
    },
    errorContainer: {
        width: '100%',
        alignItems: 'center',
    },
    successContainer: {
        width: '100%',
        alignItems: 'center',
    },
    title: {
        fontSize: FONT_SIZES['3xl'],
        fontWeight: FONT_WEIGHTS.bold,
        textAlign: 'left',
        marginBottom: SPACING.base,
        color: COLORS.text.auth.primary,
    },
    description: {
        fontSize: FONT_SIZES.base,
        textAlign: 'left',
        marginBottom: SPACING.xl,
        color: COLORS.text.auth.secondary,
        lineHeight: 22,
    },
    input: {
        marginBottom: SPACING.base,
    },
    button: {
        marginBottom: SPACING.base,
    },
    link: {
        alignSelf: 'center',
    },
    linkText: {
        color: COLORS.primary,
        fontSize: FONT_SIZES.base,
        fontWeight: FONT_WEIGHTS.medium,
    },
});