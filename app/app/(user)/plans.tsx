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
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import {
  useActivePlan,
  usePlanHistory,
  useDownloadPlanPDF,
  useTrackPlanPerformance,
  usePlanEffectiveness,
  useAdaptationHistory,
  useRegeneratePlan,
} from '@/hooks/schema/use-plans';
import { useAuth } from '@/context/auth-context';
import { PlanPerformanceModal } from '@/components/schema/plan-performance-modal';
import { PlanDetailModal } from '@/components/schema/plan-detail-modal';
import type { PlanPerformancePayload } from '@/types/schema';
import { APIError } from '@/api/client';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

const formatDate = (value?: string | null) => {
  if (!value) {
    return 'Unknown';
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return 'Unknown';
  }
  return parsed.toLocaleDateString();
};

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
  const requestRegenerationMutation = useRegeneratePlan();

  const [refreshing, setRefreshing] = React.useState(false);
  const [showHistory, setShowHistory] = React.useState(false);
  const [showPerformanceModal, setShowPerformanceModal] = React.useState(false);
  const [showPlanModal, setShowPlanModal] = React.useState(false);

  const activePlanRecord = useMemo(() => {
    if (!planHistory || planHistory.length === 0) {
      return null;
    }
    const active = planHistory.find((plan) => plan.is_active);
    return active || planHistory[0];
  }, [planHistory]);

  const { data: planEffectiveness, isLoading: isLoadingEffectiveness } = usePlanEffectiveness(activePlanRecord?.plan_id);
  const { data: adaptationHistory, isLoading: isLoadingAdaptations } = useAdaptationHistory(userID);

  const onRefresh = async () => {
    setRefreshing(true);
    await Promise.all([refetchActive(), refetchHistory()]);
    setRefreshing(false);
  };

  const handleGenerateNewPlan = () => {
    router.push('/(user)/plan-generator');
  };

  const handleViewFullPlan = () => {
    if (!activePlan) {
      Alert.alert('No Plan', 'Generate a plan to view the full workout breakdown.');
      return;
    }
    setShowPlanModal(true);
  };

  const handleRequestRegeneration = (planID: number) => {
    if (!planID || planID <= 0) {
      Alert.alert('Unavailable', 'This plan cannot be updated yet.');
      return;
    }

    const reason = 'User requested plan adjustments via the mobile app';

    Alert.alert(
      'Request Plan Adjustments',
      'We will notify your coach to review and adapt this plan.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Request',
          style: 'default',
          onPress: () =>
            requestRegenerationMutation.mutate(
              { planID, reason, userID },
              {
                onSuccess: () => {
                  Alert.alert('Request sent', 'Your coach will review the plan shortly.');
                },
                onError: (error) => {
                  let message: string | undefined;
                  if (error instanceof Error) {
                    message = error.message;
                  } else if (error && typeof (error as APIError).message === 'string') {
                    message = (error as APIError).message;
                  }
                  Alert.alert('Unable to submit request', message || 'Please try again later.');
                },
              }
            ),
        },
      ]
    );
  };

  const handleDownloadPlan = () => {
    if (!activePlanRecord) {
      Alert.alert('No Plan', 'Generate a plan first to download a PDF.');
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
  const workouts = Array.isArray(activePlan?.workouts) ? activePlan?.workouts ?? [] : [];
  const workoutCount = workouts.length;

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
                  {workoutCount} workouts • Active
                </Text>
                {isLoadingEffectiveness ? (
                  <Text style={styles.planEffectivenessLoading}>Refreshing effectiveness...</Text>
                ) : typeof planEffectiveness?.effectiveness_score === 'number' ? (
                  <Text style={styles.planEffectivenessText}>
                    Effectiveness: {Math.round(planEffectiveness.effectiveness_score)}%
                  </Text>
                ) : null}
              </View>
              {/* <TouchableOpacity style={styles.actionButton}>
                <Ionicons name="ellipsis-horizontal" size={24} color={COLORS.text.auth.secondary} />
              </TouchableOpacity> */}
            </View>

            <View style={styles.progressContainer}>
              <View style={styles.progressHeader}>
                <Text style={styles.progressLabel}>This Week's Progress</Text>
                <Text style={styles.progressPercentage}>
                  0/{workoutCount}
                </Text>
              </View>
              <View style={styles.progressBar}>
                <View style={[styles.progressFill, { width: '0%' }]} />
              </View>
            </View>

            <View style={styles.workoutsPreview}>
              {workouts.slice(0, 3).map((workout) => (
                <View key={workout.workout_id} style={styles.workoutPreviewItem}>
                  <Ionicons name="barbell" size={16} color={COLORS.primary} />
                  <Text style={styles.workoutPreviewText}>
                    Day {workout.day_of_week}: {workout.focus}
                  </Text>
                </View>
              ))}
              {workoutCount > 3 && (
                <Text style={styles.moreText}>
                  +{workoutCount - 3} more workouts
                </Text>
              )}
            </View>

            <View style={styles.planActions}>
              <TouchableOpacity style={[styles.planActionButton, styles.primaryAction]} onPress={handleViewFullPlan}>
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
                style={[styles.planActionButton, styles.secondaryAction]}
                onPress={() => handleRequestRegeneration(activePlanRecord?.plan_id ?? 0)}
                disabled={!activePlanRecord || requestRegenerationMutation.isPending}
              >
                <Ionicons
                  name={requestRegenerationMutation.isPending ? 'time' : 'refresh'}
                  size={18}
                  color={COLORS.text.auth.primary}
                />
                <Text style={styles.planActionText}>
                  {requestRegenerationMutation.isPending ? 'Submitting...' : 'Request Adjustments'}
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
              planHistory.map((plan) => {
                const isDownloadingThisPlan =
                  downloadPlanMutation.isPending && downloadPlanMutation.variables === plan.plan_id;
                const isRequestingThisPlan =
                  requestRegenerationMutation.isPending &&
                  requestRegenerationMutation.variables?.planID === plan.plan_id;

                return (
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
                          Generated {formatDate(plan.generated_at)}
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

                    <View style={styles.historyCardActions}>
                      <TouchableOpacity
                        style={styles.historyActionButton}
                        onPress={() => downloadPlanMutation.mutate(plan.plan_id)}
                        disabled={downloadPlanMutation.isPending}
                      >
                        <Ionicons
                          name={isDownloadingThisPlan ? 'cloud-download' : 'download-outline'}
                          size={16}
                          color={COLORS.primary}
                        />
                        <Text style={[styles.historyActionText, { color: COLORS.primary }]}>
                          {isDownloadingThisPlan ? 'Preparing...' : 'Download'}
                        </Text>
                      </TouchableOpacity>
                      <TouchableOpacity
                        style={styles.historyActionButton}
                        onPress={() => handleRequestRegeneration(plan.plan_id)}
                        disabled={isRequestingThisPlan}
                      >
                        <Ionicons
                          name={isRequestingThisPlan ? 'time' : 'refresh'}
                          size={16}
                          color={COLORS.text.auth.primary}
                        />
                        <Text style={styles.historyActionText}>
                          {isRequestingThisPlan ? 'Submitting...' : 'Request Updates'}
                        </Text>
                      </TouchableOpacity>
                    </View>
                  </View>
                );
              })
            ) : (
              <View style={styles.emptyHistoryState}>
                <Ionicons name="time-outline" size={40} color={COLORS.text.tertiary} />
                <Text style={styles.emptyHistoryText}>No plan history yet</Text>
              </View>
            )}
          </View>
        )}
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Adaptation History</Text>
        </View>

        {isLoadingAdaptations ? (
          <View style={styles.listLoadingRow}>
            <ActivityIndicator size="small" color={COLORS.primary} />
            <Text style={styles.historyCardMeta}>Loading adaptations...</Text>
          </View>
        ) : adaptationHistory && adaptationHistory.length > 0 ? (
          <View style={styles.adaptationList}>
            {adaptationHistory.map((adaptation) => (
              <View key={adaptation.adaptation_id} style={styles.adaptationCard}>
                <View style={styles.adaptationHeader}>
                  <Ionicons name="sparkles" size={18} color={COLORS.primary} />
                  <View style={{ flex: 1 }}>
                    <Text style={styles.adaptationReason}>{adaptation.reason}</Text>
                    <Text style={styles.adaptationMeta}>
                      {formatDate(adaptation.adaptation_date)} • Plan #{adaptation.plan_id}
                    </Text>
                  </View>
                </View>
                <Text style={styles.adaptationDetails}>
                  Triggered by: {adaptation.trigger || 'system'}
                </Text>
              </View>
            ))}
          </View>
        ) : (
          <View style={styles.emptyHistoryState}>
            <Ionicons name="leaf-outline" size={32} color={COLORS.text.tertiary} />
            <Text style={styles.emptyHistoryText}>No adaptations recorded yet</Text>
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
      <PlanDetailModal
        visible={showPlanModal}
        onClose={() => setShowPlanModal(false)}
        plan={activePlan ?? null}
        isLoading={isLoadingActive}
        effectiveness={planEffectiveness?.effectiveness_score}
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
    borderRadius: BORDER_RADIUS['2xl'],
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
  planEffectivenessLoading: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  planEffectivenessText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.auth.secondary,
    marginTop: 2,
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
    paddingHorizontal: SPACING.lg,
    borderRadius: BORDER_RADIUS.md,
    backgroundColor: COLORS.background.primary,
    gap: SPACING.xs,
  },
  primaryAction: {
    backgroundColor: COLORS.primary,
  },
  secondaryAction: {
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  performanceAction: {
    backgroundColor: COLORS.background.auth,
  },
  planActionText: {
    fontSize: 10,
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
  historyCardActions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
    marginTop: SPACING.sm,
  },
  historyActionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    paddingVertical: SPACING.xs,
    paddingHorizontal: SPACING.base,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.primary,
  },
  historyActionText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
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
  listLoadingRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.sm,
    paddingVertical: SPACING.md,
  },
  adaptationList: {
    gap: SPACING.sm,
  },
  adaptationCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    ...SHADOWS.sm,
  },
  adaptationHeader: {
    flexDirection: 'row',
    gap: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  adaptationReason: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  adaptationMeta: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  adaptationDetails: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.auth.secondary,
    lineHeight: 18,
  },
});
