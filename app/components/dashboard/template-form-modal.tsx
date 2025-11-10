import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal,
  ScrollView,
  TouchableOpacity,
  TextInput,
  Alert,
  Switch,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import type { CreateUserTemplateRequest, TemplateExercise, UserWorkoutTemplate } from '@/types/workout-template';
import { Button } from '../../components/forms/button';

interface TemplateFormModalProps {
  visible: boolean;
  onClose: () => void;
  onSave: (data: CreateUserTemplateRequest) => void;
  initialData?: UserWorkoutTemplate | null;
  isLoading?: boolean;
}

export const TemplateFormModal: React.FC<TemplateFormModalProps> = ({
  visible,
  onClose,
  onSave,
  initialData,
  isLoading = false,
}) => {
  const [name, setName] = useState(initialData?.name || '');
  const [description, setDescription] = useState(initialData?.description || '');
  const [isPublic, setIsPublic] = useState(initialData?.is_public || false);
  const [exercises, setExercises] = useState<TemplateExercise[]>(
    initialData?.exercises || []
  );
  const [currentExercise, setCurrentExercise] = useState<Partial<TemplateExercise>>({
    exercise_name: '',
    sets: 3,
    target_reps: 10,
    target_weight: undefined,
    rest_seconds: 60,
  });

  const resetForm = () => {
    setName('');
    setDescription('');
    setIsPublic(false);
    setExercises([]);
    setCurrentExercise({
      exercise_name: '',
      sets: 3,
      target_reps: 10,
      target_weight: undefined,
      rest_seconds: 60,
    });
  };

  const handleClose = () => {
    if (!initialData) {
      resetForm();
    }
    onClose();
  };

  const handleAddExercise = () => {
    if (!currentExercise.exercise_name?.trim()) {
      Alert.alert('Error', 'Please enter exercise name');
      return;
    }

    const newExercise: TemplateExercise = {
      exercise_name: currentExercise.exercise_name.trim(),
      sets: currentExercise.sets || 3,
      target_reps: currentExercise.target_reps || 10,
      target_weight: currentExercise.target_weight,
      rest_seconds: currentExercise.rest_seconds || 60,
    };

    setExercises([...exercises, newExercise]);
    setCurrentExercise({
      exercise_name: '',
      sets: 3,
      target_reps: 10,
      target_weight: undefined,
      rest_seconds: 60,
    });
  };

  const handleRemoveExercise = (index: number) => {
    setExercises(exercises.filter((_, i) => i !== index));
  };

  const handleSave = () => {
    if (!name.trim()) {
      Alert.alert('Error', 'Please enter a template name');
      return;
    }

    if (exercises.length === 0) {
      Alert.alert('Error', 'Please add at least one exercise');
      return;
    }

    const templateData: CreateUserTemplateRequest = {
      name: name.trim(),
      description: description.trim() || undefined,
      is_public: isPublic,
      exercises,
    };

    onSave(templateData);
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={handleClose}
    >
      <View style={styles.modalOverlay}>
        <View style={styles.modalContent}>
          {/* Header */}
          <View style={styles.modalHeader}>
            <Text style={styles.modalTitle}>
              {initialData ? 'Edit Template' : 'Create Template'}
            </Text>
            <TouchableOpacity onPress={handleClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.inverse} />
            </TouchableOpacity>
          </View>

          {/* Content */}
          <ScrollView style={styles.modalBody} showsVerticalScrollIndicator={false}>
            {/* Template Info */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Template Information</Text>
              
              <Text style={styles.label}>Name *</Text>
              <TextInput
                style={styles.input}
                placeholder="e.g., Upper Body Strength"
                placeholderTextColor={COLORS.text.tertiary}
                value={name}
                onChangeText={setName}
              />

              <Text style={styles.label}>Description</Text>
              <TextInput
                style={[styles.input, styles.textArea]}
                placeholder="Describe this workout template..."
                placeholderTextColor={COLORS.text.tertiary}
                value={description}
                onChangeText={setDescription}
                multiline
                numberOfLines={3}
              />

              <View style={styles.switchContainer}>
                <View>
                  <Text style={styles.label}>Make Public</Text>
                  <Text style={styles.hint}>Allow others to use this template</Text>
                </View>
                <Switch
                  value={isPublic}
                  onValueChange={setIsPublic}
                  trackColor={{ false: COLORS.border.light, true: COLORS.primary }}
                  thumbColor={COLORS.white}
                />
              </View>
            </View>

            {/* Current Exercises */}
            {exercises.length > 0 && (
              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Exercises ({exercises.length})</Text>
                {exercises.map((exercise, index) => (
                  <View key={index} style={styles.exerciseCard}>
                    <View style={styles.exerciseCardHeader}>
                      <Text style={styles.exerciseCardName}>
                        {index + 1}. {exercise.exercise_name}
                      </Text>
                      <TouchableOpacity onPress={() => handleRemoveExercise(index)}>
                        <Ionicons name="trash-outline" size={20} color={COLORS.error} />
                      </TouchableOpacity>
                    </View>
                    <Text style={styles.exerciseCardDetails}>
                      {exercise.sets} sets × {exercise.target_reps} reps
                      {exercise.target_weight ? ` @ ${exercise.target_weight} lbs` : ''}
                      {' • '}
                      {exercise.rest_seconds}s rest
                    </Text>
                  </View>
                ))}
              </View>
            )}

            {/* Add Exercise */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Add Exercise</Text>
              
              <Text style={styles.label}>Exercise Name *</Text>
              <TextInput
                style={styles.input}
                placeholder="e.g., Bench Press"
                placeholderTextColor={COLORS.text.tertiary}
                value={currentExercise.exercise_name}
                onChangeText={(text) =>
                  setCurrentExercise({ ...currentExercise, exercise_name: text })
                }
              />

              <View style={styles.row}>
                <View style={styles.rowItem}>
                  <Text style={styles.label}>Sets</Text>
                  <TextInput
                    style={styles.input}
                    placeholder="3"
                    placeholderTextColor={COLORS.text.tertiary}
                    value={currentExercise.sets?.toString()}
                    onChangeText={(text) =>
                      setCurrentExercise({
                        ...currentExercise,
                        sets: parseInt(text) || 0,
                      })
                    }
                    keyboardType="numeric"
                  />
                </View>

                <View style={styles.rowItem}>
                  <Text style={styles.label}>Reps</Text>
                  <TextInput
                    style={styles.input}
                    placeholder="10"
                    placeholderTextColor={COLORS.text.tertiary}
                    value={currentExercise.target_reps}
                    onChangeText={(text) =>
                      setCurrentExercise({ ...currentExercise, target_reps: text })
                    }
                  />
                </View>
              </View>

              <View style={styles.row}>
                <View style={styles.rowItem}>
                  <Text style={styles.label}>Weight (lbs)</Text>
                  <TextInput
                    style={styles.input}
                    placeholder="Optional"
                    placeholderTextColor={COLORS.text.tertiary}
                    value={currentExercise.target_weight?.toString()}
                    onChangeText={(text) =>
                      setCurrentExercise({
                        ...currentExercise,
                        target_weight: text ? parseInt(text) : undefined,
                      })
                    }
                    keyboardType="numeric"
                  />
                </View>

                <View style={styles.rowItem}>
                  <Text style={styles.label}>Rest (seconds)</Text>
                  <TextInput
                    style={styles.input}
                    placeholder="60"
                    placeholderTextColor={COLORS.text.tertiary}
                    value={currentExercise.rest_seconds?.toString()}
                    onChangeText={(text) =>
                      setCurrentExercise({
                        ...currentExercise,
                        rest_seconds: parseInt(text) || 60,
                      })
                    }
                    keyboardType="numeric"
                  />
                </View>
              </View>

              <TouchableOpacity style={styles.addExerciseButton} onPress={handleAddExercise}>
                <Ionicons name="add-circle-outline" size={20} color={COLORS.primary} />
                <Text style={styles.addExerciseButtonText}>Add Exercise</Text>
              </TouchableOpacity>
            </View>
          </ScrollView>

          <View style={styles.modalFooter}>
            <Button
              onPress={handleClose}
              title="Cancel"
              variant="outline"
            />
            <Button
              style={[styles.saveButton, isLoading && styles.saveButtonDisabled]}
              onPress={handleSave}
              disabled={isLoading}
              title={isLoading ? 'Saving...' : 'Save Template'}
            />
          </View>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['2xl'],
    borderTopRightRadius: BORDER_RADIUS['2xl'],
    height: '60%',
    width: '98%',
    alignSelf: 'center',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.lg,
  },
  modalTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
  },
  closeButton: {
     backgroundColor: COLORS.background.accent,
            padding: SPACING.md,
            borderRadius: BORDER_RADIUS.full,
            justifyContent: 'center',
            alignItems: 'center',
            minWidth: 40,
            minHeight: 40,
            
  },
  modalBody: {
    flex: 1,
    padding: SPACING.lg,
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  label: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
    marginTop: SPACING.sm,
  },
  hint: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  input: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  textArea: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  switchContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: SPACING.md,
  },
  row: {
    flexDirection: 'row',
    gap: SPACING.md,
  },
  rowItem: {
    flex: 1,
  },
  exerciseCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.md,
    marginBottom: SPACING.sm,
  },
  exerciseCardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.xs,
  },
  exerciseCardName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,

    flex: 1,
  },
  exerciseCardDetails: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  addExerciseButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.background.auth,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.primary,
    borderStyle: 'dashed',
    marginTop: SPACING.md,
  },
  addExerciseButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.primary,
  },
  modalFooter: {
    flexDirection: 'row',
    gap: SPACING.md,
    padding: SPACING.lg,
    paddingTop: SPACING.md,
  },
  cancelButton: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.secondary,
  },
  cancelButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.secondary,
  },
  saveButton: {
    flex: 2,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
  },
  saveButtonDisabled: {
    opacity: 0.5,
  },
  saveButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
});
