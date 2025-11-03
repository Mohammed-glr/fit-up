import React, { useCallback, useMemo, useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

import { useAuth } from '@/context/auth-context';
import { APIError } from '@/api/client';
import { useCreatePlan } from '@/hooks/schema/use-plans';
import type { EquipmentType, FitnessGoal, FitnessLevel } from '@/types/schema';
import { Button } from '@/components/forms';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

const GOAL_OPTIONS: Array<{ label: string; value: FitnessGoal; description: string }> = [
  { label: 'Muscle Gain', value: 'muscle_gain', description: 'Build lean muscle mass' },
  { label: 'Strength', value: 'strength', description: 'Increase overall strength' },
  { label: 'Fat Loss', value: 'fat_loss', description: 'Reduce body fat percentage' },
  { label: 'Endurance', value: 'endurance', description: 'Improve cardiovascular health' },
  { label: 'General Fitness', value: 'general_fitness', description: 'Balanced training approach' },
];

const EQUIPMENT_OPTIONS: Array<{ label: string; value: EquipmentType }> = [
  { label: 'Bodyweight', value: 'bodyweight' },
  { label: 'Dumbbells', value: 'dumbbell' },
  { label: 'Barbell', value: 'barbell' },
  { label: 'Machines', value: 'machine' },
  { label: 'Kettlebells', value: 'kettlebell' },
  { label: 'Resistance Bands', value: 'resistance_band' },
];

const FITNESS_LEVELS: Array<{ label: string; value: FitnessLevel; subtitle: string }> = [
  { label: 'Beginner', value: 'beginner', subtitle: 'New to structured training' },
  { label: 'Intermediate', value: 'intermediate', subtitle: '6+ months of experience' },
  { label: 'Advanced', value: 'advanced', subtitle: '2+ years of experience' },
];

const FREQUENCY_OPTIONS = [2, 3, 4, 5, 6];
const DURATION_OPTIONS = [30, 45, 60, 75, 90];

const Chip: React.FC<{
  label: string;
  selected: boolean;
  onPress: () => void;
  subtitle?: string;
}> = ({ label, selected, onPress, subtitle }) => {
  return (
    <TouchableOpacity
      onPress={onPress}
      activeOpacity={0.8}
      style={[styles.chip, selected && styles.chipSelected]}
    >
      <Text style={[styles.chipLabel, selected && styles.chipLabelSelected]}>{label}</Text>
      {subtitle ? (
        <Text style={[styles.chipSubtitle, selected && styles.chipSubtitleSelected]}>{subtitle}</Text>
      ) : null}
    </TouchableOpacity>
  );
};

export default function PlanGeneratorScreen() {
  const router = useRouter();
  const { user } = useAuth();
  const createPlanMutation = useCreatePlan();

  const userId = useMemo(() => {
    if (typeof user?.id === 'number') {
      return user.id;
    }
    if (typeof user?.id === 'string') {
      const parsed = parseInt(user.id, 10);
      return Number.isFinite(parsed) ? parsed : 0;
    }
    return 0;
  }, [user?.id]);

  const [selectedGoals, setSelectedGoals] = useState<FitnessGoal[]>(['general_fitness']);
  const [selectedEquipment, setSelectedEquipment] = useState<EquipmentType[]>(['bodyweight']);
  const [fitnessLevel, setFitnessLevel] = useState<FitnessLevel>('beginner');
  const [weeklyFrequency, setWeeklyFrequency] = useState<number>(3);
  const [timePerWorkout, setTimePerWorkout] = useState<number>(45);

  const toggleGoal = useCallback((value: FitnessGoal) => {
    setSelectedGoals((prev) =>
      prev.includes(value) ? prev.filter((goal) => goal !== value) : [...prev, value]
    );
  }, []);

  const toggleEquipment = useCallback((value: EquipmentType) => {
    setSelectedEquipment((prev) =>
      prev.includes(value) ? prev.filter((item) => item !== value) : [...prev, value]
    );
  }, []);

  const handleSubmit = async () => {
    if (!user) {
      Alert.alert('Not Ready', 'Please sign in again before generating a plan.');
      return;
    }

    const resolvedUserId = Number.isFinite(userId) && userId > 0 ? userId : 0;

    if (selectedGoals.length === 0) {
      Alert.alert('Missing Goal', 'Choose at least one training goal to personalize your plan.');
      return;
    }

    const payload = {
  user_id: resolvedUserId,
      metadata: {
        user_goals: selectedGoals,
        available_equipment: selectedEquipment,
        fitness_level: fitnessLevel,
        weekly_frequency: weeklyFrequency,
        time_per_workout: timePerWorkout,
      },
    };

    try {
      await createPlanMutation.mutateAsync(payload);
      Alert.alert('Plan Ready', 'Your new plan has been generated.', [
        {
          text: 'View Plan',
          onPress: () => router.back(),
        },
      ]);
    } catch (error: unknown) {
      let message = 'Failed to generate plan. Please try again later.';

      if (error instanceof APIError) {
        const errorCode = typeof error.data?.code === 'string' ? error.data.code : undefined;
        if (errorCode === 'PLAN_LIMIT_REACHED' || /maximum number of active plans/i.test(error.message)) {
          message = 'You already have three active plans. Delete an existing plan before generating a new one.';
        } else if (error.message) {
          message = error.message;
        }
      } else if (error instanceof Error && error.message) {
        message = error.message;
      }

      Alert.alert('Generation Failed', message);
    }
  };

  return (
    <View style={styles.container}>
      {/* <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={24} color={COLORS.text.primary} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Plan Generator</Text>
        <View style={styles.headerSpacer} />
      </View> */}

      <ScrollView contentContainerStyle={styles.content} showsVerticalScrollIndicator={false}>
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Primary Goals</Text>
          <Text style={styles.sectionDescription}>
            Pick the results you are aiming for. The generator will tune intensity and exercise selection accordingly.
          </Text>
          <View style={styles.chipGroup}>
            {GOAL_OPTIONS.map((option) => (
              <Chip
                key={option.value}
                label={option.label}
                subtitle={option.description}
                selected={selectedGoals.includes(option.value)}
                onPress={() => toggleGoal(option.value)}
              />
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Available Equipment</Text>
          <Text style={styles.sectionDescription}>
            Select everything you can reliably access. We will filter exercises to match.
          </Text>
          <View style={styles.chipRow}>
            {EQUIPMENT_OPTIONS.map((option) => (
              <TouchableOpacity
                key={option.value}
                style={[styles.smallChip, selectedEquipment.includes(option.value) && styles.smallChipSelected]}
                onPress={() => toggleEquipment(option.value)}
              >
                <Text style={[styles.smallChipLabel, selectedEquipment.includes(option.value) && styles.smallChipLabelSelected]}>
                  {option.label}
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Experience Level</Text>
          <Text style={styles.sectionDescription}>
            We balance volume and progression based on your current training background.
          </Text>
          <View style={styles.chipGroup}>
            {FITNESS_LEVELS.map((option) => (
              <Chip
                key={option.value}
                label={option.label}
                subtitle={option.subtitle}
                selected={fitnessLevel === option.value}
                onPress={() => setFitnessLevel(option.value)}
              />
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Weekly Frequency</Text>
          <Text style={styles.sectionDescription}>
            How many sessions can you consistently commit to each week?
          </Text>
          <View style={styles.smallChipRow}>
            {FREQUENCY_OPTIONS.map((option) => (
              <TouchableOpacity
                key={option}
                style={[styles.frequencyChip, weeklyFrequency === option && styles.frequencyChipSelected]}
                onPress={() => setWeeklyFrequency(option)}
              >
                <Text
                  style={[styles.frequencyChipLabel, weeklyFrequency === option && styles.frequencyChipLabelSelected]}
                >
                  {option}x
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Session Duration</Text>
          <Text style={styles.sectionDescription}>
            Choose the average time you want to spend per workout.
          </Text>
          <View style={styles.smallChipRow}>
            {DURATION_OPTIONS.map((option) => (
              <TouchableOpacity
                key={option}
                style={[styles.durationChip, timePerWorkout === option && styles.durationChipSelected]}
                onPress={() => setTimePerWorkout(option)}
              >
                <Text
                  style={[styles.durationChipLabel, timePerWorkout === option && styles.durationChipLabelSelected]}
                >
                  {option} min
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>
      </ScrollView>

      <View style={styles.footer}>
        <Button
          title="Cancel"
          onPress={() => router.back()}
          variant="outline"
          disabled={createPlanMutation.isPending}
        />
        <TouchableOpacity
          style={[styles.generateButton, createPlanMutation.isPending && styles.generateButtonDisabled]}
          onPress={handleSubmit}
          disabled={createPlanMutation.isPending}
          activeOpacity={0.85}
        >
          {createPlanMutation.isPending ? (
            <ActivityIndicator color={COLORS.text.primary} />
          ) : (
            <>
              <Ionicons name="flash" size={20} color={COLORS.text.primary} />
              <Text style={styles.generateButtonText}>Generate Plan</Text>
            </>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingTop: SPACING['4xl'],
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING.lg,
  },
  backButton: {
    width: 44,
    height: 44,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
    justifyContent: 'center',
  ...SHADOWS.sm,
  },
  headerTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  headerSpacer: {
    width: 44,
  },
  content: {
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING['5xl'],
  },
  section: {
    marginBottom: SPACING['4xl'],
  },
  sectionTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.sm,
  },
  sectionDescription: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.md,
    lineHeight: 20,
  },
  chipGroup: {
    gap: SPACING.md,
  },
  chip: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.base,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  ...SHADOWS.sm,
  },
  chipSelected: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  chipLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 2,
  },
  chipLabelSelected: {
    color: COLORS.text.primary,
  },
  chipSubtitle: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  chipSubtitleSelected: {
    color: COLORS.text.primary,
  },
  chipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  smallChipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  smallChip: {
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  smallChipSelected: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  smallChipLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.auth.primary,
    fontWeight: FONT_WEIGHTS.medium,
  },
  smallChipLabelSelected: {
    color: COLORS.text.primary,
  },
  frequencyChip: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.base,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  frequencyChipSelected: {
    backgroundColor: COLORS.background.accent,
    borderColor: COLORS.primary,
  },
  frequencyChipLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
  },
  frequencyChipLabelSelected: {
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.bold,
  },
  durationChip: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.base,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  durationChipSelected: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  durationChipLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.primary,
  },
  durationChipLabelSelected: {
    color: COLORS.text.primary,
  },
  footer: {
    position: 'absolute',
    left: SPACING.lg,
    right: SPACING.lg,
    bottom: SPACING['3xl'],
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: SPACING.md,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.full,
    padding: SPACING.md,
    ...SHADOWS.lg,
  },
  generateButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.sm,
    backgroundColor: COLORS.primary,
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: SPACING.md,
  },
  generateButtonDisabled: {
    opacity: 0.6,
  },
  generateButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
    textTransform: 'uppercase',
  },
});
