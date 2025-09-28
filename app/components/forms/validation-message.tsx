import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS, getColorWithOpacity } from '@/constants/theme';

interface ValidationMessageProps {
  message: string;
  type?: 'error' | 'success' | 'warning' | 'info';
  showIcon?: boolean;
  style?: any;
}

export const ValidationMessage: React.FC<ValidationMessageProps> = ({
  message,
  type = 'error',
  showIcon = true,
  style,
}) => {
  const getIconName = () => {
    switch (type) {
      case 'error':
        return 'close-circle';
      case 'success':
        return 'checkmark-circle';
      case 'warning':
        return 'warning';
      case 'info':
        return 'information-circle';
      default:
        return 'close-circle';
    }
  };

  const getColor = () => {
    switch (type) {
      case 'error':
        return COLORS.error;
      case 'success':
        return COLORS.success;
      case 'warning':
        return COLORS.warning;
      case 'info':
        return COLORS.info;
      default:
        return COLORS.error;
    }
  };

  const getBackgroundColor = () => {
    switch (type) {
      case 'error':
        return getColorWithOpacity(COLORS.error, 0.1);
      case 'success':
        return getColorWithOpacity(COLORS.success, 0.1);
      case 'warning':
        return getColorWithOpacity(COLORS.warning, 0.1);
      case 'info':
        return getColorWithOpacity(COLORS.info, 0.1);
      default:
        return getColorWithOpacity(COLORS.error, 0.1);
    }
  };

  if (!message) {
    return null;
  }

  return (
    <View
      style={[
        styles.container,
        { backgroundColor: getBackgroundColor() },
        style,
      ]}
    >
      {showIcon && (
        <Ionicons
          name={getIconName()}
          size={16}
          color={getColor()}
          style={styles.icon}
        />
      )}
      <Text
        style={[
          styles.message,
          { color: getColor() },
        ]}
      >
        {message}
      </Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.base,
    marginVertical: SPACING.xs,
    borderWidth: 1,
    borderColor: 'transparent',
  },
  icon: {
    marginRight: SPACING.sm,
    marginTop: 1,
  },
  message: {
    flex: 1,
    fontSize: FONT_SIZES.sm,
    lineHeight: 20,
    fontWeight: '400',
  },
});
