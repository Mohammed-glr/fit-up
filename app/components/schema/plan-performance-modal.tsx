import React, { useEffect, useMemo, useState } from 'react';
import {
  Modal,
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  TextInput,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';

import { Button } from '@/components/forms';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import type { PlanPerformancePayload } from '@/types/schema';

interface PlanPerformanceModalProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (payload: PlanPerformancePayload) => void;
  isSubmitting?: boolean;
}

interface FieldConfig {
  key: keyof PlanPerformanceModalState;
  label: string;
  helper: string;
  suffix?: string;
  keyboardType?: 'numeric' | 'default';
  max?: number;
}

type PlanPerformanceModalState = {
  completionRate: string;
  progressRate: string;
  averageRpe: string;
  satisfaction: string;
  injuryRate: string;
};

const DEFAULT_STATE: PlanPerformanceModalState = {
  completionRate: '85',
  progressRate: '70',
  averageRpe: '7',
  satisfaction: '8',
  injuryRate: '0',
};

const FIELDS: FieldConfig[] = [
  {
    key: 'completionRate',
    label: 'Completion Rate',
    helper: 'Percent of scheduled workouts you completed this week.',
    suffix: '%',
    keyboardType: 'numeric',
    max: 100,
  },
  {
    key: 'progressRate',
    label: 'Progress Rate',
    helper: 'How much progress you felt (volume, weight, energy).',
    suffix: '%',
    keyboardType: 'numeric',
    max: 100,
  },
  {
    key: 'averageRpe',
    label: 'Average RPE',
    helper: 'How hard the sessions felt on average (1 easy – 10 max).',
    keyboardType: 'numeric',
    max: 10,
  },
  {
    key: 'satisfaction',
    label: 'Satisfaction',
    helper: 'Overall satisfaction with the plan this week (1 low – 10 high).',
    keyboardType: 'numeric',
    max: 10,
  },
  {
    key: 'injuryRate',
    label: 'Injury / Pain',
    helper: 'Time spent limited by pain or injury concerns.',
    suffix: '%',
    keyboardType: 'numeric',
    max: 100,
  },
];

const clampNumber = (value: number, min: number, max: number) => {
  if (Number.isNaN(value)) {
    return min;
  }
  return Math.min(Math.max(value, min), max);
};

export const PlanPerformanceModal: React.FC<PlanPerformanceModalProps> = ({
  visible,
  onClose,
  onSubmit,
  isSubmitting = false,
}) => {
  const [formState, setFormState] = useState<PlanPerformanceModalState>(DEFAULT_STATE);

  useEffect(() => {
    if (!visible) {
      setFormState(DEFAULT_STATE);
    }
  }, [visible]);

  const handleChange = (key: keyof PlanPerformanceModalState, value: string) => {
    setFormState((prev) => ({
      ...prev,
      [key]: value.replace(/[^0-9.]/g, ''),
    }));
  };

  const performancePayload = useMemo<PlanPerformancePayload>(() => {
    const completion = clampNumber(parseFloat(formState.completionRate) / 100, 0, 1);
    const progress = clampNumber(parseFloat(formState.progressRate) / 100, 0, 1);
    const averageRpe = clampNumber(parseFloat(formState.averageRpe), 1, 10);
    const satisfaction = clampNumber(parseFloat(formState.satisfaction), 1, 10);
    const injuryRate = clampNumber(parseFloat(formState.injuryRate) / 100, 0, 1);

    return {
      completion_rate: completion,
      progress_rate: progress,
      average_rpe: averageRpe,
      user_satisfaction: satisfaction,
      injury_rate: injuryRate,
    };
  }, [formState]);

  const handleSubmit = () => {
    onSubmit(performancePayload);
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent
      onRequestClose={onClose}
    >
      <View style={styles.backdrop}>
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : undefined}
          style={styles.modalWrapper}
        >
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Log Weekly Performance</Text>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <Ionicons name="close" size={22} color={COLORS.text.auth.primary} />
              </TouchableOpacity>
            </View>

            <Text style={styles.modalSubtitle}>
              Capture how the plan felt this week. We will use this to adapt your next cycle.
            </Text>

            <ScrollView style={styles.form} contentContainerStyle={styles.formContent}>
              {FIELDS.map((field) => {
                const value = formState[field.key];
                return (
                  <View key={field.key} style={styles.fieldRow}>
                    <View style={styles.fieldHeader}>
                      <Text style={styles.fieldLabel}>{field.label}</Text>
                      {field.suffix ? <Text style={styles.fieldSuffix}>{field.suffix}</Text> : null}
                    </View>
                    <Text style={styles.fieldHelper}>{field.helper}</Text>
                    <TextInput
                      value={value}
                      onChangeText={(text) => handleChange(field.key, text)}
                      keyboardType={field.keyboardType}
                      inputMode="decimal"
                      style={styles.input}
                      placeholder={field.suffix ? `0${field.suffix}` : '0'}
                      placeholderTextColor={COLORS.text.placeholder}
                    />
                  </View>
                );
              })}
            </ScrollView>

            <View style={styles.footer}>
              <Button
                title="Cancel"
                onPress={onClose}
                variant="outline"
                disabled={isSubmitting}
              />
              <TouchableOpacity
                style={[styles.submitButton, isSubmitting && styles.submitButtonDisabled]}
                onPress={handleSubmit}
                disabled={isSubmitting}
                activeOpacity={0.85}
              >
                {isSubmitting ? (
                  <Ionicons name="time" size={20} color={COLORS.text.primary} />
                ) : (
                  <>
                    <Ionicons name="analytics" size={20} color={COLORS.text.primary} />
                    <Text style={styles.submitButtonText}>Save Metrics</Text>
                  </>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </KeyboardAvoidingView>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  backdrop: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.4)',
    justifyContent: 'flex-end',
  },
  modalWrapper: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: COLORS.background.card,
    borderTopLeftRadius: BORDER_RADIUS['3xl'],
    borderTopRightRadius: BORDER_RADIUS['3xl'],
    padding: SPACING.lg,
    paddingBottom: SPACING['4xl'],
    minHeight: '75%',
    ...SHADOWS.lg,
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
  },
  modalTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  closeButton: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  modalSubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.lg,
    lineHeight: 20,
  },
  form: {
    flex: 1,
  },
  formContent: {
    paddingBottom: SPACING['3xl'],
  },
  fieldRow: {
    marginBottom: SPACING['2xl'],
  },
  fieldHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  fieldLabel: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  fieldSuffix: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
  },
  fieldHelper: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.sm,
    lineHeight: 18,
  },
  input: {
    backgroundColor: COLORS.background.primary,
    borderRadius: BORDER_RADIUS.lg,
    paddingVertical: SPACING.base,
    paddingHorizontal: SPACING.lg,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
  },
  footer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: SPACING.md,
    marginTop: SPACING.lg,
  },
  submitButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.sm,
    backgroundColor: COLORS.primary,
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: SPACING.md,
  },
  submitButtonDisabled: {
    opacity: 0.6,
  },
  submitButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
    textTransform: 'uppercase',
  },
});
