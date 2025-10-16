import React, { useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Link } from 'expo-router';
import { FormContainer, Button, InputField } from '@/components/forms';
import { authService } from '@/api/services/auth-service';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function ForgotPassword() {
    const [email, setEmail] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const { showError, showSuccess, showInfo } = useToastMethods();

    const validateEmail = (email: string): boolean => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    const handleSubmit = async () => {
        setError('');
        
        if (!email) {
            showError('Email is required.');
            return;
        }
        
        if (!validateEmail(email)) {
            showError('Please enter a valid email address.');
            return;
        }

        try {
            setIsSubmitting(true);
            await authService.ForgetPassword(email);
            setSuccess(true);
            showSuccess(
                `Password reset instructions have been sent to ${email}`, 
                {
                    position: 'top',
                    duration: 5000,
                }
            );
        } catch (error: any) {
            let errorMessage = 'Failed to send reset email. Please try again.';
            
            if (error.response?.status === 404) {
                errorMessage = 'No account found with this email address.';
            } else if (error.response?.status === 429) {
                errorMessage = 'Too many requests. Please try again later.';
            } else if (error.response?.data?.message) {
                errorMessage = error.response.data.message;
            }
            
            showError(errorMessage, {
                position: 'top',
                duration: error.response?.status === 429 ? 8000 : 5000,
                actionButton: error.response?.status === 404 ? {
                    text: 'Sign Up',
                    onPress: () => {
                    }
                } : undefined
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    if (success) {
        return (
            <FormContainer>
                <View style={styles.successContainer}>
                    <Text style={styles.title}>Check <br />Your Email</Text>
                    <Text style={styles.description}>
                        We've sent password reset instructions to {email}
                    </Text>
                    <Text style={styles.note}>
                        Didn't receive the email? Check your spam folder or try again.
                    </Text>
                    <Button
                        title="Send Again"
                        onPress={() => {
                            setSuccess(false);
                            showInfo('Sending reset instructions again...', {
                                position: 'top',
                                duration: 2000,
                            });
                            handleSubmit();
                        }}
                        variant="outline"
                        style={styles.button}
                    />
                    <Link href="/(auth)/login" style={styles.link}>
                        <Text style={styles.linkText}>Back to Login</Text>
                    </Link>
                </View>
            </FormContainer>
        );
    }

    return (
        <FormContainer>
            <View style={styles.container}>
                <Text style={styles.title}>Forgot <br />Password?</Text>
                <Text style={styles.description}>
                    Enter your email address and we'll send you instructions to reset your password.
                </Text>
                
                <InputField
                    label="Email"
                    value={email}
                    onChangeText={setEmail}
                    keyboardType="email-address"
                    placeholder="Enter your email"
                    disabled={isSubmitting}
                    autoCapitalize="none"
                    style={styles.input}
                />
                
                <Button
                    title="Send Reset Instructions"
                    onPress={handleSubmit}
                    loading={isSubmitting}
                    disabled={!email || isSubmitting}
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
    note: {
        fontSize: FONT_SIZES.sm,
        textAlign: 'center',
        marginBottom: SPACING.xl,
        color: COLORS.text.auth.tertiary,
        lineHeight: 20,
    },
    input: {
        marginBottom: SPACING.xl,
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