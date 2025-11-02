import React, { useMemo } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  Alert,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { useActivePlan, usePlanHistory, useDownloadPlanPDF, useTrackPlanPerformance } from '@/hooks/schema/use-plans';
import { useAuth } from '@/context/auth-context';
import { WorkoutDayCard } from '@/components/schema/workout-day-card';
import { PlanPerformanceModal } from '@/components/schema/plan-performance-modal';
import type { PlanPerformancePayload } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

export default function UserPlansScreen() {
  const { user } = useAuth();
  const router = useRouter();
  const userID = useMemo<number | undefined>(() => {
    if (!user) {
      return undefined;
    }

    if (typeof user.id === 'number') {
      return user.id;
    }

    if (typeof user.id === 'string') {
      const parsed = parseInt(user.id, 10);
      return Number.isFinite(parsed) && parsed > 0 ? parsed : 0;
    }

    return 0;
  }, [user]);

  const { data: activePlan, isLoading: isLoadingActive, refetch: refetchActive } = useActivePlan(userID);
  const { data: planHistory, isLoading: isLoadingHistory, refetch: refetchHistory } = usePlanHistory(userID);
  const downloadPlanMutation = useDownloadPlanPDF();
  const trackPerformanceMutation = useTrackPlanPerformance();

  const [refreshing, setRefreshing] = React.useState(false);
  const [showHistory, setShowHistory] = React.useState(false);
  const [showPerformanceModal, setShowPerformanceModal] = React.useState(false);

  const activePlanRecord = useMemo(() => {
    if (!planHistory || planHistory.length === 0) {
      return null;
    }
    const active = planHistory.find((plan) => plan.is_active);
    return active || planHistory[0];
  }, [planHistory]);

  const onRefresh = async () => {
    setRefreshing(true);
    await Promise.all([refetchActive(), refetchHistory()]);
    setRefreshing(false);
  };

  const handleGenerateNewPlan = () => {
    router.push('/(user)/plan-generator');
  };

  const handleDownloadPlan = () => {
    if (!activePlanRecord) {
      Alert.alert('No Plan', 'Generate a plan first to download a PDF.');
      return;
    }

    if (Platform.OS !== 'web') {
      Alert.alert('Download Unavailable', 'PDF download is currently supported on web builds.');
      return;
    }

    downloadPlanMutation.mutate(activePlanRecord.plan_id);
  };

  const handleOpenPerformanceModal = () => {
    if (!activePlanRecord) {
      Alert.alert('No Active Plan', 'Generate a plan before tracking performance.');
      return;
    }
    setShowPerformanceModal(true);
  };

  const handleSubmitPerformance = async (payload: PlanPerformancePayload) => {
    if (!activePlanRecord) {
      return;
    }

    try {
      await trackPerformanceMutation.mutateAsync({
        planID: activePlanRecord.plan_id,
        data: payload,
      });
      setShowPerformanceModal(false);
      Alert.alert('Logged', 'Thanks for sharing your progress.');
    } catch (error: any) {
      Alert.alert('Unable to Save', error?.message || 'Please try again later.');
    }
  };

  const isLoading = isLoadingActive || isLoadingHistory;

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  return (
    <>
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.contentContainer}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={COLORS.primary} />
        }
      >
      <View style={styles.header}>
        <Text style={styles.headerTitle}>My Workout Plans</Text>
        <TouchableOpacity style={styles.generateButton} onPress={handleGenerateNewPlan}>
          <Ionicons name="add-circle" size={24} color={COLORS.text.primary} />
          <Text style={styles.generateButtonText}>Generate New Plan</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Active Plan</Text>
          {activePlan && (
            <View style={styles.statusBadge}>
              <View style={styles.activeDot} />
              <Text style={styles.statusText}>Active</Text>
            </View>
          )}
        </View>

        {activePlan ? (
          <View style={styles.activePlanCard}>
            <View style={styles.planHeader}>
              <View>
                <Text style={styles.planName}>
                  Week of {new Date(activePlan.week_start).toLocaleDateString()}
                </Text>
                <Text style={styles.planMeta}>
                  {activePlan.workouts.length} workouts • Active
                </Text>
              </View>
              <TouchableOpacity style={styles.actionButton}>
                <Ionicons name="ellipsis-horizontal" size={24} color={COLORS.text.auth.secondary} />
              </TouchableOpacity>
            </View>

            <View style={styles.progressContainer}>
              <View style={styles.progressHeader}>
                <Text style={styles.progressLabel}>This Week's Progress</Text>
                <Text style={styles.progressPercentage}>
                  0/{activePlan.workouts.length}
                </Text>
              </View>
              <View style={styles.progressBar}>
                <View style={[styles.progressFill, { width: '0%' }]} />
              </View>
            </View>

            <View style={styles.workoutsPreview}>
              {activePlan.workouts.slice(0, 3).map((workout) => (
                <View key={workout.workout_id} style={styles.workoutPreviewItem}>
                  <Ionicons name="barbell" size={16} color={COLORS.primary} />
                  <Text style={styles.workoutPreviewText}>
                    Day {workout.day_of_week}: {workout.focus}
                  </Text>
                </View>
              ))}
              {activePlan.workouts.length > 3 && (
                <Text style={styles.moreText}>
                  +{activePlan.workouts.length - 3} more workouts
                </Text>
              )}
            </View>

            <View style={styles.planActions}>
              <TouchableOpacity style={[styles.planActionButton, styles.primaryAction]}>
                <Ionicons name="play" size={18} color={COLORS.text.primary} />
                <Text style={styles.planActionText}>View Full Plan</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={styles.planActionButton}
                onPress={handleDownloadPlan}
                disabled={downloadPlanMutation.isPending}
              >
                <Ionicons
                  name={downloadPlanMutation.isPending ? 'cloud-download' : 'download-outline'}
                  size={18}
                  color={COLORS.primary}
                />
                <Text style={[styles.planActionText, { color: COLORS.primary }]}>
                  {downloadPlanMutation.isPending ? 'Preparing...' : 'Download PDF'}
                </Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.planActionButton, styles.performanceAction]}
                onPress={handleOpenPerformanceModal}
                disabled={trackPerformanceMutation.isPending}
              >
                <Ionicons
                  name="stats-chart"
                  size={18}
                  color={COLORS.text.auth.primary}
                />
                <Text style={styles.planActionText}>Log Performance</Text>
              </TouchableOpacity>
            </View>
          </View>
        ) : (
          <View style={styles.emptyState}>
            <Ionicons name="calendar-outline" size={60} color={COLORS.text.tertiary} />
            <Text style={styles.emptyStateTitle}>No Active Plan</Text>
            <Text style={styles.emptyStateText}>
              Generate a personalized workout plan to get started
            </Text>
            <TouchableOpacity style={styles.emptyStateButton} onPress={handleGenerateNewPlan}>
              <Ionicons name="flash" size={20} color={COLORS.text.primary} />
              <Text style={styles.emptyStateButtonText}>Generate Plan</Text>
            </TouchableOpacity>
          </View>
        )}
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Plan History</Text>
          <TouchableOpacity onPress={() => setShowHistory(!showHistory)}>
            <Ionicons
              name={showHistory ? 'chevron-up' : 'chevron-down'}
              size={24}
              color={COLORS.text.auth.secondary}
            />
          </TouchableOpacity>
        </View>

        {showHistory && (
          <View style={styles.historyList}>
            {planHistory && planHistory.length > 0 ? (
              planHistory.map((plan) => (
                <View key={plan.plan_id} style={styles.historyCard}>
                  <View style={styles.historyCardHeader}>
                    <View style={styles.historyCardIcon}>
                      <Ionicons name="document-text" size={20} color={COLORS.primary} />
                    </View>
                    <View style={styles.historyCardInfo}>
                      <Text style={styles.historyCardTitle}>
                        Plan #{plan.plan_id}
                      </Text>
                      <Text style={styles.historyCardMeta}>
                        Generated {new Date(plan.generated_at).toLocaleDateString()}
                      </Text>
                    </View>
                    <View style={styles.historyCardStats}>
                      <Text style={styles.historyCardEffectiveness}>
                        {Math.round(plan.effectiveness)}%
                      </Text>
                      <Text style={styles.historyCardEffectivenessLabel}>
                        Effectiveness
                      </Text>
                    </View>
                  </View>
                  <TouchableOpacity style={styles.historyCardAction}>
                    <Text style={styles.historyCardActionText}>View Details</Text>
                    <Ionicons name="chevron-forward" size={16} color={COLORS.primary} />
                  </TouchableOpacity>
                </View>
              ))
            ) : (
              <View style={styles.emptyHistoryState}>
                <Ionicons name="time-outline" size={40} color={COLORS.text.tertiary} />
                <Text style={styles.emptyHistoryText}>No plan history yet</Text>
              </View>
            )}
          </View>
        )}
      </View>

        <View style={styles.statsContainer}>
          <View style={styles.statCard}>
            <Ionicons name="trophy" size={24} color={COLORS.primary} />
            <Text style={styles.statValue}>{planHistory?.length || 0}</Text>
            <Text style={styles.statLabel}>Total Plans</Text>
          </View>
          <View style={styles.statCard}>
            <Ionicons name="calendar" size={24} color={COLORS.primary} />
            <Text style={styles.statValue}>0</Text>
            <Text style={styles.statLabel}>Weeks Active</Text>
          </View>
          <View style={styles.statCard}>
            <Ionicons name="trending-up" size={24} color={COLORS.primary} />
            <Text style={styles.statValue}>-</Text>
            <Text style={styles.statLabel}>Avg. Performance</Text>
          </View>
        </View>
      </ScrollView>

      <PlanPerformanceModal
        visible={showPerformanceModal}
        onClose={() => setShowPerformanceModal(false)}
        onSubmit={handleSubmitPerformance}
        isSubmitting={trackPerformanceMutation.isPending}
      />
    </>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.background.auth,
  },
  contentContainer: {
    padding: SPACING.base,
  },
  header: {
    marginBottom: SPACING.lg,
  },
  headerTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.md,
  },
  generateButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.primary,
    paddingVertical: SPACING.md,
    paddingHorizontal: SPACING.lg,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.sm,
    ...SHADOWS.base,
  },
  generateButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.success,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    gap: 4,
  },
  activeDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: COLORS.text.primary,
  },
  statusText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  activePlanCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.full,
    padding: SPACING.base,
    ...SHADOWS.base,
  },
  planHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: SPACING.md,
  },
  planName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: 4,
  },
  planMeta: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  actionButton: {
    padding: SPACING.xs,
  },
  progressContainer: {
    marginBottom: SPACING.base,
  },
  progressHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: SPACING.xs,
  },
  progressLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.secondary,
  },
  progressPercentage: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  progressBar: {
    height: 8,
    backgroundColor: COLORS.background.primary,
    borderRadius: BORDER_RADIUS.full,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.primary,
    borderRadius: BORDER_RADIUS.full,
  },
  workoutsPreview: {
    marginBottom: SPACING.base,
  },
  workoutPreviewItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginBottom: SPACING.xs,
  },
  workoutPreviewText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.secondary,
  },
  moreText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    fontStyle: 'italic',
    marginTop: SPACING.xs,
  },
  planActions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  planActionButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.sm,
    paddingHorizontal: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    backgroundColor: COLORS.background.primary,
    gap: SPACING.xs,
  },
  primaryAction: {
    backgroundColor: COLORS.primary,
  },
  performanceAction: {
    backgroundColor: COLORS.background.primary,
  },
  planActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  emptyState: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.xl,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  emptyStateTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.md,
    marginBottom: SPACING.xs,
  },
  emptyStateText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginBottom: SPACING.lg,
  },
  emptyStateButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingVertical: SPACING.md,
    paddingHorizontal: SPACING.lg,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.xs,
  },
  emptyStateButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  historyList: {
    gap: SPACING.sm,
  },
  historyCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    ...SHADOWS.sm,
  },
  historyCardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  historyCardIcon: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS['2xl'],
    backgroundColor: COLORS.background.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.sm,
  },
  historyCardInfo: {
    flex: 1,
  },
  historyCardTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 2,
  },
  historyCardMeta: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  historyCardStats: {
    alignItems: 'flex-end',
  },
  historyCardEffectiveness: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.success,
  },
  historyCardEffectivenessLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  historyCardAction: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.sm,
    gap: 4,
  },
  historyCardActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
  emptyHistoryState: {
    alignItems: 'center',
    paddingVertical: SPACING.xl,
  },
  emptyHistoryText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: SPACING.sm,
  },
  statsContainer: {
    flexDirection: 'row',
    gap: SPACING.sm,
    marginBottom: SPACING.xl,
  },
  statCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.xs,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
    textAlign: 'center',
  },
});
