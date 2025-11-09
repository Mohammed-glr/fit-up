import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity } from 'react-native';
import type { FoodLogEntryWithRecipe } from '@/types/food-tracker';
import { FoodLogEntryCard } from './food-log-entry-card';
import { BORDER_RADIUS, COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING } from '@/constants/theme';

type FoodLogListProps = {
  entries?: FoodLogEntryWithRecipe[];
  onPressEntry?: (entry: FoodLogEntryWithRecipe) => void;
};

export function FoodLogList({ entries, onPressEntry }: FoodLogListProps) {
  const [showAll, setShowAll] = React.useState(false);

  if (!entries || entries.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyTitle}>No meals logged</Text>
        <Text style={styles.emptySubtitle}>Log a meal to see it appear here and track your macros.</Text>
      </View>
    );
  }

  const visibleEntries = showAll ? entries : entries.slice(0, 3);

  return (
    <View style={styles.listContainer}>
      {visibleEntries.map((entry) => (
        <FoodLogEntryCard key={entry.id} entry={entry} onPress={onPressEntry} />
      ))}
      {entries.length > 3 && (
        <TouchableOpacity style={styles.viewAllButton} onPress={() => setShowAll(!showAll)}>
          <Text style={styles.viewAllText}>{showAll ? 'Show Less' : 'View All'}</Text>
        </TouchableOpacity>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  listContainer: {
    gap: SPACING.md,
  },
  emptyContainer: {
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    backgroundColor: COLORS.darkGray,
    gap: SPACING.xs,
  },
  emptyTitle: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  emptySubtitle: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
  },
  viewAllButton: {
    marginTop: SPACING.sm,
    padding: SPACING.sm,
    alignItems: 'center',
  },
  viewAllText: {
    color: COLORS.primary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
});
