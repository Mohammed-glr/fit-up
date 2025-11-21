import React from 'react';
import { SafeAreaView, StyleSheet, ScrollView, View, Text, TouchableOpacity, RefreshControl, Platform } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useCoachDashboard } from '@/hooks/schema/use-coach';
import { DashboardGreeting } from '@/components/dashboard/dashboard-greeting';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

export default function CoachDashboardScreen() {
  const router = useRouter();
  const { data: currentUser } = useCurrentUser();
  const { data: dashboard, refetch, isRefetching } = useCoachDashboard();

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching}
              onRefresh={refetch}
              tintColor={COLORS.primary}
            />
          }
        >
          <DashboardGreeting name={`Coach ${currentUser?.name}`} />

          {/* Stats Overview */}
          <MotiView
            from={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ type: 'timing', duration: 400, delay: 100 }}
          >
            <View style={styles.statsGrid}>
              <StatCard
                icon="people"
                label="Total Clients"
                value={dashboard?.total_clients || 0}
                color={COLORS.primary}
              />
              <StatCard
                icon="checkmark-circle"
                label="Active Clients"
                value={dashboard?.active_clients || 0}
                color={COLORS.success}
              />
              <StatCard
                icon="document-text"
                label="Active Schemas"
                value={dashboard?.active_schemas || 0}
                color={COLORS.warning}
              />
              <StatCard
                icon="barbell"
                label="Total Workouts"
                value={dashboard?.total_workouts || 0}
                color={COLORS.info}
              />
            </View>
          </MotiView>

          {/* Performance Metric */}
          {/* <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 400, delay: 200 }}
          >
            <View style={styles.performanceCard}>
              <View style={styles.performanceHeader}>
                <View>
                  <Text style={styles.performanceTitle}>Client Completion Rate</Text>
                  <Text style={styles.performanceSubtitle}>Average across all clients</Text>
                </View>
                <View style={styles.performanceValueContainer}>
                  <Text style={styles.performanceValue}>
                    {Math.round((dashboard?.average_completion || 0) * 100)}%
                  </Text>
                </View>
              </View>
              <View style={styles.progressBarContainer}>
                <View
                  style={[
                    styles.progressBar,
                    { width: `${(dashboard?.average_completion || 0) * 100}%` },
                  ]}
                />
              </View>
            </View>
          </MotiView> */}

          {/* Quick Actions */}
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 400, delay: 300 }}
          >
            {/* <Text style={styles.sectionTitle}>Quick Actions</Text> */}
            <View style={styles.actionsGrid}>
              <ActionCard
                icon="person-add"
                label="Add Client"
                color={COLORS.primary}
                onPress={() => router.push('/(coach)/clients')}
              />
              <ActionCard
                icon="create"
                label="New Schema"
                color={COLORS.success}
                onPress={() => router.push('/(coach)/schema-templates')}
              />
              <ActionCard
                icon="chatbubbles"
                label="Messages"
                color={COLORS.info}
                onPress={() => router.push('/(coach)/conversations')}
              />
              <ActionCard
                icon="stats-chart"
                label="Analytics"
                color={COLORS.warning}
                onPress={() => {}}
              />
              <ActionCard
                icon="people-outline"
                label="View Clients"
                color="#10B981"
                onPress={() => router.push('/(coach)/clients')}
              />
              <ActionCard
                icon="document-text-outline"
                label="Templates"
                color="#8B5CF6"
                onPress={() => router.push('/(coach)/schema-templates')}
              />
            </View>
          </MotiView>

          {/* Recent Activity */}
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 400, delay: 400 }}
          >
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>Recent Activity</Text>
              <TouchableOpacity onPress={() => {}}>
                <Text style={styles.viewAll}>View All</Text>
              </TouchableOpacity>
            </View>
            
            {dashboard?.recent_activity && dashboard.recent_activity.length > 0 ? (
              dashboard.recent_activity.slice(0, 5).map((activity, index) => (
                <ActivityItem
                  key={index}
                  icon={getActivityIcon(activity.activity_type)}
                  title={activity.description}
                  time={new Date(activity.timestamp).toLocaleDateString()}
                  onPress={() => {}}
                />
              ))
            ) : (
              <View style={styles.emptyState}>
                <Ionicons name="pulse-outline" size={48} color={COLORS.text.tertiary} />
                <Text style={styles.emptyStateText}>No recent activity</Text>
              </View>
            )}
          </MotiView>

          <View style={styles.spacer} />
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const getActivityIcon = (type: string): keyof typeof Ionicons.glyphMap => {
  switch (type) {
    case 'workout_completed':
      return 'checkmark-circle';
    case 'client_joined':
      return 'person-add';
    case 'schema_assigned':
      return 'document-text';
    case 'message_sent':
      return 'chatbubble';
    default:
      return 'pulse';
  }
};

