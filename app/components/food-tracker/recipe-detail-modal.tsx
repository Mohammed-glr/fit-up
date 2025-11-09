import React from 'react';
import { Modal, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import Ionicons from '@expo/vector-icons/Ionicons';
import type { UserRecipeDetail } from '@/types/food-tracker';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';

type RecipeDetailModalProps = {
  visible: boolean;
  recipe?: UserRecipeDetail | null;
  isLoading?: boolean;
  onClose: () => void;
  onEdit?: (recipe: UserRecipeDetail) => void;
  onDelete?: (recipeId: number) => void;
};

export function RecipeDetailModal({ visible, recipe, isLoading, onClose, onEdit, onDelete }: RecipeDetailModalProps) {
  const handleEdit = React.useCallback(() => {
    if (recipe && onEdit) {
      onEdit(recipe);
    }
  }, [recipe, onEdit]);

  const handleDelete = React.useCallback(() => {
    if (recipe && onDelete) {
      onDelete(recipe.id);
    }
  }, [recipe, onDelete]);

  return (
    <Modal visible={visible} animationType="slide" transparent onRequestClose={onClose}>
      <View style={styles.overlay}>
        <View style={styles.sheet}>
          <View style={styles.handle} />
          <View style={styles.header}>
            <Text style={styles.title} numberOfLines={2}>
              {recipe?.name ?? 'Recipe'}
            </Text>
            <View style={styles.headerActions}>
              {onEdit && (
                <TouchableOpacity onPress={handleEdit} style={styles.actionButton}>
                  <Ionicons name="pencil" size={22} color={COLORS.primary} />
                </TouchableOpacity>
              )}
              {onDelete && (
                <TouchableOpacity onPress={handleDelete} style={styles.actionButton}>
                  <Ionicons name="trash-outline" size={22} color={COLORS.error} />
                </TouchableOpacity>
              )}
              <TouchableOpacity onPress={onClose} accessibilityRole="button" accessibilityLabel="Close recipe details">
                <Ionicons name="close" size={24} color={COLORS.text.tertiary} />
              </TouchableOpacity>
            </View>
          </View>

          <ScrollView contentContainerStyle={styles.content}>
            {isLoading && !recipe ? (
              <Text style={styles.loadingText}>Loading recipe...</Text>
            ) : (
              <>
                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Category</Text>
                  <Text style={styles.sectionValue}>{recipe?.category?.toUpperCase()}</Text>
                </View>

                <View style={[styles.section, styles.sectionRow]}>
                  <View style={[styles.chip, styles.chipSpacing]}>
                    <Ionicons name="time" size={14} color={COLORS.text.secondary} />
                    <Text style={styles.chipText}>{recipe?.prep_time ?? 0} min prep</Text>
                  </View>
                  <View style={[styles.chip, styles.chipSpacing]}>
                    <Ionicons name="flame" size={14} color={COLORS.text.secondary} />
                    <Text style={styles.chipText}>{Math.round(recipe?.calories ?? 0)} kcal</Text>
                  </View>
                  <View style={[styles.chip, styles.chipSpacing]}>
                    <Text style={styles.chipText}>{recipe?.servings ?? 0} servings</Text>
                  </View>
                </View>

                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Description</Text>
                  <Text style={styles.bodyText}>{recipe?.description ?? 'No description provided.'}</Text>
                </View>

                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Ingredients</Text>
                  {recipe?.ingredients?.map((ingredient) => (
                    <Text key={ingredient.id} style={styles.listItemText}>
                      • {ingredient.amount} {ingredient.unit} {ingredient.item}
                    </Text>
                  ))}
                  {(!recipe?.ingredients || recipe.ingredients.length === 0) && (
                    <Text style={styles.bodyText}>No ingredients listed.</Text>
                  )}
                </View>

                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Instructions</Text>
                  {recipe?.instructions?.map((instruction) => (
                    <Text key={instruction.id} style={styles.listItemText}>
                      {instruction.step_number}. {instruction.instruction}
                    </Text>
                  ))}
                  {(!recipe?.instructions || recipe.instructions.length === 0) && (
                    <Text style={styles.bodyText}>No instructions provided.</Text>
                  )}
                </View>
              </>
            )}
          </ScrollView>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.4)',
    justifyContent: 'flex-end',
  },
  sheet: {
    maxHeight: '90%',
    backgroundColor: COLORS.darkGray,
    borderTopLeftRadius: BORDER_RADIUS['3xl'],
    borderTopRightRadius: BORDER_RADIUS['3xl'],
    paddingBottom: SPACING['3xl'],
  },
  handle: {
    width: 48,
    height: 5,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.border.subtle,
    alignSelf: 'center',
    marginTop: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: SPACING.xl,
    marginBottom: SPACING.md,
  },
  headerActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
  },
  actionButton: {
    padding: SPACING.xs,
  },
  title: {
    flex: 1,
    marginRight: SPACING.lg,
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  content: {
    paddingHorizontal: SPACING.xl,
    paddingBottom: SPACING['3xl'],
  },
  loadingText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.base,
  },
  section: {
    marginBottom: SPACING.lg,
  },
  sectionRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginBottom: 0,
  },
  chipSpacing: {
    marginRight: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  sectionLabel: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  sectionValue: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
  },
  bodyText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.base,
    lineHeight: FONT_SIZES.base * 1.4,
  },
  listItemText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.base,
    lineHeight: FONT_SIZES.base * 1.5,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 6,
    paddingHorizontal: SPACING.sm,
    borderRadius: BORDER_RADIUS.base,
    backgroundColor: COLORS.background.accent,
  },
  chipText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    marginLeft: 6,
  },
});
