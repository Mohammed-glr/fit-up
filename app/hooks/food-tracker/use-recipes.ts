import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import recipeAPI from '@/api/services/recipe-service';
import type {
  CreateRecipeRequest,
  ListUserRecipesResponse,
  RecipeListParams,
  RecipeSearchParams,
  SearchRecipesResponse,
  UserRecipeDetail,
  GetFavoritesResponse,
} from '@/types/food-tracker';
import { APIError } from '@/api/client';
import { recipeKeys } from './keys';

type ToggleFavoriteVariables = { recipeId: number };
type UpdateRecipeVariables = { recipeId: number; data: CreateRecipeRequest };
type DeleteRecipeVariables = { recipeId: number };

type SearchParams = RecipeSearchParams & { favorites_only?: boolean };

const { userRecipeService } = recipeAPI;

export const useUserRecipes = (params?: RecipeListParams, options?: { enabled?: boolean }) => {
  return useQuery<ListUserRecipesResponse, APIError>({
    queryKey: recipeKeys.list(params),
    queryFn: () => userRecipeService.List(params),
    enabled: options?.enabled ?? true,
  });
};

export const useRecipeFavorites = (options?: { enabled?: boolean }) => {
  return useQuery<GetFavoritesResponse, APIError>({
    queryKey: recipeKeys.favorites,
    queryFn: () => userRecipeService.GetFavorites(),
    enabled: options?.enabled ?? true,
  });
};

export const useRecipeSearch = (params?: SearchParams, options?: { enabled?: boolean }) => {
  return useQuery<SearchRecipesResponse, APIError>({
    queryKey: recipeKeys.search(params),
    queryFn: () => userRecipeService.Search(params),
    enabled: options?.enabled ?? Boolean(params && Object.keys(params).length > 0),
  });
};

export const useUserRecipeDetail = (recipeId?: number | null) => {
  const enabled = typeof recipeId === 'number' && recipeId > 0;

  return useQuery<UserRecipeDetail, APIError>({
    queryKey: recipeKeys.detail(recipeId ?? null),
    queryFn: () => userRecipeService.Retrieve(recipeId as number),
    enabled,
  });
};

export const useCreateUserRecipe = () => {
  const queryClient = useQueryClient();

  return useMutation<UserRecipeDetail, APIError, CreateRecipeRequest>({
    mutationFn: (data) => userRecipeService.Create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: recipeKeys.all });
      queryClient.invalidateQueries({ queryKey: recipeKeys.favorites });
    },
  });
};

export const useUpdateUserRecipe = () => {
  const queryClient = useQueryClient();

  return useMutation<UserRecipeDetail, APIError, UpdateRecipeVariables>({
    mutationFn: ({ recipeId, data }) => userRecipeService.Update(recipeId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: recipeKeys.all });
      queryClient.invalidateQueries({ queryKey: recipeKeys.detail(variables.recipeId) });
      queryClient.invalidateQueries({ queryKey: recipeKeys.favorites });
    },
  });
};

export const useDeleteUserRecipe = () => {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, DeleteRecipeVariables>({
    mutationFn: ({ recipeId }) => userRecipeService.Delete(recipeId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: recipeKeys.all });
      queryClient.invalidateQueries({ queryKey: recipeKeys.detail(variables.recipeId) });
      queryClient.invalidateQueries({ queryKey: recipeKeys.favorites });
    },
  });
};

export const useToggleFavoriteRecipe = () => {
  const queryClient = useQueryClient();

  return useMutation<string | undefined, APIError, ToggleFavoriteVariables>({
    mutationFn: ({ recipeId }) => userRecipeService.ToggleFavorite(recipeId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: recipeKeys.all });
      queryClient.invalidateQueries({ queryKey: recipeKeys.detail(variables.recipeId) });
      queryClient.invalidateQueries({ queryKey: recipeKeys.favorites });
    },
  });
};
