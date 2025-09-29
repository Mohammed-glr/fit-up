import React, { useState, useEffect, useRef } from 'react';
import { Text, StyleSheet, View, Animated } from 'react-native';
import { InputField } from '@/components/forms';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';

interface PasswordInputProps {
    label?: string;
    value: string;
    onChangeText: (text: string) => void;
    error?: string;
    placeholder?: string;
    disabled?: boolean;
    style?: any;
    showStrengthIndicator?: boolean;
    showRequirements?: boolean;
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
    showRequirements = false,
}: PasswordInputProps) {
    const [strength, setStrength] = useState(0);
    const fadeAnim = useRef(new Animated.Value(0)).current;
    const scaleAnims = useRef([1, 2, 3, 4, 5].map(() => new Animated.Value(0.8))).current;
    const requirementAnims = useRef([0, 1, 2, 3, 4].map(() => new Animated.Value(0))).current;
    const pulseAnim = useRef(new Animated.Value(1)).current;

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
            const newStrength = calculatePasswordStrength(text);
            setStrength(newStrength);
            
            scaleAnims.forEach((anim, index) => {
                if (newStrength > index) {
                    Animated.spring(anim, {
                        toValue: 1,
                        useNativeDriver: true,
                        tension: 200,
                        friction: 7,
                    }).start();
                } else {
                    Animated.spring(anim, {
                        toValue: 0.8,
                        useNativeDriver: true,
                        tension: 200,
                        friction: 7,
                    }).start();
                }
            });

            if (newStrength >= 4) {
                Animated.sequence([
                    Animated.timing(pulseAnim, {
                        toValue: 1.1,
                        duration: 200,
                        useNativeDriver: true,
                    }),
                    Animated.timing(pulseAnim, {
                        toValue: 1,
                        duration: 200,
                        useNativeDriver: true,
                    })
                ]).start();
            }
        }
        
        if (showRequirements) {
            const requirements = getPasswordRequirements(text);
            requirements.forEach((req, index) => {
                Animated.spring(requirementAnims[index], {
                    toValue: req.met ? 1 : 0,
                    useNativeDriver: true,
                    tension: 200,
                    friction: 7,
                }).start();
            });
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

    const getPasswordRequirements = (password: string) => {
        return [
            { text: 'At least 8 characters', met: password.length >= 8 },
            { text: 'Lowercase letter', met: /[a-z]/.test(password) },
            { text: 'Uppercase letter', met: /[A-Z]/.test(password) },
            { text: 'Number', met: /[0-9]/.test(password) },
            { text: 'Special character', met: /[^A-Za-z0-9]/.test(password) },
        ];
    };

    const getStrengthColor = (strength: number): string => {
        switch (strength) {
            case 0:
            case 1:
                return COLORS.error;
            case 2:
            case 3:
                return COLORS.warning;
            case 4:
            case 5:
                return COLORS.success;
            default:
                return COLORS.border.light;
        }
    };

    useEffect(() => {
        if (showStrengthIndicator && value.length > 0) {
            Animated.timing(fadeAnim, {
                toValue: 1,
                duration: 300,
                useNativeDriver: true,
            }).start();
        } else {
            Animated.timing(fadeAnim, {
                toValue: 0,
                duration: 300,
                useNativeDriver: true,
            }).start();
        }
    }, [showStrengthIndicator, value.length, fadeAnim]);

    return (
        <>
            <InputField
                label={label}
                value={value}
                onChangeText={handlePasswordChange}
                error={error}
                isPassword
                placeholder={placeholder}
                disabled={disabled}
                style={style}
            />
            {showStrengthIndicator && value.length > 0 && (
                <Animated.View style={[
                    styles.strengthContainer,
                    {
                        opacity: fadeAnim,
                        transform: [{
                            translateY: fadeAnim.interpolate({
                                inputRange: [0, 1],
                                outputRange: [10, 0],
                            })
                        }]
                    }
                ]}>
                    <View style={styles.strengthBarContainer}>
                        {[1, 2, 3, 4, 5].map((segment, index) => (
                            <Animated.View
                                key={segment}
                                style={[
                                    styles.strengthSegment,
                                    {
                                        backgroundColor: strength >= segment 
                                            ? getStrengthColor(strength) 
                                            : COLORS.border.light,
                                        transform: [{ scaleY: scaleAnims[index] }]
                                    }
                                ]}
                            />
                        ))}
                    </View>
                    {/* <View style={styles.strengthTextContainer}>
                        <Animated.Text style={[
                            styles.strengthIndicator,
                            { 
                                color: getStrengthColor(strength),
                                transform: [
                                    {
                                        scale: fadeAnim.interpolate({
                                            inputRange: [0, 1],
                                            outputRange: [0.8, 1],
                                        })
                                    },
                                    { scale: strength >= 4 ? pulseAnim : 1 }
                                ]
                            }
                        ]}>
                            Password Strength: {getStrengthText(strength)}
                        </Animated.Text>
                    </View> */}
                </Animated.View>
            )}
            {showRequirements && value.length > 0 && (
                <Animated.View style={[
                    styles.requirementsContainer,
                    {
                        opacity: fadeAnim,
                        transform: [{
                            translateY: fadeAnim.interpolate({
                                inputRange: [0, 1],
                                outputRange: [15, 0],
                            })
                        }]
                    }
                ]}>
                    {/* <Animated.Text style={[
                        styles.requirementsTitle,
                        {
                            transform: [{
                                scale: fadeAnim.interpolate({
                                    inputRange: [0, 1],
                                    outputRange: [0.9, 1],
                                })
                            }]
                        }
                    ]}>
                        Password Requirements:
                    </Animated.Text> */}
                    {getPasswordRequirements(value).map((req, index) => (
                        <Animated.View 
                            key={index} 
                            style={[
                                styles.requirementRow,
                                {
                                    transform: [
                                        {
                                            translateX: requirementAnims[index].interpolate({
                                                inputRange: [0, 1],
                                                outputRange: [10, 0],
                                            })
                                        },
                                        {
                                            scale: requirementAnims[index].interpolate({
                                                inputRange: [0, 1],
                                                outputRange: [0.95, 1],
                                            })
                                        }
                                    ]
                                }
                            ]}
                        >
                            <Animated.Text style={[
                                styles.requirementIcon,
                                { 
                                    color: req.met ? COLORS.success : COLORS.text.secondary,
                                    transform: [{
                                        rotate: requirementAnims[index].interpolate({
                                            inputRange: [0, 1],
                                            outputRange: ['0deg', '360deg'],
                                        })
                                    }]
                                }
                            ]}>
                                {req.met ? '✓' : '○'}
                            </Animated.Text>
                            <Animated.Text style={[
                                styles.requirementText,
                                { 
                                    color: req.met ? COLORS.success : COLORS.text.secondary,
                                    opacity: requirementAnims[index].interpolate({
                                        inputRange: [0, 1],
                                        outputRange: [0.6, 1],
                                    })
                                }
                            ]}>
                                {req.text}
                            </Animated.Text>
                        </Animated.View>
                    ))}
                </Animated.View>
            )}
        </>
    );
}

