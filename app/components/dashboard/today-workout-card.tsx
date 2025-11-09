import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, ActivityIndicator } from 'react-native';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { TodayWorkout } from '@/hooks/user/use-today-workout';

interface TodayWorkoutCardProps {
  workout: TodayWorkout | null | undefined;
  isLoading: boolean;
}

export const TodayWorkoutCard: React.FC<TodayWorkoutCardProps> = ({ workout, isLoading }) => {
  const router = useRouter();

  if (isLoading) {
    return (
      <MotiView
        from={{ opacity: 0, translateY: 20 }}
        animate={{ opacity: 1, translateY: 0 }}
        transition={{ type: 'timing', duration: 400 }}
        style={styles.container}
      >
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={COLORS.primary} />
          <Text style={styles.loadingText}>Loading today's workout...</Text>
        </View>
      </MotiView>
    );
  }

  if (!workout) {
    return (
      <MotiView
        from={{ opacity: 0, translateY: 20 }}
        animate={{ opacity: 1, translateY: 0 }}
        transition={{ type: 'timing', duration: 400 }}
        style={styles.container}
      >
        <View style={styles.emptyContainer}>
          <Ionicons name="calendar-outline" size={48} color={COLORS.text.tertiary} />
          <Text style={styles.emptyTitle}>No Workout Scheduled</Text>
          <Text style={styles.emptySubtitle}>
            Get started by creating a workout plan
          </Text>
          <TouchableOpacity
            style={styles.createPlanButton}
            onPress={() => router.push('/(user)/plan-generator')}
            activeOpacity={0.7}
          >
            <Text style={styles.createPlanButtonText}>Create Plan</Text>
          </TouchableOpacity>
        </View>
      </MotiView>
    );
  }

  if (workout.is_rest) {
    return (
      <MotiView
        from={{ opacity: 0, translateY: 20 }}
        animate={{ opacity: 1, translateY: 0 }}
        transition={{ type: 'timing', duration: 400 }}
        style={[styles.container, styles.restDayContainer]}
      >
        <View style={styles.restIconContainer}>
          <Ionicons name="bed" size={40} color={COLORS.info} />
        </View>
        <View style={styles.restContent}>
          <Text style={styles.restTitle}>Rest Day</Text>
          <Text style={styles.restSubtitle}>{workout.day_title}</Text>
          <Text style={styles.restDescription}>
            Recovery is essential for progress. Take today to let your muscles heal and come back stronger! 💪
          </Text>
        </View>
      </MotiView>
    );
  }

  const handleStartWorkout = () => {
    router.push('/(user)/workout-session');
  };

  return (
    <MotiView
      from={{ opacity: 0, translateY: 20 }}
      animate={{ opacity: 1, translateY: 0 }}
      transition={{ type: 'timing', duration: 400 }}
      style={styles.container}
    >
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Ionicons name="barbell" size={24} color={COLORS.primary} />
          <View style={styles.headerTextContainer}>
            <Text style={styles.headerTitle}>Today's Workout</Text>
            <Text style={styles.headerSubtitle}>{workout.day_title}</Text>
          </View>
        </View>
        {workout.is_completed && (
          <View style={styles.completedBadge}>
            <Ionicons name="checkmark-circle" size={20} color={COLORS.success} />
            <Text style={styles.completedText}>Completed</Text>
          </View>
        )}
      </View>

      {/* Focus and Stats */}
      <View style={styles.statsRow}>
        <View style={styles.statItem}>
          <Ionicons name="fitness" size={16} color={COLORS.text.secondary} />
          <Text style={styles.statText}>{workout.focus}</Text>
        </View>
        <View style={styles.statItem}>
          <Ionicons name="list" size={16} color={COLORS.text.secondary} />
          <Text style={styles.statText}>{workout.total_exercises} exercises</Text>
        </View>
        <View style={styles.statItem}>
          <Ionicons name="time" size={16} color={COLORS.text.secondary} />
          <Text style={styles.statText}>{workout.estimated_minutes} min</Text>
        </View>
      </View>

      {/* Exercise Preview */}
      <Text style={styles.exercisesTitle}>Exercises</Text>
      <ScrollView
        style={styles.exercisesList}
        showsVerticalScrollIndicator={false}
        nestedScrollEnabled={true}
      >
        {workout.exercises?.slice(0, 5).map((exercise, index) => (
          <View key={index} style={styles.exerciseItem}>
            <View style={styles.exerciseNumber}>
              <Text style={styles.exerciseNumberText}>{index + 1}</Text>
            </View>
            <View style={styles.exerciseDetails}>
              <Text style={styles.exerciseName}>{exercise.name}</Text>
              <Text style={styles.exerciseMeta}>
                {exercise.sets} sets × {exercise.reps} reps
              </Text>
            </View>
          </View>
        ))}
        {workout.exercises?.length > 5 && (
          <Text style={styles.moreExercises}>
            +{workout.exercises.length - 5} more exercises
          </Text>
        )}
      </ScrollView>

      {/* Action Button */}
      <TouchableOpacity
        style={[styles.startButton, workout.is_completed && styles.startButtonCompleted]}
        onPress={handleStartWorkout}
        activeOpacity={0.8}
      >
        <Ionicons
          name={workout.is_completed ? 'checkmark-circle' : 'play-circle'}
          size={24}
          color={COLORS.white}
        />
        <Text style={styles.startButtonText}>
          {workout.is_completed ? 'View Workout' : 'Start Workout'}
        </Text>
      </TouchableOpacity>
    </MotiView>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: COLORS.primarySoft,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginTop: SPACING.lg,
    ...SHADOWS.sm,
  },
  loadingContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.xl,
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  emptyContainer: {
    alignItems: 'center',
    paddingVertical: SPACING.xl,
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
  createPlanButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    marginTop: SPACING.lg,
  },
  createPlanButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
  restDayContainer: {
    backgroundColor: `${COLORS.info}10`,
    borderWidth: 2,
    borderColor: `${COLORS.info}30`,
  },
  restIconContainer: {
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  restContent: {
    alignItems: 'center',
  },
  restTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.xs,
  },
  restSubtitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.info,
    marginBottom: SPACING.md,
  },
  restDescription: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
    textAlign: 'center',
    lineHeight: 20,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  headerTextContainer: {
    marginLeft: SPACING.sm,
    flex: 1,
  },
  headerTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.secondary,
  },
  headerSubtitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
  },
  completedBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: `${COLORS.success}15`,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.md,
  },
  completedText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.success,
    marginLeft: 4,
  },
  statsRow: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginBottom: SPACING.lg,
  },
  statItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  statText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
  },
  exercisesTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.sm,
  },
  exercisesList: {
    maxHeight: 200,
    marginBottom: SPACING.md,
  },
  exerciseItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: SPACING.sm,
  },
  exerciseNumber: {
    width: 28,
    height: 28,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: `${COLORS.primary}15`,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: SPACING.sm,
  },
  exerciseNumberText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.primary,
  },
  exerciseDetails: {
    flex: 1,
  },
  exerciseName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.primary,
    marginBottom: 2,
  },
  exerciseMeta: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
  },
  moreExercises: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    paddingVertical: SPACING.sm,
    fontStyle: 'italic',
  },
  startButton: {
    backgroundColor: COLORS.primary,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
    gap: SPACING.sm,
    ...SHADOWS.sm,
  },
  startButtonCompleted: {
    backgroundColor: COLORS.success,
  },
  startButtonText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.white,
  },
});