interface StatCardProps {
  icon: keyof typeof Ionicons.glyphMap;
  label: string;
  value: number;
  color: string;
}

const StatCard: React.FC<StatCardProps> = ({ icon, label, value, color }) => (
  <View style={[styles.statCard, { backgroundColor: color + '15' }]}>
    <View style={[styles.statIconContainer, { backgroundColor: color + '15' }]}>
      <Ionicons name={icon} size={24} color={color} />
    </View>
    <View style={styles.statContent}>
      <Text style={styles.statLabel}>{label}</Text>
      <Text style={styles.statValue}>{value}</Text>
    </View>
  </View>
);

interface ActionCardProps {
  icon: keyof typeof Ionicons.glyphMap;
  label: string;
  color: string;
  onPress: () => void;
}

const ActionCard: React.FC<ActionCardProps> = ({ icon, label, color, onPress }) => (
  <MotiView
    from={{ opacity: 0, scale: 0.8 }}
    animate={{ opacity: 1, scale: 1 }}
    transition={{ type: 'spring', damping: 15 }}
    style={styles.actionCardContainer}
  >
    <TouchableOpacity
      style={[styles.actionCard, { backgroundColor: `${color}10` }]}
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View style={[styles.actionIconContainer, { backgroundColor: color }]}>
        <Ionicons name={icon} size={24} color={COLORS.white} />
      </View>
      <Text style={styles.actionLabel} numberOfLines={2}>
        {label}
      </Text>
    </TouchableOpacity>
  </MotiView>
);

interface ActivityItemProps {
  icon: keyof typeof Ionicons.glyphMap;
  title: string;
  time: string;
  onPress: () => void;
}

const ActivityItem: React.FC<ActivityItemProps> = ({ icon, title, time, onPress }) => (
  <TouchableOpacity style={styles.activityItem} onPress={onPress}>
    <View style={styles.activityIconContainer}>
      <Ionicons name={icon} size={20} color={COLORS.primary} />
    </View>
    <View style={styles.activityContent}>
      <Text style={styles.activityTitle}>{title}</Text>
      <Text style={styles.activityTime}>{time}</Text>
    </View>
    <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
  </TouchableOpacity>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
    overflow: 'hidden',
  },
  safeArea: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.lg,
    gap: SPACING.md,
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
    marginTop: SPACING.lg,
  },
  statCard: {
    flex: 1,
    minWidth: '47%',
    borderRadius: BORDER_RADIUS.full,
    padding: SPACING.md,
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    ...SHADOWS.sm,
  },
  statIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
  },
  statContent: {
    flex: 1,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    marginBottom: 4,
  },
  statValue: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  performanceCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginTop: SPACING.lg,
    ...SHADOWS.sm,
  },
  performanceHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  performanceTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  performanceSubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  performanceValueContainer: {
    backgroundColor: COLORS.primary + '20',
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
  },
  performanceValue: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  progressBarContainer: {
    height: 8,
    backgroundColor: COLORS.background.accent,
    borderRadius: BORDER_RADIUS.full,
    overflow: 'hidden',
  },
  progressBar: {
    height: '100%',
    backgroundColor: COLORS.primary,
    borderRadius: BORDER_RADIUS.full,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginTop: SPACING.xl,
    marginBottom: SPACING.md,
  },
  viewAll: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
  actionsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  actionCardContainer: {
    width: '31%',
    minWidth: 100,
  },
  actionCard: {
    aspectRatio: 1,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.sm,
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    ...SHADOWS.sm,
  },
  actionIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: SPACING.xs,
  },
  actionLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
    textAlign: 'center',
  },
  activityItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.xl,
    padding: SPACING.base,
    marginBottom: SPACING.sm,
    ...SHADOWS.sm,
  },
  activityIconContainer: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary + '20',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: SPACING.base,
  },
  activityContent: {
    flex: 1,
  },
  activityTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.inverse,
    marginBottom: 2,
  },
  activityTime: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: SPACING['3xl'],
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    ...SHADOWS.sm,
  },
  emptyStateText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: SPACING.sm,
  },
  spacer: {
    height: SPACING.xl,
  },
});
