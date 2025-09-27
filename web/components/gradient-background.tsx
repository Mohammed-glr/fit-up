import React from 'react';
import { View, StyleSheet } from 'react-native';
import { COLORS } from '@/constants/theme';

interface FullScreenBackgroundProps {
  children: React.ReactNode;
  style?: any;
  variant?: 'primary' | 'secondary' | 'accent';
}

export const FullScreenBackground: React.FC<FullScreenBackgroundProps> = ({
  children,
  style,
  variant = 'secondary',
}) => {
  const getBackgroundColor = () => {
    switch (variant) {
      case 'primary':
        return COLORS.background.primary;
      case 'secondary':
        return COLORS.background.secondary;
      case 'accent':
        return COLORS.background.accent;
      default:
        return COLORS.background.secondary;
    }
  };

  return (
    <View style={[styles.container, { backgroundColor: getBackgroundColor() }, style]}>
      {children}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    width: '100%',
  },
});