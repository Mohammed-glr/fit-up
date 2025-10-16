import React from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import { router } from 'expo-router';
import { Button } from '@/components/forms';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';

interface EmailVerificationNoticeProps {
    email?: string;
    onResend: () => void;
    message?: string | null;
    error?: string | null;
    isResending?: boolean;
    disabled?: boolean;
}

const EmailVerificationNotice: React.FC<EmailVerificationNoticeProps> = ({
    email,
    onResend,
    message,
    error,
    isResending = false,
    disabled = false,
}) => {
    const handleOpenVerify = () => {
        router.push({ pathname: '/(auth)/verify-email' } as never);
    };

    return (
        <View style={styles.container}>
            <Text style={styles.title}>Verify your email</Text>
            <Text style={styles.description}>
                We sent a verification link to your email address. Please confirm your email to continue.
            </Text>
            {email ? (
                <Text style={styles.emailHint}>
                    Current email: <Text style={styles.emailValue}>{email}</Text>
                </Text>
            ) : (
                <Text style={styles.emailHint}>
                    Enter the email associated with your account and tap resend below.
                </Text>
            )}
            {message ? <Text style={styles.success}>{message}</Text> : null}
            {error ? <Text style={styles.error}>{error}</Text> : null}
            <Button
                title={isResending ? 'Sending…' : 'Resend verification email'}
                onPress={onResend}
                loading={isResending}
                disabled={disabled}
                variant="outline"
                style={styles.button}
            />
            <Pressable onPress={handleOpenVerify} style={styles.link}>
                <Text style={styles.linkText}>Already have a code? Verify it here.</Text>
            </Pressable>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        backgroundColor: COLORS.background.secondary,
        borderRadius: BORDER_RADIUS.lg,
        padding: SPACING.lg,
        marginTop: SPACING.lg,
        borderWidth: 1,
        borderColor: COLORS.border.light,
    },
    title: {
        fontSize: FONT_SIZES.lg,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.primary,
        marginBottom: SPACING.xs,
    },
    description: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.secondary,
        marginBottom: SPACING.sm,
        lineHeight: 20,
    },
    emailHint: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.secondary,
        marginBottom: SPACING.sm,
    },
    emailValue: {
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.primary,
    },
    success: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.success || '#1DB954',
        marginBottom: SPACING.sm,
    },
    error: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.error || '#FF3B30',
        marginBottom: SPACING.sm,
    },
    button: {
        marginTop: SPACING.sm,
    },
    link: {
        marginTop: SPACING.base,
    },
    linkText: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.primary,
        fontWeight: FONT_WEIGHTS.medium,
    },
});

export default EmailVerificationNotice;
