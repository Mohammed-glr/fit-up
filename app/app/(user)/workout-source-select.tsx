import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Modal,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Stack, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { useTodayWorkout } from '@/hooks/user/use-today-workout';
import { useUserSchemas, useSchemaWithWorkouts } from '@/hooks/schema/use-user-schemas';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export default function WorkoutSourceSelectScreen() {
  const router = useRouter();
  const { data: currentUser } = useCurrentUser();
  const { data: todayWorkout, isLoading: loadingPlan } = useTodayWorkout();
  const { data: schemas, isLoading: loadingSchemas } = useUserSchemas(currentUser?.id || '');

  const [showDaySelector, setShowDaySelector] = useState(false);
  const [selectedSchemaId, setSelectedSchemaId] = useState<number | null>(null);

  const { data: selectedSchemaData, isLoading: loadingSchemaData } = useSchemaWithWorkouts(
    selectedSchemaId || 0
  );

  const activeSchema = schemas?.find(s => s.active);

  const handleSelectAIPlan = () => {
    router.push('/(user)/workout-session?source=ai');
  };

  const handleSelectSchema = (schemaId: number) => {
    setSelectedSchemaId(schemaId);
    setShowDaySelector(true);
  };

  const handleSelectDay = (dayIndex: number) => {
    if (selectedSchemaId) {
      setShowDaySelector(false);
      router.push(`/(user)/workout-session?source=schema&schemaId=${selectedSchemaId}&dayIndex=${dayIndex}`);
    }
  };

  if (loadingPlan || loadingSchemas) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
        <Text style={styles.loadingText}>Loading workout options...</Text>
      </View>
    );
  }

  const hasAIPlan = todayWorkout && todayWorkout.exercises && todayWorkout.exercises.length > 0;
  const hasActiveSchema = !!activeSchema;

  if (!hasAIPlan && !hasActiveSchema) {
    return (
      <SafeAreaView style={styles.container}>
        <Stack.Screen options={{ title: 'Start Workout', headerShown: true }} />
        <View style={styles.emptyContainer}>
          <Ionicons name="barbell-outline" size={64} color={COLORS.text.tertiary} />
          <Text style={styles.emptyTitle}>No Workouts Available</Text>
          <Text style={styles.emptySubtitle}>
            Create an AI plan or ask your coach to assign a training schema
          </Text>
          <TouchableOpacity
            style={styles.createPlanButton}
            onPress={() => router.push('/(user)/plan-generator')}
          >
            <Text style={styles.createPlanButtonText}>Create Plan</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <Stack.Screen options={{ title: 'Choose Workout Source', headerShown: true }} />

      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <Text style={styles.headerTitle}>Select Your Workout</Text>
        <Text style={styles.headerSubtitle}>
          Choose between your generated plan or coach-created schema
        </Text>

        {hasAIPlan && (
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 300 }}
          >
            <TouchableOpacity
              style={styles.sourceCard}
              onPress={handleSelectAIPlan}
              activeOpacity={0.7}
            >
              <View style={styles.sourceIconContainer}>
                <Ionicons name="flash" size={32} color={COLORS.primary} />
              </View>

              <View style={styles.sourceContent}>
                <View style={styles.sourceHeader}>
                  <Text style={styles.sourceTitle}>Plan</Text>
                  <View style={styles.badge}>
                    <Text style={styles.badgeText}>Today's Workout</Text>
                  </View>
                </View>

                <Text style={styles.sourcePlanName}>{todayWorkout.plan_name}</Text>
                <Text style={styles.sourceDay}>{todayWorkout.day_title}</Text>

                <View style={styles.sourceStats}>
                  <View style={styles.statBadge}>
                    <Ionicons name="barbell-outline" size={14} color={COLORS.text.tertiary} />
                    <Text style={styles.statText}>
                      {todayWorkout.total_exercises} exercises
                    </Text>
                  </View>
                  <View style={styles.statBadge}>
                    <Ionicons name="time-outline" size={14} color={COLORS.text.tertiary} />
                    <Text style={styles.statText}>
                      ~{todayWorkout.estimated_minutes} min
                    </Text>
                  </View>
                  <View style={styles.statBadge}>
                    <Ionicons name="fitness-outline" size={14} color={COLORS.text.tertiary} />
                    <Text style={styles.statText}>{todayWorkout.focus}</Text>
                  </View>
                </View>
              </View>

              <Ionicons name="chevron-forward" size={24} color={COLORS.text.tertiary} />
            </TouchableOpacity>
          </MotiView>
        )}

        {hasActiveSchema && (
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 300, delay: 100 }}
          >
            <TouchableOpacity
              style={styles.sourceCard}
              onPress={() => handleSelectSchema(activeSchema.schema_id)}
              activeOpacity={0.7}
            >
              <View style={[styles.sourceIconContainer, styles.coachIcon]}>
                <Ionicons name="people" size={32} color={COLORS.success} />
              </View>

              <View style={styles.sourceContent}>
                <View style={styles.sourceHeader}>
                  <Text style={styles.sourceTitle}>Coach Schema</Text>
                  <View style={[styles.badge, styles.coachBadge]}>
                    <Text style={styles.badgeText}>Active</Text>
                  </View>
                </View>

                <Text style={styles.sourcePlanName}>
                  Week of {new Date(activeSchema.week_start).toLocaleDateString()}
                </Text>
                <Text style={styles.sourceDay}>Professional training program</Text>

                <View style={styles.sourceStats}>
                  <View style={styles.statBadge}>
                    <Ionicons name="calendar-outline" size={14} color={COLORS.text.tertiary} />
                    <Text style={styles.statText}>
                      Weekly plan
                    </Text>
                  </View>
                  <View style={styles.statBadge}>
                    <Ionicons name="trophy-outline" size={14} color={COLORS.text.tertiary} />
                    <Text style={styles.statText}>Coach designed</Text>
                  </View>
                </View>
              </View>

              <Ionicons name="chevron-forward" size={24} color={COLORS.text.tertiary} />
            </TouchableOpacity>
          </MotiView>
        )}

        {!hasAIPlan && hasActiveSchema && (
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 300, delay: 200 }}
          >
            <TouchableOpacity
              style={styles.createPlanCard}
              onPress={() => router.push('/(user)/plan-generator')}
              activeOpacity={0.7}
            >
              <Ionicons name="add-circle" size={24} color={COLORS.primary} />
              <Text style={styles.createPlanCardText}>Create Plan</Text>
              <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
            </TouchableOpacity>
          </MotiView>
        )}
      </ScrollView>

      {/* Day Selector Modal */}
      <Modal
        visible={showDaySelector}
        transparent
        animationType="slide"
        onRequestClose={() => setShowDaySelector(false)}
      >
        <View style={styles.modalOverlay}>
          <MotiView
            from={{ translateY: 300, opacity: 0 }}
            animate={{ translateY: 0, opacity: 1 }}
            transition={{ type: 'spring', damping: 20 }}
            style={styles.daySelectorModal}
          >
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Select Workout Day</Text>
              <TouchableOpacity
                onPress={() => setShowDaySelector(false)}
                style={styles.modalCloseButton}
              >
                <Ionicons name="close" size={24} color={COLORS.text.inverse} />
              </TouchableOpacity>
            </View>

            {loadingSchemaData ? (
              <ActivityIndicator size="large" color={COLORS.primary} style={{ marginVertical: SPACING.xl }} />
            ) : (
              <ScrollView style={styles.daysList} showsVerticalScrollIndicator={false}>
                {selectedSchemaData?.workouts.map((workout, index) => (
                  <TouchableOpacity
                    key={workout.workout_id}
                    style={styles.dayCard}
                    onPress={() => handleSelectDay(index)}
                    activeOpacity={0.7}
                  >
                    <View style={styles.dayIconContainer}>
                      <Text style={styles.dayNumber}>{index + 1}</Text>
                    </View>
                    <View style={styles.dayContent}>
                      <Text style={styles.dayName}>{DAY_NAMES[workout.day_of_week]}</Text>
                      <Text style={styles.dayFocus}>{workout.focus}</Text>
                      <View style={styles.dayStats}>
                        <Ionicons name="barbell-outline" size={12} color={COLORS.text.tertiary} />
                        <Text style={styles.dayStatsText}>
                          {workout.exercises.length} exercises
                        </Text>
                      </View>
                    </View>
                    <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
                  </TouchableOpacity>
                ))}
              </ScrollView>
            )}
          </MotiView>
        </View>
      </Modal>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingTop: 0,
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.background.auth,
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginTop: SPACING.lg,
    marginBottom: SPACING.sm,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    textAlign: 'center',
    marginBottom: SPACING.xl,
  },
  createPlanButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
  },
  createPlanButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.white,
  },
  scrollContent: {
    padding: SPACING.base,
  },
  headerTitle: {
    fontSize: 32,
    fontWeight: '700',
    color: COLORS.text.inverse,
    marginBottom: SPACING.xs,
  },
  headerSubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    marginBottom: SPACING.xl,
  },
  sourceCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS['3xl'],
    padding: SPACING.md,
    marginBottom: SPACING.md,
    ...SHADOWS.base,
  },
  sourceIconContainer: {
    width: 60,
    height: 60,
    borderRadius: BORDER_RADIUS['3xl'],
    backgroundColor: `${COLORS.primary}15`,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  coachIcon: {
    backgroundColor: `${COLORS.success}15`,
  },
  sourceContent: {
    flex: 1,
  },
  sourceHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.xs,
  },
  sourceTitle: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.tertiary,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginRight: SPACING.sm,
  },
  badge: {
    backgroundColor: `${COLORS.primary}20`,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 2,
    borderRadius: BORDER_RADIUS.sm,
  },
  coachBadge: {
    backgroundColor: `${COLORS.success}20`,
  },
  badgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  sourcePlanName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginBottom: 2,
  },
  sourceDay: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    marginBottom: SPACING.sm,
  },
  sourceStats: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  statBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  statText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  createPlanCard: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS.full,
    padding: SPACING.lg,
    borderWidth: 1,
    borderColor: COLORS.border.light,
    borderStyle: 'dashed',
  },
  createPlanCardText: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
    marginLeft: SPACING.md,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    justifyContent: 'flex-end',
  },
  daySelectorModal: {
    backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['2xl'],
    borderTopRightRadius: BORDER_RADIUS['2xl'],
    maxHeight: '80%',
    paddingBottom: SPACING.xl,
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.lg,
  },
  modalTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  modalCloseButton: {
    padding: SPACING.xs,
  },
  daysList: {
    padding: SPACING.lg,
  },
  dayCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS.full,
    padding: SPACING.md,
    marginBottom: SPACING.sm,
  },
  dayIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.lg,
    backgroundColor: COLORS.primaryDark,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  dayNumber: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  dayContent: {
    flex: 1,
  },
  dayName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
    marginBottom: 2,
  },
  dayFocus: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    marginBottom: 4,
  },
  dayStats: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  dayStatsText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
});
