import React from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { UserStats } from '@/hooks/user/use-user-stats';

interface MetricsGridProps {
  stats?: UserStats;
  isLoading: boolean;
}

interface MetricCardProps {
  icon: keyof typeof Ionicons.glyphMap;
  label: string;
  value: string | number;
  color: string;
  delay: number;
  isLoading?: boolean;
}

const MetricCard: React.FC<MetricCardProps> = ({ icon, label, value, color, delay, isLoading }) => (
  <MotiView
    from={{ opacity: 0, scale: 0.9 }}
    animate={{ opacity: 1, scale: 1 }}
    transition={{ type: 'timing', duration: 400, delay }}
    style={[styles.metricCard]}
  >
    <View style={[styles.iconContainer, { backgroundColor: `${color}15` }]}>
      <Ionicons name={icon} size={24} color={color} />
    </View>
    <View style={styles.metricContent}>
      <Text style={styles.metricLabel}>{label}</Text>
      {isLoading ? (
        <ActivityIndicator size="small" color={color} />
      ) : (
        <Text style={styles.metricValue}>{value}</Text>
      )}
    </View>
  </MotiView>
);

export const MetricsGrid: React.FC<MetricsGridProps> = ({ stats, isLoading }) => {
  return (
    <View style={styles.container}>
      <View style={styles.row}>
        <MetricCard
          icon="flame"
          label="Current Streak"
          value={`${stats?.current_streak || 0} days`}
          color={COLORS.error}
          delay={0}
          isLoading={isLoading}
        />
        <MetricCard
          icon="trophy"
          label="Longest Streak"
          value={`${stats?.longest_streak || 0} days`}
          color={COLORS.warning}
          delay={100}
          isLoading={isLoading}
        />
      </View>
      <View style={styles.row}>
        <MetricCard
          icon="calendar"
          label="Total Weeks"
          value={stats?.total_weeks || 0}
          color={COLORS.info}
          delay={200}
          isLoading={isLoading}
        />
        <MetricCard
          icon="checkmark-circle"
          label="Completion Rate"
          value={`${Math.round(stats?.completion_rate || 0)}%`}
          color={COLORS.success}
          delay={300}
          isLoading={isLoading}
        />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: SPACING.lg,
  },
  row: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginBottom: SPACING.md,
  },
  metricCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.md,
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    ...SHADOWS.sm,
  },
  iconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.lg,
    justifyContent: 'center',
    alignItems: 'center',
  },
  metricContent: {
    flex: 1,
  },
  metricLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.tertiary,
    marginBottom: 4,
  },
  metricValue: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
  },
});
