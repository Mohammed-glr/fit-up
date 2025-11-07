import React from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  StyleSheet,
  ActivityIndicator,
  Alert,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import type {
  CreateRecipeRequest,
  RecipeCategory,
  RecipeDifficulty,
  UserRecipeDetail,
} from '@/types/food-tracker';

interface RecipeFormModalProps {
  visible: boolean;
  recipe?: UserRecipeDetail | null;
  onClose: () => void;
  onSubmit: (data: CreateRecipeRequest) => void;
  isSubmitting?: boolean;
}

const categories: RecipeCategory[] = ['breakfast', 'lunch', 'dinner', 'snack', 'dessert'];
const difficulties: RecipeDifficulty[] = ['easy', 'medium', 'hard'];

export function RecipeFormModal({
  visible,
  recipe,
  onClose,
  onSubmit,
  isSubmitting = false,
}: RecipeFormModalProps) {
  const isEditing = Boolean(recipe);

  const [formData, setFormData] = React.useState<CreateRecipeRequest>({
    name: '',
    description: '',
    category: 'breakfast',
    difficulty: 'easy',
    calories: 0,
    protein: 0,
    carbs: 0,
    fat: 0,
    fiber: 0,
    prep_time: 0,
    cook_time: 0,
    servings: 1,
    image_url: '',
    ingredients: [],
    instructions: [],
    tags: [],
  });

  React.useEffect(() => {
    if (recipe) {
      setFormData({
        name: recipe.name,
        description: recipe.description,
        category: recipe.category,
        difficulty: recipe.difficulty,
        calories: recipe.calories,
        protein: recipe.protein,
        carbs: recipe.carbs,
        fat: recipe.fat,
        fiber: recipe.fiber,
        prep_time: recipe.prep_time,
        cook_time: recipe.cook_time,
        servings: recipe.servings,
        image_url: recipe.image_url || '',
        ingredients: recipe.ingredients.map((ing) => ({
          ingredient_id: ing.id,
          item: ing.item,
          amount: ing.amount,
          unit: ing.unit,
          order_index: ing.order_index,
        })),
        instructions: recipe.instructions.map((inst) => ({
          instruction_id: inst.id,
          step_number: inst.step_number,
          instruction: inst.instruction,
        })),
        tags: recipe.tags.map((tag) => tag.tag_name),
      });
    } else {
      setFormData({
        name: '',
        description: '',
        category: 'breakfast',
        difficulty: 'easy',
        calories: 0,
        protein: 0,
        carbs: 0,
        fat: 0,
        fiber: 0,
        prep_time: 0,
        cook_time: 0,
        servings: 1,
        image_url: '',
        ingredients: [],
        instructions: [],
        tags: [],
      });
    }
  }, [recipe]);

  const handleAddIngredient = () => {
    setFormData((prev) => ({
      ...prev,
      ingredients: [
        ...prev.ingredients,
        {
          item: '',
          amount: 0,
          unit: '',
          order_index: prev.ingredients.length + 1,
        },
      ],
    }));
  };

  const handleRemoveIngredient = (index: number) => {
    setFormData((prev) => ({
      ...prev,
      ingredients: prev.ingredients.filter((_, i) => i !== index),
    }));
  };

  const handleUpdateIngredient = (
    index: number,
    field: 'item' | 'amount' | 'unit',
    value: string | number
  ) => {
    setFormData((prev) => ({
      ...prev,
      ingredients: prev.ingredients.map((ing, i) =>
        i === index ? { ...ing, [field]: value } : ing
      ),
    }));
  };

  const handleAddInstruction = () => {
    setFormData((prev) => ({
      ...prev,
      instructions: [
        ...prev.instructions,
        {
          step_number: prev.instructions.length + 1,
          instruction: '',
        },
      ],
    }));
  };

  const handleRemoveInstruction = (index: number) => {
    setFormData((prev) => ({
      ...prev,
      instructions: prev.instructions.filter((_, i) => i !== index).map((inst, i) => ({
        ...inst,
        step_number: i + 1,
      })),
    }));
  };

  const handleUpdateInstruction = (index: number, value: string) => {
    setFormData((prev) => ({
      ...prev,
      instructions: prev.instructions.map((inst, i) =>
        i === index ? { ...inst, instruction: value } : inst
      ),
    }));
  };

  const handleSubmit = () => {
    if (!formData.name.trim()) {
      Alert.alert('Error', 'Recipe name is required');
      return;
    }
    if (formData.servings <= 0) {
      Alert.alert('Error', 'Servings must be greater than 0');
      return;
    }
    if (formData.ingredients.length === 0) {
      Alert.alert('Error', 'At least one ingredient is required');
      return;
    }
    if (formData.instructions.length === 0) {
      Alert.alert('Error', 'At least one instruction is required');
      return;
    }

    onSubmit(formData);
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <View style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={onClose} disabled={isSubmitting}>
            <Ionicons name="close" size={28} color={COLORS.text.inverse} />
          </TouchableOpacity>
          <Text style={styles.title}>{isEditing ? 'Edit Recipe' : 'New Recipe'}</Text>
          <TouchableOpacity onPress={handleSubmit} disabled={isSubmitting}>
            {isSubmitting ? (
              <ActivityIndicator color={COLORS.primary} />
            ) : (
              <Text style={styles.saveButton}>Save</Text>
            )}
          </TouchableOpacity>
        </View>

        <ScrollView style={styles.scrollView} contentContainerStyle={styles.scrollContent}>
          {/* Basic Info */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Basic Information</Text>

            <TextInput
              style={styles.input}
              placeholder="Recipe Name"
              placeholderTextColor={COLORS.text.placeholder}
              value={formData.name}
              onChangeText={(text) => setFormData((prev) => ({ ...prev, name: text }))}
            />

            <TextInput
              style={[styles.input, styles.textArea]}
              placeholder="Description"
              placeholderTextColor={COLORS.text.placeholder}
              value={formData.description}
              onChangeText={(text) => setFormData((prev) => ({ ...prev, description: text }))}
              multiline
              numberOfLines={3}
            />

            <TextInput
              style={styles.input}
              placeholder="Image URL (optional)"
              placeholderTextColor={COLORS.text.placeholder}
              value={formData.image_url}
              onChangeText={(text) => setFormData((prev) => ({ ...prev, image_url: text }))}
            />

            <View style={styles.row}>
              <View style={styles.halfWidth}>
                <Text style={styles.label}>Category</Text>
                <View style={styles.pickerWrapper}>
                  {categories.map((cat) => (
                    <TouchableOpacity
                      key={cat}
                      style={[
                        styles.chip,
                        formData.category === cat && styles.chipActive,
                      ]}
                      onPress={() => setFormData((prev) => ({ ...prev, category: cat }))}
                    >
                      <Text
                        style={[
                          styles.chipText,
                          formData.category === cat && styles.chipTextActive,
                        ]}
                      >
                        {cat}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
              </View>

              <View style={styles.halfWidth}>
                <Text style={styles.label}>Difficulty</Text>
                <View style={styles.pickerWrapper}>
                  {difficulties.map((diff) => (
                    <TouchableOpacity
                      key={diff}
                      style={[
                        styles.chip,
                        formData.difficulty === diff && styles.chipActive,
                      ]}
                      onPress={() => setFormData((prev) => ({ ...prev, difficulty: diff }))}
                    >
                      <Text
                        style={[
                          styles.chipText,
                          formData.difficulty === diff && styles.chipTextActive,
                        ]}
                      >
                        {diff}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
              </View>
            </View>
          </View>

          {/* Nutrition */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Nutrition (per serving)</Text>
            <View style={styles.row}>
              <View style={styles.halfWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Calories"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.calories.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, calories: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
              <View style={styles.halfWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Protein (g)"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.protein.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, protein: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
            </View>

            <View style={styles.row}>
              <View style={styles.halfWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Carbs (g)"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.carbs.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, carbs: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
              <View style={styles.halfWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Fat (g)"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.fat.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, fat: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
            </View>

            <TextInput
              style={styles.input}
              placeholder="Fiber (g)"
              placeholderTextColor={COLORS.text.placeholder}
              value={formData.fiber.toString()}
              onChangeText={(text) =>
                setFormData((prev) => ({ ...prev, fiber: parseInt(text) || 0 }))
              }
              keyboardType="numeric"
            />
          </View>

          {/* Time & Servings */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Time & Servings</Text>
            <View style={styles.row}>
              <View style={styles.thirdWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Prep (min)"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.prep_time.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, prep_time: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
              <View style={styles.thirdWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Cook (min)"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.cook_time.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, cook_time: parseInt(text) || 0 }))
                  }
                  keyboardType="numeric"
                />
              </View>
              <View style={styles.thirdWidth}>
                <TextInput
                  style={styles.input}
                  placeholder="Servings"
                  placeholderTextColor={COLORS.text.placeholder}
                  value={formData.servings.toString()}
                  onChangeText={(text) =>
                    setFormData((prev) => ({ ...prev, servings: parseInt(text) || 1 }))
                  }
                  keyboardType="numeric"
                />
              </View>
            </View>
          </View>

          {/* Ingredients */}
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>Ingredients</Text>
              <TouchableOpacity onPress={handleAddIngredient}>
                <Ionicons name="add-circle" size={28} color={COLORS.primary} />
              </TouchableOpacity>
            </View>

            {formData.ingredients.map((ingredient, index) => (
              <View key={index} style={styles.ingredientRow}>
                <View style={styles.ingredientInputs}>
                  <TextInput
                    style={[styles.input, { flex: 2 }]}
                    placeholder="Item"
                    placeholderTextColor={COLORS.text.placeholder}
                    value={ingredient.item}
                    onChangeText={(text) => handleUpdateIngredient(index, 'item', text)}
                  />
                  <TextInput
                    style={[styles.input, { flex: 1 }]}
                    placeholder="Amount"
                    placeholderTextColor={COLORS.text.placeholder}
                    value={ingredient.amount.toString()}
                    onChangeText={(text) =>
                      handleUpdateIngredient(index, 'amount', parseFloat(text) || 0)
                    }
                    keyboardType="decimal-pad"
                  />
                  <TextInput
                    style={[styles.input, { flex: 1 }]}
                    placeholder="Unit"
                    placeholderTextColor={COLORS.text.placeholder}
                    value={ingredient.unit}
                    onChangeText={(text) => handleUpdateIngredient(index, 'unit', text)}
                  />
                </View>
                <TouchableOpacity onPress={() => handleRemoveIngredient(index)}>
                  <Ionicons name="trash-outline" size={24} color={COLORS.error} />
                </TouchableOpacity>
              </View>
            ))}
          </View>

          {/* Instructions */}
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>Instructions</Text>
              <TouchableOpacity onPress={handleAddInstruction}>
                <Ionicons name="add-circle" size={28} color={COLORS.primary} />
              </TouchableOpacity>
            </View>

            {formData.instructions.map((instruction, index) => (
              <View key={index} style={styles.instructionRow}>
                <Text style={styles.stepNumber}>{index + 1}</Text>
                <TextInput
                  style={[styles.input, styles.instructionInput]}
                  placeholder={`Step ${index + 1}`}
                  placeholderTextColor={COLORS.text.placeholder}
                  value={instruction.instruction}
                  onChangeText={(text) => handleUpdateInstruction(index, text)}
                  multiline
                />
                <TouchableOpacity onPress={() => handleRemoveInstruction(index)}>
                  <Ionicons name="trash-outline" size={24} color={COLORS.error} />
                </TouchableOpacity>
              </View>
            ))}
          </View>
        </ScrollView>
      </View>
    </Modal>
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
    padding: SPACING.lg,
    paddingTop: Platform.OS === 'ios' ? SPACING['3xl'] : SPACING.lg,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.dark,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  saveButton: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.lg,
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  label: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.sm,
  },
  input: {
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    borderRadius: BORDER_RADIUS.md,
    padding: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  textArea: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  row: {
    flexDirection: 'row',
    gap: SPACING.md,
  },
  halfWidth: {
    flex: 1,
  },
  thirdWidth: {
    flex: 1,
  },
  pickerWrapper: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  chip: {
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  chipActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  chipText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    textTransform: 'capitalize',
  },
  chipTextActive: {
    color: COLORS.text.primary,
  },
  ingredientRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  ingredientInputs: {
    flex: 1,
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  instructionRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: SPACING.sm,
    marginBottom: SPACING.md,
  },
  stepNumber: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
    marginTop: SPACING.md,
  },
  instructionInput: {
    flex: 1,
    minHeight: 60,
  },
});
