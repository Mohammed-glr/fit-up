import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { Button } from '../forms/button';

interface WorkoutSession {
  session_id: number;
  workout_title: string;
  completed_at: string;
  duration_minutes: number;
  total_exercises: number;
  total_volume_lbs: number;
}

interface WorkoutAttachmentPickerProps {
  visible: boolean;
  onClose: () => void;
  onSelectWorkout: (sessionId: number) => void;
  recentWorkouts?: WorkoutSession[];
  isLoading?: boolean;
}

export const WorkoutAttachmentPicker: React.FC<WorkoutAttachmentPickerProps> = ({
  visible,
  onClose,
  onSelectWorkout,
  recentWorkouts = [],
  isLoading = false,
}) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const handleSelect = (sessionId: number) => {
    onSelectWorkout(sessionId);
    onClose();
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="slide"
      onRequestClose={onClose}
    >
      <View style={styles.overlay}>
        <View style={styles.modalContainer}>
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.title}>Share Workout</Text>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.inverse} />
            </TouchableOpacity>
          </View>

          {/* Content */}
          {isLoading ? (
            <View style={styles.loadingContainer}>
              <ActivityIndicator size="large" color={COLORS.primary} />
              <Text style={styles.loadingText}>Loading recent workouts...</Text>
            </View>
          ) : recentWorkouts.length === 0 ? (
            <View style={styles.emptyContainer}>
              <Ionicons name="barbell-outline" size={64} color={COLORS.text.tertiary} />
              <Text style={styles.emptyTitle}>No Recent Workouts</Text>
              <Text style={styles.emptySubtitle}>
                Complete a workout to share it with your coach
              </Text>
            </View>
          ) : (
            <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
              <Text style={styles.sectionTitle}>Recent Workouts</Text>
              {recentWorkouts.map((workout) => (
                <TouchableOpacity
                  key={workout.session_id}
                  style={styles.workoutCard}
                  onPress={() => handleSelect(workout.session_id)}
                >
                  <View style={styles.workoutHeader}>
                    <View style={styles.workoutIconContainer}>
                      <Ionicons name="fitness" size={24} color={COLORS.primary} />
                    </View>
                    <View style={styles.workoutInfo}>
                      <Text style={styles.workoutTitle}>{workout.workout_title}</Text>
                      <Text style={styles.workoutDate}>{formatDate(workout.completed_at)}</Text>
                    </View>
                    <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
                  </View>
                  
                  <View style={styles.workoutStats}>
                    <View style={styles.statBadge}>
                      <Ionicons name="time-outline" size={16} color={COLORS.text.secondary} />
                      <Text style={styles.statText}>{workout.duration_minutes} min</Text>
                    </View>
                    <View style={styles.statBadge}>
                      <Ionicons name="barbell-outline" size={16} color={COLORS.text.secondary} />
                      <Text style={styles.statText}>{workout.total_exercises} exercises</Text>
                    </View>
                    <View style={styles.statBadge}>
                      <Ionicons name="trending-up-outline" size={16} color={COLORS.text.secondary} />
                      <Text style={styles.statText}>{workout.total_volume_lbs.toFixed(0)} lbs</Text>
                    </View>
                  </View>
                </TouchableOpacity>
              ))}
            </ScrollView>
          )}

          {/* Footer */}
          {!isLoading && (
            <View style={styles.footer}>
              <Button
                onPress={onClose}
                title="Cancel"
                variant="outline"
              />
            </View>
          )}
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  modalContainer: {
    backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['2xl'],
    borderTopRightRadius: BORDER_RADIUS['2xl'],
    height: '80%',
    ...SHADOWS.lg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.lg,
  },
  title: {
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
  loadingContainer: {
    padding: SPACING['3xl'],
    alignItems: 'center',
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
  },
  emptyContainer: {
    padding: SPACING['3xl'],
    alignItems: 'center',
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginTop: SPACING.lg,
    marginBottom: SPACING.sm,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    textAlign: 'center',
  },
  content: {
    flex: 1,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.placeholder,
    textTransform: 'uppercase',
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.lg,
    paddingBottom: SPACING.md,
  },
  workoutCard: {
    marginHorizontal: SPACING.lg,
    marginBottom: SPACING.md,
    padding: SPACING.md,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    ...SHADOWS.sm,
  },
  workoutHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  workoutIconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: `${COLORS.primary}20`,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: SPACING.md,
  },
  workoutInfo: {
    flex: 1,
  },
  workoutTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    marginBottom: SPACING.xs / 2,
  },
  workoutDate: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
  },
  workoutStats: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  statBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs / 2,
    backgroundColor: COLORS.background.secondary,
    borderRadius: BORDER_RADIUS.md,
    gap: SPACING.xs / 2,
  },
  statText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
  },
  footer: {
    padding: SPACING.lg,
  },
});
