import React, { useState } from 'react';
import { TextInput, View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS, FONT_WEIGHTS, SHADOWS } from '@/constants/theme';

interface InputFieldProps {
  label?: string;
  placeholder?: string;
  value: string;
  onChangeText: (text: string) => void;
  error?: string;
  isPassword?: boolean;
  leftIcon?: keyof typeof Ionicons.glyphMap;
  rightIcon?: keyof typeof Ionicons.glyphMap;
  onRightIconPress?: () => void;
  keyboardType?: 'default' | 'email-address' | 'numeric' | 'phone-pad';
  autoCapitalize?: 'none' | 'sentences' | 'words' | 'characters';
  style?: any;
  disabled?: boolean;
}

export const InputField: React.FC<InputFieldProps> = ({
  label,
  placeholder,
  value,
  onChangeText,
  error,
  isPassword = false,
  leftIcon,
  rightIcon,
  onRightIconPress,
  keyboardType = 'default',
  autoCapitalize = 'none',
  style,
  disabled = false,
}) => {
  const [isPasswordVisible, setIsPasswordVisible] = useState(false);
  const [isFocused, setIsFocused] = useState(false);

  const togglePasswordVisibility = () => {
    setIsPasswordVisible(!isPasswordVisible);
  };

  const getSecureTextEntry = () => {
    if (isPassword) {
      return !isPasswordVisible;
    }
    return false;
  };

  const getRightIcon = () => {
    if (isPassword) {
      return isPasswordVisible ? 'eye-off' : 'eye';
    }
    return rightIcon;
  };

  const handleRightIconPress = () => {
    if (isPassword) {
      togglePasswordVisibility();
    } else if (onRightIconPress) {
      onRightIconPress();
    }
  };

  return (
    <View style={[styles.container, style]}>
      {label && (
        <Text style={styles.label}>
          {label}
        </Text>
      )}
      
      <View
        style={[
          styles.inputContainer,
          isFocused && styles.focused,
          error && styles.error,
          disabled && styles.disabled,
        ]}
      >
        {leftIcon && (
          <Ionicons
            name={leftIcon}
            size={20}
            color={COLORS.text.tertiary}
            style={styles.leftIcon}
          />
        )}
        
        <TextInput
          value={value}
          onChangeText={onChangeText}
          placeholder={placeholder}
          secureTextEntry={getSecureTextEntry()}
          keyboardType={keyboardType}
          autoCapitalize={autoCapitalize}
          editable={!disabled}
          style={styles.input}
          placeholderTextColor={COLORS.text.placeholder}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        
        {(getRightIcon() || onRightIconPress) && (
          <TouchableOpacity
            onPress={handleRightIconPress}
            style={styles.rightIconContainer}
            activeOpacity={0.7}
          >
            <Ionicons
              name={getRightIcon()}
              size={20}
              color={COLORS.text.tertiary}
            />
          </TouchableOpacity>
        )}
      </View>
      
      {error && (
        <Text style={styles.errorText}>
          {error}
        </Text>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: SPACING.base,
  },
  label: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    marginBottom: SPACING.sm,
    color: COLORS.text.primary,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 0.5,
    borderColor: COLORS.border.accent || COLORS.border.accent,
    borderRadius: BORDER_RADIUS.full,
    paddingHorizontal: SPACING.base,
    minHeight: 52,
    backgroundColor: COLORS.background.surface || COLORS.background.primary,
    ...SHADOWS.sm,
    transitionDuration: '200ms',
    transitionProperty: 'all',
    transitionTimingFunction: 'ease-in-out',
    
  },
  focused: {
    borderColor: COLORS.primary,
    borderWidth: 1.5,
  },
  error: {
    borderColor: COLORS.error,
  },
  disabled: {
    backgroundColor: COLORS.background.secondary,
    opacity: 0.6,
  },
  input: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    paddingVertical: SPACING.md,
    color: COLORS.text.primary,
  },
  leftIcon: {
    marginRight: SPACING.md,
    padding: SPACING.xs,
    color: COLORS.text.placeholder,
  },
  rightIconContainer: {
    marginLeft: SPACING.md,
    padding: SPACING.xs,
  },
  errorText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.error,
    marginTop: SPACING.xs,
    marginLeft: SPACING.xs,
    fontWeight: FONT_WEIGHTS.medium,
    backgroundColor: COLORS.background.errorSoft,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    overflow: 'hidden',
    width : '100%', 
  },
});