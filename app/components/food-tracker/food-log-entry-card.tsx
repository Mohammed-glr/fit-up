import React from 'react';
import { StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import Ionicons from '@expo/vector-icons/Ionicons';
import type { FoodLogEntryWithRecipe } from '@/types/food-tracker';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';

const capitalize = (value?: string | null) => {
  if (!value) {
    return '';
  }
  return value.charAt(0).toUpperCase() + value.slice(1);
};

type FoodLogEntryCardProps = {
  entry: FoodLogEntryWithRecipe;
  onPress?: (entry: FoodLogEntryWithRecipe) => void;
};

export function FoodLogEntryCard({ entry, onPress }: FoodLogEntryCardProps) {
  const handlePress = React.useCallback(() => {
    if (onPress) {
      onPress(entry);
    }
  }, [entry, onPress]);

  return (
    <TouchableOpacity activeOpacity={0.85} style={styles.container} onPress={handlePress}>
      
      <View style={styles.header}>
        <View style={styles.mealTypeBadge}>
          <Text style={styles.mealTypeText}>{capitalize(entry.meal_type)}</Text>
        </View>
        <View style={styles.calorieRow}>
          <Ionicons name="flame" size={16} color={COLORS.primary} />
          <Text style={styles.calorieText}>{Math.round(entry.calories)} kcal</Text>
        </View>
      </View>

      <Text style={styles.title} numberOfLines={1}>
        {entry.recipe_name ?? 'Logged meal'}
      </Text>

      <View style={styles.macrosRow}>
        <Text style={styles.macroText}>P {Math.round(entry.protein)}g</Text>
        <Text style={styles.divider}>•</Text>
        <Text style={styles.macroText}>C {Math.round(entry.carbs)}g</Text>
        <Text style={styles.divider}>•</Text>
        <Text style={styles.macroText}>F {Math.round(entry.fat)}g</Text>
        <Text style={styles.divider}>•</Text>
        <Text style={styles.macroText}>Servings {entry.servings}</Text>
      </View>

      <View style={styles.footerRow}>
        <Text style={styles.sourceText}>
          {entry.recipe_source === 'system' ? 'FitUp recipe' : 'Your recipe'}
        </Text>
        <Text style={styles.timeText}>{new Date(entry.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: COLORS.darkGray,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    gap: SPACING.sm,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  mealTypeBadge: {
    backgroundColor: COLORS.primarySoft,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
    borderRadius: BORDER_RADIUS.base,
  },
  mealTypeText: {
    color: COLORS.primary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    letterSpacing: 0.2,
  },
  calorieRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  calorieText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  title: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  macrosRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  macroText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  divider: {
    color: COLORS.text.placeholder,
  },
  footerRow: {
    marginTop: SPACING.xs,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  sourceText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.xs,
  },
  timeText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.xs,
  },
});
