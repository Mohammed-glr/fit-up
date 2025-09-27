import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { 
    FormContainer,
    Button,
    InputField,
    ValidationMessage
} from '@/components/forms';
interface LoginFormData {
    identifier: string;
    password: string;
}

interface LoginFormError {
    identifier?: string;
    password?: string;
    general?: string;
}

export default function LoginForm() {
    const [formData, setFormData] = useState<LoginFormData>({
        identifier: "",
        password: "",
    });
    const [formError, setFormError] = useState<LoginFormError>({});
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const { login } = useAuth();

    const validate = (): boolean => {
        const errors: LoginFormError = {};

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        
        if (!formData.identifier) {
            errors.identifier = "Email or username is required.";
        } else if (emailRegex.test(formData.identifier)) {
            // It's an email format - email regex already validated it's valid
            // No additional validation needed
        } else {
            // It's username format, validate as username
            if (formData.identifier.length < 3) {
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
            await login({
                identifier: formData.identifier,
                password: formData.password,
            });
        } catch (error: any) {
            let errorMessage = "Login failed. Please try again.";

            // Handle Axios error response
            if (error.response?.status) {
                switch (error.response.status) {
                    case 401:
                        errorMessage = "Invalid email or password.";
                        break;
                    case 403:
                        errorMessage = "Account is disabled or not verified.";
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

            setFormError({ general: errorMessage });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <FormContainer>
            {formError.general && <ValidationMessage message={formError.general} />}
            <InputField
                label="Email"
                value={formData.identifier}
                onChangeText={(value) => handleChange("identifier", value)}
                error={formError.identifier}
                keyboardType="email-address"
                leftIcon="mail"
                autoCapitalize="none"
                placeholder="Enter your email"
                disabled={isSubmitting}
                style={{ marginBottom: 16 }}
            />
            <InputField
                label="Password"
                value={formData.password}
                onChangeText={(value) => handleChange("password", value)}
                error={formError.password}
                leftIcon="alert"
                isPassword
                placeholder="Enter your password"
                disabled={isSubmitting}
                style={{ marginBottom: 24 }}
                
            />
            <Button
                title="Log In"
                onPress={handleSubmit}
                loading={isSubmitting}
            />
        </FormContainer>    
    );
}