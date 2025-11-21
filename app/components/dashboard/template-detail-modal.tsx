import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal,
  ScrollView,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import type { UserWorkoutTemplate } from '@/types/workout-template';

interface TemplateDetailModalProps {
  visible: boolean;
  template: UserWorkoutTemplate | null;
  onClose: () => void;
  onStartWorkout?: (template: UserWorkoutTemplate) => void;
  onEdit?: (template: UserWorkoutTemplate) => void;
  canEdit?: boolean;
}

export const TemplateDetailModal: React.FC<TemplateDetailModalProps> = ({
  visible,
  template,
  onClose,
  onStartWorkout,
  onEdit,
  canEdit = false,
}) => {
  if (!template) return null;

  const handleStartWorkout = () => {
    onClose();
    if (onStartWorkout) {
      onStartWorkout(template);
    } else {
      Alert.alert(
        'Start Workout',
        'Workout session integration coming soon!',
        [{ text: 'OK' }]
      );
    }
  };

  const handleEdit = () => {
    onClose();
    if (onEdit) {
      onEdit(template);
    }
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={onClose}
    >
      <View style={styles.modalOverlay}>
        <View style={styles.modalContent}>
          {/* Header */}
          <View style={styles.modalHeader}>
            <View style={styles.headerLeft}>
              <Text style={styles.modalTitle}>{template.name}</Text>
              {template.is_public && (
                <View style={styles.publicBadge}>
                  <Ionicons name="globe-outline" size={12} color={COLORS.primary} />
                  <Text style={styles.publicBadgeText}>Public</Text>
                </View>
              )}
            </View>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.inverse} />
            </TouchableOpacity>
          </View>

          {/* Content */}
          <ScrollView style={styles.modalBody} showsVerticalScrollIndicator={false}>
            {template.description && (
              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Description</Text>
                <Text style={styles.description}>{template.description}</Text>
              </View>
            )}

            <View style={styles.section}>
              <Text style={styles.sectionTitle}>
                Exercises ({template.exercises.length})
              </Text>
              {template.exercises.map((exercise, index) => (
                <View key={index} style={styles.exerciseCard}>
                  <View style={styles.exerciseHeader}>
                    <Text style={styles.exerciseName}>
                      {index + 1}. {exercise.exercise_name}
                    </Text>
                  </View>
                  <View style={styles.exerciseDetails}>
                    {exercise.sets != null && (
                      <View style={styles.exerciseDetail}>
                        <Ionicons name="repeat-outline" size={16} color={COLORS.text.tertiary} />
                        <Text style={styles.exerciseDetailText}>
                          {exercise.sets} sets
                        </Text>
                      </View>
                    )}
                    {exercise.target_reps != null && (
                      <View style={styles.exerciseDetail}>
                        <Ionicons name="fitness-outline" size={16} color={COLORS.text.tertiary} />
                        <Text style={styles.exerciseDetailText}>
                          {exercise.target_reps} reps
                        </Text>
                      </View>
                    )}
                    {exercise.target_weight != null && exercise.target_weight > 0 && (
                      <View style={styles.exerciseDetail}>
                        <Ionicons name="barbell-outline" size={16} color={COLORS.text.tertiary} />
                        <Text style={styles.exerciseDetailText}>
                          {exercise.target_weight} lbs
                        </Text>
                      </View>
                    )}
                    {exercise.rest_seconds != null && (
                      <View style={styles.exerciseDetail}>
                        <Ionicons name="time-outline" size={16} color={COLORS.text.tertiary} />
                        <Text style={styles.exerciseDetailText}>
                          {exercise.rest_seconds}s rest
                        </Text>
                      </View>
                    )}
                  </View>
                </View>
              ))}
            </View>

            <View style={styles.metaSection}>
              <Text style={styles.metaText}>
                Created: {new Date(template.created_at).toLocaleDateString()}
              </Text>
              <Text style={styles.metaText}>
                Updated: {new Date(template.updated_at).toLocaleDateString()}
              </Text>
            </View>
          </ScrollView>

          {/* Footer Actions */}
          <View style={styles.modalFooter}>
            {canEdit && onEdit && (
              <TouchableOpacity style={styles.editButton} onPress={handleEdit}>
                <Ionicons name="create-outline" size={20} color={COLORS.primary} />
                <Text style={styles.editButtonText}>Edit</Text>
              </TouchableOpacity>
            )}
            <TouchableOpacity style={styles.startButton} onPress={handleStartWorkout}>
              <Ionicons name="play-circle" size={20} color={COLORS.white} />
              <Text style={styles.startButtonText}>Start Workout</Text>
            </TouchableOpacity>
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

    paddingBottom: SPACING.xl,
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    padding: SPACING.lg,
  },
  headerLeft: {
    flex: 1,
    gap: SPACING.xs,
  },
  modalTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
  },
  publicBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: `${COLORS.primary}15`,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.md,
    alignSelf: 'flex-start',
  },
  publicBadgeText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.semibold as any,
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
  description: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    lineHeight: 22,
  },
  exerciseCard: {
    backgroundColor: COLORS.primarySoft,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.md,
    marginBottom: SPACING.sm,
  },
  exerciseHeader: {
    marginBottom: SPACING.sm,
  },
  exerciseName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
  },
  exerciseDetails: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  exerciseDetail: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  exerciseDetailText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  metaSection: {
    paddingTop: SPACING.md,
    gap: SPACING.xs,
  },
  metaText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  modalFooter: {
    flexDirection: 'row',
    gap: SPACING.md,
    padding: SPACING.lg,
    paddingTop: SPACING.md,
  },
  editButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.background.secondary,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.primary,
  },
  editButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.primary,
  },
  startButton: {
    flex: 2,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.primary,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
  },
  startButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
});
