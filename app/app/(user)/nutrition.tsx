import React from 'react';
import { ActivityIndicator, Alert, RefreshControl, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { useRouter } from 'expo-router';
import { DateToggle } from '@/components/food-tracker/date-toggle';
import { NutritionSummaryCard } from '@/components/food-tracker/nutrition-summary-card';
import { FoodLogList } from '@/components/food-tracker/food-log-list';
import { RecipeCarousel } from '@/components/food-tracker/recipe-carousel';
import { FoodLogModal } from '@/components/food-tracker/food-log-modal';
import { useDailyNutritionSummary, useNutritionGoals } from '@/hooks/food-tracker/use-nutrition';
import { useFoodLogsByDate, useLogFood } from '@/hooks/food-tracker/use-food-logs';
import { useToggleFavoriteRecipe, useUserRecipeDetail, useUserRecipes } from '@/hooks/food-tracker/use-recipes';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { UserRecipe, CreateFoodLogRequest } from '@/types/food-tracker';
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
  const [isLogModalVisible, setIsLogModalVisible] = React.useState(false);
  const [selectedRecipeForLog, setSelectedRecipeForLog] = React.useState<UserRecipe | null>(null);

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
  const logFoodMutation = useLogFood();
  const {
    data: selectedRecipe,
    isLoading: selectedRecipeLoading,
  } = useUserRecipeDetail(selectedRecipeId);

  const favoriteRecipes = favoritesResponse?.recipes ?? [];
  const foodLogs = logsResponse?.logs ?? [];

  const isRefreshing = summaryRefetching || logsRefetching || goalsRefetching || (favoritesFetching && !favoritesLoading);
  const summaryCardLoading = summaryLoading || goalsLoading;
  const favoritesAreLoading = favoritesLoading && favoriteRecipes.length === 0;

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

  const handleBrowseRecipes = React.useCallback(() => {
    router.push('/(user)/system-recipes');
  }, [router]);

  const handleLogMeal = React.useCallback(() => {
    setSelectedRecipeForLog(null);
    setIsLogModalVisible(true);
  }, []);

  const handleLogRecipe = React.useCallback((recipe: UserRecipe) => {
    setSelectedRecipeForLog(recipe);
    setIsLogModalVisible(true);
  }, []);

  const handleCloseLogModal = React.useCallback(() => {
    setIsLogModalVisible(false);
    setSelectedRecipeForLog(null);
  }, []);

  const handleSubmitFoodLog = React.useCallback((data: CreateFoodLogRequest) => {
    logFoodMutation.mutate(data, {
      onSuccess: () => {
        handleCloseLogModal();
        Alert.alert('Success', 'Food logged successfully');
      },
      onError: (error) => {
        Alert.alert('Error', error.message || 'Failed to log food');
      },
    });
  }, [logFoodMutation, handleCloseLogModal]);

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
          isLoading={summaryCardLoading}
          onPressSetGoals={() => router.push('/(user)/profile')}
        />
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Food log</Text>
          <TouchableOpacity style={styles.ctaButton} onPress={handleLogMeal}>
            <Text style={styles.ctaText}>Log meal</Text>
          </TouchableOpacity>
        </View>
        {logsLoading ? (
          <ActivityIndicator color={COLORS.primary} />
        ) : (
          <FoodLogList entries={foodLogs} />
        )}
      </View>

      <View style={styles.section}>
        {favoritesAreLoading ? (
          <ActivityIndicator color={COLORS.primary} />
        ) : (
          <RecipeCarousel
            title="Favorite recipes"
            recipes={favoriteRecipes}
            emptyMessage="Mark recipes as favorites to see them here."
            onPressRecipe={handleSelectRecipe}
            onToggleFavorite={handleToggleFavorite}
          />
        )}
        <View style={styles.recipeButtonsRow}>
          <TouchableOpacity style={[styles.viewAllButton, styles.primaryButton]} onPress={handleViewAllRecipes}>
            <Text style={[styles.viewAllText, styles.primaryButtonText]}>My Recipes</Text>
          </TouchableOpacity>
          <TouchableOpacity style={[styles.viewAllButton, styles.secondaryButton]} onPress={handleBrowseRecipes}>
            <Text style={[styles.viewAllText, styles.secondaryButtonText]}>Browse All</Text>
          </TouchableOpacity>
        </View>
      </View>

      <RecipeDetailModal
        visible={selectedRecipeId !== null}
        recipe={selectedRecipe}
        isLoading={selectedRecipeLoading}
        onClose={handleCloseRecipeModal}
      />

      <FoodLogModal
        visible={isLogModalVisible}
        date={selectedDate}
        selectedRecipe={selectedRecipeForLog}
        onClose={handleCloseLogModal}
        onSubmit={handleSubmitFoodLog}
        isSubmitting={logFoodMutation.isPending}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  contentContainer: {
    padding: SPACING.xl,
    paddingBottom: SPACING['6xl'],
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
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.semibold,
    marginBottom: SPACING.md,
  },
  ctaButton: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    marginBottom: SPACING.md,
  },
  ctaText: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  recipeButtonsRow: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginTop: SPACING.base,
  },
  viewAllButton: {
    flex: 1,
    marginTop: SPACING.base,
    paddingVertical: SPACING.md,
    alignItems: 'center',
    borderRadius: BORDER_RADIUS.full,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  primaryButton: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  secondaryButton: {
    backgroundColor: 'transparent',
  },
  viewAllText: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  primaryButtonText: {
    color: COLORS.text.primary,
  },
  secondaryButtonText: {
    color: COLORS.text.inverse,
  },
});
