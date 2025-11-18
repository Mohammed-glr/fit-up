import React from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface WorkoutExerciseDetail {
  we_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
  exercise: {
    exercise_id: number;
    name: string;
    muscle_groups: string[];
    difficulty: string;
    equipment: string;
    type: string;
  };
}

interface WorkoutWithExercises {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
  exercises: WorkoutExerciseDetail[];
}

interface WeeklySchema {
  schema_id: number;
  user_id: string;
  week_start: string;
  active: boolean;
  metadata?: {
    custom_data?: {
      title?: string;
      notes?: string;
    };
  };
  workouts?: WorkoutWithExercises[];
}

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export default function CoachSchemaViewScreen() {
  const router = useRouter();
  const { schemaId } = useLocalSearchParams<{ schemaId: string }>();

  const { data: schema, isLoading } = useQuery({
    queryKey: ['schema-details', schemaId],
    queryFn: async () => {
      if (!schemaId) throw new Error('Schema ID required');
      const response = await httpClient.get(`schemas/${schemaId}/workouts`);
      return response.data as WeeklySchema;
    },
    enabled: !!schemaId,
  });

  const handleEditSchema = () => {
    if (!schema) return;
    router.push({
      pathname: '/(coach)/schema-create',
      params: {
        userId: schema.user_id,
        schemaId: schema.schema_id.toString(),
      },
    });
  };

  const handleDeleteSchema = () => {
    Alert.alert(
      'Delete Schema',
      'Are you sure you want to delete this schema? This action cannot be undone.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await httpClient.delete(`schemas/${schemaId}`);
              Alert.alert('Success', 'Schema deleted successfully');
              router.back();
            } catch (error) {
              Alert.alert('Error', 'Failed to delete schema');
            }
          },
        },
      ]
    );
  };

  const handleToggleActive = async () => {
    if (!schema) return;
    try {
      await httpClient.patch(`schemas/${schemaId}/active`, {
        active: !schema.active,
      });
      Alert.alert('Success', `Schema ${schema.active ? 'deactivated' : 'activated'} successfully`);
    } catch (error) {
      Alert.alert('Error', 'Failed to update schema status');
    }
  };

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (!schema) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.errorText}>Schema not found</Text>
        <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
          <Text style={styles.backButtonText}>Go Back</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const workouts = schema.workouts || [];
  const schemaTitle = schema.metadata?.custom_data?.title || `Schema #${schema.schema_id}`;
  const schemaNotes = schema.metadata?.custom_data?.notes;

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.contentContainer}>
      <View style={styles.header}>
        <View style={styles.headerInfo}>
          <Text style={styles.headerTitle}>{schemaTitle}</Text>
          <Text style={styles.headerSubtitle}>
            Week of {new Date(schema.week_start).toLocaleDateString()}
          </Text>
        </View>
        {schema.active && (
          <View style={styles.activeBadge}>
            <View style={styles.activeDot} />
            <Text style={styles.activeText}>Active</Text>
          </View>
        )}
      </View>

      {schemaNotes && (
        <View style={styles.notesCard}>
          <View style={styles.notesHeader}>
            <Ionicons name="document-text" size={20} color={COLORS.primary} />
            <Text style={styles.notesTitle}>Notes</Text>
          </View>
          <Text style={styles.notesText}>{schemaNotes}</Text>
        </View>
      )}

      <View style={styles.actionsRow}>
        <TouchableOpacity style={styles.actionButton} onPress={handleEditSchema}>
          <Ionicons name="create-outline" size={20} color={COLORS.primary} />
          <Text style={styles.actionButtonText}>Edit</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.actionButton, styles.actionButtonSecondary]}
          onPress={handleToggleActive}
        >
          <Ionicons
            name={schema.active ? 'pause-outline' : 'play-outline'}
            size={20}
            color={COLORS.warning}
          />
          <Text style={[styles.actionButtonText, styles.actionButtonTextSecondary]}>
            {schema.active ? 'Deactivate' : 'Activate'}
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.actionButton, styles.actionButtonDanger]}
          onPress={handleDeleteSchema}
        >
          <Ionicons name="trash-outline" size={20} color={COLORS.error} />
          <Text style={[styles.actionButtonText, styles.actionButtonTextDanger]}>Delete</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.statsRow}>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>{workouts.length}</Text>
          <Text style={styles.statLabel}>Workouts</Text>
        </View>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>
            {workouts.reduce((sum, w) => sum + (w.exercises?.length || 0), 0)}
          </Text>
          <Text style={styles.statLabel}>Exercises</Text>
        </View>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>
            {workouts.reduce(
              (sum, w) => sum + (w.exercises?.reduce((s, e) => s + e.sets, 0) || 0),
              0
            )}
          </Text>
          <Text style={styles.statLabel}>Total Sets</Text>
        </View>
      </View>

      {workouts.length === 0 ? (
        <View style={styles.emptyState}>
          <Ionicons name="calendar-outline" size={60} color={COLORS.text.tertiary} />
          <Text style={styles.emptyStateText}>No workouts in this schema</Text>
        </View>
      ) : (
        <View style={styles.workoutsList}>
          <Text style={styles.sectionTitle}>Workouts</Text>
          {workouts.map((workout) => (
            <View key={workout.workout_id} style={styles.workoutCard}>
              <View style={styles.workoutHeader}>
                <View style={styles.dayBadge}>
                  <Ionicons name="calendar" size={16} color={COLORS.primary} />
                  <Text style={styles.dayText}>
                    {DAY_NAMES[workout.day_of_week] || `Day ${workout.day_of_week}`}
                  </Text>
                </View>
                <Text style={styles.focusText}>{workout.focus}</Text>
              </View>

              <View style={styles.exercisesList}>
                {workout.exercises?.map((exercise, idx) => (
                  <View key={exercise.we_id} style={styles.exerciseRow}>
                    <Text style={styles.exerciseNumber}>{idx + 1}</Text>
                    <View style={styles.exerciseInfo}>
                      <Text style={styles.exerciseName}>{exercise.exercise.name}</Text>
                      <Text style={styles.exerciseDetails}>
                        {exercise.sets} sets × {exercise.reps} reps
                        {exercise.rest_seconds > 0 && ` • ${exercise.rest_seconds}s rest`}
                      </Text>
                      {exercise.exercise.muscle_groups &&
                        exercise.exercise.muscle_groups.length > 0 && (
                          <Text style={styles.muscleGroups}>
                            {exercise.exercise.muscle_groups.join(', ')}
                          </Text>
                        )}
                    </View>
                  </View>
                ))}
              </View>
            </View>
          ))}
        </View>
      )}
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
  contentContainer: {
    padding: SPACING.base,
    paddingBottom: SPACING['6xl'],
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.lg,
  },
  headerInfo: {
    flex: 1,
  },
  headerTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  headerSubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  activeBadge: {
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
    backgroundColor: COLORS.primarySoft,
  },
  activeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  notesCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    marginBottom: SPACING.md,
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
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.secondary,
    lineHeight: 22,
  },
  actionsRow: {
    flexDirection: 'row',
    gap: SPACING.sm,
    marginBottom: SPACING.lg,
  },
  actionButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.xl,
    paddingVertical: SPACING.sm,
    borderWidth: 1,
    borderColor: COLORS.primary,
  },
  actionButtonSecondary: {
    borderColor: COLORS.warning,
  },
  actionButtonDanger: {
    borderColor: COLORS.error,
  },
  actionButtonText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  actionButtonTextSecondary: {
    color: COLORS.warning,
  },
  actionButtonTextDanger: {
    color: COLORS.error,
  },
  statsRow: {
    flexDirection: 'row',
    gap: SPACING.sm,
    marginBottom: SPACING.lg,
  },
  statCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.xl,
    padding: SPACING.base,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  statLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.md,
  },
  workoutsList: {
    gap: SPACING.md,
  },
  workoutCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    ...SHADOWS.base,
  },
  workoutHeader: {
    marginBottom: SPACING.md,
  },
  dayBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.background.auth,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    alignSelf: 'flex-start',
    marginBottom: SPACING.xs,
  },
  dayText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  focusText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  exercisesList: {
    gap: SPACING.sm,
  },
  exerciseRow: {
    flexDirection: 'row',
    gap: SPACING.sm,
    paddingVertical: SPACING.xs,
  },
  exerciseNumber: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
    width: 24,
  },
  exerciseInfo: {
    flex: 1,
  },
  exerciseName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 2,
  },
  exerciseDetails: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.secondary,
  },
  muscleGroups: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  emptyState: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.xl,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  emptyStateText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: SPACING.md,
  },
  errorText: {
    fontSize: FONT_SIZES.lg,
    color: COLORS.error,
    marginBottom: SPACING.md,
  },
  backButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.xl,
  },
  backButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.white,
  },
});
