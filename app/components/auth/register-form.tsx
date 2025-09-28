import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { View, Text, StyleSheet, Keyboard } from 'react-native';
import { Link } from 'expo-router';
import { MotiView } from 'moti';
import { 
    FormContainer,
    Button,
    InputField,
    ValidationMessage
} from '@/components/forms';
import OAuthButtons from './oauth-buttons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
interface RegisterFormData {
    username: string;
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
}

interface RegisterFormError {
    username?: string;
    name?: string;
    email?: string;
    password?: string;
    confirmPassword?: string;
    general?: string;
}

export default function RegisterForm() {
    const [currentStep, setCurrentStep] = useState(1);
    const [prevStep, setPrevStep] = useState(1);
    const [formData, setFormData] = useState<RegisterFormData>({
        username: "",
        name: "",
        email: "",
        password: "",
        confirmPassword: "",
    });
    const [formError, setFormError] = useState<RegisterFormError>({});
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const { register } = useAuth();



    const validateCurrentStep = (): boolean => {
        const errors: RegisterFormError = {};

        switch (currentStep) {
            case 1:
                if (!formData.name.trim()) {
                    errors.name = "Name is required.";
                }
                if (!formData.username.trim()) {
                    errors.username = "Username is required.";
                }
                break;
            case 2:
                if (!formData.email.trim()) {
                    errors.email = "Email is required.";
                } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
                    errors.email = "Email is invalid.";
                }
                break;
            case 3:
                if (!formData.password) {
                    errors.password = "Password is required.";
                } else if (formData.password.length < 6) {
                    errors.password = "Password must be at least 6 characters.";
                }
                if (formData.confirmPassword !== formData.password) {
                    errors.confirmPassword = "Passwords do not match.";
                }
                break;
            default:
                break;
        }
        setFormError(errors);
        return Object.keys(errors).length === 0;
    }

    const handleNext = () => {
        if (validateCurrentStep()) {
            Keyboard.dismiss();
            setPrevStep(currentStep);
            setCurrentStep((prev) => prev + 1);
        }
    }

    const handleBack = () => {
        setCurrentStep((prev) => prev - 1);
        setFormError({});
    }

    const handleChange = (field: keyof RegisterFormData, value: string) => {
        setFormData({
            ...formData,
            [field]: value,
        });

        if (formError[field]) {
            setFormError(prev => ({
                ...prev,
                [field]: undefined
            }));
        }
    }

    const renderStepIndicator = () => {
        const totalSteps = 3;
        
        return (
            <MotiView 
                style={styles.stepIndicatorContainer}
                from={{ opacity: 0, translateY: -10 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 300 }}
            >
                <View style={styles.progressBarsContainer}>
                    {Array.from({ length: totalSteps }, (_, index) => {
                        const stepNumber = index + 1;
                        const isActive = stepNumber === currentStep;
                        const isCompleted = stepNumber < currentStep;
                        
                        return (
                            <MotiView 
                                key={stepNumber}
                                style={styles.progressBarStep}
                                animate={{
                                    backgroundColor: isCompleted 
                                        ? COLORS.success || COLORS.primary
                                        : isActive 
                                        ? COLORS.primary
                                        : COLORS.border.light,
                                    scaleY: isActive ? 1.2 : 1,
                                }}
                                transition={{
                                    type: 'timing',
                                    duration: 400,
                                    delay: index * 100,
                                }}
                            />
                        );
                    })}
                </View>
            </MotiView>
        );
    }

    const renderStepContent = () => {
        switch (currentStep) {
            case 1:
                return (
                    <View>
                        <MotiView
                            from={{ opacity: 0, translateY: -10 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 300 }}
                        >
                            <Text style={styles.stepTitle}>Personal Information</Text>
                            <Text style={styles.stepDescription}>Let's start with your basic details</Text>
                        </MotiView>
                        
                        <MotiView
                            from={{ opacity: 0, translateY: 20 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 400, delay: 200 }}
                        >
                            <InputField
                                label="Full Name"
                                value={formData.name}
                                onChangeText={(value) => handleChange('name', value)}
                                error={formError.name}
                                autoCapitalize="words"
                                placeholder="Enter your full name"
                                disabled={isSubmitting}
                            />
                            <InputField
                                label="Username"
                                value={formData.username}
                                onChangeText={(value) => handleChange('username', value)}   
                                error={formError.username}
                                autoCapitalize="none"
                                placeholder="Choose a unique username"
                                disabled={isSubmitting}
                            />
                        </MotiView>
                    </View>
                );
            
            case 2:
                return (
                    <View>
                        <MotiView
                            from={{ opacity: 0, translateY: -10 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 300 }}
                        >
                            <Text style={styles.stepTitle}>Email Address</Text>
                            <Text style={styles.stepDescription}>We'll use this to verify your account</Text>
                        </MotiView>
                        
                        <MotiView
                            from={{ opacity: 0, translateY: 20 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 400, delay: 200 }}
                        >
                            <InputField
                                label="Email"
                                value={formData.email}
                                onChangeText={(value) => handleChange('email', value)}
                                error={formError.email}
                                autoCapitalize="none"
                                placeholder="Enter your email address"
                                disabled={isSubmitting}
                                keyboardType="email-address"
                            />
                        </MotiView>
                    </View>
                );
            
            case 3:
                return (
                    <View>
                        <MotiView
                            from={{ opacity: 0, translateY: -10 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 300 }}
                        >
                            <Text style={styles.stepTitle}>Create Password</Text>
                            <Text style={styles.stepDescription}>Choose a strong password to secure your account</Text>
                        </MotiView>
                        
                        <MotiView
                            from={{ opacity: 0, translateY: 20 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 400, delay: 200 }}
                        >
                            <InputField
                                label="Password"
                                value={formData.password}
                                onChangeText={(value) => handleChange('password', value)}
                                error={formError.password}
                                isPassword
                                placeholder="Enter your password"
                                disabled={isSubmitting}
                                style={{ marginBottom: SPACING.base }}
                            />

                            <InputField
                                label="Confirm Password"
                                value={formData.confirmPassword}
                                onChangeText={(value) => handleChange('confirmPassword', value)}
                                error={formError.confirmPassword}
                                isPassword
                                placeholder="Confirm your password"
                                disabled={isSubmitting}
                            />
                        </MotiView>
                    </View>
                );
            
            default:
                return null;
        }
    }

    const renderNavigationButtons = () => {
        return (
            <MotiView 
                style={styles.navigationContainer}
                from={{ opacity: 0, translateY: 20 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 400, delay: 200 }}
            >
                {currentStep > 1 && (
                    <MotiView
                        from={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.9 }}
                        transition={{ type: 'spring', damping: 15, stiffness: 150 }}
                        style={styles.backButton}
                    >
                        <Button
                            title="Back"
                            onPress={handleBack}
                            disabled={isSubmitting}
                            variant="outline"
                        />
                    </MotiView>
                )}
                
                <MotiView
                    animate={{ 
                        scale: isSubmitting ? 0.95 : 1,
                        opacity: isSubmitting ? 0.7 : 1 
                    }}
                    transition={{ type: 'spring', damping: 15, stiffness: 150 }}
                    style={[styles.nextButton, currentStep === 1 && currentStep <= 1 && styles.fullWidthButton]}
                >
                    {currentStep < 3 ? (
                        <Button
                            title="Next"
                            onPress={handleNext}
                            disabled={isSubmitting}
                        />
                    ) : (
                        <Button
                            title="Create Account"
                            onPress={handleSubmit}
                            loading={isSubmitting}
                            disabled={isSubmitting}
                        />
                    )}
                </MotiView>
            </MotiView>
        );
    }

    const handleSubmit = async () => {
        if (!validateCurrentStep()) {
            return;
        }

        try {
            setIsSubmitting(true);
            await register({
                username: formData.username,
                name: formData.name,
                email: formData.email,
                password: formData.password,
                confirmPassword: formData.confirmPassword,
            });

        } catch (error: any) {
            let errorMessage = "Registration failed. Please try again.";

            if (error.response?.status) {
                switch (error.response.status) {
                    case 400:
                        errorMessage = "Invalid registration data. Please check your input.";
                        break;
                    case 409:
                        errorMessage = "Username or email already exists.";
                        break;
                    case 422:
                        errorMessage = error.response.data?.message || "Validation failed.";
                        break;
                    case 500:
                        errorMessage = "Server error. Please try again later.";
                        break;
                    default:
                        errorMessage = error.response.data?.message || "Registration failed. Please try again.";
                }
            } else if (error.request) {
                errorMessage = "Network error. Please check your connection.";
            } else if (error.message) {
                errorMessage = error.message;
            }

            setFormError({ general: errorMessage });
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <FormContainer>
            {renderStepIndicator()}
            
            {formError.general && <ValidationMessage message={formError.general} />}
            
            <MotiView 
                key={currentStep}
                style={styles.stepContent}
                from={{ opacity: 0, translateX: currentStep > prevStep ? 50 : -50 }}
                animate={{ opacity: 1, translateX: 0 }}
                exit={{ opacity: 0, translateX: currentStep > prevStep ? -50 : 50 }}
                transition={{ type: 'timing', duration: 300 }}
            >
                {renderStepContent()}
            </MotiView>

            {renderNavigationButtons()}

            {currentStep === 1 && (
                <MotiView
                    from={{ opacity: 0, translateY: 30 }}
                    animate={{ opacity: 1, translateY: 0 }}
                    exit={{ opacity: 0, translateY: 30 }}
                    transition={{ type: 'timing', duration: 400, delay: 300 }}
                >
                    <View style={styles.dividerContainer}>
                        <View style={styles.divider} />
                        <Text style={styles.dividerText}>Or continue with</Text>
                        <View style={styles.divider} />
                    </View>

                    <OAuthButtons disabled={isSubmitting} />
                </MotiView>
            )}

            <View style={styles.loginContainer}>
                <Text style={styles.loginText}>Already have an account? </Text>
                <Link href="/(auth)/login" style={styles.link}>
                    <Text style={styles.linkText}>Sign In</Text>
                </Link>
            </View>
        </FormContainer>
    )
}

