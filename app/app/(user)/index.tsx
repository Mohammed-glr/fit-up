import React from 'react';
import { SafeAreaView, StyleSheet, ScrollView, View } from 'react-native';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useUserStats } from '@/hooks/user/use-user-stats';
import { useTodayWorkout } from '@/hooks/user/use-today-workout';
import { DashboardGreeting } from '@/components/dashboard/dashboard-greeting';
import { TodayWorkoutCard } from '@/components/dashboard/today-workout-card';
import { MetricsGrid } from '@/components/dashboard/metrics-grid';
import { QuickActionsGrid } from '@/components/dashboard/quick-actions-grid';
import { ActivityFeed } from '@/components/dashboard/activity-feed';
import { WorkoutHistorySummary } from '@/components/dashboard/workout-history-summary';
import { COLORS, SPACING } from '@/constants/theme';

export default function DashboardScreen() {
  const { data: currentUser } = useCurrentUser();
  const { data: userStats, isLoading: isLoadingStats } = useUserStats();
  const { data: todayWorkout, isLoading: isLoadingWorkout } = useTodayWorkout();

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <ScrollView 
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          <DashboardGreeting name={currentUser?.name} />
          
          <TodayWorkoutCard workout={todayWorkout} isLoading={isLoadingWorkout} />
          
          <WorkoutHistorySummary />

          <MetricsGrid stats={userStats} isLoading={isLoadingStats} />
          
          <QuickActionsGrid />
          
          <ActivityFeed limit={10} />
          
          
          <View style={styles.spacer} />
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
    overflow: 'hidden',
    paddingBottom: SPACING['6xl'],
  },
  safeArea: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.lg,
  },
  spacer: {
    height: SPACING.xl,
  },
});

