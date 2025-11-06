import React from 'react';
import {
  ActivityIndicator,
  FlatList,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { useToggleFavoriteRecipe, useUserRecipeDetail, useUserRecipes } from '@/hooks/food-tracker/use-recipes';
import { RecipeCard } from '@/components/food-tracker/recipe-card';
import { RecipeDetailModal } from '@/components/food-tracker/recipe-detail-modal';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { RecipeCategory, UserRecipe } from '@/types/food-tracker';

const categories: Array<{ label: string; value?: RecipeCategory }> = [
  { label: 'All' },
  { label: 'Breakfast', value: 'breakfast' },
  { label: 'Lunch', value: 'lunch' },
  { label: 'Dinner', value: 'dinner' },
  { label: 'Snacks', value: 'snack' },
  { label: 'Dessert', value: 'dessert' },
];

const useDebouncedValue = (value: string, delay = 300) => {
  const [debouncedValue, setDebouncedValue] = React.useState(value);

  React.useEffect(() => {
    const handle = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => clearTimeout(handle);
  }, [value, delay]);

  return debouncedValue;
};

export default function RecipesScreen() {
  const [searchTerm, setSearchTerm] = React.useState('');
  const [selectedCategory, setSelectedCategory] = React.useState<RecipeCategory | undefined>();
  const [favoritesOnly, setFavoritesOnly] = React.useState(false);
  const [selectedRecipeId, setSelectedRecipeId] = React.useState<number | null>(null);

  const debouncedSearch = useDebouncedValue(searchTerm.trim());

  const queryParams = React.useMemo(
    () => ({
      limit: 40,
      offset: 0,
      search: debouncedSearch || undefined,
      category: selectedCategory,
      favorites_only: favoritesOnly || undefined,
      sort_by: 'updated_at',
      sort_order: 'desc' as const,
    }),
    [debouncedSearch, selectedCategory, favoritesOnly],
  );

  const {
    data: recipesResponse,
    isLoading,
    isRefetching,
    refetch,
  } = useUserRecipes(queryParams);

  const toggleFavoriteMutation = useToggleFavoriteRecipe();
  const {
    data: selectedRecipe,
    isLoading: selectedRecipeLoading,
  } = useUserRecipeDetail(selectedRecipeId);

  const recipes = recipesResponse?.recipes ?? [];
  const isRefreshing = isRefetching && !isLoading;

  const handleToggleFavorite = React.useCallback((recipe: UserRecipe) => {
    toggleFavoriteMutation.mutate({ recipeId: recipe.id });
  }, [toggleFavoriteMutation]);

  const handleSelectRecipe = React.useCallback((recipe: UserRecipe) => {
    setSelectedRecipeId(recipe.id);
  }, []);

  const handleCloseRecipeModal = React.useCallback(() => {
    setSelectedRecipeId(null);
  }, []);

  const renderRecipe = React.useCallback(({ item }: { item: UserRecipe }) => (
    <RecipeCard
      recipe={item}
      onPress={handleSelectRecipe}
      onToggleFavorite={handleToggleFavorite}
      style={styles.recipeCard}
    />
  ), [handleSelectRecipe, handleToggleFavorite]);

  const keyExtractor = React.useCallback((item: UserRecipe) => item.id.toString(), []);

  return (
    <View style={styles.container}>
      <FlatList
        data={recipes}
        keyExtractor={keyExtractor}
        renderItem={renderRecipe}
        ItemSeparatorComponent={() => <View style={{ height: SPACING.lg }} />}
        ListHeaderComponent={(
          <View style={styles.header}>
            <Text style={styles.heading}>Your recipes</Text>
            <Text style={styles.subheading}>Personal recipes you have created or saved.</Text>

            <View style={styles.searchWrapper}>
              <TextInput
                value={searchTerm}
                onChangeText={setSearchTerm}
                placeholder="Search recipes..."
                placeholderTextColor={COLORS.text.placeholder}
                style={styles.searchInput}
                returnKeyType="search"
              />
            </View>

            <View style={styles.filterRow}>
              {categories.map((category) => {
                const isActive = (!category.value && !selectedCategory) || category.value === selectedCategory;
                return (
                  <TouchableOpacity
                    key={category.label}
                    style={[styles.filterChip, isActive && styles.filterChipActive]}
                    onPress={() =>
                      setSelectedCategory((prev) =>
                        category.value ? (prev === category.value ? undefined : category.value) : undefined
                      )
                    }
                  >
                    <Text style={[styles.filterChipText, isActive && styles.filterChipTextActive]}>{category.label}</Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            <TouchableOpacity
              style={[styles.favoritesToggle, favoritesOnly && styles.favoritesToggleActive]}
              onPress={() => setFavoritesOnly((prev) => !prev)}
            >
              <Text style={[styles.favoritesToggleText, favoritesOnly && styles.favoritesToggleTextActive]}>
                {favoritesOnly ? 'Showing favorites' : 'Show favorites only'}
              </Text>
            </TouchableOpacity>
          </View>
        )}
        ListEmptyComponent={
          isLoading ? (
            <ActivityIndicator color={COLORS.primary} style={styles.loader} />
          ) : (
            <View style={styles.emptyState}>
              <Text style={styles.emptyTitle}>No recipes found</Text>
              <Text style={styles.emptySubtitle}>
                Try adjusting your filters or create a new recipe to get started.
              </Text>
            </View>
          )
        }
        contentContainerStyle={styles.listContent}
        refreshing={isRefreshing}
        onRefresh={refetch}
      />

      <RecipeDetailModal
        visible={selectedRecipeId !== null}
        recipe={selectedRecipe}
        isLoading={selectedRecipeLoading}
        onClose={handleCloseRecipeModal}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.primary,
  },
  listContent: {
    padding: SPACING.xl,
    paddingBottom: SPACING['3xl'],
  },
  header: {
    gap: SPACING.md,
    marginBottom: SPACING.lg,
  },
  heading: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
  },
  subheading: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.base,
  },
  searchWrapper: {
    marginTop: SPACING.sm,
  },
  searchInput: {
    height: 48,
    borderRadius: BORDER_RADIUS.base,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    paddingHorizontal: SPACING.lg,
    backgroundColor: COLORS.background.surface,
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.base,
  },
  filterRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginTop: SPACING.sm,
  },
  filterChip: {
    paddingHorizontal: SPACING.md,
    paddingVertical: 8,
    borderRadius: BORDER_RADIUS.full,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    marginRight: SPACING.sm,
    marginBottom: SPACING.sm,
    backgroundColor: COLORS.background.surface,
  },
  filterChipActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  filterChipText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  filterChipTextActive: {
    color: COLORS.background.surface,
  },
  favoritesToggle: {
    marginTop: SPACING.sm,
    paddingVertical: SPACING.sm,
    paddingHorizontal: SPACING.lg,
    borderRadius: BORDER_RADIUS.base,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    backgroundColor: COLORS.background.surface,
  },
  favoritesToggleActive: {
    borderColor: COLORS.primary,
    backgroundColor: COLORS.primarySoft,
  },
  favoritesToggleText: {
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    textAlign: 'center',
  },
  favoritesToggleTextActive: {
    color: COLORS.primary,
  },
  recipeCard: {
    width: '100%',
  },
  loader: {
    marginTop: SPACING.lg,
  },
  emptyState: {
    marginTop: SPACING['2xl'],
    padding: SPACING['2xl'],
    borderRadius: BORDER_RADIUS.xl,
    backgroundColor: COLORS.background.secondary,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
    alignItems: 'center',
  },
  emptyTitle: {
    color: COLORS.text.primary,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  emptySubtitle: {
    marginTop: SPACING.sm,
    color: COLORS.text.secondary,
    fontSize: FONT_SIZES.sm,
    textAlign: 'center',
  },
});
