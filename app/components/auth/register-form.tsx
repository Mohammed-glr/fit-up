import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { View, Text, StyleSheet, Keyboard, Pressable } from 'react-native';
import { Link, router } from 'expo-router';
import { MotiView } from 'moti';
import { 
    FormContainer,
    Button,
    InputField
} from '@/components/forms';
import OAuthButtons from './oauth-buttons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';
import PasswordInput from "./password-input";
interface RegisterFormData {
    username: string;
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
    role: 'client' | 'coach';
}

interface RegisterFormError {
    username?: string;
    name?: string;
    email?: string;
    password?: string;
    confirmPassword?: string;
    role?: string;
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
        role: "client",
    });
    const [formError, setFormError] = useState<RegisterFormError>({});
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const { register } = useAuth();
    const { showError, showSuccess } = useToastMethods();



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
        const totalSteps = 4;
        
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
                            <Text style={styles.stepTitle}>
                                Personal Information
                            </Text>
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
                            <Text style={styles.stepTitle}>
                                Email Address
                            </Text>
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
                            <Text style={styles.stepTitle}>
                                Create Password
                            </Text>
                            <Text style={styles.stepDescription}>Choose a strong password to secure your account</Text>
                        </MotiView>
                        
                        <MotiView
                            from={{ opacity: 0, translateY: 20 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 400, delay: 200 }}
                        >
                            <PasswordInput
                                label="Password"
                                value={formData.password}
                                onChangeText={(value) => handleChange('password', value)}
                                error={formError.password}
                                placeholder="Enter your password"
                                disabled={isSubmitting}
                                style={{ marginBottom: SPACING.base }}
                                showStrengthIndicator={true}
                                showRequirements={true}
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
            
            case 4:
                return (
                    <View>
                        <MotiView
                            from={{ opacity: 0, translateY: -10 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 300 }}
                        >
                            <Text style={styles.stepTitle}>
                                Choose Your Role
                            </Text>
                            <Text style={styles.stepDescription}>Select how you'll be using FitUp</Text>
                        </MotiView>
                        
                        <MotiView
                            from={{ opacity: 0, translateY: 20 }}
                            animate={{ opacity: 1, translateY: 0 }}
                            transition={{ type: 'timing', duration: 400, delay: 200 }}
                            style={styles.roleContainer}
                        >
                            <RoleOption
                                title="Client"
                                description="Track workouts, nutrition, and fitness progress"
                                icon="🏃"
                                isSelected={formData.role === 'client'}
                                onSelect={() => handleChange('role', 'client')}
                            />
                            <RoleOption
                                title="Coach"
                                description="Create programs and manage multiple clients"
                                icon="🏋️"
                                isSelected={formData.role === 'coach'}
                                onSelect={() => handleChange('role', 'coach')}
                            />
                        </MotiView>
                    </View>
                );
            
            default:
                return null;
        }
    }
    
    const RoleOption = ({ title, description, icon, isSelected, onSelect }: {
        title: string;
        description: string;
        icon: string;
        isSelected: boolean;
        onSelect: () => void;
    }) => (
        <MotiView
            animate={{
                backgroundColor: isSelected ? COLORS.primarySoft : COLORS.background.secondary,
                borderColor: isSelected ? COLORS.primary : COLORS.border.light,
                scale: isSelected ? 1.02 : 1,
            }}
            transition={{ type: 'spring', damping: 15, stiffness: 150 }}
            style={[styles.roleOption, isSelected && styles.roleOptionSelected]}
        >
            <Pressable onPress={onSelect} style={styles.roleOptionContent}>
                <Text style={styles.roleIcon}>{icon}</Text>
                <View style={styles.roleTextContainer}>
                    <Text style={[styles.roleTitle, isSelected && styles.roleTextSelected]}>{title}</Text>
                    <Text style={[styles.roleDescription, isSelected && styles.roleDescriptionSelected]}>{description}</Text>
                </View>
            </Pressable>
        </MotiView>
    );

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
                    {currentStep < 4 ? (
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
                role: formData.role,
            });
            
            showSuccess('Account created! Check your email to verify before logging in.', {
                position: 'top',
                duration: 3000,
            });
            
            router.replace('/(auth)/login');

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

            showError(errorMessage, {
                position: 'top',
                duration: 5000,
                actionButton: error.response?.status === 409 ? {
                    text: 'Try Login',
                    onPress: () => {
                    }
                } : undefined
            });
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <FormContainer>
            {renderStepIndicator()}
            
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
                    {/* <View style={styles.dividerContainer}>
                        <View style={styles.divider} />
                        <Text style={styles.dividerText}>Or continue with</Text>
                        <View style={styles.divider} />
                    </View>

                    <OAuthButtons disabled={isSubmitting} /> */}
                </MotiView>
            )}

            <MotiView
                from={{ opacity: 0, translateY: 20 }}
                animate={{ opacity: 1, translateY: 0 }}
                transition={{ type: 'timing', duration: 400, delay: 600 }}

            >
                <View style={styles.loginContainer}>
                    <Text style={styles.loginText}>Already have an account? </Text>
                    <Link href="/(auth)/login" style={styles.link}>
                        <Text style={styles.linkText}>Sign In</Text>
                    </Link>
                </View>
            </MotiView>
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

    stepContent: {
        marginBottom: SPACING.xl,
    },
    stepTitle: {
        fontSize: FONT_SIZES["3xl"],
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.auth.primary,
        textAlign: 'left',
        marginBottom: SPACING.xs,
    },
    stepDescription: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.auth.secondary,
        textAlign: 'left',
        marginBottom: SPACING.xl,
    },

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
        backgroundColor: COLORS.text.auth.tertiary,
        opacity: 0.3,
    },
    dividerText: {
        marginHorizontal: SPACING.base,
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.auth.tertiary,
    },
    loginContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        marginTop: SPACING.xl,
    },
    loginText: {
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
    
    roleContainer: {
        gap: SPACING.base,
    },
    roleOption: {
        borderRadius: 32,
        borderWidth: 0.5,
        padding: SPACING.base,
        marginBottom: SPACING.base,
        borderColor: COLORS.border.accent,
    },
    roleOptionSelected: {
        borderWidth: 2,
    },
    roleOptionContent: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: SPACING.base,
    },
    roleIcon: {
        fontSize: 32,
    },
    roleTextContainer: {
        flex: 1,
    },
    roleTitle: {
        fontSize: FONT_SIZES.lg,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.auth.placeholder,
        marginBottom: 4,
    },
    roleTextSelected: {
        color: COLORS.primaryDark,
    },
    roleDescription: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.auth.tertiary,
        lineHeight: 18,
    },
    roleDescriptionSelected: {
        color: COLORS.primary,
    },
});
