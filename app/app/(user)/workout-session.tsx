import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  Alert,
  ActivityIndicator,
  Modal,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { GestureDetector, Gesture } from 'react-native-gesture-handler';
import * as Haptics from 'expo-haptics';
import { useTodayWorkout } from '@/hooks/user/use-today-workout';
import { useWorkoutCompletion, ExerciseSetLog } from '@/hooks/user/use-workout-completion';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface ExerciseSet {
  setNumber: number;
  reps: number;
  weight: number;
  completed: boolean;
  notes?: string;
}

interface ExerciseProgress {
  exerciseId: number;
  exerciseName: string;
  sets: ExerciseSet[];
  targetSets: number;
  targetReps: string;
  restSeconds: number;
}

export default function WorkoutSessionScreen() {
  const router = useRouter();
  const { data: todayWorkout, isLoading } = useTodayWorkout();
  const { mutate: saveWorkout, isPending: isSaving } = useWorkoutCompletion();
  
  const [currentExerciseIndex, setCurrentExerciseIndex] = useState(0);
  const [exerciseProgress, setExerciseProgress] = useState<ExerciseProgress[]>([]);
  const [startTime, setStartTime] = useState<Date>(new Date());
  const [isRestTimerActive, setIsRestTimerActive] = useState(false);
  const [restTimeRemaining, setRestTimeRemaining] = useState(0);
  const [showRestModal, setShowRestModal] = useState(false);
  const [workoutStarted, setWorkoutStarted] = useState(false);

  useEffect(() => {
    if (todayWorkout && todayWorkout.exercises && !workoutStarted) {
      const initialProgress: ExerciseProgress[] = todayWorkout.exercises.map((exercise, index) => ({
        exerciseId: exercise.exercise_id || index,
        exerciseName: exercise.name,
        targetSets: exercise.sets,
        targetReps: exercise.reps,
        restSeconds: exercise.rest_seconds,
        sets: Array.from({ length: exercise.sets }, (_, i) => ({
          setNumber: i + 1,
          reps: 0,
          weight: 0,
          completed: false,
        })),
      }));
      setExerciseProgress(initialProgress);
      setWorkoutStarted(true);
      setStartTime(new Date());
    }
  }, [todayWorkout, workoutStarted]);

  useEffect(() => {
    let interval: ReturnType<typeof setInterval> | undefined;
    if (isRestTimerActive && restTimeRemaining > 0) {
      interval = setInterval(() => {
        setRestTimeRemaining((prev) => {
          if (prev <= 1) {
            setIsRestTimerActive(false);
            setShowRestModal(false);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [isRestTimerActive, restTimeRemaining]);

  const currentExercise = exerciseProgress[currentExerciseIndex];
  const totalExercises = exerciseProgress.length;
  const completedExercises = exerciseProgress.filter((ex) => 
    ex.sets.every((set) => set.completed)
  ).length;

  const handleCompleteSet = useCallback((setIndex: number) => {
    const exercise = exerciseProgress[currentExerciseIndex];
    const isCompleting = !exercise.sets[setIndex].completed;

    setExerciseProgress((prev) => {
      const updated = [...prev];
      const ex = updated[currentExerciseIndex];
      ex.sets[setIndex].completed = !ex.sets[setIndex].completed;
      return updated;
    });

    if (isCompleting) {
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
      
      const allSetsCompleted = exercise.sets.every((set, idx) => 
        idx === setIndex || set.completed
      );
      
      if (allSetsCompleted) {
        setTimeout(() => {
          Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
        }, 200);
      }
    } else {
      Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    }

    if (
      isCompleting &&
      setIndex < exercise.sets.length - 1 &&
      exercise.restSeconds > 0
    ) {
      setRestTimeRemaining(exercise.restSeconds);
      setIsRestTimerActive(true);
      setShowRestModal(true);
    }
  }, [currentExerciseIndex, exerciseProgress]);

  const handleUpdateSet = useCallback((setIndex: number, field: 'reps' | 'weight', value: number) => {
    setExerciseProgress((prev) => {
      const updated = [...prev];
      updated[currentExerciseIndex].sets[setIndex][field] = value;
      return updated;
    });
  }, [currentExerciseIndex]);

  const handleUpdateSetNotes = useCallback((setIndex: number, notes: string) => {
    setExerciseProgress((prev) => {
      const updated = [...prev];
      updated[currentExerciseIndex].sets[setIndex].notes = notes;
      return updated;
    });
  }, [currentExerciseIndex]);

  const handleNextExercise = useCallback(() => {
    if (currentExerciseIndex < totalExercises - 1) {
      setCurrentExerciseIndex((prev) => prev + 1);
    }
  }, [currentExerciseIndex, totalExercises]);

  const handlePreviousExercise = useCallback(() => {
    if (currentExerciseIndex > 0) {
      setCurrentExerciseIndex((prev) => prev - 1);
    }
  }, [currentExerciseIndex]);

  const handleSkipExercise = useCallback(() => {
    Alert.alert(
      'Skip Exercise',
      'Are you sure you want to skip this exercise?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Skip',
          style: 'destructive',
          onPress: handleNextExercise,
        },
      ]
    );
  }, [handleNextExercise]);

  const handleFinishWorkout = useCallback(() => {
    if (!todayWorkout) return;

    const totalSets = exerciseProgress.reduce((acc, ex) => acc + ex.sets.length, 0);
    const completedSets = exerciseProgress.reduce(
      (acc, ex) => acc + ex.sets.filter((s) => s.completed).length,
      0
    );
    const completionRate = totalSets > 0 ? Math.round((completedSets / totalSets) * 100) : 0;

    Alert.alert(
      'Finish Workout',
      `You completed ${completedSets} out of ${totalSets} sets (${completionRate}%).\n\nFinish workout?`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Finish',
          onPress: () => {
            const exerciseLogs: ExerciseSetLog[] = [];
            exerciseProgress.forEach((exercise) => {
              exercise.sets.forEach((set) => {
                exerciseLogs.push({
                  exercise_id: exercise.exerciseId,
                  exercise_name: exercise.exerciseName,
                  set_number: set.setNumber,
                  reps: set.reps,
                  weight: set.weight,
                  completed: set.completed,
                  notes: set.notes,
                });
              });
            });

            const now = new Date();
            const durationSeconds = Math.floor((now.getTime() - startTime.getTime()) / 1000);

            saveWorkout(
              {
                plan_id: todayWorkout.plan_id,
                day_index: todayWorkout.day_index,
                duration_seconds: durationSeconds,
                completed_at: now.toISOString(),
                exercises: exerciseLogs,
              },
              {
                onSuccess: (response) => {
                  Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
                  setTimeout(() => {
                    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
                  }, 150);

                  Alert.alert(
                    'Workout Complete! 🎉',
                    `Great job! You completed ${response.completed_sets} sets with a ${response.completion_rate.toFixed(0)}% completion rate.\n\n` +
                    `Total Volume: ${response.total_volume.toFixed(0)} lbs\n` +
                    `Duration: ${response.duration_minutes} minutes\n` +
                    `Current Streak: ${response.new_streak} days${response.is_personal_best ? '\n\n🏆 New Personal Best!' : ''}`,
                    [
                      {
                        text: 'Done',
                        onPress: () => router.back(),
                      },
                    ]
                  );
                },
                onError: (error) => {
                  Alert.alert(
                    'Error',
                    'Failed to save workout. Please try again.',
                    [{ text: 'OK' }]
                  );
                },
              }
            );
          },
        },
      ]
    );
  }, [exerciseProgress, todayWorkout, startTime, saveWorkout, router]);

  const handleSkipRest = useCallback(() => {
    setIsRestTimerActive(false);
    setShowRestModal(false);
    setRestTimeRemaining(0);
  }, []);

  const swipeGesture = Gesture.Pan()
    .onEnd((event) => {
      const SWIPE_THRESHOLD = 50;
      const SWIPE_VELOCITY_THRESHOLD = 500;

      if (
        event.translationX < -SWIPE_THRESHOLD || 
        event.velocityX < -SWIPE_VELOCITY_THRESHOLD
      ) {
        if (currentExerciseIndex < totalExercises - 1) {
          Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
          handleNextExercise();
        }
      }
      else if (
        event.translationX > SWIPE_THRESHOLD || 
        event.velocityX > SWIPE_VELOCITY_THRESHOLD
      ) {
        if (currentExerciseIndex > 0) {
          Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
          handlePreviousExercise();
        }
      }
    })
    .runOnJS(true);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const getElapsedTime = () => {
    const now = new Date();
    const elapsed = Math.floor((now.getTime() - startTime.getTime()) / 1000);
    return formatTime(elapsed);
  };

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
        <Text style={styles.loadingText}>Loading workout...</Text>
      </View>
    );
  }

  if (!todayWorkout || !todayWorkout.exercises || todayWorkout.exercises.length === 0) {
    return (
      <SafeAreaView style={styles.container}>
        <Stack.Screen options={{ title: 'Workout Session', headerShown: true }} />
        <View style={styles.emptyContainer}>
          <Ionicons name="barbell-outline" size={64} color={COLORS.text.tertiary} />
          <Text style={styles.emptyTitle}>No Workout Available</Text>
          <Text style={styles.emptySubtitle}>
            Create a workout plan to get started
          </Text>
          <TouchableOpacity
            style={styles.createPlanButton}
            onPress={() => router.push('/plan-generator')}
          >
            <Text style={styles.createPlanButtonText}>Create Plan</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  if (!currentExercise) {
    return null;
  }

  return (
    <SafeAreaView style={styles.container} edges={['bottom']}>
      <Stack.Screen 
        options={{ 
          title: todayWorkout.day_title,
          headerShown: true,
        //   headerLeft: () => (
        //     // <TouchableOpacity onPress={() => router.back()}>
        //     //   <Ionicons name="close" size={28} color={COLORS.text.inverse} style={{ marginLeft: SPACING.md, padding: SPACING.md, backgroundColor: COLORS.primaryDark, borderRadius: BORDER_RADIUS.full}} />
        //     // </TouchableOpacity>
        //   ),
        }} 
      />

      {/* Progress Header */}
      <View style={styles.progressHeader}>
        <View style={styles.progressInfo}>
          <Text style={styles.progressText}>
            Exercise {currentExerciseIndex + 1} of {totalExercises}
          </Text>
          <Text style={styles.timeText}>
            <Ionicons name="time-outline" size={16} /> {getElapsedTime()}
          </Text>
        </View>
        <View style={styles.progressBarContainer}>
          <View 
            style={[
              styles.progressBarFill, 
              { width: `${((currentExerciseIndex + 1) / totalExercises) * 100}%` }
            ]} 
          />
        </View>
      </View>

      <ScrollView 
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <GestureDetector gesture={swipeGesture}>
          <View style={styles.gestureContainer}>
            <View style={styles.swipeIndicator}>
              <View style={styles.swipeHint}>
                {currentExerciseIndex > 0 && (
                  <View style={styles.swipeHintLeft}>
                    <Ionicons name="chevron-back" size={16} color={COLORS.text.tertiary} />
                    <Text style={styles.swipeHintText}>Swipe</Text>
                  </View>
                )}
                {currentExerciseIndex < totalExercises - 1 && (
                  <View style={styles.swipeHintRight}>
                    <Text style={styles.swipeHintText}>Swipe</Text>
                    <Ionicons name="chevron-forward" size={16} color={COLORS.text.tertiary} />
                  </View>
                )}
              </View>
            </View>

            <MotiView
              from={{ opacity: 0, translateY: -20 }}
              animate={{ opacity: 1, translateY: 0 }}
              transition={{ type: 'timing', duration: 300 }}
              key={currentExerciseIndex}
            >
          <View style={styles.exerciseHeader}>
            <View style={styles.exerciseIconContainer}>
              <Ionicons name="fitness" size={32} color={COLORS.primary} />
            </View>
            <Text style={styles.exerciseName}>{currentExercise.exerciseName}</Text>
            <Text style={styles.exerciseTarget}>
              {currentExercise.targetSets} sets × {currentExercise.targetReps} reps
            </Text>
          </View>
        </MotiView>

        <View style={styles.setsContainer}>
          {currentExercise.sets.map((set, index) => (
            <MotiView
              key={index}
              from={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ type: 'timing', duration: 200, delay: index * 50 }}
            >
              <View style={[styles.setCard, set.completed && styles.setCardCompleted]}>
                <View style={styles.setHeader}>
                  <View style={styles.setNumberContainer}>
                    <Text style={styles.setNumber}>Set {set.setNumber}</Text>
                    {set.completed && (
                      <Ionicons name="checkmark-circle" size={20} color={COLORS.success} />
                    )}
                  </View>
                  <TouchableOpacity
                    style={[
                      styles.completeButton,
                      set.completed && styles.completeButtonActive,
                    ]}
                    onPress={() => handleCompleteSet(index)}
                  >
                    <Text style={[
                      styles.completeButtonText,
                      set.completed && styles.completeButtonTextActive,
                    ]}>
                      {set.completed ? 'Completed' : 'Mark Complete'}
                    </Text>
                  </TouchableOpacity>
                </View>

                <View style={styles.setInputsRow}>
                  <View style={styles.inputGroup}>
                    <Text style={styles.inputLabel}>Reps</Text>
                    <View style={styles.inputWithButtons}>
                      <TouchableOpacity
                        style={styles.inputButton}
                        onPress={() => handleUpdateSet(index, 'reps', Math.max(0, set.reps - 1))}
                      >
                        <Ionicons name="remove" size={20} color={COLORS.text.sc.error} />
                      </TouchableOpacity>
                      <TextInput
                        style={styles.input}
                        value={set.reps.toString()}
                        onChangeText={(text) => {
                          const value = parseInt(text) || 0;
                          handleUpdateSet(index, 'reps', value);
                        }}
                        keyboardType="number-pad"
                        selectTextOnFocus
                      />
                      <TouchableOpacity
                        style={styles.inputButton}
                        onPress={() => handleUpdateSet(index, 'reps', set.reps + 1)}
                      >
                        <Ionicons name="add" size={20} color={COLORS.text.sc.success} />
                      </TouchableOpacity>
                    </View>
                  </View>

                  <View style={styles.inputGroup}>
                    <Text style={styles.inputLabel}>Weight (lbs)</Text>
                    <View style={styles.inputWithButtons}>
                      <TouchableOpacity
                        style={styles.inputButton}
                        onPress={() => handleUpdateSet(index, 'weight', Math.max(0, set.weight - 5))}
                      >
                        <Ionicons name="remove" size={20} color={COLORS.text.sc.error} />
                      </TouchableOpacity>
                      <TextInput
                        style={styles.input}
                        value={set.weight.toString()}
                        onChangeText={(text) => {
                          const value = parseInt(text) || 0;
                          handleUpdateSet(index, 'weight', value);
                        }}
                        keyboardType="number-pad"
                        selectTextOnFocus
                      />
                      <TouchableOpacity
                        style={styles.inputButton}
                        onPress={() => handleUpdateSet(index, 'weight', set.weight + 5)}
                      >
                        <Ionicons name="add" size={20} color={COLORS.text.sc.success} />
                      </TouchableOpacity>
                    </View>
                  </View>
                </View>

                <View style={styles.notesContainer}>
                  <View style={styles.notesHeader}>
                    <Ionicons name="document-text-outline" size={16} color={COLORS.text.secondary} />
                    <Text style={styles.notesLabel}>Notes (optional)</Text>
                  </View>
                  <TextInput
                    style={styles.notesInput}
                    value={set.notes || ''}
                    onChangeText={(text) => handleUpdateSetNotes(index, text)}
                    placeholder="Add notes about form, feeling, etc..."
                    placeholderTextColor={COLORS.text.tertiary}
                    multiline
                    numberOfLines={2}
                    maxLength={200}
                  />
                </View>
              </View>
            </MotiView>
          ))}
        </View>
          </View>
        </GestureDetector>

        <View style={styles.navigationContainer}>
          <TouchableOpacity
            style={[styles.navButton, currentExerciseIndex === 0 && styles.navButtonDisabled]}
            onPress={handlePreviousExercise}
            disabled={currentExerciseIndex === 0}
          >
            <Ionicons 
              name="chevron-back" 
              size={24} 
              color={currentExerciseIndex === 0 ? COLORS.text.tertiary : COLORS.primary} 
            />
            <Text style={[
              styles.navButtonText,
              currentExerciseIndex === 0 && styles.navButtonTextDisabled,
            ]}>
              Previous
            </Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.skipButton}
            onPress={handleSkipExercise}
          >
            <Text style={styles.skipButtonText}>Skip Exercise</Text>
          </TouchableOpacity>

          {currentExerciseIndex < totalExercises - 1 ? (
            <TouchableOpacity
              style={styles.navButton}
              onPress={handleNextExercise}
            >
              <Text style={styles.navButtonText}>Next</Text>
              <Ionicons name="chevron-forward" size={24} color={COLORS.primary} />
            </TouchableOpacity>
          ) : (
            <TouchableOpacity
              style={styles.finishButton}
              onPress={handleFinishWorkout}
            >
              <Ionicons name="checkmark-circle" size={24} color={COLORS.white} />
              <Text style={styles.finishButtonText}>Finish</Text>
            </TouchableOpacity>
          )}
        </View>
      </ScrollView>

      <Modal
        visible={showRestModal}
        transparent
        animationType="fade"
        onRequestClose={handleSkipRest}
      >
        <View style={styles.modalOverlay}>
          <MotiView
            from={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ type: 'spring', damping: 15 }}
            style={styles.restModal}
          >
            <View style={styles.restIconContainer}>
              <Ionicons name="timer" size={64} color={COLORS.primary} />
            </View>
            <Text style={styles.restTitle}>Rest Time</Text>
            <Text style={styles.restTimer}>{formatTime(restTimeRemaining)}</Text>
            <Text style={styles.restSubtitle}>Take a breather before the next set</Text>
            
            <TouchableOpacity
              style={styles.skipRestButton}
              onPress={handleSkipRest}
            >
              <Text style={styles.skipRestButtonText}>Skip Rest</Text>
            </TouchableOpacity>
          </MotiView>
        </View>
      </Modal>
    </SafeAreaView>
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
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginTop: SPACING.lg,
    marginBottom: SPACING.sm,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    textAlign: 'center',
    marginBottom: SPACING.xl,
  },
  createPlanButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
  },
  createPlanButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
  progressHeader: {
    borderRadius: BORDER_RADIUS.lg,
    margin: SPACING.md,
    padding: SPACING.lg,
    backgroundColor: COLORS.darkGray,
  },
  progressInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  progressText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.placeholder,
  },
  timeText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
  },
  progressBarContainer: {
    height: 4,
    backgroundColor: COLORS.border.light,
    borderRadius: BORDER_RADIUS.full,
    overflow: 'hidden',
  },
  progressBarFill: {
    height: '100%',
    backgroundColor: COLORS.primary,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.lg,
  },
  exerciseHeader: {
    alignItems: 'center',
    marginBottom: SPACING.xl,
  },
  exerciseIconContainer: {
    width: 80,
    height: 80,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: `${COLORS.primary}15`,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  exerciseName: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.primarySoft,
    textAlign: 'center',
    marginBottom: SPACING.xs,
  },
  exerciseTarget: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
  },
  setsContainer: {
    marginBottom: SPACING.xl,
  },
  setCard: {
    backgroundColor: COLORS.darkGray,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING.md,
    borderWidth: 2,
    borderColor: 'transparent',
    ...SHADOWS.sm,
  },
  setCardCompleted: {
    borderColor: COLORS.success,
    backgroundColor: `${COLORS.success}10`,
  },
  setHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  setNumberContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  setNumber: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
  },
  completeButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
  },
  completeButtonActive: {
    backgroundColor: COLORS.success,
  },
  completeButtonText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
  completeButtonTextActive: {
    color: COLORS.white,
  },
  setInputsRow: {
    flexDirection: 'row',
    gap: SPACING.md,
  },
  inputGroup: {
    flex: 1,
  },
  inputLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
  },
  inputWithButtons: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    shadowColor: COLORS.white,
    ...SHADOWS.sm,
    borderRadius: BORDER_RADIUS.md,
  },
  inputButton: {
    padding: SPACING.sm,
  },
  input: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    textAlign: 'center',
    paddingVertical: SPACING.sm,
  },
  notesContainer: {
    marginTop: SPACING.md,
    paddingTop: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.light,
  },
  notesHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginBottom: SPACING.xs,
  },
  notesLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
  },
  notesInput: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.sm,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.inverse,
    minHeight: 60,
    textAlignVertical: 'top',
  },
  gestureContainer: {
    flex: 1,
  },
  swipeIndicator: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  swipeHint: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',
  },
  swipeHintLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  swipeHintRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  swipeHintText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    fontWeight: '500' as any,
  },
  navigationContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: SPACING.md,
  },
  navButton: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
    backgroundColor: COLORS.darkGray,
    gap: SPACING.xs,
  },
  navButtonDisabled: {
    opacity: 0.4,
  },
  navButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.primary,
  },
  navButtonTextDisabled: {
    color: COLORS.text.tertiary,
  },
  skipButton: {
    padding: SPACING.md,

  },
  skipButtonText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  finishButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primaryDark,
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
    gap: SPACING.xs,
  },
  finishButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  restModal: {
    backgroundColor: COLORS.background.secondary,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING['2xl'],
    alignItems: 'center',
    width: '100%',
    maxWidth: 400,
    ...SHADOWS.lg,
  },
  restIconContainer: {
    marginBottom: SPACING.lg,
  },
  restTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.md,
  },
  restTimer: {
    fontSize: FONT_SIZES['6xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.primary,
    marginBottom: SPACING.sm,
  },
  restSubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.secondary,
    textAlign: 'center',
    marginBottom: SPACING.xl,
  },
  skipRestButton: {
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.light,
  },
  skipRestButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.secondary,
  },
});