const styles = StyleSheet.create({
    strengthContainer: {
        alignItems: 'center',
        justifyContent: 'center',
        marginTop: SPACING.xs,
    },
    strengthBarContainer: {
        alignItems: 'center',
        justifyContent: 'center',
        width: '90%',
        flexDirection: 'row',
        height: 4,
        borderRadius: 4,
        overflow: 'hidden',
        marginBottom: SPACING.xs,
        gap: 2,
    },
    strengthSegment: {
        flex: 1,
        height: '100%',
        borderRadius: 1,
        backgroundColor: COLORS.border.light,
    },
    strengthTextContainer: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    strengthLabel: {
        fontSize: FONT_SIZES.xs,
        fontWeight: FONT_WEIGHTS.medium,
        color: COLORS.text.secondary,
    },
    strengthIndicator: {
        fontSize: FONT_SIZES.xs,
        fontWeight: FONT_WEIGHTS.medium,
    },
    requirementsContainer: {
        marginTop: SPACING.sm,
        paddingTop: SPACING.xs,
       
    },
    requirementsTitle: {
        fontSize: FONT_SIZES.xs,
        fontWeight: FONT_WEIGHTS.medium,
        color: COLORS.text.secondary,
        marginBottom: SPACING.xs,
    },
    requirementRow: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 2,
    },
    requirementIcon: {
        fontSize: FONT_SIZES.sm,
        marginRight: SPACING.xs,
        width: 16,
        textAlign: 'center',
    },
    requirementText: {
        fontSize: FONT_SIZES.xs,
        flex: 1,
    },
});