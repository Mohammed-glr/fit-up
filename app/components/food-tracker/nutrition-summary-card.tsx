import React from 'react';
import { ActivityIndicator, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { DailyNutritionSummary, NutritionGoals } from '@/types/food-tracker';

const macroDefinitions = [
  { key: 'calories', label: 'Calories', unit: 'kcal' },
  { key: 'protein', label: 'Protein', unit: 'g' },
  { key: 'carbs', label: 'Carbs', unit: 'g' },
  { key: 'fat', label: 'Fat', unit: 'g' },
  { key: 'fiber', label: 'Fiber', unit: 'g' },
] as const;

type NutritionSummaryCardProps = {
  summary?: DailyNutritionSummary;
  goals?: NutritionGoals | null;
  isLoading?: boolean;
  onPressSetGoals?: () => void;
};

const resolveMacroValue = (summary: DailyNutritionSummary | undefined, key: typeof macroDefinitions[number]['key']) => {
  switch (key) {
    case 'calories':
      return summary?.total_calories ?? 0;
    case 'protein':
      return summary?.total_protein ?? 0;
    case 'carbs':
      return summary?.total_carbs ?? 0;
    case 'fat':
      return summary?.total_fat ?? 0;
    case 'fiber':
      return summary?.total_fiber ?? 0;
    default:
      return 0;
  }
};

const resolveGoalValue = (goals: NutritionGoals | undefined | null, key: typeof macroDefinitions[number]['key']) => {
  if (!goals) {
    return undefined;
  }

  switch (key) {
    case 'calories':
      return goals.calories_goal;
    case 'protein':
      return goals.protein_goal;
    case 'carbs':
      return goals.carbs_goal;
    case 'fat':
      return goals.fat_goal;
    case 'fiber':
      return goals.fiber_goal;
    default:
      return undefined;
  }
};

const formatNumber = (value: number) => {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}k`;
  }
  return Math.round(value).toString();
};

const getProgressRatio = (value: number, goal?: number) => {
  if (!goal || goal <= 0) {
    return 0;
  }
  return Math.min(value / goal, 1);
};

const getProgressColor = (value: number, goal?: number) => {
  if (!goal || goal <= 0) {
    return COLORS.primaryLight;
  }

  const ratio = value / goal;
  if (ratio < 0.7) {
    return COLORS.primaryLight;
  }
  if (ratio <= 1.1) {
    return COLORS.primary;
  }
  return '#F97316';
};

export function NutritionSummaryCard({ summary, goals, isLoading, onPressSetGoals }: NutritionSummaryCardProps) {
  if (isLoading) {
    return (
      <View style={[styles.container, styles.loadingContainer]}>
        <ActivityIndicator size="small" color={COLORS.primary} />
        <Text style={styles.loadingText}>Loading nutrition summary...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.headerRow}>
        <Text style={styles.heading}>Nutrition Overview</Text>
        {!goals && (
          <TouchableOpacity onPress={onPressSetGoals} style={styles.setGoalsButton}>
            <Text style={styles.setGoalsText}>Set goals</Text>
          </TouchableOpacity>
        )}
      </View>

      <View style={styles.caloriesRow}>
        <View>
          <Text style={styles.caloriesLabel}>Total calories</Text>
          <Text style={styles.caloriesValue}>{formatNumber(summary?.total_calories ?? 0)}</Text>
        </View>
        {goals?.calories_goal ? (
          <Text style={styles.goalText}>Goal: {formatNumber(goals.calories_goal)}</Text>
        ) : (
          <Text style={styles.goalTextMuted}>No calorie goal</Text>
        )}
      </View>

      {macroDefinitions.map((macro) => {
        const value = resolveMacroValue(summary, macro.key);
        const goalValue = resolveGoalValue(goals ?? undefined, macro.key);
        const progressRatio = getProgressRatio(value, goalValue);
        const progressColor = getProgressColor(value, goalValue);

        return (
          <View key={macro.key} style={styles.macroRow}>
            <View style={styles.macroHeader}>
              <Text style={styles.macroLabel}>{macro.label}</Text>
              <Text style={styles.macroValue}>{formatNumber(value)} {macro.unit}</Text>
            </View>
            <View style={styles.progressTrack}>
              <View style={[styles.progressFill, { flex: progressRatio || 0, backgroundColor: progressColor }]} />
              <View style={{ flex: 1 - (progressRatio || 0) }} />
            </View>
            <Text style={styles.goalCaption}>
              {goalValue ? `Goal ${formatNumber(goalValue)} ${macro.unit}` : 'No goal'}
            </Text>
          </View>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: COLORS.background.surface,
    borderRadius: BORDER_RADIUS.xl,
    padding: SPACING.xl,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    gap: SPACING.lg,
  },
  loadingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  loadingText: {
    marginLeft: SPACING.sm,
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
  },
  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  heading: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  setGoalsButton: {
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    backgroundColor: COLORS.primarySoft,
    borderRadius: BORDER_RADIUS.base,
  },
  setGoalsText: {
    color: COLORS.primary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  caloriesRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  caloriesLabel: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
  },
  caloriesValue: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
  },
  goalText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  goalTextMuted: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
  },
  macroRow: {
    gap: SPACING.xs,
  },
  macroHeader: {
    flexDirection: 'row',
    alignItems: 'baseline',
    justifyContent: 'space-between',
  },
  macroLabel: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
  },
  macroValue: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
  },
  progressTrack: {
    flexDirection: 'row',
    height: 10,
    backgroundColor: COLORS.primarySoft,
    borderRadius: BORDER_RADIUS.full,
    overflow: 'hidden',
  },
  progressFill: {
    borderTopLeftRadius: BORDER_RADIUS.full,
    borderBottomLeftRadius: BORDER_RADIUS.full,
  },
  goalCaption: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.xs,
  },
});
