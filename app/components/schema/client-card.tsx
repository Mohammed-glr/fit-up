import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { ClientSummary } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface ClientCardProps {
  client: ClientSummary;
  onPress?: () => void;
}

export const ClientCard: React.FC<ClientCardProps> = ({ client, onPress }) => {
  const fullName = `${client.first_name} ${client.last_name}`;
  const completionPercentage = Math.round(client.completion_rate * 100);
  const hasActiveSchema = !!client.current_schema_id;
  const lastWorkoutDate = client.last_workout_date
    ? new Date(client.last_workout_date).toLocaleDateString()
    : null;

  return (
    <TouchableOpacity
      style={styles.card}
      onPress={onPress}
      activeOpacity={0.7}
      disabled={!onPress}
    >
      <View style={styles.header}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>
            {client.first_name[0]}{client.last_name[0]}
          </Text>
        </View>
        
        <View style={styles.info}>
          <Text style={styles.name}>{fullName}</Text>
          <Text style={styles.email}>{client.email}</Text>
          
          <View style={styles.badges}>
            <View style={[styles.badge, hasActiveSchema ? styles.badgeActive : styles.badgeInactive]}>
              <Text style={[styles.badgeText, hasActiveSchema && styles.badgeTextActive]}>
                {hasActiveSchema ? 'Active' : 'Inactive'}
              </Text>
            </View>
            <View style={styles.badge}>
              <Text style={styles.badgeText}>{client.fitness_level}</Text>
            </View>
          </View>
        </View>

        <TouchableOpacity style={styles.menuButton} onPress={(e) => {
          e.stopPropagation();
          // TODO: Show options menu
        }}>
          <Ionicons name="ellipsis-vertical" size={20} color={COLORS.text.auth.secondary} />
        </TouchableOpacity>
      </View>

      <View style={styles.stats}>
        <View style={styles.statItem}>
          <Ionicons name="barbell" size={16} color={COLORS.primary} />
          <Text style={styles.statLabel}>{client.total_workouts}</Text>
          <Text style={styles.statSubtext}>workouts</Text>
        </View>
        <View style={styles.statDivider} />
        <View style={styles.statItem}>
          <Ionicons name="flame" size={16} color={COLORS.warning} />
          <Text style={styles.statLabel}>{client.current_streak}</Text>
          <Text style={styles.statSubtext}>day streak</Text>
        </View>
        <View style={styles.statDivider} />
        <View style={styles.statItem}>
          <Ionicons name="trending-up" size={16} color={COLORS.success} />
          <Text style={styles.statLabel}>{completionPercentage}%</Text>
          <Text style={styles.statSubtext}>completion</Text>
        </View>
      </View>

      {lastWorkoutDate && (
        <View style={styles.footer}>
          <Ionicons name="calendar-outline" size={14} color={COLORS.text.tertiary} />
          <Text style={styles.lastWorkout}>Last workout: {lastWorkoutDate}</Text>
        </View>
      )}

      {client.active_goals > 0 && (
        <View style={styles.goalsTag}>
          <Ionicons name="trophy" size={14} color={COLORS.primary} />
          <Text style={styles.goalsText}>{client.active_goals} active goals</Text>
        </View>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    marginBottom: SPACING.sm,
    ...SHADOWS.sm,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: SPACING.md,
  },
  avatar: {
    width: 52,
    height: 52,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  avatarText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  info: {
    flex: 1,
  },
  name: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 2,
  },
  email: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
  },
  badges: {
    flexDirection: 'row',
    gap: SPACING.xs,
  },
  badge: {
    paddingHorizontal: SPACING.sm,
    paddingVertical: 2,
    borderRadius: BORDER_RADIUS.sm,
    backgroundColor: COLORS.background.primary,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
  },
  badgeActive: {
    backgroundColor: COLORS.success,
    borderColor: COLORS.success,
  },
  badgeInactive: {
    backgroundColor: 'transparent',
    borderColor: COLORS.border.subtle,
  },
  badgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    textTransform: 'capitalize',
  },
  badgeTextActive: {
    color: COLORS.text.primary,
  },
  menuButton: {
    padding: SPACING.xs,
  },
  stats: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-around',
    paddingVertical: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.subtle,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.subtle,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
    gap: 2,
  },
  statDivider: {
    width: 1,
    height: 28,
    backgroundColor: COLORS.border.subtle,
  },
  statLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  statSubtext: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  footer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginTop: SPACING.md,
  },
  lastWorkout: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  goalsTag: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
    backgroundColor: COLORS.background.primary,
    borderRadius: BORDER_RADIUS.sm,
    alignSelf: 'flex-start',
    marginTop: SPACING.sm,
  },
  goalsText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
});
