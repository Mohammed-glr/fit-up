import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { View, Text, StyleSheet } from 'react-native';
import { Link } from 'expo-router';
import { 
    FormContainer,
    Button,
    InputField,
    ValidationMessage
} from '@/components/forms';
import OAuthButtons from './oauth-buttons';
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

    const validate = (): boolean => {
        const errors: RegisterFormError = {};
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (!formData.username) {
            errors.username = "Username is required.";
        } else if (formData.username.length < 3) {
            errors.username = "Username must be at least 3 characters long.";
        }
        if (!formData.name) {
            errors.name = "Name is required.";
        } else if (formData.name.length < 3) {
            errors.name = "Name must be at least 3 characters long.";
        }

        if (!formData.email) {
            errors.email = "Email is required.";
        } else if (!emailRegex.test(formData.email)) {
            errors.email = "Please enter a valid email address.";
        }

        if (!formData.password) {
            errors.password = "Password is required.";
        }
        else if (formData.password.length < 6) {
            errors.password = "Password must be at least 6 characters long.";
        }

        if (!formData.confirmPassword) {
            errors.confirmPassword = "Please confirm your password.";
        } else if (formData.confirmPassword !== formData.password) {
            errors.confirmPassword = "Passwords do not match.";
        }
        setFormError(errors);
        return Object.keys(errors).length === 0;
    }

    const handleChange = (field: keyof RegisterFormData, value: string) => {
        setFormData({
            ...formData,
            [field]: value,
        });

        validate();
    }

    const handleSubmit = async () => {
        if (!validate()) {
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

    return(
        <FormContainer>
            {formError.general && <ValidationMessage message={formError.general} />}
            <InputField
                label="Username"
                value={formData.username}
                onChangeText={(value) => handleChange('username', value)}   
                error={formError.username}
                autoCapitalize="none"
                placeholder="Enter your username"
                disabled={isSubmitting}
                leftIcon="person"
            />
            <InputField
                label="Name"
                value={formData.name}
                onChangeText={(value) => handleChange('name', value)}
                error={formError.name}
                autoCapitalize="words"
                placeholder="Enter your full name"
                disabled={isSubmitting}
                leftIcon="person"
            />
            <InputField
                label="Email"
                value={formData.email}
                onChangeText={(value) => handleChange('email', value)}
                error={formError.email}
                autoCapitalize="none"
                placeholder="Enter your email"
                disabled={isSubmitting}
                leftIcon="mail"
                keyboardType="email-address"
                style={{ marginBottom: 16 }}
            />

            <InputField
                label="Password"
                value={formData.password}
                onChangeText={(value) => handleChange('password', value)}
                error={formError.password}
                isPassword
                leftIcon="lock-closed"
                placeholder="Enter your password"
                disabled={isSubmitting}
                style={{ marginBottom: 16 }}
            />

            <InputField
                label="Confirm Password"
                value={formData.confirmPassword}
                onChangeText={(value) => handleChange('confirmPassword', value)}
                error={formError.confirmPassword}
                isPassword
                leftIcon="lock-closed"
                placeholder="Confirm your password"
                disabled={isSubmitting}
                style={{ marginBottom: 24 }}
            />

            <Button
                title="Create Account"
                onPress={handleSubmit}
                loading={isSubmitting}
                disabled={isSubmitting}
                style={styles.registerButton}
            />

            <View style={styles.dividerContainer}>
                <View style={styles.divider} />
                <Text style={styles.dividerText}>Or continue with</Text>
                <View style={styles.divider} />
            </View>

            <OAuthButtons disabled={isSubmitting} />

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
    registerButton: {
        marginBottom: 24,
    },
    dividerContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        marginVertical: 24,
    },
    divider: {
        flex: 1,
        height: 1,
        backgroundColor: '#E1E5E9',
    },
    dividerText: {
        marginHorizontal: 16,
        fontSize: 14,
        color: '#666',
    },
    loginContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        marginTop: 24,
    },
    loginText: {
        fontSize: 16,
        color: '#666',
    },
    link: {
        // Link styles handled by expo-router
    },
    linkText: {
        fontSize: 16,
        color: '#007AFF',
        fontWeight: '500',
    },
});
