import React from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  Alert,
  SafeAreaView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { useAuth } from '@/context/auth-context';
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
  workouts?: WorkoutWithExercises[];
}

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export default function UserSchemaScreen() {
  const { user } = useAuth();
  const router = useRouter();
  const [refreshing, setRefreshing] = React.useState(false);
  const [selectedSchema, setSelectedSchema] = React.useState<WeeklySchema | null>(null);

  const userID = user?.id;

  const { data: schemas, isLoading, refetch } = useQuery({
    queryKey: ['user-schemas', userID],
    queryFn: async () => {
      if (!userID) throw new Error('User ID required');
      const response = await httpClient.get(`schemas/user/${userID}`);
      return response.data as { schemas: WeeklySchema[] };
    },
    enabled: !!userID,
  });

  const { data: schemaDetails, isLoading: isLoadingDetails } = useQuery({
    queryKey: ['schema-details', selectedSchema?.schema_id],
    queryFn: async () => {
      if (!selectedSchema?.schema_id) throw new Error('Schema ID required');
      const response = await httpClient.get(`schemas/${selectedSchema.schema_id}/workouts`);
      return response.data as WeeklySchema;
    },
    enabled: !!selectedSchema?.schema_id,
  });

  const onRefresh = async () => {
    setRefreshing(true);
    await refetch();
    setRefreshing(false);
  };

  const handleViewSchema = (schema: WeeklySchema) => {
    setSelectedSchema(schema);
  };

  const handleCloseDetails = () => {
    setSelectedSchema(null);
  };

  const activeSchema = schemas?.schemas?.find(s => s.active);
  const archivedSchemas = schemas?.schemas?.filter(s => !s.active) || [];

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (selectedSchema) {
    const displaySchema = schemaDetails || selectedSchema;
    const workouts = displaySchema.workouts || [];

    return (
      <ScrollView style={styles.container} contentContainerStyle={styles.scrollContent}>
        <View style={styles.header}>
          <View style={styles.headerInfo}>
            <Text style={styles.headerTitle}>
              Week of {new Date(displaySchema.week_start).toLocaleDateString()}
            </Text>
            <Text style={styles.headerSubtitle}>
              {displaySchema.active ? 'Active Schema' : 'Past Schema'}
            </Text>
          </View>
          {displaySchema.active && (
            <View style={styles.activeBadge}>
              <View style={styles.activeDot} />
              <Text style={styles.activeText}>Active</Text>
            </View>
          )}
        </View>

        {isLoadingDetails ? (
          <View style={styles.loadingSection}>
            <ActivityIndicator size="small" color={COLORS.primary} />
            <Text style={styles.loadingText}>Loading workouts...</Text>
          </View>
        ) : workouts.length === 0 ? (
          <View style={styles.emptyState}>
            <Ionicons name="calendar-outline" size={60} color={COLORS.text.tertiary} />
            <Text style={styles.emptyStateText}>No workouts in this schema</Text>
          </View>
        ) : (
          <View style={styles.workoutsList}>
            {workouts.map((workout) => (
              <View key={workout.workout_id} style={styles.workoutCard}>
                <View style={styles.workoutHeader}>
                  <View style={styles.dayBadge}>
                    <Ionicons name="calendar" size={16} color={COLORS.primary} />
                    <Text style={styles.dayText}>{DAY_NAMES[workout.day_of_week] || `Day ${workout.day_of_week}`}</Text>
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
                        {exercise.exercise.muscle_groups && exercise.exercise.muscle_groups.length > 0 && (
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

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          refreshControl={
            <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={COLORS.primary} />
          }
          showsVerticalScrollIndicator={false}
        >
          <View style={styles.header}>
            <View style={styles.titleContainer}>
              <Ionicons name="calendar" size={32} color={COLORS.primary} />
              <Text style={styles.title}>Training Schema</Text>
            </View>
            <Text style={styles.subtitle}>Coach-designed workout programs</Text>
          </View>

      {activeSchema ? (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Active Schema</Text>
          <TouchableOpacity
            style={styles.schemaCard}
            onPress={() => handleViewSchema(activeSchema)}
          >
            <View style={styles.schemaHeader}>
              <View style={styles.schemaIcon}>
                <Ionicons name="calendar" size={24} color={COLORS.primary} />
              </View>
              <View style={styles.schemaInfo}>
                <Text style={styles.schemaTitle}>
                  Week of {new Date(activeSchema.week_start).toLocaleDateString()}
                </Text>
                <Text style={styles.schemaMeta}>
                  {activeSchema.workouts?.length || 0} workouts
                </Text>
              </View>
              <View style={styles.activeBadge}>
                <View style={styles.activeDot} />
                <Text style={styles.activeText}>Active</Text>
              </View>
            </View>
          </TouchableOpacity>
        </View>
      ) : (
        <View style={styles.emptyState}>
          <Ionicons name="calendar-outline" size={60} color="#444444" />
          <Text style={styles.emptyStateTitle}>No Active Schema</Text>
          <Text style={styles.emptyStateText}>
            Your coach hasn't assigned a training schema yet
          </Text>
        </View>
      )}

      {archivedSchemas.length > 0 && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Past Schemas</Text>
          <View style={styles.schemasList}>
            {archivedSchemas.map((schema) => (
              <TouchableOpacity
                key={schema.schema_id}
                style={styles.schemaCard}
                onPress={() => handleViewSchema(schema)}
              >
                <View style={styles.schemaHeader}>
                  <View style={styles.schemaIcon}>
                    <Ionicons name="document-text" size={24} color={COLORS.text.auth.secondary} />
                  </View>
                  <View style={styles.schemaInfo}>
                    <Text style={styles.schemaTitle}>
                      Week of {new Date(schema.week_start).toLocaleDateString()}
                    </Text>
                    <Text style={styles.schemaMeta}>
                      {schema.workouts?.length || 0} workouts
                    </Text>
                  </View>
                  <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
                </View>
              </TouchableOpacity>
            ))}
          </View>
        </View>
      )}

          <View style={styles.bottomSpacer} />
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A0A0A',
  },
  safeArea: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: 20,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#0A0A0A',
  },
  header: {
    marginBottom: 24,
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
  section: {
    marginBottom: 32,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
    marginBottom: 16,
  },
  schemasList: {
    gap: 12,
  },
  schemaCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    marginBottom: 12,
  },
  schemaHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  schemaIcon: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.xl,
    backgroundColor: COLORS.background.auth,
    justifyContent: 'center',
    alignItems: 'center',
  },
  schemaInfo: {
    flex: 1,
  },
  schemaTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 4,
  },
  schemaMeta: {
    fontSize: 14,
    color: '#888888',
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
  emptyState: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: 32,
    alignItems: 'center',
    marginBottom: 24,
  },
  emptyStateTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
    marginTop: 16,
    marginBottom: 8,
  },
  emptyStateText: {
    fontSize: 14,
    color: '#888888',
    textAlign: 'center',
  },
  workoutsList: {
    gap: 12,
  },
  workoutCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: 12,
  },
  workoutHeader: {
    marginBottom: 16,
  },
  dayBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    backgroundColor: '#2A2A2A',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: BORDER_RADIUS.full,
    alignSelf: 'flex-start',
    marginBottom: 8,
  },
  dayText: {
    fontSize: 12,
    fontWeight: '500',
    color: COLORS.primary,
  },
  focusText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  exercisesList: {
    gap: 12,
  },
  exerciseRow: {
    flexDirection: 'row',
    gap: 12,
    paddingVertical: 8,
  },
  exerciseNumber: {
    fontSize: 16,
    fontWeight: '700',
    color: COLORS.primary,
    width: 24,
  },
  exerciseInfo: {
    flex: 1,
  },
  exerciseName: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 4,
  },
  exerciseDetails: {
    fontSize: 14,
    color: '#888888',
  },
  muscleGroups: {
    fontSize: 12,
    color: '#666666',
    marginTop: 4,
  },
  loadingSection: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 12,
    paddingVertical: 32,
  },
  loadingText: {
    fontSize: 14,
    color: '#888888',
  },
  bottomSpacer: {
    height: 40,
  },
});
