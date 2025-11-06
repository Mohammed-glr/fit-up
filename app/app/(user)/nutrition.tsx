import React from 'react';
import { RefreshControl, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { useRouter } from 'expo-router';
import { DateToggle } from '@/components/food-tracker/date-toggle';
import { NutritionSummaryCard } from '@/components/food-tracker/nutrition-summary-card';
import { FoodLogList } from '@/components/food-tracker/food-log-list';
import { RecipeCarousel } from '@/components/food-tracker/recipe-carousel';
import { useDailyNutritionSummary, useNutritionGoals } from '@/hooks/food-tracker/use-nutrition';
import { useFoodLogsByDate } from '@/hooks/food-tracker/use-food-logs';
import { useToggleFavoriteRecipe, useUserRecipeDetail, useUserRecipes } from '@/hooks/food-tracker/use-recipes';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { UserRecipe } from '@/types/food-tracker';
import { RecipeDetailModal } from '@/components/food-tracker/recipe-detail-modal';

const todayKey = (() => {
  const now = new Date();
  const year = now.getFullYear();
  const month = `${now.getMonth() + 1}`.padStart(2, '0');
  const day = `${now.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
})();

export default function NutritionScreen() {
  const router = useRouter();
  const [selectedDate, setSelectedDate] = React.useState<string>(todayKey);
  const [selectedRecipeId, setSelectedRecipeId] = React.useState<number | null>(null);

  const {
    data: summary,
    isLoading: summaryLoading,
    isRefetching: summaryRefetching,
    refetch: refetchSummary,
  } = useDailyNutritionSummary(selectedDate);

  const {
    data: goals,
    isLoading: goalsLoading,
    isRefetching: goalsRefetching,
    refetch: refetchGoals,
  } = useNutritionGoals();

  const {
    data: logsResponse,
    isLoading: logsLoading,
    isRefetching: logsRefetching,
    refetch: refetchLogs,
  } = useFoodLogsByDate(selectedDate);

  const {
    data: favoritesResponse,
    isLoading: favoritesLoading,
    isFetching: favoritesFetching,
    refetch: refetchFavorites,
  } = useUserRecipes({ favorites_only: true, limit: 10 });

  const toggleFavoriteMutation = useToggleFavoriteRecipe();
  const {
    data: selectedRecipe,
    isLoading: selectedRecipeLoading,
  } = useUserRecipeDetail(selectedRecipeId);

  const favoriteRecipes = favoritesResponse?.recipes ?? [];
  const foodLogs = logsResponse?.logs ?? [];

  const isRefreshing = summaryRefetching || logsRefetching || goalsRefetching || (favoritesFetching && !favoritesLoading);
  const isLoading = summaryLoading || logsLoading || goalsLoading;

  const handleRefresh = React.useCallback(async () => {
    await Promise.all([
      refetchSummary(),
      refetchLogs(),
      refetchGoals(),
      refetchFavorites(),
    ]);
  }, [refetchSummary, refetchLogs, refetchGoals, refetchFavorites]);

  const handleToggleFavorite = React.useCallback((recipe: UserRecipe) => {
    toggleFavoriteMutation.mutate({ recipeId: recipe.id });
  }, [toggleFavoriteMutation]);

  const handleSelectRecipe = React.useCallback((recipe: UserRecipe) => {
    setSelectedRecipeId(recipe.id);
  }, []);

  const handleCloseRecipeModal = React.useCallback(() => {
    setSelectedRecipeId(null);
  }, []);

  const handleViewAllRecipes = React.useCallback(() => {
    router.push('/(user)/recipes');
  }, [router]);

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.contentContainer}
      refreshControl={(
        <RefreshControl refreshing={isRefreshing} onRefresh={handleRefresh} tintColor={COLORS.primary} />
      )}
    >
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Daily nutrition</Text>
        <DateToggle date={selectedDate} onChange={setSelectedDate} />
      </View>

      <View style={styles.section}>
        <NutritionSummaryCard
          summary={summary}
          goals={goals}
          isLoading={isLoading}
          onPressSetGoals={() => router.push('/(user)/profile')}
        />
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Food log</Text>
          <TouchableOpacity style={styles.ctaButton} onPress={() => router.push('/(user)/plans')}>
            <Text style={styles.ctaText}>Log meal</Text>
          </TouchableOpacity>
        </View>
        <FoodLogList entries={foodLogs} />
      </View>

      <View style={styles.section}>
        <RecipeCarousel
          title="Favorite recipes"
          recipes={favoriteRecipes}
          emptyMessage="Mark recipes as favorites to see them here."
          onPressRecipe={handleSelectRecipe}
          onToggleFavorite={handleToggleFavorite}
        />
        <TouchableOpacity style={styles.viewAllButton} onPress={handleViewAllRecipes}>
          <Text style={styles.viewAllText}>Browse all recipes</Text>
        </TouchableOpacity>
      </View>

      <RecipeDetailModal
        visible={selectedRecipeId !== null}
        recipe={selectedRecipe}
        isLoading={selectedRecipeLoading}
        onClose={handleCloseRecipeModal}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.primary,
  },
  contentContainer: {
    padding: SPACING.xl,
    paddingBottom: SPACING['3xl'],
  },
  section: {
    marginBottom: SPACING['2xl'],
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  sectionTitle: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  ctaButton: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.base,
    backgroundColor: COLORS.primary,
  },
  ctaText: {
    color: COLORS.background.surface,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  viewAllButton: {
    marginTop: SPACING.base,
    paddingVertical: SPACING.sm,
    alignItems: 'center',
    borderRadius: BORDER_RADIUS.base,
    borderWidth: 1,
    borderColor: COLORS.border.light,
  },
  viewAllText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
});