const styles = StyleSheet.create({
    stepIndicatorContainer: {
        marginBottom: SPACING.xl,
        paddingHorizontal: SPACING.base,
    },
    progressBarsContainer: {
        flexDirection: 'row',
        gap: SPACING.sm,
        justifyContent: 'space-between',
    },
    progressBarStep: {
        flex: 1,
        height: 4,
        borderRadius: 2,
    },

    // Step Content Styles
    stepContent: {
        marginBottom: SPACING.xl,
    },
    stepTitle: {
        fontSize: FONT_SIZES.xl,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.primary,
        textAlign: 'center',
        marginBottom: SPACING.xs,
    },
    stepDescription: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.secondary,
        textAlign: 'center',
        marginBottom: SPACING.xl,
    },

    // Navigation Styles
    navigationContainer: {
        flexDirection: 'row',
        gap: SPACING.base,
        marginBottom: SPACING.xl,
    },
    backButton: {
        flex: 1,
    },
    nextButton: {
        flex: 1, 
    },
    fullWidthButton: {
        flex: 1,
    },

    // Original Styles
    registerButton: {
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
        backgroundColor: COLORS.border.light,
    },
    dividerText: {
        marginHorizontal: SPACING.base,
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.tertiary,
    },
    loginContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        marginTop: SPACING.xl,
    },
    loginText: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.tertiary,
    },
    link: {
        // Link styles handled by expo-router
    },
    linkText: {
        fontSize: FONT_SIZES.base,
        color: COLORS.primary,
        fontWeight: FONT_WEIGHTS.medium,
    },
});
