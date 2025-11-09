import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING } from '@/constants/theme';

/**
 * Workout Templates Screen
 * TODO: Complete implementation with:
 * - User templates list
 * - Public templates browser
 * - Create/Edit/Delete templates
 * - Start workout from template
 */
export default function TemplatesScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Workout Templates</Text>
      <Text style={styles.message}>
        Coming Soon: Save and reuse workout configurations
      </Text>
      <Text style={styles.subMessage}>
        Create custom templates, browse community templates, and quick-start workouts
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.primary,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  title: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.md,
    textAlign: 'center',
  },
  message: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
    marginBottom: SPACING.sm,
    textAlign: 'center',
  },
  subMessage: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
});
