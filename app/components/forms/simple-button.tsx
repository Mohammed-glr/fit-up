import React from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
  ViewStyle,
  TextStyle,
} from 'react-native';

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
  const getButtonStyle = (): (ViewStyle | Partial<ViewStyle>)[] => {
    let buttonStyle: (ViewStyle | Partial<ViewStyle>)[] = [styles.button];
    
    if (disabled || loading) {
      buttonStyle.push(styles.disabled);
    } else if (variant === 'primary') {
      buttonStyle.push(styles.primary);
    } else if (variant === 'secondary') {
      buttonStyle.push(styles.secondary);
    } else if (variant === 'outline') {
      buttonStyle.push(styles.outline);
    }
    
    return buttonStyle;
  }

const getTextStyle = (): (TextStyle | Partial<TextStyle>)[] => {
  let textStyle: (TextStyle | Partial<TextStyle>)[] = [styles.text];
  if (disabled || loading) {
    textStyle.push(styles.disabledText);
  } else if (variant === 'primary') {
    textStyle.push(styles.primaryText);
  } else if (variant === 'secondary') {
    textStyle.push(styles.secondaryText);
  } else if (variant === 'outline') {
    textStyle.push(styles.outlineText);
  } 
  
  return textStyle;
};

  return (
    <TouchableOpacity
      style={[...getButtonStyle(), style]}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator 
          color={variant === 'primary' ? '#FFFFFF' : '#007AFF'} 
          size="small" 
        />
      ) : (
        <Text style={getTextStyle()}>
          {title}
        </Text>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    minHeight: 44,
  },
  primary: {
    backgroundColor: '#007AFF',
  },
  secondary: {
    backgroundColor: '#F2F2F7',
  },
  outline: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: '#007AFF',
  },
  disabled: {
    backgroundColor: '#C7C7CC',
    opacity: 0.6,
  },
  text: {
    fontSize: 16,
    fontWeight: '600',
  },
  primaryText: {
    color: '#FFFFFF',
  },
  secondaryText: {
    color: '#000000',
  },
  outlineText: {
    color: '#007AFF',
  },
  disabledText: {
    color: '#8E8E93',
  },
});