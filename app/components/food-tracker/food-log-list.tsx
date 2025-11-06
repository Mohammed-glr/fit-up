import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import type { FoodLogEntryWithRecipe } from '@/types/food-tracker';
import { FoodLogEntryCard } from './food-log-entry-card';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING } from '@/constants/theme';

type FoodLogListProps = {
  entries?: FoodLogEntryWithRecipe[];
  onPressEntry?: (entry: FoodLogEntryWithRecipe) => void;
};

export function FoodLogList({ entries, onPressEntry }: FoodLogListProps) {
  if (!entries || entries.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyTitle}>No meals logged</Text>
        <Text style={styles.emptySubtitle}>Log a meal to see it appear here and track your macros.</Text>
      </View>
    );
  }

  return (
    <View style={styles.listContainer}>
      {entries.map((entry) => (
        <FoodLogEntryCard key={entry.id} entry={entry} onPress={onPressEntry} />
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  listContainer: {
    gap: SPACING.md,
  },
  emptyContainer: {
    padding: SPACING.lg,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    borderRadius: 20,
    backgroundColor: COLORS.background.secondary,
    gap: SPACING.xs,
  },
  emptyTitle: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  emptySubtitle: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
  },
});
