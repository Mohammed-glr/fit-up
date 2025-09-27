import React, { useState } from 'react';
import { Text, StyleSheet } from 'react-native';
import { InputField } from '@/components/forms';

interface PasswordInputProps {
    label?: string;
    value: string;
    onChangeText: (text: string) => void;
    error?: string;
    placeholder?: string;
    disabled?: boolean;
    style?: any;
    showStrengthIndicator?: boolean;
}

export default function PasswordInput({
    label = "Password",
    value,
    onChangeText,
    error,
    placeholder = "Enter your password",
    disabled = false,
    style,
    showStrengthIndicator = false,
}: PasswordInputProps) {
    const [strength, setStrength] = useState(0);

    const calculatePasswordStrength = (password: string): number => {
        let score = 0;
        if (password.length >= 8) score += 1;
        if (/[a-z]/.test(password)) score += 1;
        if (/[A-Z]/.test(password)) score += 1;
        if (/[0-9]/.test(password)) score += 1;
        if (/[^A-Za-z0-9]/.test(password)) score += 1;
        return score;
    };

    const handlePasswordChange = (text: string) => {
        onChangeText(text);
        if (showStrengthIndicator) {
            setStrength(calculatePasswordStrength(text));
        }
    };

    const getStrengthText = (strength: number): string => {
        switch (strength) {
            case 0:
            case 1:
                return 'Weak';
            case 2:
            case 3:
                return 'Medium';
            case 4:
            case 5:
                return 'Strong';
            default:
                return '';
        }
    };

    const getStrengthColor = (strength: number): string => {
        switch (strength) {
            case 0:
            case 1:
                return '#FF6B6B';
            case 2:
            case 3:
                return '#FF9800';
            case 4:
            case 5:
                return '#4CAF50';
            default:
                return '#E1E5E9';
        }
    };

    return (
        <>
            <InputField
                label={label}
                value={value}
                onChangeText={handlePasswordChange}
                error={error}
                isPassword
                leftIcon="lock-closed"
                placeholder={placeholder}
                disabled={disabled}
                style={style}
            />
            {showStrengthIndicator && value.length > 0 && (
                <Text style={[
                    styles.strengthIndicator,
                    { color: getStrengthColor(strength) }
                ]}>
                    Password Strength: {getStrengthText(strength)}
                </Text>
            )}
        </>
    );
}

const styles = StyleSheet.create({
    strengthIndicator: {
        marginTop: 4,
        fontSize: 12,
        fontWeight: '500',
    },
});