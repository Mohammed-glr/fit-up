import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { ManualSchemaRequest, ManualWorkoutRequest, ManualExerciseRequest } from '@/types/schema';
import { WorkoutDayCard } from './workout-day-card';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { Button } from '../forms';
import { useRouter } from 'expo-router';
import { useWorkoutEditorContext } from '@/context/workout-editor-context';

interface SchemaFormProps {
  initialData?: ManualSchemaRequest;
  onSubmit: (data: ManualSchemaRequest) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

const EMPTY_WORKOUT: ManualWorkoutRequest = {
  day_of_week: 1,
  workout_name: '',
  focus: '',
  exercises: [],
  notes: '',
  estimated_minutes: 0,
};

export const SchemaForm: React.FC<SchemaFormProps> = ({
  initialData,
  onSubmit,
  onCancel,
  isLoading = false,
}) => {
  const router = useRouter();
  const { currentWorkout, setCurrentWorkout } = useWorkoutEditorContext();
  const [schemaName, setSchemaName] = useState(initialData?.name || '');
  const [description, setDescription] = useState(initialData?.description || '');
  const [workouts, setWorkouts] = useState<ManualWorkoutRequest[]>(
    initialData?.workouts || Array.from({ length: 7 }, (_, i) => ({
      ...EMPTY_WORKOUT,
      day_of_week: i + 1,
    }))
  );

  const handleWorkoutChange = (dayOfWeek: number, updates: Partial<ManualWorkoutRequest>) => {
    setWorkouts((prev) =>
      prev.map((w) => (w.day_of_week === dayOfWeek ? { ...w, ...updates } : w))
    );
  };

  const handleAddExercises = (dayOfWeek: number, exerciseIds: number[]) => {
    const newExercises: ManualExerciseRequest[] = exerciseIds.map((id) => ({
      exercise_id: id,
      sets: 3,
      reps: '10',
      rest_seconds: 60,
    }));

    setWorkouts((prev) =>
      prev.map((w) =>
        w.day_of_week === dayOfWeek
          ? { ...w, exercises: [...(w.exercises || []), ...newExercises] }
          : w
      )
    );
  };

  const handleRemoveExercise = (dayOfWeek: number, exerciseIndex: number) => {
    setWorkouts((prev) =>
      prev.map((w) =>
        w.day_of_week === dayOfWeek
          ? {
              ...w,
              exercises: w.exercises?.filter((_, i) => i !== exerciseIndex) || [],
            }
          : w
      )
    );
  };

  const handleSubmit = () => {
    if (!schemaName.trim()) {
      Alert.alert('Validation Error', 'Please enter a schema name');
      return;
    }

    // Filter out empty workouts and validate
    const validWorkouts = workouts.filter(
      (w) => w.exercises && w.exercises.length > 0 && w.workout_name.trim() && w.focus.trim()
    );

    if (validWorkouts.length === 0) {
      Alert.alert('Validation Error', 'Please add at least one workout with a name, focus, and exercises');
      return;
    }

    const data: ManualSchemaRequest = {
      user_id: initialData?.user_id || 0,
      coach_id: initialData?.coach_id || '',
      name: schemaName.trim(),
      description: description.trim(),
      start_date: initialData?.start_date || new Date().toISOString().split('T')[0],
      workouts: validWorkouts,
    };

    onSubmit(data);
  };

  const openDayEditor = (dayOfWeek: number) => {
    const workout = workouts.find((w) => w.day_of_week === dayOfWeek);
    if (workout) {
      router.push({
        pathname: '/(coach)/workout-editor',
        params: {
          dayOfWeek: dayOfWeek.toString(),
          workoutData: JSON.stringify(workout),
        },
      });
    }
  };

  // Listen for returned workout data from workout-editor
  useEffect(() => {
    if (currentWorkout) {
      setWorkouts((prev) =>
        prev.map((w) =>
          w.day_of_week === currentWorkout.day_of_week ? currentWorkout : w
        )
      );
      setCurrentWorkout(null);
    }
  }, [currentWorkout, setCurrentWorkout]);

  return (
    <View style={styles.container}>
      <ScrollView style={styles.scrollView} contentContainerStyle={styles.scrollContent}>
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Schema Information</Text>
          <View style={styles.inputContainer}>
            <Text style={styles.label}>Schema Name *</Text>
            <TextInput
              style={styles.input}
              value={schemaName}
              onChangeText={setSchemaName}
              placeholder="e.g., Full Body Workout"
              placeholderTextColor={COLORS.text.tertiary}
            />
          </View>
          <View style={styles.inputContainer}>
            <Text style={styles.label}>Description</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              value={description}
              onChangeText={setDescription}
              placeholder="Brief description of the program..."
              placeholderTextColor={COLORS.text.tertiary}
              multiline
              numberOfLines={3}
            />
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Weekly Schedule</Text>
          {workouts.map((workout) => (
            <WorkoutDayCard
              key={workout.day_of_week}
              dayOfWeek={workout.day_of_week}
              workout={workout.exercises && workout.exercises.length > 0 ? workout : undefined}
              onEdit={() => openDayEditor(workout.day_of_week)}
              isRestDay={!workout.exercises || workout.exercises.length === 0}
            />
          ))}
        </View>
      </ScrollView>

      <View style={styles.footer}>
        <Button
          title={isLoading ? 'Canceling...' : 'Cancel'}
          onPress={onCancel}
          disabled={isLoading}
          variant="outline"
          
        />
        <Button
          title={isLoading ? 'Saving...' : 'Save Schema'}
          onPress={handleSubmit}
          disabled={isLoading}
          variant="primary"
          style={{ flex: 1 }}
        />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
  },
  section: {
    marginBottom: SPACING['6xl'],
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.base,
  },
  inputContainer: {
    marginBottom: SPACING.base,
  },
  label: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.secondary,
    marginBottom: SPACING.xs,
  },
  input: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  textArea: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  footer: {
       position: 'absolute',
        left: 0,
        right: 0,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        paddingHorizontal: SPACING.md,
        bottom: Platform.OS === 'ios' ? 20 : 16,
        marginHorizontal: 16,
        height: Platform.OS === 'ios' ? 76 : 76,
        borderRadius: BORDER_RADIUS.full,
        backgroundColor: COLORS.background.card,
        elevation: 20,
        shadowColor: '#000',
        shadowOffset: { 
          width: 0, 
          height: 12 
        },
        gap: SPACING.md,
        shadowOpacity: 0.35,
        shadowRadius: 24,
        overflow: 'visible',
        backdropFilter: 'blur(30px)',

  },
});
