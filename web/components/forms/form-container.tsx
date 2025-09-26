import React from 'react';
import { View, StyleSheet, ViewStyle } from 'react-native';
import { ThemedView } from '@/components/themed-view';

interface FormContainerProps {
  children: React.ReactNode;
  style?: ViewStyle;
}

export const FormContainer: React.FC<FormContainerProps> = ({ 
  children, 
  style 
}) => {
  return (
    <ThemedView style={[styles.container, style]}>
      {children}
    </ThemedView>
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 20,
    width: '100%',
  },
});
