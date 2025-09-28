import React, { useState } from 'react';
import { View, StyleSheet, Text, TouchableOpacity, Animated } from 'react-native';
import { Button } from '@/components/forms';
import { SPACING, COLORS, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { useToastMethods } from '@/components/ui/toast-provider';
import Ionicons from '@expo/vector-icons/Ionicons';
import { MotiView } from 'moti';

interface OAuthButtonsProps {
    disabled?: boolean;
}

export default function OAuthButtons({ disabled = false }: OAuthButtonsProps) {
    const { showInfo, showError } = useToastMethods();
    const [isExpanded, setIsExpanded] = useState(false);

    const toggleExpanded = () => {
        setIsExpanded(!isExpanded);
    };

    const handleGoogleLogin = async () => {
        setIsExpanded(false); // Close popover
        try {
            // TODO: Implement Google OAuth login
            showInfo('Google login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('Google login failed:', error);
            showError('Google login failed. Please try again.');
        }
    };

    const handleGitHubLogin = async () => {
        setIsExpanded(false); // Close popover
        try {
            // TODO: Implement GitHub OAuth login
            showInfo('GitHub login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('GitHub login failed:', error);
            showError('GitHub login failed. Please try again.');
        }
    };

    const handleFacebookLogin = async () => {
        setIsExpanded(false); // Close popover
        try {
            // TODO: Implement Facebook OAuth login
            showInfo('Facebook login will be implemented soon!', {
                position: 'top',
                duration: 4000,
            });
        } catch (error) {
            console.error('Facebook login failed:', error);
            showError('Facebook login failed. Please try again.');
        }
    };

    return (
        <View style={styles.container}>
            {/* Main OAuth Button */}
            <TouchableOpacity
                style={[styles.mainButton, disabled && styles.disabled]}
                onPress={toggleExpanded}
                disabled={disabled}
                activeOpacity={0.7}
            >
                <View style={styles.mainButtonContent}>
                    <Ionicons 
                        name="logo-google" 
                        size={20} 
                        color={COLORS.text.secondary} 
                        style={styles.mainButtonIcon}
                    />
                    <Text style={styles.mainButtonText}>Continue with Social</Text>
                    <MotiView
                        animate={{ rotate: isExpanded ? '180deg' : '0deg' }}
                        transition={{ type: 'timing', duration: 200 }}
                    >
                        <Ionicons 
                            name="chevron-down" 
                            size={16} 
                            color={COLORS.text.tertiary}
                        />
                    </MotiView>
                </View>
            </TouchableOpacity>

            {/* Popover Container */}
            {isExpanded && (
                <>
                    {/* Backdrop */}
                    <TouchableOpacity
                        style={styles.backdrop}
                        onPress={() => setIsExpanded(false)}
                        activeOpacity={1}
                    />
                    
                    <MotiView
                        from={{ opacity: 0, scale: 0.95, translateY: -10 }}
                        animate={{ opacity: 1, scale: 1, translateY: 0 }}
                        exit={{ opacity: 0, scale: 0.95, translateY: -10 }}
                        transition={{ type: 'spring', damping: 15, stiffness: 150 }}
                        style={[styles.popover, { backgroundColor: '#FFFFFF' }]}
                    >
                    <View style={styles.popoverContent}>
                        <TouchableOpacity
                            style={[styles.oauthOption, disabled && styles.disabled]}
                            onPress={handleGoogleLogin}
                            disabled={disabled}
                            activeOpacity={0.7}
                        >
                            <Ionicons name="logo-google" size={20} color="#DB4437" />
                            <Text style={styles.oauthOptionText}>Continue with Google</Text>
                        </TouchableOpacity>

                        <View style={styles.separator} />

                        <TouchableOpacity
                            style={[styles.oauthOption, disabled && styles.disabled]}
                            onPress={handleGitHubLogin}
                            disabled={disabled}
                            activeOpacity={0.7}
                        >
                            <Ionicons name="logo-github" size={20} color="#333" />
                            <Text style={styles.oauthOptionText}>Continue with GitHub</Text>
                        </TouchableOpacity>

                        <View style={styles.separator} />

                        <TouchableOpacity
                            style={[styles.oauthOption, disabled && styles.disabled]}
                            onPress={handleFacebookLogin}
                            disabled={disabled}
                            activeOpacity={0.7}
                        >
                            <Ionicons name="logo-facebook" size={20} color="#1877F2" />
                            <Text style={styles.oauthOptionText}>Continue with Facebook</Text>
                        </TouchableOpacity>
                    </View>
                </MotiView>
                </>
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        width: '100%',
        position: 'relative',
    },
    mainButton: {
        backgroundColor: COLORS.background.surface,
        borderWidth: 1,
        borderColor: COLORS.primarySoft,
        borderRadius: BORDER_RADIUS.full,
        paddingVertical: SPACING.md,
        paddingHorizontal: SPACING.base,
        ...SHADOWS.sm,
    },
    mainButtonContent: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
    },
    mainButtonIcon: {
        marginRight: SPACING.sm,
        color: COLORS.primaryDark,
    },
    mainButtonText: {
        flex: 1,
        textAlign: 'center',
        fontSize: FONT_SIZES.base,
        fontWeight: FONT_WEIGHTS.medium,
        color: COLORS.primaryDark,
    },
    popover: {
        position: 'absolute',
        top: '100%',
        left: 0,
        right: 0,
        marginTop: SPACING.sm,
        zIndex: 9999,
        backgroundColor: '#FFFFFF',
        borderRadius: BORDER_RADIUS.lg,
        borderWidth: 1,
        borderColor: COLORS.border.light,
        ...SHADOWS.lg,
        overflow: 'hidden',
    },
    popoverContent: {
    },
    oauthOption: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingVertical: SPACING.md,
        paddingHorizontal: SPACING.base,
        backgroundColor: COLORS.background.surface,
    },
    oauthOptionText: {
        marginLeft: SPACING.md,
        fontSize: FONT_SIZES.base,
        fontWeight: FONT_WEIGHTS.medium,
        color: COLORS.text.primary,
        flex: 1,
    },
    separator: {
        height: 1,
        backgroundColor: COLORS.border.subtle,
    },
    disabled: {
        opacity: 0.5,
    },
    backdrop: {
        position: 'absolute',
        top: -1000, 
        left: -1000,
        right: -1000,
        bottom: -1000,
        backgroundColor: COLORS.surface.backdrop,
        zIndex: 9998,
    },
});