import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Modal,
  ScrollView,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { ManualWorkoutRequest, ManualExerciseRequest } from '@/types/schema';
import { ExercisePicker } from './exercise-picker';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { Button } from '../forms';

interface WorkoutEditorModalProps {
  visible: boolean;
  workout: ManualWorkoutRequest;
  onSave: (workout: ManualWorkoutRequest) => void;
  onClose: () => void;
}

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export const WorkoutEditorModal: React.FC<WorkoutEditorModalProps> = ({
  visible,
  workout,
  onSave,
  onClose,
}) => {
  const [workoutName, setWorkoutName] = useState(workout.workout_name || '');
  const [focus, setFocus] = useState(workout.focus || '');
  const [notes, setNotes] = useState(workout.notes || '');
  const [exercises, setExercises] = useState<ManualExerciseRequest[]>(workout.exercises || []);
  const [showExercisePicker, setShowExercisePicker] = useState(false);

  useEffect(() => {
    if (visible) {
      setWorkoutName(workout.workout_name || '');
      setFocus(workout.focus || '');
      setNotes(workout.notes || '');
      setExercises(workout.exercises || []);
    }
  }, [visible, workout]);

  const dayName = workout.day_of_week >= 1 && workout.day_of_week <= 7
    ? DAY_NAMES[(workout.day_of_week + 6) % 7]
    : `Day ${workout.day_of_week}`;

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
      day_of_week: workout.day_of_week,
      workout_name: workoutName.trim(),
      focus: focus.trim(),
      notes: notes.trim(),
      estimated_minutes: workout.estimated_minutes || 0,
      exercises,
    };

    onSave(updatedWorkout);
    onClose();
  };

  return (
    <Modal visible={visible} animationType="slide" onRequestClose={onClose}>
      <View style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={onClose} style={styles.closeButton}>
            <Ionicons name="close" size={28} color={COLORS.text.auth.primary} />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>{dayName} Workout</Text>
          <View style={styles.placeholder} />
        </View>

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

        <View style={styles.footer}>
          <Button title="Cancel" onPress={onClose} variant="outline" />
          <Button title="Save Workout" onPress={handleSave} />
        </View>

        <ExercisePicker
          visible={showExercisePicker}
          onClose={() => setShowExercisePicker(false)}
          onSelect={handleAddExercise}
        />
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.dark,
  },
  closeButton: {
    padding: SPACING.xs,
  },
  headerTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  placeholder: {
    width: 44,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
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
    borderRadius: BORDER_RADIUS.md,
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
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: COLORS.primary,
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.bold,
    textAlign: 'center',
    lineHeight: 28,
    marginRight: SPACING.sm,
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
    backgroundColor: COLORS.background.auth,
    borderRadius: BORDER_RADIUS.sm,
    padding: SPACING.sm,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.primary,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    textAlign: 'center',
  },
  footer: {
    flexDirection: 'row',
    gap: SPACING.md,
    padding: SPACING.base,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.dark,
  },
});
