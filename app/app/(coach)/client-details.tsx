import React, { useMemo, useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

import {
  useClientDetails,
  useClientProgress,
  useClientWorkouts,
  useClientSchemas,
} from '@/hooks/schema/use-coach';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

const DAY_LABELS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

const StatPill: React.FC<{ icon: keyof typeof Ionicons.glyphMap; label: string; value: string }> = ({
  icon,
  label,
  value,
}) => (
  <View style={styles.statPill}>
    <Ionicons name={icon} size={16} color={COLORS.primary} />
    <View style={styles.statPillText}>
      <Text style={styles.statPillLabel}>{label}</Text>
      <Text style={styles.statPillValue}>{value}</Text>
    </View>
  </View>
);

const DayDistributionChart: React.FC<{ data: number[] }> = ({ data }) => {
  const max = Math.max(...data, 1);
  const BAR_MAX_HEIGHT = 120;
  return (
    <View style={styles.chartContainer}>
      {data.map((count, index) => {
        const scaledHeight = Math.max((count / max) * BAR_MAX_HEIGHT, 4);
        return (
          <View key={DAY_LABELS[index]} style={styles.chartBarWrapper}>
            <View style={[styles.chartBar, { height: scaledHeight }]} />
            <Text style={styles.chartLabel}>{DAY_LABELS[index]}</Text>
          </View>
        );
      })}
    </View>
  );
};

export default function ClientDetailsScreen() {
  const params = useLocalSearchParams<{ userId?: string }>();
  const router = useRouter();
  const [refreshing, setRefreshing] = useState(false);

  const clientId = useMemo(() => {
    if (!params.userId) {
      return 0;
    }
    const parsed = parseInt(params.userId, 10);
    return Number.isFinite(parsed) ? parsed : 0;
  }, [params.userId]);

  const clientDetailsQuery = useClientDetails(clientId);
  const clientProgressQuery = useClientProgress(clientId);
  const clientWorkoutsQuery = useClientWorkouts(clientId);
  const clientSchemasQuery = useClientSchemas(clientId);

  const refetchAll = async () => {
    setRefreshing(true);
    await Promise.all([
      clientDetailsQuery.refetch(),
      clientProgressQuery.refetch(),
      clientWorkoutsQuery.refetch(),
      clientSchemasQuery.refetch(),
    ]);
    setRefreshing(false);
  };

  const isLoading =
    clientDetailsQuery.isLoading ||
    clientProgressQuery.isLoading ||
    clientWorkoutsQuery.isLoading ||
    clientSchemasQuery.isLoading;

  const client = clientDetailsQuery.data;
  const progress = clientProgressQuery.data;
  const workouts = clientWorkoutsQuery.data?.workouts || [];
  const schemas = clientSchemasQuery.data?.schemas || [];

  const dayDistribution = useMemo(() => {
    const base = new Array(7).fill(0);
    workouts.forEach((workout) => {
      const day = workout.day_of_week ?? 0;
      if (day >= 0 && day < base.length) {
        base[day] += 1;
      }
    });
    return base;
  }, [workouts]);

  if (!clientId) {
    return (
      <View style={styles.centered}> 
        <Text style={styles.errorText}>Missing client identifier.</Text>
      </View>
    );
  }

  if (isLoading) {
    return (
      <View style={styles.centered}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (!client) {
    return (
      <View style={styles.centered}>
        <Text style={styles.errorText}>Unable to load client details. Try again later.</Text>
      </View>
    );
  }

  const fullName = `${client.first_name} ${client.last_name}`;
  const completionPercentage = Math.round((client.completion_rate || 0) * 100);
  const lastWorkout = client.last_workout_date
    ? new Date(client.last_workout_date).toLocaleDateString()
    : 'No recent activity';

  const handleCreateSchema = () => {
    router.push({
      pathname: '/(coach)/schema-create',
      params: { userId: client.user_id.toString() },
    });
  };

  const activeSchema = schemas.find((schema) => schema.active);

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.contentContainer}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={refetchAll}
          tintColor={COLORS.primary}
        />
      }
    >
      {/* <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
          <Ionicons name="arrow-back" size={22} color={COLORS.text.primary} />
        </TouchableOpacity>
        <View style={styles.headerInfo}>
          <Text style={styles.clientName}>{fullName}</Text>
          <Text style={styles.clientSubtitle}>{client}</Text>
        </View>
        <TouchableOpacity style={styles.assignButton} onPress={handleCreateSchema}>
          <Ionicons name="add-circle" size={24} color={COLORS.text.primary} />
          <Text style={styles.assignButtonText}>Assign Schema</Text>
        </TouchableOpacity>
      </View> */}

      <View style={styles.overviewCard}>
        <Text style={styles.sectionTitle}>Overview</Text>
        <View style={styles.quickStats}>
          <StatPill icon="barbell" label="Total Workouts" value={`${client.total_workouts}`} />
          <StatPill icon="flame" label="Current Streak" value={`${client.current_streak} days`} />
          <StatPill icon="checkmark-circle" label="Completion" value={`${completionPercentage}%`} />
        </View>
        <View style={styles.overviewMeta}>
          <View style={styles.metaRow}>
            <Text style={styles.metaLabel}>Fitness Level</Text>
            <Text style={styles.metaValue}>{client.fitness_level || 'Unspecified'}</Text>
          </View>
          <View style={styles.metaRow}>
            <Text style={styles.metaLabel}>Last Workout</Text>
            <Text style={styles.metaValue}>{lastWorkout}</Text>
          </View>
          <View style={styles.metaRow}>
            <Text style={styles.metaLabel}>Active Goals</Text>
            <Text style={styles.metaValue}>{client.active_goals}</Text>
          </View>
        </View>
      </View>

      <View style={styles.sectionCard}>
        <View style={styles.sectionHeader}>
          <View>
            <Text style={styles.sectionTitle}>Weekly Distribution</Text>
            <Text style={styles.sectionSubtitle}>Completed sessions by day</Text>
          </View>
          <Ionicons name="bar-chart" size={20} color={COLORS.primary} />
        </View>
        <DayDistributionChart data={dayDistribution} />
      </View>

      <View style={styles.sectionCard}>
        <View style={styles.sectionHeader}>
          <View>
            <Text style={styles.sectionTitle}>Performance Snapshot</Text>
            <Text style={styles.sectionSubtitle}>Live adherence insights</Text>
          </View>
          <Ionicons name="pulse" size={20} color={COLORS.success} />
        </View>
        <View style={styles.performanceGrid}>
          <View style={styles.performanceTile}>
            <Text style={styles.performanceValue}>{completionPercentage}%</Text>
            <Text style={styles.performanceLabel}>Plan Completion</Text>
          </View>
          <View style={styles.performanceTile}>
            <Text style={styles.performanceValue}>{progress?.current_streak || 0}</Text>
            <Text style={styles.performanceLabel}>Current Streak</Text>
          </View>
          <View style={styles.performanceTile}>
            <Text style={styles.performanceValue}>{progress?.total_workouts || 0}</Text>
            <Text style={styles.performanceLabel}>Workouts Logged</Text>
          </View>
        </View>
      </View>

      <View style={styles.sectionCard}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Assigned Schemas</Text>
          <TouchableOpacity onPress={() => Alert.alert('Coming Soon', 'Schema comparison view is coming soon.') }>
            <Text style={styles.viewAll}>Compare</Text>
          </TouchableOpacity>
        </View>

        {schemas.length === 0 ? (
          <View style={styles.emptyState}>
            <Ionicons name="document-text" size={48} color={COLORS.text.tertiary} />
            <Text style={styles.emptyTitle}>No schemas assigned</Text>
            <Text style={styles.emptySubtitle}>Create a custom plan tailored to this client.</Text>
          </View>
        ) : (
          schemas.slice(0, 4).map((schema) => {
            const startDate = schema.week_start
              ? new Date(schema.week_start).toLocaleDateString()
              : 'Unknown';
            const workoutCount = schema.workouts?.length || 0;
            const isActive = activeSchema?.schema_id === schema.schema_id;
            const rawTitle = (schema.metadata as any)?.custom_data?.title;
            const schemaTitle = typeof rawTitle === 'string' && rawTitle.trim().length > 0
              ? rawTitle
              : `Schema #${schema.schema_id}`;

            return (
              <View key={schema.schema_id} style={[styles.schemaCard, isActive && styles.schemaCardActive]}>
                <View style={styles.schemaHeader}>
                  <Text style={styles.schemaName}>{schemaTitle}</Text>
                  {isActive ? (
                    <View style={styles.activeBadge}>
                      <View style={styles.activeDot} />
                      <Text style={styles.activeBadgeText}>Active</Text>
                    </View>
                  ) : null}
                </View>
                <View style={styles.schemaMeta}>
                  <Text style={styles.schemaMetaText}>{workoutCount} workouts • Starts {startDate}</Text>
                  <TouchableOpacity
                    style={styles.schemaAction}
                    onPress={() => {
                      Alert.alert(
                        'Clone Schema',
                        'Clone and assign this schema to another client?',
                        [
                          { text: 'Cancel', style: 'cancel' },
                          {
                            text: 'Clone',
                            onPress: () =>
                              router.push({
                                pathname: '/(coach)/schema-create',
                                params: { userId: client.user_id.toString(), schemaId: schema.schema_id.toString() },
                              }),
                          },
                        ]
                      );
                    }}
                  >
                    <Text style={styles.schemaActionText}>Clone</Text>
                    <Ionicons name="chevron-forward" size={16} color={COLORS.primary} />
                  </TouchableOpacity>
                </View>
              </View>
            );
          })
        )}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  contentContainer: {
    padding: SPACING.lg,
    paddingBottom: SPACING['5xl'],
  },
  centered: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.background.auth,
  },
  errorText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING['3xl'],
  },
  backButton: {
    width: 44,
    height: 44,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: SPACING.md,
    ...SHADOWS.sm,
  },
  overviewCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING['3xl'],
    ...SHADOWS.sm,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  quickStats: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
    marginTop: SPACING.md,
  },
  statPill: {
    flex: 1,
    minWidth: 100,
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    backgroundColor: COLORS.background.card,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 12 
    },
    shadowOpacity: 0.35,
    shadowRadius: 24,
    borderRadius: BORDER_RADIUS.full,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
  },
  statPillText: {
    flex: 1,
  },
  statPillLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  statPillValue: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  overviewMeta: {
    marginTop: SPACING['2xl'],
    gap: SPACING.md,
  },
  metaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  metaLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  metaValue: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
  },
  sectionCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING['3xl'],
    ...SHADOWS.sm,
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
  },
  sectionSubtitle: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 4,
  },
  chartContainer: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    justifyContent: 'space-between',
    gap: SPACING.md,
    paddingVertical: SPACING.base,
  },
  chartBarWrapper: {
    flex: 1,
    alignItems: 'center',
  },
  chartBar: {
    width: 18,
    borderRadius: BORDER_RADIUS.md,
    backgroundColor: COLORS.primary,
  },
  chartLabel: {
    marginTop: SPACING.sm,
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  performanceGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  performanceTile: {
    flex: 1,
    minWidth: 100,
    backgroundColor: COLORS.background.card,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 12 
    },
    shadowOpacity: 0.35,
    shadowRadius: 24,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    alignItems: 'center',
  },
  performanceValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  performanceLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 4,
  },
  viewAll: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: SPACING['3xl'],
    gap: SPACING.sm,
  },
  emptyTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
  schemaCard: {
    backgroundColor: COLORS.background.primary,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  schemaCardActive: {
    borderColor: COLORS.primary,
  },
  schemaHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: SPACING.sm,
  },
  schemaName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  activeBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.success,
    borderRadius: BORDER_RADIUS.full,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
  },
  activeDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: COLORS.text.primary,
  },
  activeBadgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  schemaMeta: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  schemaMetaText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  schemaAction: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  schemaActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
});
