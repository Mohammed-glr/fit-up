import React from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import {
  useToggleFavoriteRecipe,
  useUserRecipeDetail,
  useUserRecipes,
  useCreateUserRecipe,
  useUpdateUserRecipe,
  useDeleteUserRecipe,
} from '@/hooks/food-tracker/use-recipes';
import { RecipeCard } from '@/components/food-tracker/recipe-card';
import { RecipeDetailModal } from '@/components/food-tracker/recipe-detail-modal';
import { RecipeFormModal } from '@/components/food-tracker/recipe-form-modal';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { RecipeCategory, UserRecipe, SystemRecipe, CreateRecipeRequest } from '@/types/food-tracker';
import { Ionicons } from '@expo/vector-icons';
import { useRecipeContext } from '@/context/recipe-context';

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
  const [isFormVisible, setIsFormVisible] = React.useState(false);
  const [editingRecipe, setEditingRecipe] = React.useState<UserRecipe | null>(null);

  const { setOnCreateRecipe } = useRecipeContext();

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
    error,
    refetch,
  } = useUserRecipes(queryParams);

  const toggleFavoriteMutation = useToggleFavoriteRecipe();
  const createRecipeMutation = useCreateUserRecipe();
  const updateRecipeMutation = useUpdateUserRecipe();
  const deleteRecipeMutation = useDeleteUserRecipe();

  const {
    data: selectedRecipe,
    isLoading: selectedRecipeLoading,
  } = useUserRecipeDetail(selectedRecipeId);

  React.useEffect(() => {
    if (__DEV__) {
      console.log('[RecipesScreen] Query params:', queryParams);
      console.log('[RecipesScreen] Response:', recipesResponse);
      console.log('[RecipesScreen] Error:', error);
      console.log('[RecipesScreen] IsLoading:', isLoading);
    }
  }, [recipesResponse, error, isLoading, queryParams]);

  const recipes = recipesResponse?.recipes ?? [];
  const isRefreshing = isRefetching && !isLoading;

  const handleToggleFavorite = React.useCallback(
    (recipe: UserRecipe | SystemRecipe) => {
      if ('is_favorite' in recipe) {
        toggleFavoriteMutation.mutate({ recipeId: recipe.id });
      }
    },
    [toggleFavoriteMutation]
  );

  const handleSelectRecipe = React.useCallback((recipe: UserRecipe | SystemRecipe) => {
    if ('user_id' in recipe) {
      setSelectedRecipeId(recipe.id);
    }
  }, []);

  const handleCloseRecipeModal = React.useCallback(() => {
    setSelectedRecipeId(null);
  }, []);

  const handleCreateRecipe = React.useCallback(() => {
    setEditingRecipe(null);
    setIsFormVisible(true);
  }, []);

  const handleEditRecipe = React.useCallback((recipe: UserRecipe) => {
    setEditingRecipe(recipe);
    setIsFormVisible(true);
    setSelectedRecipeId(null);
  }, []);

  const handleDeleteRecipe = React.useCallback(
    (recipeId: number) => {
      Alert.alert(
        'Delete Recipe',
        'Are you sure you want to delete this recipe? This action cannot be undone.',
        [
          { text: 'Cancel', style: 'cancel' },
          {
            text: 'Delete',
            style: 'destructive',
            onPress: () => {
              deleteRecipeMutation.mutate(
                { recipeId },
                {
                  onSuccess: () => {
                    setSelectedRecipeId(null);
                    Alert.alert('Success', 'Recipe deleted successfully');
                  },
                  onError: (error) => {
                    Alert.alert('Error', error.message || 'Failed to delete recipe');
                  },
                }
              );
            },
          },
        ]
      );
    },
    [deleteRecipeMutation]
  );

  const handleFormSubmit = React.useCallback(
    (data: CreateRecipeRequest) => {
      if (editingRecipe) {
        updateRecipeMutation.mutate(
          { recipeId: editingRecipe.id, data },
          {
            onSuccess: () => {
              setIsFormVisible(false);
              setEditingRecipe(null);
              Alert.alert('Success', 'Recipe updated successfully');
            },
            onError: (error) => {
              Alert.alert('Error', error.message || 'Failed to update recipe');
            },
          }
        );
      } else {
        createRecipeMutation.mutate(data, {
          onSuccess: () => {
            setIsFormVisible(false);
            Alert.alert('Success', 'Recipe created successfully');
          },
          onError: (error) => {
            Alert.alert('Error', error.message || 'Failed to create recipe');
          },
        });
      }
    },
    [editingRecipe, createRecipeMutation, updateRecipeMutation]
  );

  const handleCloseForm = React.useCallback(() => {
    setIsFormVisible(false);
    setEditingRecipe(null);
  }, []);

  React.useEffect(() => {
    setOnCreateRecipe(handleCreateRecipe);
    return () => setOnCreateRecipe(undefined);
  }, [handleCreateRecipe, setOnCreateRecipe]);

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
            <View style={styles.titleContainer}>
              <Ionicons
                name="book"
                size={28}
                color={COLORS.primary}
              />
              <Text style={styles.headerTitle}>My Recipes</Text>
            </View>
            <Text style={styles.headerSubtitle}>Personal recipes you have created or saved.</Text>

            <View style={styles.searchWrapper}>
              <Ionicons
                name="search"
                size={20}
                color={COLORS.text.tertiary}
                style={styles.searchIcon}
              />
              <TextInput
                value={searchTerm}
                onChangeText={setSearchTerm}
                placeholder="Search recipes..."
                placeholderTextColor={COLORS.text.placeholder}
                style={styles.searchInput}
                returnKeyType="search"
              />
              {searchTerm.length > 0 && (
                <TouchableOpacity onPress={() => setSearchTerm('')}>
                  <Ionicons name="close-circle" size={20} color={COLORS.text.tertiary} />
                </TouchableOpacity>
              )}
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
          ) : error ? (
            <View style={styles.emptyState}>
              <Text style={styles.emptyTitle}>Error loading recipes</Text>
              <Text style={styles.emptySubtitle}>
                {error.message || 'Something went wrong. Please try again.'}
              </Text>
            </View>
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
        onEdit={handleEditRecipe}
        onDelete={handleDeleteRecipe}
      />

      <RecipeFormModal
        visible={isFormVisible}
        recipe={editingRecipe ? selectedRecipe : null}
        onClose={handleCloseForm}
        onSubmit={handleFormSubmit}
        isSubmitting={createRecipeMutation.isPending || updateRecipeMutation.isPending}
      />
    </View>
  );
}

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  listContent: {
    padding: SPACING.base,
    paddingBottom: SPACING['3xl'],
  },
  header: {
    gap: SPACING.md,
    marginBottom: SPACING.lg,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 4,
  },
  title: {
    fontSize: 32,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  subtitle: {
    fontSize: 16,
    color: '#888888',
  },
  headerInfo: {
    flex: 1,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  headerSubtitle: {
    fontSize: 14,
    color: '#888888',
    marginTop: 2,
  },
  searchWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    marginTop: SPACING.md,
    paddingHorizontal: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  searchInput: {
    flex: 1,
    paddingVertical: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
  },
  searchIcon: {
    marginRight: SPACING.sm,
  },
  filterRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginTop: SPACING.sm,
  },
  filterChip: {
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    marginRight: SPACING.sm,
    marginBottom: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  filterChipActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  filterChipText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
  },
  filterChipTextActive: {
    color: COLORS.text.primary,
  },
  favoritesToggle: {
    marginTop: SPACING.sm,
    paddingVertical: SPACING.sm,
    paddingHorizontal: SPACING.lg,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
  },
  favoritesToggleActive: {
    backgroundColor: COLORS.primary,
    
  },
  favoritesToggleText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    textAlign: 'center',
  },
  favoritesToggleTextActive: {
    color: COLORS.text.primary,
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
    borderRadius: BORDER_RADIUS['2xl'],
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
  },
  emptyTitle: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  emptySubtitle: {
    marginTop: SPACING.sm,
    color: COLORS.text.placeholder,
    fontSize: FONT_SIZES.sm,
    textAlign: 'center',
  },
});
