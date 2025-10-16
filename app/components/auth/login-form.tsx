import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { View, Text, StyleSheet } from 'react-native';
import { Link, router } from 'expo-router';
import { MotiView } from 'moti';
import { 
    FormContainer,
    Button,
    InputField
} from '@/components/forms';
import OAuthButtons from './oauth-buttons';
import EmailVerificationNotice from './email-verification-notice';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';
interface LoginFormData {
    identifier: string;
    password: string;
}

interface LoginFormError {
    identifier?: string;
    password?: string;
    general?: string;
}

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export default function LoginForm() {
    const [formData, setFormData] = useState<LoginFormData>({
        identifier: "",
        password: "",
    });
    const [formError, setFormError] = useState<LoginFormError>({});
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const [showVerificationPrompt, setShowVerificationPrompt] = useState<boolean>(false);
    const [resendEmail, setResendEmail] = useState<string | null>(null);
    const [isResending, setIsResending] = useState<boolean>(false);
    const { 
        login,
        resendVerification,
        verificationMessage,
        verificationError,
    } = useAuth();
    const { showError, showSuccess } = useToastMethods();

    const validate = (): boolean => {
        const errors: LoginFormError = {};

        const identifier = formData.identifier.trim();
        
        if (!identifier) {
            errors.identifier = "Email or username is required.";
        } else if (!emailRegex.test(identifier)) {
            if (identifier.length < 3) {
                errors.identifier = "Username must be at least 3 characters long.";
            }
        }

        if (!formData.password) {
            errors.password = "Password is required.";
        } else if (formData.password.length < 6) {
            errors.password = "Password must be at least 6 characters long.";
        }

        setFormError(errors);
        return Object.keys(errors).length === 0;
    }

    const handleChange = (field: keyof LoginFormData, value: string) => {
        if (field === "identifier") {
            setShowVerificationPrompt(false);
            setResendEmail(null);
        }
        setFormData({
            ...formData,
            [field]: value,
        });
    }

    const handleSubmit = async () => {
        if (!validate()) {
            return;
        }

        try {
            setIsSubmitting(true);
            setShowVerificationPrompt(false);
            setResendEmail(null);
            const user = await login({
                identifier: formData.identifier,
                password: formData.password,
            });
            
            showSuccess('Welcome back! Successfully logged in.', {
                position: 'top',
                duration: 2000,
            });
            
            setTimeout(() => {
                if (user.role === 'coach') {
                    router.replace('/(coach)');
                } else {
                    router.replace('/(user)');
                }
            }, 500);
        } catch (error: any) {
            let errorMessage = "Login failed. Please try again.";

            if (error.response?.status) {
                switch (error.response.status) {
                    case 401:
                        errorMessage = "Invalid email or password.";
                        break;
                    case 403:
                        errorMessage = "Your email is not verified yet. Please verify to continue.";
                        setShowVerificationPrompt(true);
                        if (emailRegex.test(formData.identifier.trim())) {
                            setResendEmail(formData.identifier.trim());
                        } else {
                            setResendEmail(null);
                            setFormError(prev => ({
                                ...prev,
                                identifier: 'Enter your email address so we can resend verification.',
                            }));
                        }
                        break;
                    case 429:
                        errorMessage = "Too many login attempts. Please try again later.";
                        break;
                    case 500:
                        errorMessage = "Server error. Please try again later.";
                        break;
                    default:
                        errorMessage = error.response.data?.message || "Login failed. Please try again.";
                }
            } else if (error.request) {
                errorMessage = "Network error. Please check your connection.";
            } else if (error.message) {
                errorMessage = error.message;
            }

            showError(errorMessage, {
                position: 'top',
                duration: 5000,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleResendVerification = async () => {
        const targetEmail = (resendEmail || formData.identifier).trim();
        if (!targetEmail || !emailRegex.test(targetEmail)) {
            setFormError(prev => ({
                ...prev,
                identifier: 'Enter a valid email address to resend verification.',
            }));
            return;
        }

        setFormError(prev => ({
            ...prev,
            identifier: undefined,
        }));

        setIsResending(true);
        try {
            const message = await resendVerification(targetEmail);
            showSuccess(message || 'Verification email sent. Please check your inbox.', {
                position: 'top',
                duration: 3000,
            });
        } catch (error: any) {
            const errorMessage = error?.response?.data?.message || 'Unable to resend verification email.';
            showError(errorMessage, {
                position: 'top',
                duration: 5000,
            });
        } finally {
            setIsResending(false);
        }
    };

    return (
        <FormContainer>
            <MotiView
                from={{ opacity: 0, translateY: -10 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 300 }}
            >
                <Text style={styles.title}>Welcome <br />Back</Text>
                <Text style={styles.subtitle}>Sign in to continue to FitUp</Text>
            </MotiView>

            <MotiView
                from={{ opacity: 0, translateY: 20 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 400, delay: 200 }}
            >
                <InputField
                    label="Email"
                    value={formData.identifier}
                    onChangeText={(value) => handleChange("identifier", value)}
                    error={formError.identifier}
                    keyboardType="email-address"
                    autoCapitalize="none"
                    placeholder="Enter your email"
                    disabled={isSubmitting}
                    style={{ marginBottom: SPACING.base }}
                />
                <InputField
                    label="Password"
                    value={formData.password}
                    onChangeText={(value) => handleChange("password", value)}
                    error={formError.password}
                    isPassword
                    placeholder="Enter your password"
                    disabled={isSubmitting}
                    style={{ marginBottom: SPACING.xl }}
                />
                {(showVerificationPrompt || verificationMessage || verificationError) && (
                    <EmailVerificationNotice
                        email={(() => {
                            const candidate = (resendEmail || formData.identifier).trim();
                            return emailRegex.test(candidate) ? candidate : undefined;
                        })()}
                        onResend={handleResendVerification}
                        message={verificationMessage}
                        error={verificationError}
                        isResending={isResending}
                        disabled={isSubmitting || isResending}
                    />
                )}
            </MotiView>

            <MotiView
                from={{ opacity: 0, translateY: 20 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 400, delay: 400 }}
            >
                <View style={styles.forgotPasswordContainer}>
                    <Link href="/(auth)/forgot-password" style={styles.link}>
                        <Text style={styles.linkText}>Forgot Password?</Text>
                    </Link>
                </View>

                <Button
                    title="Log In"
                    onPress={handleSubmit}
                    loading={isSubmitting}
                    disabled={isSubmitting}
                    style={styles.loginButton}
                />
            </MotiView>

            <MotiView
                from={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ type: 'spring', damping: 15, stiffness: 150, delay: 600 }}
            >
                <View style={styles.dividerContainer}>
                    <View style={styles.divider} />
                    <Text style={styles.dividerText}>Or continue with</Text>
                    <View style={styles.divider} />
                </View>

                <OAuthButtons disabled={isSubmitting} />
            </MotiView>

            <MotiView
                from={{ opacity: 0, translateY: 20 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 400, delay: 800 }}
            >
                <View style={styles.signupContainer}>
                    <Text style={styles.signupText}>Don't have an account? </Text>
                    <Link href="/(auth)/register" style={styles.link}>
                        <Text style={styles.linkText}>Sign Up</Text>
                    </Link>
                </View>
            </MotiView>
        </FormContainer>    
    );
}

const styles = StyleSheet.create({
    title: {
        fontSize: FONT_SIZES['3xl'],
        fontWeight: FONT_WEIGHTS.bold,
        textAlign: 'left',
        marginBottom: SPACING.xs,
        color: COLORS.text.auth.primary,
    },
    subtitle: {
        fontSize: FONT_SIZES.base,
        textAlign: 'left',
        marginBottom: SPACING.xl,
        color: COLORS.text.auth.secondary,
    },
    forgotPasswordContainer: {
        alignItems: 'flex-start',
        marginBottom: SPACING.xl,
    },
    loginButton: {
        marginBottom: SPACING.xl,
    },
    dividerContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        marginVertical: SPACING.xl,
    },
    divider: {
        flex: 1,
        height: 1,
        backgroundColor: COLORS.text.auth.tertiary,
        opacity: 0.3,
    },
    dividerText: {
        marginHorizontal: SPACING.base,
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.auth.tertiary,
    },
    signupContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        marginTop: SPACING.xl,
    },
    signupText: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.auth.tertiary,
    },
    link: {
    },
    linkText: {
        fontSize: FONT_SIZES.base,
        color: COLORS.primary,
        fontWeight: FONT_WEIGHTS.medium,
    },
});