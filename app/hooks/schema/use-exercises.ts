import { useQuery } from '@tanstack/react-query';
import { exerciseService } from '@/api/services/schema-service';
import type { Exercise, PaginatedResponse } from '@/types/schema';
import { APIError } from '@/api/client';

export const exerciseKeys = {
  all: ['exercises'] as const,
  lists: () => [...exerciseKeys.all, 'list'] as const,
  list: (filters: string) => [...exerciseKeys.lists(), { filters }] as const,
  details: () => [...exerciseKeys.all, 'detail'] as const,
  detail: (id: number) => [...exerciseKeys.details(), id] as const,
  recommended: () => [...exerciseKeys.all, 'recommended'] as const,
  mostUsed: () => [...exerciseKeys.all, 'most-used'] as const,
  byMuscleGroup: (group: string) => [...exerciseKeys.all, 'muscle-group', group] as const,
  byEquipment: (equipment: string) => [...exerciseKeys.all, 'equipment', equipment] as const,
  search: (query: string) => [...exerciseKeys.all, 'search', query] as const,
  usageStats: (id: number) => [...exerciseKeys.detail(id), 'stats'] as const,
};

export const useExercises = (params?: { limit?: number; offset?: number }) => {
  return useQuery<Exercise[], APIError>({
    queryKey: exerciseKeys.list(JSON.stringify(params || {})),
    queryFn: () => exerciseService.List(),
    staleTime: 10 * 60 * 1000, 
  });
};

export const useExercise = (id: number) => {
  return useQuery({
    queryKey: exerciseKeys.detail(id),
    queryFn: () => exerciseService.ListById(id),
    enabled: !!id,
    staleTime: 10 * 60 * 1000,
  });
};

export const useFilterExercises = (filters: {
  muscle_groups?: string[];
  difficulty?: string;
  equipment?: string[];
  type?: string[];
  search?: string;
}) => {
  const filterKey = JSON.stringify(filters);
  
  return useQuery<PaginatedResponse<Exercise>, APIError>({
    queryKey: exerciseKeys.list(filterKey),
    queryFn: () => exerciseService.Filter({
      muscleGroups: filters.muscle_groups,
      equipment: filters.equipment,
      exerciseTypes: filters.type,
      fitnessLevels: filters.difficulty ? [filters.difficulty] : undefined,
    }),
    enabled: Object.keys(filters).length > 0,
    staleTime: 5 * 60 * 1000,
  });
};

export const useSearchExercises = (query: string) => {
  return useQuery({
    queryKey: exerciseKeys.search(query),
    queryFn: () => exerciseService.Search(),
    enabled: query.length > 0,
    staleTime: 5 * 60 * 1000,
  });
};

export const useExercisesByMuscleGroup = (muscleGroup: string) => {
  return useQuery({
    queryKey: exerciseKeys.byMuscleGroup(muscleGroup),
    queryFn: () => exerciseService.ListByMuscleGroup(muscleGroup),
    enabled: !!muscleGroup,
    staleTime: 10 * 60 * 1000,
  });
};

export const useExercisesByEquipment = (equipment: string) => {
  return useQuery({
    queryKey: exerciseKeys.byEquipment(equipment),
    queryFn: () => exerciseService.GetEquipment(equipment),
    enabled: !!equipment,
    staleTime: 10 * 60 * 1000,
  });
};

export const useRecommendedExercises = () => {
  return useQuery({
    queryKey: exerciseKeys.recommended(),
    queryFn: () => exerciseService.GetRecommended(),
    staleTime: 15 * 60 * 1000,
  });
};

export const useMostUsedExercises = () => {
  return useQuery({
    queryKey: exerciseKeys.mostUsed(),
    queryFn: () => exerciseService.GetMostUsed(),
    staleTime: 15 * 60 * 1000,
  });
};

export const useExerciseUsageStats = (id: number) => {
  return useQuery({
    queryKey: exerciseKeys.usageStats(id),
    queryFn: async () => {
      // This method doesn't exist yet in the service, return empty for now
      return { stats: {} };
    },
    enabled: !!id,
    staleTime: 10 * 60 * 1000,
  });
};
