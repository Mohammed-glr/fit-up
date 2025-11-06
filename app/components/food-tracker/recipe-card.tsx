import React from 'react';
import { Image, StyleSheet, Text, TouchableOpacity, View, type StyleProp, type ViewStyle } from 'react-native';
import Ionicons from '@expo/vector-icons/Ionicons';
import type { UserRecipe } from '@/types/food-tracker';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';

type RecipeCardProps = {
  recipe: UserRecipe;
  onPress?: (recipe: UserRecipe) => void;
  onToggleFavorite?: (recipe: UserRecipe) => void;
  style?: StyleProp<ViewStyle>;
};

const resolveImageSource = (imageUrl?: string) => {
  if (imageUrl && imageUrl.length > 0) {
    return { uri: imageUrl };
  }
  return null;
};

const formatMacro = (label: string, value: number, unit: string) => `${label} ${Math.round(value)}${unit}`;

export function RecipeCard({ recipe, onPress, onToggleFavorite, style }: RecipeCardProps) {
  const handlePress = React.useCallback(() => {
    onPress?.(recipe);
  }, [onPress, recipe]);

  const handleToggleFavorite = React.useCallback(() => {
    onToggleFavorite?.(recipe);
  }, [onToggleFavorite, recipe]);

  const imageSource = resolveImageSource(recipe.image_url);

  return (
    <TouchableOpacity style={[styles.container, style]} activeOpacity={0.85} onPress={handlePress}>
      <View style={styles.imageWrapper}>
        {imageSource ? (
          <Image source={imageSource} style={styles.image} resizeMode="cover" />
        ) : (
          <View style={styles.placeholder}>
            <Ionicons name="fast-food" size={32} color={COLORS.primary} />
          </View>
        )}

        <TouchableOpacity
          style={[styles.favoriteButton, recipe.is_favorite && styles.favoriteButtonActive]}
          onPress={handleToggleFavorite}
          accessibilityRole="button"
          accessibilityLabel={recipe.is_favorite ? 'Remove from favorites' : 'Add to favorites'}
        >
          <Ionicons
            name={recipe.is_favorite ? 'heart' : 'heart-outline'}
            size={18}
            color={recipe.is_favorite ? COLORS.background.surface : COLORS.primary}
          />
        </TouchableOpacity>
      </View>

      <View style={styles.content}>
        <Text style={styles.name} numberOfLines={2}>
          {recipe.name}
        </Text>
        <Text style={styles.subtitle} numberOfLines={1}>
          {recipe.category.toUpperCase()} • {recipe.difficulty.toUpperCase()}
        </Text>

        <View style={styles.macroRow}>
          <Text style={styles.macroText}>{formatMacro('P', recipe.protein, 'g')}</Text>
          <Text style={styles.dot}>•</Text>
          <Text style={styles.macroText}>{formatMacro('C', recipe.carbs, 'g')}</Text>
          <Text style={styles.dot}>•</Text>
          <Text style={styles.macroText}>{formatMacro('F', recipe.fat, 'g')}</Text>
        </View>

        <View style={styles.footerRow}>
          <View style={styles.caloriesBadge}>
            <Ionicons name="flame" size={14} color={COLORS.background.surface} />
            <Text style={styles.calorieText}>{Math.round(recipe.calories)} kcal</Text>
          </View>
          <Text style={styles.servingsText}>{recipe.servings} servings</Text>
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: {
    width: 220,
    backgroundColor: COLORS.background.surface,
    borderRadius: BORDER_RADIUS.xl,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    overflow: 'hidden',
  },
  imageWrapper: {
    height: 140,
    backgroundColor: COLORS.background.secondary,
    position: 'relative',
  },
  image: {
    width: '100%',
    height: '100%',
  },
  placeholder: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  favoriteButton: {
    position: 'absolute',
    top: SPACING.sm,
    right: SPACING.sm,
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: COLORS.background.surface,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
  },
  favoriteButtonActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  content: {
    padding: SPACING.lg,
    gap: SPACING.sm,
  },
  name: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  subtitle: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
  },
  macroRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  macroText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  dot: {
    color: COLORS.text.placeholder,
  },
  footerRow: {
    marginTop: SPACING.xs,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  caloriesBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    backgroundColor: COLORS.primary,
    paddingVertical: 4,
    paddingHorizontal: SPACING.sm,
    borderRadius: BORDER_RADIUS.base,
  },
  calorieText: {
    color: COLORS.background.surface,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  servingsText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.xs,
  },
});
