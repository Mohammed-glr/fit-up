import React from 'react';
import { View, StyleSheet, ViewStyle, ScrollView } from 'react-native';
import { ThemedView } from '@/components/themed-view';
import { SPACING } from '@/constants/theme';

interface FormContainerProps {
  children: React.ReactNode;
  style?: ViewStyle;
  fullScreen?: boolean;
  scrollable?: boolean;
}

export const FormContainer: React.FC<FormContainerProps> = ({ 
  children, 
  style,
  fullScreen = true,
  scrollable = true
}) => {
  const content = (
    <ThemedView style={[styles.container, fullScreen && styles.fullScreen, style]} fullScreen={fullScreen}>
      <View style={styles.content}>
        {children}
      </View>
    </ThemedView>
  );

  if (scrollable && fullScreen) {
    return (
      <ScrollView 
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        {content}
      </ScrollView>
    );
  }

  return content;
};

const styles = StyleSheet.create({
  container: {
    width: '100%',
  },
  fullScreen: {
    flex: 1,
  },
  content: {
    padding: SPACING.xl,
    paddingTop: SPACING['3xl'],
    paddingBottom: SPACING['4xl'],
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
  },
});
