import React, { useState } from 'react';
import {
  Modal,
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  TextInput,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { MealType, UserRecipe, SystemRecipe, CreateFoodLogRequest } from '@/types/food-tracker';
import { Button } from '../forms/button';

interface FoodLogModalProps {
  visible: boolean;
  date: string;
  onClose: () => void;
  onSubmit: (data: CreateFoodLogRequest) => void;
  isSubmitting?: boolean;
  selectedRecipe?: UserRecipe | SystemRecipe | null;
}

const mealTypes: Array<{ label: string; value: MealType }> = [
  { label: 'Breakfast', value: 'breakfast' },
  { label: 'Lunch', value: 'lunch' },
  { label: 'Dinner', value: 'dinner' },
  { label: 'Snack', value: 'snack' },
];

export const FoodLogModal: React.FC<FoodLogModalProps> = ({
  visible,
  date,
  onClose,
  onSubmit,
  isSubmitting = false,
  selectedRecipe,
}) => {
  const [mealType, setMealType] = useState<MealType>('breakfast');
  const [servings, setServings] = useState('1');
  const [calories, setCalories] = useState('');
  const [protein, setProtein] = useState('');
  const [carbs, setCarbs] = useState('');
  const [fat, setFat] = useState('');
  const [fiber, setFiber] = useState('');

  React.useEffect(() => {
    if (selectedRecipe && visible) {
      const servingsNum = parseFloat(servings) || 1;
      setCalories(Math.round(selectedRecipe.calories * servingsNum).toString());
      setProtein(Math.round(selectedRecipe.protein * servingsNum).toString());
      setCarbs(Math.round(selectedRecipe.carbs * servingsNum).toString());
      setFat(Math.round(selectedRecipe.fat * servingsNum).toString());
      setFiber(Math.round(selectedRecipe.fiber * servingsNum).toString());
    }
  }, [selectedRecipe, servings, visible]);

  const handleSubmit = () => {
    const servingsNum = parseFloat(servings);
    const caloriesNum = parseInt(calories);
    const proteinNum = parseInt(protein);
    const carbsNum = parseInt(carbs);
    const fatNum = parseInt(fat);
    const fiberNum = parseInt(fiber);

    if (isNaN(servingsNum) || servingsNum <= 0) {
      Alert.alert('Error', 'Please enter a valid servings amount');
      return;
    }

    if (isNaN(caloriesNum) || isNaN(proteinNum) || isNaN(carbsNum) || isNaN(fatNum)) {
      Alert.alert('Error', 'Please enter valid nutrition values');
      return;
    }

    const payload: CreateFoodLogRequest = {
      log_date: date,
      meal_type: mealType,
      calories: caloriesNum,
      protein: proteinNum,
      carbs: carbsNum,
      fat: fatNum,
      fiber: fiberNum || 0,
      servings: servingsNum,
    };

    if (selectedRecipe) {
      if ('user_id' in selectedRecipe) {
        payload.user_recipe_id = selectedRecipe.id;
      } else {
        payload.system_recipe_id = selectedRecipe.id;
      }
    }

    onSubmit(payload);
  };

  const handleClose = () => {
    setMealType('breakfast');
    setServings('1');
    setCalories('');
    setProtein('');
    setCarbs('');
    setFat('');
    setFiber('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent
      onRequestClose={handleClose}
    >
      <View style={styles.overlay}>
        <View style={styles.modal}>
          <View style={styles.header}>
            <Text style={styles.title}>
              {selectedRecipe ? `Log: ${selectedRecipe.name}` : 'Log Food'}
            </Text>
            <TouchableOpacity onPress={handleClose} hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}>
              <Ionicons name="close" size={24} color={COLORS.text.inverse} style={styles.cancelBTNN} />
            </TouchableOpacity>
          </View>

          <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
            <Text style={styles.sectionTitle}>Meal Type</Text>
            <View style={styles.mealTypeContainer}>
              {mealTypes.map((type) => (
                <TouchableOpacity
                  key={type.value}
                  style={[
                    styles.mealTypeButton,
                    mealType === type.value && styles.mealTypeButtonActive,
                  ]}
                  onPress={() => setMealType(type.value)}
                >
                  <Text
                    style={[
                      styles.mealTypeText,
                      mealType === type.value && styles.mealTypeTextActive,
                    ]}
                  >
                    {type.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text style={styles.sectionTitle}>Servings</Text>
            <TextInput
              style={styles.input}
              value={servings}
              onChangeText={setServings}
              keyboardType="decimal-pad"
              placeholder="1"
              placeholderTextColor={COLORS.text.placeholder}
            />

            <Text style={styles.sectionTitle}>Nutrition (per serving)</Text>
            <View style={styles.nutritionGrid}>
              <View style={styles.nutritionInput}>
                <Text style={styles.inputLabel}>Calories</Text>
                <TextInput
                  style={styles.input}
                  value={calories}
                  onChangeText={setCalories}
                  keyboardType="number-pad"
                  placeholder="0"
                  placeholderTextColor={COLORS.text.placeholder}
                  editable={!selectedRecipe}
                />
              </View>

              <View style={styles.nutritionInput}>
                <Text style={styles.inputLabel}>Protein (g)</Text>
                <TextInput
                  style={styles.input}
                  value={protein}
                  onChangeText={setProtein}
                  keyboardType="number-pad"
                  placeholder="0"
                  placeholderTextColor={COLORS.text.placeholder}
                  editable={!selectedRecipe}
                />
              </View>

              <View style={styles.nutritionInput}>
                <Text style={styles.inputLabel}>Carbs (g)</Text>
                <TextInput
                  style={styles.input}
                  value={carbs}
                  onChangeText={setCarbs}
                  keyboardType="number-pad"
                  placeholder="0"
                  placeholderTextColor={COLORS.text.placeholder}
                  editable={!selectedRecipe}
                />
              </View>

              <View style={styles.nutritionInput}>
                <Text style={styles.inputLabel}>Fat (g)</Text>
                <TextInput
                  style={styles.input}
                  value={fat}
                  onChangeText={setFat}
                  keyboardType="number-pad"
                  placeholder="0"
                  placeholderTextColor={COLORS.text.placeholder}
                  editable={!selectedRecipe}
                />
              </View>

              <View style={styles.nutritionInput}>
                <Text style={styles.inputLabel}>Fiber (g)</Text>
                <TextInput
                  style={styles.input}
                  value={fiber}
                  onChangeText={setFiber}
                  keyboardType="number-pad"
                  placeholder="0"
                  placeholderTextColor={COLORS.text.placeholder}
                  editable={!selectedRecipe}
                />
              </View>
            </View>
          </ScrollView>

          <View style={styles.footer}>
            <Button
              onPress={handleClose}
              title={isSubmitting ? '' : 'Cancel'}
              variant="outline"
            />
            <Button
              style={[styles.button, styles.submitButton, isSubmitting && styles.buttonDisabled]}
              onPress={handleSubmit}
              disabled={isSubmitting}
              title={isSubmitting ? 'Logging...' : 'Log Food'}
            />
          </View>
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
  modal: {
    backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['2xl'],
    borderTopRightRadius: BORDER_RADIUS['2xl'],
    maxHeight: '90%',
    width: '98%',
    alignSelf: 'center',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  cancelBTNN: {
     backgroundColor: COLORS.background.accent,
        padding: SPACING.md,
        borderRadius: BORDER_RADIUS.full,
        justifyContent: 'center',
        alignItems: 'center',
        minWidth: 40,
        minHeight: 40,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    flex: 1,
  },
  content: {
    padding: SPACING.xl,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
    marginTop: SPACING.lg,
  },
  mealTypeContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  mealTypeButton: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
  },
  mealTypeButtonActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  mealTypeText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
  },
  mealTypeTextActive: {
    color: COLORS.text.primary,
  },
  input: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  nutritionGrid: {
    gap: SPACING.md,
  },
  nutritionInput: {
    flex: 1,
  },
  inputLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
  },
  footer: {
    flexDirection: 'row',
    gap: SPACING.md,
    padding: SPACING.xl,
  },
  button: {
    flex: 1,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
  },
  cancelButton: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  submitButton: {
    backgroundColor: COLORS.primary,
  },
  buttonDisabled: {
    opacity: 0.5,
  },
  cancelButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
  },
  submitButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
});
