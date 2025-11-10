import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Stack } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { useWorkoutHistory } from '@/hooks/user/use-workout-history';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    weekday: 'short', 
    month: 'short', 
    day: 'numeric',
    year: 'numeric'
  });
};

const formatDateShort = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric'
  });
};

export default function WorkoutHistoryScreen() {
  const [page, setPage] = useState(1);
  const [refreshing, setRefreshing] = useState(false);
  const { data, isLoading, error, refetch } = useWorkoutHistory({ page, pageSize: 20 });

  const handleLoadMore = () => {
    if (data?.has_more) {
      setPage((prev) => prev + 1);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    setPage(1);
    await refetch();
    setRefreshing(false);
  };

  if (isLoading && page === 1) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
        <Text style={styles.loadingText}>Loading workout history...</Text>
      </View>
    );
  }

  if (error || !data) {
    return (
      <SafeAreaView style={styles.container}>
        <Stack.Screen options={{ title: 'Workout History' }} />
        <View style={styles.errorContainer}>
          <Ionicons name="alert-circle-outline" size={64} color={COLORS.error} />
          <Text style={styles.errorTitle}>Failed to Load History</Text>
          <Text style={styles.errorSubtitle}>Please try again later</Text>
          <TouchableOpacity style={styles.retryButton} onPress={handleRefresh}>
            <Text style={styles.retryButtonText}>Retry</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  if (data.workouts.length === 0) {
    return (
      <SafeAreaView style={styles.container}>
        <Stack.Screen options={{ title: 'Workout History' }} />
        <View style={styles.emptyContainer}>
          <Ionicons name="calendar-outline" size={64} color={COLORS.text.tertiary} />
          <Text style={styles.emptyTitle}>No Workouts Yet</Text>
          <Text style={styles.emptySubtitle}>
            Complete your first workout to see it here
          </Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <Stack.Screen options={{ title: 'Workout History' }} />
      
      {/* Stats Header */}
      <View style={styles.statsHeader}>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>{data.total_count}</Text>
          <Text style={styles.statLabel}>Total Workouts</Text>
        </View>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>
            {data.workouts.reduce((acc, w) => acc + w.completed_sets, 0)}
          </Text>
          <Text style={styles.statLabel}>Sets Completed</Text>
        </View>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>
            {Math.round(data.workouts.reduce((acc, w) => acc + w.total_volume, 0) / 1000)}k
          </Text>
          <Text style={styles.statLabel}>Total Volume</Text>
        </View>
      </View>

      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={handleRefresh}
            tintColor={COLORS.primary}
          />
        }
      >
        {data.workouts.map((workout, index) => (
          <MotiView
            key={workout.date}
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 300, delay: index * 50 }}
            style={styles.workoutCard}
          >
            {/* Date Header */}
            <View style={styles.dateHeader}>
              <View style={styles.dateIconContainer}>
                <Ionicons name="calendar" size={20} color={COLORS.primary} />
              </View>
              <View style={styles.dateInfo}>
                <Text style={styles.dateText}>{formatDate(workout.date)}</Text>
                {workout.day_title && (
                  <Text style={styles.dayTitle}>{workout.day_title}</Text>
                )}
              </View>
            </View>

            {/* Workout Stats */}
            <View style={styles.workoutStats}>
              <View style={styles.statItem}>
                <Ionicons name="fitness" size={16} color={COLORS.text.secondary} />
                <Text style={styles.statItemText}>
                  {workout.total_exercises} exercises
                </Text>
              </View>
              <View style={styles.statItem}>
                <Ionicons name="checkmark-circle" size={16} color={COLORS.success} />
                <Text style={styles.statItemText}>
                  {workout.completed_sets} sets
                </Text>
              </View>
              <View style={styles.statItem}>
                <Ionicons name="barbell" size={16} color={COLORS.warning} />
                <Text style={styles.statItemText}>
                  {Math.round(workout.total_volume)} lbs
                </Text>
              </View>
              <View style={styles.statItem}>
                <Ionicons name="time" size={16} color={COLORS.info} />
                <Text style={styles.statItemText}>
                  {workout.duration_minutes} min
                </Text>
              </View>
            </View>

            {/* Exercises List */}
            {workout.exercises && workout.exercises.length > 0 && (
              <View style={styles.exercisesList}>
                <Text style={styles.exercisesLabel}>Exercises:</Text>
                <View style={styles.exerciseTags}>
                  {workout.exercises.slice(0, 3).map((exercise, idx) => (
                    <View key={idx} style={styles.exerciseTag}>
                      <Text style={styles.exerciseTagText} numberOfLines={1}>
                        {exercise}
                      </Text>
                    </View>
                  ))}
                  {workout.exercises.length > 3 && (
                    <View style={styles.exerciseTag}>
                      <Text style={styles.exerciseTagText}>
                        +{workout.exercises.length - 3} more
                      </Text>
                    </View>
                  )}
                </View>
              </View>
            )}
          </MotiView>
        ))}

        {/* Load More Button */}
        {data.has_more && (
          <TouchableOpacity
            style={styles.loadMoreButton}
            onPress={handleLoadMore}
            disabled={isLoading}
          >
            {isLoading ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.loadMoreText}>Load More</Text>
                <Ionicons name="chevron-down" size={20} color={COLORS.primary} />
              </>
            )}
          </TouchableOpacity>
        )}

        <View style={styles.spacer} />
      </ScrollView>
    </SafeAreaView>
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
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  errorTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginTop: SPACING.md,
  },
  errorSubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: SPACING.sm,
    textAlign: 'center',
  },
  retryButton: {
    marginTop: SPACING.xl,
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
  },
  retryButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginTop: SPACING.md,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    marginTop: SPACING.sm,
    textAlign: 'center',
  },
  statsHeader: {
    flexDirection: 'row',
    padding: SPACING.lg,
    gap: SPACING.md,
  },
  statCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.md,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  statValue: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.primary,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
    textAlign: 'center',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.lg,
    paddingTop: 0,
  },
  workoutCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING.md,
    ...SHADOWS.sm,
  },
  dateHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  dateIconContainer: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: `${COLORS.primary}15`,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  dateInfo: {
    flex: 1,
  },
  dateText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
  },
  dayTitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  workoutStats: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
    marginBottom: SPACING.md,
  },
  statItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  statItemText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  exercisesList: {
  },
  exercisesLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
    marginBottom: SPACING.xs,
  },
  exerciseTags: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.xs,
  },
  exerciseTag: {
    backgroundColor: `${COLORS.primary}10`,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs / 2,
    borderRadius: BORDER_RADIUS.sm,
    maxWidth: 150,
  },
  exerciseTagText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.primary,
  },
  loadMoreButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING.md,
    marginTop: SPACING.md,
    gap: SPACING.xs,
  },
  loadMoreText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.primary,
  },
  spacer: {
    height: SPACING.xl,
  },
});
