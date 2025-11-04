import React from 'react';
import {
  Modal,
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { GeneratedPlan, GeneratedPlanWorkout } from '@/types/schema';
import { WorkoutDayCard } from '@/components/schema/workout-day-card';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';


interface PlanDetailModalProps {
  visible: boolean;
  plan: GeneratedPlan | null;
  onClose: () => void;
  isLoading?: boolean;
  effectiveness?: number;
}

export const PlanDetailModal: React.FC<PlanDetailModalProps> = ({
  visible,
  plan,
  onClose,
  isLoading = false,
  effectiveness,
}) => {
  const workouts = React.useMemo<GeneratedPlanWorkout[]>(() => {
    if (!plan || !Array.isArray(plan.workouts)) {
      return [];
    }
    return plan.workouts;
  }, [plan]);

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent
      onRequestClose={onClose}
    >
      <View style={styles.backdrop}>
        <View style={styles.modalContainer}>
          <View style={styles.modalHeader}>
            <View style={styles.titleGroup}>
              <Text style={styles.title}>Weekly Plan Overview</Text>
              {plan?.week_start ? (
                <Text style={styles.subtitle}>
                  Starting {new Date(plan.week_start).toLocaleDateString()}
                </Text>
              ) : null}
            </View>
            <TouchableOpacity onPress={onClose} style={styles.closeButton} accessibilityRole="button">
              <Ionicons name="close" size={22} color={COLORS.text.auth.primary} />
            </TouchableOpacity>
          </View>

          {typeof effectiveness === 'number' ? (
            <View style={styles.effectivenessBadge}>
              <Ionicons name="sparkles" size={18} color={COLORS.text.primary} />
              <Text style={styles.effectivenessText}>
                Effectiveness score: {Math.round(effectiveness)}%
              </Text>
            </View>
          ) : null}

          {isLoading ? (
            <View style={styles.loadingState}>
              <ActivityIndicator size="large" color={COLORS.primary} />
              <Text style={styles.loadingLabel}>Loading plan details...</Text>
            </View>
          ) : plan && workouts.length > 0 ? (
            <ScrollView style={styles.scrollArea} contentContainerStyle={styles.scrollContent}>
              {workouts.map((workout) => (
                <WorkoutDayCard
                  key={`${workout.plan_id}-${workout.day_index}-${workout.workout_id ?? workout.plan_id}`}
                  dayOfWeek={workout.day_index}
                  workout={workout}
                  isRestDay={workout.is_rest}
                />
              ))}
            </ScrollView>
          ) : (
            <View style={styles.emptyState}>
              <Ionicons name="calendar-outline" size={60} color={COLORS.text.tertiary} />
              <Text style={styles.emptyTitle}>No workouts found</Text>
              <Text style={styles.emptySubtitle}>
                Generate a plan to review the complete workout breakdown.
              </Text>
            </View>
          )}
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  backdrop: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.45)',
    justifyContent: 'flex-end',
  },
  modalContainer: {
    maxHeight: '90%',
    backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['3xl'],
    borderTopRightRadius: BORDER_RADIUS['3xl'],
    paddingTop: SPACING.lg,
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING['3xl'],
    ...SHADOWS.lg,
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
  },
  titleGroup: {
    gap: 4,
  },
  title: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  subtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  closeButton: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },

  scrollContent: {
    overflow: 'visible',
    paddingBottom: SPACING.md,
  },
  effectivenessBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.xs,
    marginBottom: SPACING.md,
  },
  effectivenessText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  scrollArea: {
    maxHeight: '85%',
  },
  
  loadingState: {
    alignItems: 'center',
    gap: SPACING.md,
    paddingVertical: SPACING['2xl'],
  },
  loadingLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  emptyState: {
    alignItems: 'center',
    gap: SPACING.sm,
    paddingVertical: SPACING['2xl'],
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
    paddingHorizontal: SPACING.lg,
  },
});
