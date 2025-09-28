import React from 'react';
import { TouchableOpacity, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS, FONT_WEIGHTS, SHADOWS } from '@/constants/theme';

interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: 'primary' | 'secondary' | 'outline';
  disabled?: boolean;
  loading?: boolean;
  style?: any;
}

export const Button: React.FC<ButtonProps> = ({
  title,
  onPress,
  variant = 'primary',
  disabled = false,
  loading = false,
  style,
}) => {
  return (
    <TouchableOpacity
      style={[
        styles.button,
        variant === 'primary' && styles.primary,
        variant === 'secondary' && styles.secondary,
        variant === 'outline' && styles.outline,
        (disabled || loading) && styles.disabled,
        style,
      ]}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator 
          color={variant === 'primary' ? COLORS.white : COLORS.primary} 
          size="small" 
        />
      ) : (
        <Text
          style={[
            styles.text,
            variant === 'primary' && styles.primaryText,
            variant === 'secondary' && styles.secondaryText,
            variant === 'outline' && styles.outlineText,
            (disabled || loading) && styles.disabledText,
          ]}
        >
          {title}
        </Text>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.base,
    minHeight: 48,
  },
  primary: {
    backgroundColor: COLORS.primary,
    ...SHADOWS.sm,
  },
  secondary: {
    backgroundColor: COLORS.background.secondary,
    borderWidth: 1,
    borderColor: COLORS.border.subtle || COLORS.border.light,
  },
  outline: {
    backgroundColor: 'transparent',
    borderWidth: 1.5,
    borderColor: COLORS.primary,
  },
  disabled: {
    backgroundColor: COLORS.lightGray,
    opacity: 0.6,
  },
  text: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  primaryText: {
    color: COLORS.text.inverse,
  },
  secondaryText: {
    color: COLORS.text.primary,
  },
  outlineText: {
    color: COLORS.primary,
  },
  disabledText: {
    color: COLORS.text.placeholder,
  },
});