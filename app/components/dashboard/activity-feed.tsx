import React from 'react';
import { View, Text, StyleSheet, ScrollView, ActivityIndicator } from 'react-native';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';
import { useActivityFeed, ActivityFeedItem, ActivityType } from '@/hooks/user/use-activity-feed';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface ActivityFeedProps {
  limit?: number;
}

const getActivityIcon = (type: ActivityType): keyof typeof Ionicons.glyphMap => {
  switch (type) {
    case 'workout_completed':
      return 'fitness';
    case 'pr_achieved':
      return 'trophy';
    case 'streak_milestone':
      return 'flame';
    case 'coach_message':
      return 'chatbubble';
    case 'new_plan':
      return 'calendar';
    case 'goal_achieved':
      return 'star';
    case 'plan_completed':
      return 'checkmark-circle';
    default:
      return 'information-circle';
  }
};

const getActivityColor = (type: ActivityType): string => {
  switch (type) {
    case 'workout_completed':
      return COLORS.success;
    case 'pr_achieved':
      return COLORS.warning;
    case 'streak_milestone':
      return COLORS.error;
    case 'coach_message':
      return COLORS.info;
    case 'new_plan':
      return COLORS.primary;
    case 'goal_achieved':
      return COLORS.warning;
    case 'plan_completed':
      return COLORS.success;
    default:
      return COLORS.text.secondary;
  }
};

const formatTimestamp = (timestamp: string): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return 'Just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
};

const ActivityItem: React.FC<{ activity: ActivityFeedItem; index: number }> = ({ activity, index }) => {
  const icon = getActivityIcon(activity.type);
  const color = getActivityColor(activity.type);

  return (
    <MotiView
      from={{ opacity: 0, translateX: -20 }}
      animate={{ opacity: 1, translateX: 0 }}
      transition={{ type: 'timing', duration: 300, delay: index * 50 }}
      style={styles.activityItem}
    >
      <View style={[styles.iconContainer, { backgroundColor: `${color}15` }]}>
        <Ionicons name={icon} size={24} color={color} />
      </View>
      
      <View style={styles.activityContent}>
        <View style={styles.activityHeader}>
          <Text style={styles.activityTitle}>{activity.title}</Text>
          <Text style={styles.activityTime}>{formatTimestamp(activity.timestamp)}</Text>
        </View>
        <Text style={styles.activityDescription} numberOfLines={2}>
          {activity.description}
        </Text>
      </View>
    </MotiView>
  );
};

export const ActivityFeed: React.FC<ActivityFeedProps> = ({ limit = 10 }) => {
  const { data, isLoading, error } = useActivityFeed(limit);

  if (isLoading) {
    return (
      <View style={styles.container}>
        <Text style={styles.sectionTitle}>Recent Activity</Text>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={COLORS.primary} />
          <Text style={styles.loadingText}>Loading activities...</Text>
        </View>
      </View>
    );
  }

  if (error || !data || data.activities.length === 0) {
    return (
      <View style={styles.container}>
        <Text style={styles.sectionTitle}>Recent Activity</Text>
        <View style={styles.emptyContainer}>
          <Ionicons name="time-outline" size={48} color={COLORS.text.tertiary} />
          <Text style={styles.emptyTitle}>No Recent Activity</Text>
          <Text style={styles.emptySubtitle}>
            Complete workouts to see your activity here
          </Text>
        </View>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.sectionTitle}>Recent Activity</Text>
      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        nestedScrollEnabled={true}
      >
        {data.activities.map((activity, index) => (
          <ActivityItem key={activity.id} activity={activity} index={index} />
        ))}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: SPACING.lg,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  loadingContainer: {
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.xl,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
  },
  emptyContainer: {
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.xl,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  emptyTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginTop: SPACING.md,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    marginTop: SPACING.xs,
    textAlign: 'center',
  },
  scrollView: {
    maxHeight: 400,
  },
  scrollContent: {
    paddingBottom: SPACING.sm,
  },
  activityItem: {
    flexDirection: 'row',
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.md,
    marginBottom: SPACING.md,
    ...SHADOWS.sm,
  },
  iconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  activityContent: {
    flex: 1,
  },
  activityHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.xs,
  },
  activityTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    flex: 1,
  },
  activityTime: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginLeft: SPACING.sm,
    marginRight: SPACING.sm,
  },
  activityDescription: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    lineHeight: 20,
  },
});
