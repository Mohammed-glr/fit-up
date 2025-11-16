import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
} from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import type { ManualWorkoutRequest, ManualExerciseRequest } from '@/types/schema';
import { ExercisePicker } from '@/components/schema/exercise-picker';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import { Button } from '@/components/forms';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useWorkoutEditorContext } from '@/context/workout-editor-context';

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export default function WorkoutEditorScreen() {
  const { setOnSaveWorkout, setIsSavingWorkout, setCurrentWorkout } = useWorkoutEditorContext();
  const params = useLocalSearchParams<{
    dayOfWeek?: string;
    workoutData?: string;
  }>();

  const dayOfWeek = params.dayOfWeek ? parseInt(params.dayOfWeek) : 1;
  const initialWorkout: ManualWorkoutRequest = params.workoutData
    ? JSON.parse(params.workoutData)
    : {
        day_of_week: dayOfWeek,
        workout_name: '',
        focus: '',
        notes: '',
        estimated_minutes: 0,
        exercises: [],
      };

  const [workoutName, setWorkoutName] = useState(initialWorkout.workout_name || '');
  const [focus, setFocus] = useState(initialWorkout.focus || '');
  const [notes, setNotes] = useState(initialWorkout.notes || '');
  const [exercises, setExercises] = useState<ManualExerciseRequest[]>(initialWorkout.exercises || []);
  const [showExercisePicker, setShowExercisePicker] = useState(false);

  const dayName = dayOfWeek >= 1 && dayOfWeek <= 7
    ? DAY_NAMES[(dayOfWeek + 6) % 7]
    : `Day ${dayOfWeek}`;

  const handleAddExercise = (exercise: any) => {
    const newExercise: ManualExerciseRequest = {
      exercise_id: exercise.exercise_id,
      sets: 3,
      reps: '10',
      rest_seconds: 60,
      weight: '',
      tempo: '',
      notes: '',
      order_index: exercises.length,
      is_superset: false,
      superset_group: 0,
    };
    setExercises([...exercises, newExercise]);
    setShowExercisePicker(false);
  };

  const handleRemoveExercise = (index: number) => {
    setExercises(exercises.filter((_, i) => i !== index));
  };

  const handleUpdateExercise = (index: number, updates: Partial<ManualExerciseRequest>) => {
    setExercises(exercises.map((ex, i) => (i === index ? { ...ex, ...updates } : ex)));
  };

  const handleSave = () => {
    if (!workoutName.trim()) {
      Alert.alert('Validation Error', 'Please enter a workout name');
      return;
    }

    if (!focus.trim()) {
      Alert.alert('Validation Error', 'Please enter a workout focus');
      return;
    }

    if (exercises.length === 0) {
      Alert.alert('Validation Error', 'Please add at least one exercise');
      return;
    }

    const updatedWorkout: ManualWorkoutRequest = {
      day_of_week: dayOfWeek,
      workout_name: workoutName.trim(),
      focus: focus.trim(),
      notes: notes.trim(),
      estimated_minutes: initialWorkout.estimated_minutes || 0,
      exercises,
    };

    setCurrentWorkout(updatedWorkout);
    router.back();
  };

  useEffect(() => {
    setOnSaveWorkout(() => handleSave);
    return () => {
      setOnSaveWorkout(undefined);
    };
  }, [workoutName, focus, notes, exercises, setOnSaveWorkout]);

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['bottom']}>
        <ScrollView style={styles.scrollView} contentContainerStyle={styles.scrollContent}>
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Workout Details</Text>
            
            <View style={styles.inputContainer}>
              <Text style={styles.label}>Workout Name *</Text>
              <TextInput
                style={styles.input}
                value={workoutName}
                onChangeText={setWorkoutName}
                placeholder="e.g., Upper Body Strength"
                placeholderTextColor={COLORS.text.tertiary}
              />
            </View>

            <View style={styles.inputContainer}>
              <Text style={styles.label}>Focus *</Text>
              <TextInput
                style={styles.input}
                value={focus}
                onChangeText={setFocus}
                placeholder="e.g., Chest & Triceps"
                placeholderTextColor={COLORS.text.tertiary}
              />
            </View>

            <View style={styles.inputContainer}>
              <Text style={styles.label}>Notes (Optional)</Text>
              <TextInput
                style={[styles.input, styles.textArea]}
                value={notes}
                onChangeText={setNotes}
                placeholder="Additional notes or instructions..."
                placeholderTextColor={COLORS.text.tertiary}
                multiline
                numberOfLines={3}
              />
            </View>
          </View>

          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>
                Exercises ({exercises.length})
              </Text>
              <TouchableOpacity
                style={styles.addButton}
                onPress={() => setShowExercisePicker(true)}
              >
                <Ionicons name="add-circle" size={24} color={COLORS.primary} />
                <Text style={styles.addButtonText}>Add Exercise</Text>
              </TouchableOpacity>
            </View>

            {exercises.length === 0 ? (
              <View style={styles.emptyState}>
                <Ionicons name="barbell-outline" size={48} color={COLORS.text.tertiary} />
                <Text style={styles.emptyText}>No exercises added yet</Text>
                <Text style={styles.emptySubtext}>Tap "Add Exercise" to get started</Text>
              </View>
            ) : (
              <View style={styles.exercisesList}>
                {exercises.map((exercise, index) => (
                  <View key={index} style={styles.exerciseCard}>
                    <View style={styles.exerciseHeader}>
                      <Text style={styles.exerciseNumber}>{index + 1}</Text>
                      <Text style={styles.exerciseName}>
                        Exercise ID: {exercise.exercise_id}
                      </Text>
                      <TouchableOpacity
                        onPress={() => handleRemoveExercise(index)}
                        style={styles.deleteButton}
                      >
                        <Ionicons name="trash-outline" size={20} color={COLORS.error} />
                      </TouchableOpacity>
                    </View>

                    <View style={styles.exerciseDetails}>
                      <View style={styles.detailRow}>
                        <View style={styles.detailInput}>
                          <Text style={styles.detailLabel}>Sets</Text>
                          <TextInput
                            style={styles.smallInput}
                            value={String(exercise.sets)}
                            onChangeText={(text) =>
                              handleUpdateExercise(index, { sets: parseInt(text) || 3 })
                            }
                            keyboardType="numeric"
                            placeholderTextColor={COLORS.text.tertiary}
                          />
                        </View>

                        <View style={styles.detailInput}>
                          <Text style={styles.detailLabel}>Reps</Text>
                          <TextInput
                            style={styles.smallInput}
                            value={exercise.reps}
                            onChangeText={(text) =>
                              handleUpdateExercise(index, { reps: text })
                            }
                            placeholder="10"
                            placeholderTextColor={COLORS.text.tertiary}
                          />
                        </View>

                        <View style={styles.detailInput}>
                          <Text style={styles.detailLabel}>Rest (sec)</Text>
                          <TextInput
                            style={styles.smallInput}
                            value={String(exercise.rest_seconds)}
                            onChangeText={(text) =>
                              handleUpdateExercise(index, { rest_seconds: parseInt(text) || 60 })
                            }
                            keyboardType="numeric"
                            placeholderTextColor={COLORS.text.tertiary}
                          />
                        </View>
                      </View>
                    </View>
                  </View>
                ))}
              </View>
            )}
          </View>
        </ScrollView>

        {/* <View style={styles.footer}>
          <Button title="Cancel" onPress={() => router.back()} variant="outline" />
          <Button title="Save Workout" onPress={handleSave} />
        </View> */}

        <ExercisePicker
          visible={showExercisePicker}
          onClose={() => setShowExercisePicker(false)}
          onSelect={handleAddExercise}
        />
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  safeArea: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
    paddingBottom: 100,
  },
  section: {
    marginBottom: SPACING['4xl'],
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
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
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  textArea: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  addButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  addButtonText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  emptyState: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['4xl'],
  },
  emptyText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.secondary,
    marginTop: SPACING.md,
  },
  emptySubtext: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
  },
  exercisesList: {
    gap: SPACING.md,
  },
  exerciseCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  exerciseHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  exerciseNumber: {
    borderRadius: 14,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
    marginRight: SPACING.sm,
    backgroundColor: COLORS.primary,
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.bold,
  },
  exerciseName: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
  },
  deleteButton: {
    padding: SPACING.xs,
  },
  exerciseDetails: {
    gap: SPACING.sm,
  },
  detailRow: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  detailInput: {
    flex: 1,
  },
  detailLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
  },
  smallInput: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.sm,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.primary,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 12 
    },
    shadowOpacity: 0.35,
    shadowRadius: 24,
    textAlign: 'center',
  },
  footer: {
    position: 'absolute',
    left: 16,
    right: 16,
    bottom: 16,
    flexDirection: 'row',
    gap: SPACING.md,
    padding: SPACING.base,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderColor: 'rgba(255, 255, 255, 0.1)',
    borderWidth: 1,
    elevation: 20,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 12 
    },
    shadowOpacity: 0.35,
    shadowRadius: 24,
  },
});
