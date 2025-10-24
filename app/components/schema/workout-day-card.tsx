import React from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { WorkoutDetail, ManualWorkoutRequest } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface WorkoutDayCardProps {
  dayOfWeek: number;
  workout?: WorkoutDetail | ManualWorkoutRequest;
  onEdit?: () => void;
  onDelete?: () => void;
  isRestDay?: boolean;
}

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export const WorkoutDayCard: React.FC<WorkoutDayCardProps> = ({
  dayOfWeek,
  workout,
  onEdit,
  onDelete,
  isRestDay = false,
}) => {
  const dayName = DAY_NAMES[dayOfWeek] || 'Unknown';

  const handleDelete = () => {
    Alert.alert(
      'Delete Workout',
      `Are you sure you want to delete the ${dayName} workout?`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: onDelete,
        },
      ]
    );
  };

  if (isRestDay || !workout) {
    return (
      <View style={styles.card}>
        <View style={styles.cardHeader}>
          <View style={styles.dayBadge}>
            <Text style={styles.dayNumber}>{dayOfWeek}</Text>
          </View>
          <View style={styles.headerInfo}>
            <Text style={styles.dayName}>{dayName}</Text>
            <Text style={styles.restDayLabel}>Rest Day</Text>
          </View>
          {onEdit && (
            <TouchableOpacity onPress={onEdit} style={styles.iconButton}>
              <Ionicons name="add-circle-outline" size={24} color={COLORS.primary} />
            </TouchableOpacity>
          )}
        </View>
      </View>
    );
  }

  const exercises = 'exercises' in workout ? workout.exercises : [];
  const exerciseCount = exercises?.length || 0;
  const workoutName = 'workout_name' in workout ? workout.workout_name : ('focus' in workout ? workout.focus : 'Workout');
  
  let estimatedTime: number | undefined;
  if ('estimated_minutes' in workout) {
    estimatedTime = workout.estimated_minutes;
  } else if ('estimated_min' in workout) {
    estimatedTime = workout.estimated_min as number | undefined;
  }
  
  const focus = 'focus' in workout ? workout.focus : '';

  return (
    <View style={styles.card}>
      <View style={styles.cardHeader}>
        <View style={styles.dayBadge}>
          <Text style={styles.dayNumber}>{dayOfWeek}</Text>
        </View>
        <View style={styles.headerInfo}>
          <Text style={styles.dayName}>{dayName}</Text>
          <Text style={styles.workoutName}>{workoutName}</Text>
          {focus && <Text style={styles.focusText}>{focus}</Text>}
        </View>
        <View style={styles.actions}>
          {onEdit && (
            <TouchableOpacity onPress={onEdit} style={styles.iconButton}>
              <Ionicons name="create-outline" size={24} color={COLORS.text.tertiary} />
            </TouchableOpacity>
          )}
          {onDelete && (
            <TouchableOpacity onPress={handleDelete} style={styles.iconButton}>
              <Ionicons name="trash-outline" size={24} color={COLORS.error} />
            </TouchableOpacity>
          )}
        </View>
      </View>

      <View style={styles.cardBody}>
        <View style={styles.statsRow}>
          <View style={styles.stat}>
            <Ionicons name="barbell" size={16} color={COLORS.primary} />
            <Text style={styles.statText}>{exerciseCount} exercises</Text>
          </View>
          {estimatedTime && (
            <View style={styles.stat}>
              <Ionicons name="time" size={16} color={COLORS.primary} />
              <Text style={styles.statText}>{estimatedTime} min</Text>
            </View>
          )}
        </View>

        {'exercises' in workout && workout.exercises && workout.exercises.length > 0 && (
          <View style={styles.exercisesList}>
            {workout.exercises.slice(0, 3).map((exercise, index) => {
              const exerciseName = 'exercise' in exercise ? exercise.exercise.name : 'Exercise';
              const sets = exercise.sets;
              const reps = exercise.reps;
              
              return (
                <View key={index} style={styles.exerciseItem}>
                  <View style={styles.exerciseMarker} />
                  <View style={styles.exerciseInfo}>
                    <Text style={styles.exerciseName}>{exerciseName}</Text>
                    <Text style={styles.exerciseDetails}>
                      {sets} sets × {reps} reps
                    </Text>
                  </View>
                </View>
              );
            })}
            {workout.exercises.length > 3 && (
              <Text style={styles.moreText}>
                +{workout.exercises.length - 3} more exercises
              </Text>
            )}
          </View>
        )}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.full,
    marginBottom: SPACING.md,
    overflow: 'hidden',
    ...SHADOWS.sm,
  },
  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: SPACING.base,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.dark,
  },
  dayBadge: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  dayNumber: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  headerInfo: {
    flex: 1,
  },
  dayName: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.placeholder,
    marginBottom: 2,
  },
  workoutName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  focusText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  restDayLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    fontStyle: 'italic',
  },
  actions: {
    flexDirection: 'row',
    gap: SPACING.xs,
  },
  iconButton: {
    padding: SPACING.xs,
  },
  cardBody: {
    padding: SPACING.base,
  },
  statsRow: {
    flexDirection: 'row',
    gap: SPACING.lg,
    marginBottom: SPACING.md,
  },
  stat: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  statText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.secondary,
  },
  exercisesList: {
    gap: SPACING.sm,
  },
  exerciseItem: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  exerciseMarker: {
    width: 4,
    height: 4,
    borderRadius: 2,
    backgroundColor: COLORS.primary,
    marginRight: SPACING.sm,
  },
  exerciseInfo: {
    flex: 1,
  },
  exerciseName: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
    marginBottom: 2,
  },
  exerciseDetails: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  moreText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.primary,
    marginTop: SPACING.xs,
    marginLeft: SPACING.base,
  },
});
