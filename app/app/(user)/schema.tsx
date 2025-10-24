import React from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { WorkoutDayCard } from '@/components/schema/workout-day-card';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

export default function UserSchemaScreen() {
  const [refreshing, setRefreshing] = React.useState(false);
  const [selectedWeek, setSelectedWeek] = React.useState(0);
  const isLoading = false;
  const hasSchema = false;

  const onRefresh = async () => {
    setRefreshing(true);
    // TODO: Refetch schema data
    setTimeout(() => setRefreshing(false), 1000);
  };

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (!hasSchema) {
    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="calendar-outline" size={80} color={COLORS.text.tertiary} />
        <Text style={styles.emptyTitle}>No Workout Schema</Text>
        <Text style={styles.emptySubtitle}>
          Your coach hasn't assigned a workout program yet.
        </Text>
        <Text style={styles.emptySubtitle}>
          Check back later or contact your coach.
        </Text>
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.contentContainer}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={COLORS.primary} />
      }
    >
      <View style={styles.header}>
        <View style={styles.headerTop}>
          <View>
            <Text style={styles.schemaName}>Full Body Strength</Text>
            <Text style={styles.coachName}>by Coach Sarah</Text>
          </View>
          <View style={styles.statusBadge}>
            <View style={styles.activeDot} />
            <Text style={styles.statusText}>Active</Text>
          </View>
        </View>
        
        <View style={styles.schemaInfo}>
          <View style={styles.infoItem}>
            <Ionicons name="calendar" size={16} color={COLORS.primary} />
            <Text style={styles.infoText}>Week 1 of 12</Text>
          </View>
          <View style={styles.infoItem}>
            <Ionicons name="barbell" size={16} color={COLORS.primary} />
            <Text style={styles.infoText}>5 days/week</Text>
          </View>
          <View style={styles.infoItem}>
            <Ionicons name="time" size={16} color={COLORS.primary} />
            <Text style={styles.infoText}>45-60 min</Text>
          </View>
        </View>

        <View style={styles.progressContainer}>
          <View style={styles.progressHeader}>
            <Text style={styles.progressLabel}>Week Progress</Text>
            <Text style={styles.progressPercentage}>3/5 complete</Text>
          </View>
          <View style={styles.progressBar}>
            <View style={[styles.progressFill, { width: '60%' }]} />
          </View>
        </View>
      </View>

      <View style={styles.weekSelector}>
        <TouchableOpacity
          style={styles.weekButton}
          onPress={() => setSelectedWeek(Math.max(0, selectedWeek - 1))}
          disabled={selectedWeek === 0}
        >
          <Ionicons 
            name="chevron-back" 
            size={24} 
            color={selectedWeek === 0 ? COLORS.text.tertiary : COLORS.primary} 
          />
        </TouchableOpacity>
        <Text style={styles.weekLabel}>Week {selectedWeek + 1}</Text>
        <TouchableOpacity
          style={styles.weekButton}
          onPress={() => setSelectedWeek(selectedWeek + 1)}
        >
          <Ionicons name="chevron-forward" size={24} color={COLORS.primary} />
        </TouchableOpacity>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>This Week's Workouts</Text>
        
        {/* TODO: Map through actual workout data */}
        {[0, 1, 2, 3, 4, 5, 6].map((day) => (
          <WorkoutDayCard
            key={day}
            dayOfWeek={day}
            isRestDay={day === 3 || day === 6}
          />
        ))}
      </View>

      <View style={styles.notesContainer}>
        <View style={styles.notesHeader}>
          <Ionicons name="document-text" size={20} color={COLORS.primary} />
          <Text style={styles.notesTitle}>Coach Notes</Text>
        </View>
        <Text style={styles.notesText}>
          Focus on progressive overload this week. Increase weight by 5-10% from last week if you completed all sets with good form.
        </Text>
      </View>
    </ScrollView>
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
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
    backgroundColor: COLORS.background.auth,
  },
  emptyTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.lg,
    marginBottom: SPACING.xs,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginBottom: SPACING.xs,
  },
  contentContainer: {
    padding: SPACING.base,
  },
  header: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    marginBottom: SPACING.base,
    ...SHADOWS.sm,
  },
  headerTop: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: SPACING.md,
  },
  schemaName: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: 4,
  },
  coachName: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.success,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.xs,
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
  schemaInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
  },
  infoItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  infoText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.auth.secondary,
  },
  progressContainer: {
    marginTop: SPACING.sm,
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
  weekSelector: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.md,
    marginBottom: SPACING.base,
    ...SHADOWS.sm,
  },
  weekButton: {
    padding: SPACING.xs,
  },
  weekLabel: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  section: {
    marginBottom: SPACING.base,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.md,
  },
  notesContainer: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    marginBottom: SPACING.xl,
    borderLeftWidth: 4,
    borderLeftColor: COLORS.primary,
    ...SHADOWS.sm,
  },
  notesHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginBottom: SPACING.sm,
  },
  notesTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  notesText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.secondary,
    lineHeight: 20,
  },
});

