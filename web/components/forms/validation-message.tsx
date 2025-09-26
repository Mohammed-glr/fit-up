import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

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
        return '#FF6B6B';
      case 'success':
        return '#4CAF50';
      case 'warning':
        return '#FF9800';
      case 'info':
        return '#2196F3';
      default:
        return '#FF6B6B';
    }
  };

  const getBackgroundColor = () => {
    switch (type) {
      case 'error':
        return '#FFF5F5';
      case 'success':
        return '#F5FFF5';
      case 'warning':
        return '#FFF8F0';
      case 'info':
        return '#F0F8FF';
      default:
        return '#FFF5F5';
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
    padding: 12,
    borderRadius: 6,
    marginVertical: 4,
    borderWidth: 1,
    borderColor: 'transparent',
  },
  icon: {
    marginRight: 8,
    marginTop: 1,
  },
  message: {
    flex: 1,
    fontSize: 14,
    lineHeight: 20,
    fontWeight: '400',
  },
});
