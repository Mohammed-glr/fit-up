import React from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import type { UserRecipe } from '@/types/food-tracker';
import { RecipeCard } from './recipe-card';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING } from '@/constants/theme';

type RecipeCarouselProps = {
  title: string;
  recipes?: UserRecipe[];
  onPressRecipe?: (recipe: UserRecipe) => void;
  onToggleFavorite?: (recipe: UserRecipe) => void;
  emptyMessage?: string;
};

export function RecipeCarousel({ title, recipes, onPressRecipe, onToggleFavorite, emptyMessage }: RecipeCarouselProps) {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{title}</Text>
      {recipes && recipes.length > 0 ? (
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          {recipes.map((recipe, index) => (
            <View key={recipe.id} style={[styles.cardWrapper, index === recipes.length - 1 && styles.lastCardWrapper]}>
              <RecipeCard
                recipe={recipe}
                onPress={onPressRecipe}
                onToggleFavorite={onToggleFavorite}
              />
            </View>
          ))}
        </ScrollView>
      ) : (
        <Text style={styles.emptyText}>{emptyMessage ?? 'No recipes yet.'}</Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: SPACING.lg,
  },
  title: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  scrollContent: {
    paddingLeft: SPACING.md,
    paddingRight: SPACING.md,
  },
  emptyText: {
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
  },
  cardWrapper: {
    marginRight: SPACING.md,
  },
  lastCardWrapper: {
    marginRight: 0,
  },
});
