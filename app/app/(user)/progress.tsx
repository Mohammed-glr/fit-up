import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
  SafeAreaView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useExerciseProgress } from '@/hooks/user/use-exercise-progress';
import { useMostUsedExercises } from '@/hooks/user/use-most-used-exercises';
import { LineChart } from '@/components/charts';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS } from '@/constants/theme';

type TimeRange = '7d' | '30d' | '90d' | 'all';

const TIME_RANGES: { value: TimeRange; label: string }[] = [
  { value: '7d', label: '7 Days' },
  { value: '30d', label: '30 Days' },
  { value: '90d', label: '90 Days' },
  { value: 'all', label: 'All Time' },
];

// Mock exercise list - in real app, this would come from an API
const EXERCISES = [
  { id: 1, name: 'Bench Press' },
  { id: 2, name: 'Squat' },
  { id: 3, name: 'Deadlift' },
  { id: 4, name: 'Overhead Press' },
  { id: 5, name: 'Barbell Row' },
  { id: 6, name: 'Pull Up' },
];

export default function ProgressScreen() {
  const [selectedExercise, setSelectedExercise] = useState<string>('Bench Press');
  const [timeRange, setTimeRange] = useState<TimeRange>('30d');

  const getDateRange = () => {
    const endDate = new Date();
    const startDate = new Date();

    switch (timeRange) {
      case '7d':
        startDate.setDate(endDate.getDate() - 7);
        break;
      case '30d':
        startDate.setDate(endDate.getDate() - 30);
        break;
      case '90d':
        startDate.setDate(endDate.getDate() - 90);
        break;
      case 'all':
        return { startDate: undefined, endDate: undefined };
    }

    return {
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
    };
  };

  const { startDate, endDate } = getDateRange();

  const { data: progressData, isLoading, error } = useExerciseProgress({
    exerciseName: selectedExercise,
    startDate,
    endDate,
  });

  const volumeData = progressData?.data_points.map((point, index) => ({
    x: index,
    y: point.volume,
    label: new Date(point.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
  })) || [];

  const weightData = progressData?.data_points.map((point, index) => ({
    x: index,
    y: point.weight,
    label: new Date(point.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
  })) || [];

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.container}>
        <View style={styles.header}>
                <View style={styles.titleContainer}>
                    <Ionicons name="barbell" size={32} color={COLORS.primary} />
                <Text style={styles.headerTitle}>Progress Tracker</Text>
                </View>
                <Text style={styles.headerSubtitle}>
                    Monitor your exercise progress over time!
                </Text>
            </View>

        <ScrollView 
          style={styles.scrollView}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Select Exercise</Text>
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.exerciseList}
            >
              {EXERCISES.map((exercise) => (
                <TouchableOpacity
                  key={exercise.id}
                  style={[
                    styles.exerciseChip,
                    selectedExercise === exercise.name && styles.exerciseChipActive,
                  ]}
                  onPress={() => setSelectedExercise(exercise.name)}
                >
                  <Text
                    style={[
                      styles.exerciseChipText,
                      selectedExercise === exercise.name && styles.exerciseChipTextActive,
                    ]}
                  >
                    {exercise.name}
                  </Text>
                </TouchableOpacity>
              ))}
            </ScrollView>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Time Range</Text>
            <View style={styles.timeRangeContainer}>
              {TIME_RANGES.map((range) => (
                <TouchableOpacity
                  key={range.value}
                  style={[
                    styles.timeRangeButton,
                    timeRange === range.value && styles.timeRangeButtonActive,
                  ]}
                  onPress={() => setTimeRange(range.value)}
                >
                  <Text
                    style={[
                      styles.timeRangeText,
                      timeRange === range.value && styles.timeRangeTextActive,
                    ]}
                  >
                    {range.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Loading State */}
          {isLoading && (
            <View style={styles.loadingContainer}>
              <ActivityIndicator size="large" color={COLORS.primary} />
              <Text style={styles.loadingText}>Loading progress data...</Text>
            </View>
          )}

          {/* Error State */}
          {error && (
            <View style={styles.errorContainer}>
              <Ionicons name="alert-circle" size={48} color={COLORS.error} />
              <Text style={styles.errorText}>Failed to load progress data</Text>
            </View>
          )}

          {/* Stats Summary */}
          {progressData && !isLoading && (
            <>
              <View style={styles.statsGrid}>
                <View style={styles.statCard}>
                  <Ionicons name="trending-up" size={24} color={COLORS.primary} />
                  <Text style={styles.statValue}>
                    {progressData.max_weight.toFixed(1)} kg
                  </Text>
                  <Text style={styles.statLabel}>Max Weight</Text>
                </View>

                <View style={styles.statCard}>
                  <Ionicons name="fitness" size={24} color={COLORS.primary} />
                  <Text style={styles.statValue}>
                    {progressData.max_volume.toFixed(0)} kg
                  </Text>
                  <Text style={styles.statLabel}>Max Volume</Text>
                </View>

                <View style={styles.statCard}>
                  <Ionicons name="repeat" size={24} color={COLORS.primary} />
                  <Text style={styles.statValue}>{progressData.total_sets}</Text>
                  <Text style={styles.statLabel}>Total Sets</Text>
                </View>

                <View style={styles.statCard}>
                  <Ionicons name="calendar" size={24} color={COLORS.primary} />
                  <Text style={styles.statValue}>
                    {progressData.data_points.length}
                  </Text>
                  <Text style={styles.statLabel}>Workouts</Text>
                </View>
              </View>

              {/* Volume Chart */}
              {volumeData.length > 0 && (
                <View style={styles.chartSection}>
                  <Text style={styles.chartTitle}>Volume Over Time</Text>
                  <Text style={styles.chartSubtitle}>
                    Total weight lifted (kg) per session
                  </Text>
                  <LineChart
                    data={volumeData}
                    color={COLORS.primary}
                    yAxisLabel="Volume (kg)"
                    showGradient={true}
                  />
                </View>
              )}

              {/* Weight Chart */}
              {weightData.length > 0 && (
                <View style={styles.chartSection}>
                  <Text style={styles.chartTitle}>Max Weight Over Time</Text>
                  <Text style={styles.chartSubtitle}>
                    Heaviest weight used (kg) per session
                  </Text>
                  <LineChart
                    data={weightData}
                    color={COLORS.success}
                    yAxisLabel="Weight (kg)"
                    showGradient={true}
                  />
                </View>
              )}

              {/* Personal Records */}
              {progressData.data_points.some(p => p.is_personal_record) && (
                <View style={styles.section}>
                  <Text style={styles.sectionTitle}>
                    <Ionicons name="trophy" size={20} color={COLORS.warning} /> Personal Records
                  </Text>
                  <View style={styles.prList}>
                    {progressData.data_points
                      .filter(p => p.is_personal_record)
                      .map((pr, index) => (
                        <View key={index} style={styles.prCard}>
                          <View style={styles.prIcon}>
                            <Ionicons name="trophy" size={20} color={COLORS.warning} />
                          </View>
                          <View style={styles.prInfo}>
                            <Text style={styles.prDate}>
                              {new Date(pr.date).toLocaleDateString('en-US', {
                                month: 'long',
                                day: 'numeric',
                                year: 'numeric',
                              })}
                            </Text>
                            <Text style={styles.prDetails}>
                              {pr.weight} kg × {pr.reps} reps ({pr.sets} sets)
                            </Text>
                          </View>
                          <View style={styles.prBadge}>
                            <Text style={styles.prBadgeText}>PR</Text>
                          </View>
                        </View>
                      ))}
                  </View>
                </View>
              )}
            </>
          )}

          {/* Empty State */}
          {progressData && progressData.data_points.length === 0 && !isLoading && (
            <View style={styles.emptyContainer}>
              <Ionicons name="bar-chart-outline" size={64} color={COLORS.text.tertiary} />
              <Text style={styles.emptyTitle}>No Progress Data</Text>
              <Text style={styles.emptyText}>
                Start tracking {selectedExercise} to see your progress over time
              </Text>
            </View>
          )}
        </ScrollView>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  container: {
    flex: 1,
  },
   header: {
    marginBottom: 24,
    margin: SPACING.base,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 4,
  },
  title: {
    fontSize: 32,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  subtitle: {
    fontSize: 16,
    color: '#888888',
  },
  headerInfo: {
    flex: 1,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  headerSubtitle: {
    fontSize: 14,
    color: '#888888',
    marginTop: 2,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
    paddingBottom: SPACING['6xl'],
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: '600',
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  exerciseList: {
    flexDirection: 'row',
    gap: SPACING.sm,
    paddingRight: SPACING.lg,
  },
  exerciseChip: {
    backgroundColor: COLORS.background.card,
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.lg,
  },
  exerciseChipActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  exerciseChipText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.text.placeholder,
  },
  exerciseChipTextActive: {
    color: COLORS.black,
  },
  timeRangeContainer: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  timeRangeButton: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    alignItems: 'center',
  },
  timeRangeButtonActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  timeRangeText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.text.placeholder,
  },
  timeRangeTextActive: {
    color: COLORS.black,
  },
  loadingContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['4xl'],
  },
  loadingText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
    marginTop: SPACING.md,
  },
  errorContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['4xl'],
  },
  errorText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.error,
    marginTop: SPACING.md,
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
    marginBottom: SPACING.xl,
  },
  statCard: {
    flex: 1,
    minWidth: '45%',
    backgroundColor: COLORS.background.card,
    padding: SPACING.base,
    borderRadius: BORDER_RADIUS['2xl'],
    alignItems: 'center',
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: '700',
    color: COLORS.text.inverse,
    marginTop: SPACING.sm,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.placeholder,
    marginTop: SPACING.xs,
  },
  chartSection: {
    marginBottom: SPACING.xl,
  },
  chartTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: '600',
    color: COLORS.text.inverse,
    marginBottom: SPACING.xs,
  },
  chartSubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    marginBottom: SPACING.md,
  },
  prList: {
    gap: SPACING.sm,
  },
  prCard: {
    backgroundColor: COLORS.background.card,
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS['2xl'],
    flexDirection: 'row',
    alignItems: 'center',
  },
  prIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: COLORS.background.warningSoft,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: SPACING.md,
  },
  prInfo: {
    flex: 1,
  },
  prDate: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.text.inverse,
    marginBottom: SPACING.xs / 2,
  },
  prDetails: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
  },
  prBadge: {
    backgroundColor: COLORS.warning,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs / 2,
    borderRadius: BORDER_RADIUS.sm,
  },
  prBadgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: '700',
    color: COLORS.black,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['4xl'],
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: '600',
    color: COLORS.text.inverse,
    marginTop: SPACING.base,
    marginBottom: SPACING.xs,
  },
  emptyText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    textAlign: 'center',
  },
});
