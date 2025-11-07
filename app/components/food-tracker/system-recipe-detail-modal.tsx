import React from 'react';
import {
  Modal,
  ScrollView,
  View,
  Text,
  TouchableOpacity,
  Image,
  StyleSheet,
  ActivityIndicator,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import type { SystemRecipeDetail } from '@/types/food-tracker';

interface SystemRecipeDetailModalProps {
  visible: boolean;
  recipe?: SystemRecipeDetail | null;
  isLoading?: boolean;
  onClose: () => void;
}

export function SystemRecipeDetailModal({
  visible,
  recipe,
  isLoading,
  onClose,
}: SystemRecipeDetailModalProps) {
  if (!visible) return null;

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <View style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={onClose}>
            <Ionicons name="close" size={28} color={COLORS.text.inverse} />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Recipe Details</Text>
          <View style={{ width: 28 }} />
        </View>

        {isLoading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color={COLORS.primary} />
          </View>
        ) : recipe ? (
          <ScrollView style={styles.scrollView} contentContainerStyle={styles.scrollContent}>
            {recipe.image_url && (
              <Image source={{ uri: recipe.image_url }} style={styles.image} resizeMode="cover" />
            )}

            <View style={styles.content}>
              <Text style={styles.name}>{recipe.name}</Text>
              <Text style={styles.description}>{recipe.description}</Text>

              <View style={styles.metaRow}>
                <View style={styles.metaBadge}>
                  <Text style={styles.metaText}>{recipe.category}</Text>
                </View>
                <View style={styles.metaBadge}>
                  <Text style={styles.metaText}>{recipe.difficulty}</Text>
                </View>
                <View style={styles.metaBadge}>
                  <Ionicons name="time-outline" size={14} color={COLORS.text.tertiary} />
                  <Text style={styles.metaText}>
                    {recipe.prep_time + recipe.cook_time} min
                  </Text>
                </View>
              </View>

              {recipe.tags && recipe.tags.length > 0 && (
                <View style={styles.tagsContainer}>
                  {recipe.tags.map((tag, index) => (
                    <View key={index} style={styles.tag}>
                      <Text style={styles.tagText}>{tag.tag_name}</Text>
                    </View>
                  ))}
                </View>
              )}

              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Nutrition (per serving)</Text>
                <View style={styles.nutritionGrid}>
                  <View style={styles.nutritionItem}>
                    <Text style={styles.nutritionValue}>{recipe.calories}</Text>
                    <Text style={styles.nutritionLabel}>Calories</Text>
                  </View>
                  <View style={styles.nutritionItem}>
                    <Text style={styles.nutritionValue}>{recipe.protein}g</Text>
                    <Text style={styles.nutritionLabel}>Protein</Text>
                  </View>
                  <View style={styles.nutritionItem}>
                    <Text style={styles.nutritionValue}>{recipe.carbs}g</Text>
                    <Text style={styles.nutritionLabel}>Carbs</Text>
                  </View>
                  <View style={styles.nutritionItem}>
                    <Text style={styles.nutritionValue}>{recipe.fat}g</Text>
                    <Text style={styles.nutritionLabel}>Fat</Text>
                  </View>
                </View>
              </View>

              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Ingredients</Text>
                {recipe.ingredients.map((ingredient, index) => (
                  <View key={index} style={styles.ingredientRow}>
                    <Ionicons name="checkmark-circle" size={20} color={COLORS.primary} />
                    <Text style={styles.ingredientText}>
                      {ingredient.amount} {ingredient.unit} {ingredient.item}
                    </Text>
                  </View>
                ))}
              </View>

              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Instructions</Text>
                {recipe.instructions.map((instruction, index) => (
                  <View key={index} style={styles.instructionRow}>
                    <View style={styles.stepNumber}>
                      <Text style={styles.stepNumberText}>{instruction.step_number}</Text>
                    </View>
                    <Text style={styles.instructionText}>{instruction.instruction}</Text>
                  </View>
                ))}
              </View>
            </View>
          </ScrollView>
        ) : (
          <View style={styles.emptyState}>
            <Text style={styles.emptyText}>Recipe not found</Text>
          </View>
        )}
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
  headerTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: SPACING['3xl'],
  },
  image: {
    width: '100%',
    height: 250,
    backgroundColor: COLORS.background.secondary,
  },
  content: {
    padding: SPACING.lg,
  },
  name: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginBottom: SPACING.sm,
  },
  description: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.lg,
    lineHeight: 24,
  },
  metaRow: {
    flexDirection: 'row',
    gap: SPACING.sm,
    marginBottom: SPACING.lg,
  },
  metaBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    paddingVertical: SPACING.xs,
    paddingHorizontal: SPACING.md,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.full,
  },
  metaText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    textTransform: 'capitalize',
  },
  tagsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
    marginBottom: SPACING.lg,
  },
  tag: {
    paddingVertical: SPACING.xs,
    paddingHorizontal: SPACING.md,
    backgroundColor: COLORS.primary,
    borderRadius: BORDER_RADIUS.full,
  },
  tagText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.primary,
  },
  section: {
    marginTop: SPACING.xl,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  nutritionGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  nutritionItem: {
    flex: 1,
    minWidth: '45%',
    padding: SPACING.lg,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    alignItems: 'center',
  },
  nutritionValue: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
    marginBottom: SPACING.xs,
  },
  nutritionLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  ingredientRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
    paddingVertical: SPACING.sm,
  },
  ingredientText: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  instructionRow: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginBottom: SPACING.lg,
  },
  stepNumber: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: COLORS.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  stepNumberText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  instructionText: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
    lineHeight: 24,
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
  },
  emptyText: {
    fontSize: FONT_SIZES.lg,
    color: COLORS.text.tertiary,
  },
});
