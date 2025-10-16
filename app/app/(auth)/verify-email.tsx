import { useEffect, useMemo, useState } from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { InputField, Button, FormContainer } from '@/components/forms';
import { useAuth } from '@/context/auth-context';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const VerifyEmailScreen = () => {
    const params = useLocalSearchParams<{ token?: string; email?: string }>();
    const {
        verifyEmail,
        resendVerification,
        verificationMessage,
        verificationError,
        isLoading,
        isEmailVerified,
    } = useAuth();

    const initialToken = useMemo(() => (typeof params.token === 'string' ? params.token : ''), [params.token]);
    const initialEmail = useMemo(() => (typeof params.email === 'string' ? params.email : ''), [params.email]);

    const [tokenInput, setTokenInput] = useState(initialToken);
    const [emailInput, setEmailInput] = useState(initialEmail);
    const [tokenError, setTokenError] = useState<string | undefined>();
    const [emailError, setEmailError] = useState<string | undefined>();
    const [autoAttempted, setAutoAttempted] = useState(false);
    const [isResending, setIsResending] = useState(false);

    useEffect(() => {
        if (initialToken && !autoAttempted) {
            verifyEmail(initialToken);
            setAutoAttempted(true);
        }
    }, [initialToken, autoAttempted, verifyEmail]);

    useEffect(() => {
        setTokenInput(initialToken);
    }, [initialToken]);

    useEffect(() => {
        setEmailInput(initialEmail);
    }, [initialEmail]);

    const handleVerify = async () => {
        if (!tokenInput.trim()) {
            setTokenError('Enter the verification token sent to your email.');
            return;
        }
        setTokenError(undefined);
        try {
            await verifyEmail(tokenInput.trim());
        } catch (error: any) {
            setTokenError(error?.response?.data?.message || 'Unable to verify email.');
        }
    };

    const handleResend = async () => {
        if (!emailInput.trim() || !emailPattern.test(emailInput.trim())) {
            setEmailError('Enter a valid email to resend the verification link.');
            return;
        }
        setEmailError(undefined);
        setIsResending(true);
        try {
            await resendVerification(emailInput.trim());
        } catch (error: any) {
            setEmailError(error?.response?.data?.message || 'Unable to resend verification email.');
        } finally {
            setIsResending(false);
        }
    };

    return (
        <FormContainer>
            <View style={styles.header}>
                <Text style={styles.title}>Verify your email</Text>
                <Text style={styles.subtitle}>
                    Confirm your email address to unlock your FitUp account. Paste the token from your email or open this page via the link we sent you.
                </Text>
            </View>

            <View style={styles.section}>
                <InputField
                    label="Verification token"
                    placeholder="Paste token"
                    value={tokenInput}
                    onChangeText={setTokenInput}
                    error={tokenError}
                    autoCapitalize="none"
                    disabled={isLoading}
                />
                <Button
                    title={isLoading ? 'Verifying…' : 'Verify email'}
                    onPress={handleVerify}
                    loading={isLoading}
                    disabled={isLoading}
                    style={styles.primaryButton}
                />
            </View>

            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Need a new email?</Text>
                <InputField
                    label="Account email"
                    placeholder="name@example.com"
                    value={emailInput}
                    onChangeText={setEmailInput}
                    error={emailError}
                    keyboardType="email-address"
                    autoCapitalize="none"
                    disabled={isLoading || isResending}
                />
                <Button
                    title={isResending ? 'Sending…' : 'Resend verification email'}
                    onPress={handleResend}
                    loading={isResending}
                    disabled={isLoading || isResending}
                    variant="outline"
                />
            </View>

            <View style={styles.feedback}>
                {isLoading ? <ActivityIndicator color={COLORS.primary} /> : null}
                {verificationMessage ? (
                    <Text style={styles.success}>{verificationMessage}</Text>
                ) : null}
                {verificationError ? <Text style={styles.error}>{verificationError}</Text> : null}
                {isEmailVerified ? (
                    <Button
                        title="Back to login"
                        onPress={() => router.replace('/(auth)/login')}
                        style={styles.backButton}
                    />
                ) : null}
            </View>
        </FormContainer>
    );
};

const styles = StyleSheet.create({
    header: {
        marginBottom: SPACING.lg,
    },
    title: {
        fontSize: FONT_SIZES['2xl'],
        fontWeight: FONT_WEIGHTS.bold,
        color: COLORS.text.primary,
        marginBottom: SPACING.xs,
    },
    subtitle: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.secondary,
        lineHeight: 20,
    },
    section: {
        marginBottom: SPACING.xl,
    },
    sectionTitle: {
        fontSize: FONT_SIZES.base,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.primary,
        marginBottom: SPACING.sm,
    },
    primaryButton: {
        marginTop: SPACING.sm,
    },
    feedback: {
        marginTop: SPACING.xl,
        gap: SPACING.sm,
    },
    success: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.success || '#1DB954',
    },
    error: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.error || '#FF3B30',
    },
    backButton: {
        marginTop: SPACING.base,
    },
});

export default VerifyEmailScreen;
