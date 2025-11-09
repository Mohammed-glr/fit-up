import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ActivityIndicator } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { useWorkoutHistory } from '@/hooks/user/use-workout-history';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS } from '@/constants/theme';

// Simple date formatter
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  return `${months[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}`;
};

export const WorkoutHistorySummary = () => {
  const router = useRouter();
  const { data: historyData, isLoading } = useWorkoutHistory({
    page: 1,
    pageSize: 3, // Just show last 3 workouts
  });

  if (isLoading) {
    return (
      <View style={styles.container}>
        <View style={styles.header}>
          <Text style={styles.title}>Recent Workouts</Text>
        </View>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="small" color={COLORS.primary} />
        </View>
      </View>
    );
  }

  if (!historyData || historyData.workouts.length === 0) {
    return null; // Don't show if no workout history
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Recent Workouts</Text>
        <TouchableOpacity
          onPress={() => router.push('/(user)/workout-history')}
          style={styles.viewAllButton}
        >
          <Text style={styles.viewAllText}>View All</Text>
          <Ionicons name="chevron-forward" size={16} color={COLORS.primary} />
        </TouchableOpacity>
      </View>

      <View style={styles.workoutsContainer}>
        {historyData.workouts.map((workout, index) => (
          <MotiView
            key={workout.date}
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{
              type: 'timing',
              duration: 400,
              delay: index * 100,
            }}
          >
            <TouchableOpacity
              style={styles.workoutCard}
              onPress={() => router.push('/(user)/workout-history')}
              activeOpacity={0.7}
            >
              <View style={styles.workoutHeader}>
                <View style={styles.dateContainer}>
                  <Ionicons 
                    name="calendar-outline" 
                    size={16} 
                    color={COLORS.text.secondary} 
                  />
                  <Text style={styles.dateText}>
                    {formatDate(workout.date)}
                  </Text>
                </View>
                {workout.day_title && (
                  <Text style={styles.dayTitle}>{workout.day_title}</Text>
                )}
              </View>

              <View style={styles.statsRow}>
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>{workout.total_exercises}</Text>
                  <Text style={styles.statLabel}>Exercises</Text>
                </View>
                <View style={styles.statDivider} />
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>{workout.completed_sets}</Text>
                  <Text style={styles.statLabel}>Sets</Text>
                </View>
                <View style={styles.statDivider} />
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>
                    {workout.total_volume.toFixed(0)}
                  </Text>
                  <Text style={styles.statLabel}>kg</Text>
                </View>
              </View>

              {workout.exercises && workout.exercises.length > 0 && (
                <View style={styles.exercisesContainer}>
                  {workout.exercises.slice(0, 2).map((exercise, idx) => (
                    <View key={idx} style={styles.exerciseTag}>
                      <Text style={styles.exerciseText} numberOfLines={1}>
                        {exercise}
                      </Text>
                    </View>
                  ))}
                  {workout.exercises.length > 2 && (
                    <View style={styles.exerciseTag}>
                      <Text style={styles.exerciseText}>
                        +{workout.exercises.length - 2} more
                      </Text>
                    </View>
                  )}
                </View>
              )}
            </TouchableOpacity>
          </MotiView>
        ))}
      </View>

      <TouchableOpacity
        style={styles.viewAllCard}
        onPress={() => router.push('/(user)/workout-history')}
        activeOpacity={0.7}
      >
        <Text style={styles.viewAllCardText}>View Full History</Text>
        <Ionicons name="arrow-forward" size={20} color={COLORS.primary} />
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: SPACING.xl,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: '700',
    color: COLORS.text.inverse,
  },
  viewAllButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  viewAllText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.primary,
  },
  loadingContainer: {
    padding: SPACING.xl,
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
  },
  workoutsContainer: {
    gap: SPACING.md,
  },
  workoutCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.md,
  },
  workoutHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  dateContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  dateText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    fontWeight: '500',
  },
  dayTitle: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.primary,
  },
  statsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    paddingVertical: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
  },
  statValue: {
    fontSize: FONT_SIZES.lg,
    fontWeight: '700',
    color: COLORS.text.inverse,
    marginBottom: 2,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    fontWeight: '500',
  },
  statDivider: {
    width: 4,
    height: 24,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primaryDark,
  },
  exercisesContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.xs,
  },
  exerciseTag: {
    backgroundColor: COLORS.background.accent,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs / 2,
    borderRadius: BORDER_RADIUS.sm,
    maxWidth: 120,
  },
  exerciseText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
    fontWeight: '500',
  },
  viewAllCard: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.md,
    marginTop: SPACING.md,
    gap: SPACING.sm,
  },
  viewAllCardText: {
    fontSize: FONT_SIZES.base,
    fontWeight: '600',
    color: COLORS.primary,
  },
});
