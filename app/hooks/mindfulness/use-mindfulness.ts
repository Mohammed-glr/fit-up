import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import * as mindfulnessService from '../../api/services/mindfulness-service';
import type {
  CreateMindfulnessSessionRequest,
  CreateBreathingExerciseRequest,
  CreateGratitudeEntryRequest,
  CreateReflectionResponseRequest,
} from '../../types/mindfulness';

// Query keys
export const mindfulnessKeys = {
  all: ['mindfulness'] as const,
  sessions: () => [...mindfulnessKeys.all, 'sessions'] as const,
  sessionsList: (limit: number) =>
    [...mindfulnessKeys.sessions(), { limit }] as const,
  sessionsStats: () => [...mindfulnessKeys.sessions(), 'stats'] as const,
  breathing: () => [...mindfulnessKeys.all, 'breathing'] as const,
  breathingList: (limit: number) =>
    [...mindfulnessKeys.breathing(), { limit }] as const,
  breathingStats: () => [...mindfulnessKeys.breathing(), 'stats'] as const,
  gratitude: () => [...mindfulnessKeys.all, 'gratitude'] as const,
  gratitudeList: (limit: number) =>
    [...mindfulnessKeys.gratitude(), { limit }] as const,
  reflections: () => [...mindfulnessKeys.all, 'reflections'] as const,
  reflectionPrompts: (category?: string) =>
    [...mindfulnessKeys.reflections(), 'prompts', { category }] as const,
  reflectionResponses: (limit: number) =>
    [...mindfulnessKeys.reflections(), 'responses', { limit }] as const,
  streak: () => [...mindfulnessKeys.all, 'streak'] as const,
};

export const useMindfulnessSessions = (limit = 20) => {
  return useQuery({
    queryKey: mindfulnessKeys.sessionsList(limit),
    queryFn: () => mindfulnessService.getMindfulnessSessions(limit),
  });
};

export const useMindfulnessStats = () => {
  return useQuery({
    queryKey: mindfulnessKeys.sessionsStats(),
    queryFn: () => mindfulnessService.getMindfulnessStats(),
  });
};

export const useCreateMindfulnessSession = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMindfulnessSessionRequest) =>
      mindfulnessService.createMindfulnessSession(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.sessions() });
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.streak() });
    },
  });
};

export const useBreathingExercises = (limit = 20) => {
  return useQuery({
    queryKey: mindfulnessKeys.breathingList(limit),
    queryFn: () => mindfulnessService.getBreathingExercises(limit),
  });
};

export const useBreathingStats = () => {
  return useQuery({
    queryKey: mindfulnessKeys.breathingStats(),
    queryFn: () => mindfulnessService.getBreathingStats(),
  });
};

export const useCreateBreathingExercise = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateBreathingExerciseRequest) =>
      mindfulnessService.createBreathingExercise(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.breathing() });
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.streak() });
    },
  });
};

export const useGratitudeEntries = (limit = 20) => {
  return useQuery({
    queryKey: mindfulnessKeys.gratitudeList(limit),
    queryFn: () => mindfulnessService.getGratitudeEntries(limit),
  });
};

export const useCreateGratitudeEntry = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateGratitudeEntryRequest) =>
      mindfulnessService.createGratitudeEntry(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.gratitude() });
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.streak() });
    },
  });
};

export const useDeleteGratitudeEntry = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (entryId: number) =>
      mindfulnessService.deleteGratitudeEntry(entryId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.gratitude() });
    },
  });
};

export const useReflectionPrompts = (category?: string) => {
  return useQuery({
    queryKey: mindfulnessKeys.reflectionPrompts(category),
    queryFn: () => mindfulnessService.getReflectionPrompts(category),
  });
};

export const useReflectionResponses = (limit = 20) => {
  return useQuery({
    queryKey: mindfulnessKeys.reflectionResponses(limit),
    queryFn: () => mindfulnessService.getReflectionResponses(limit),
  });
};

export const useCreateReflectionResponse = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateReflectionResponseRequest) =>
      mindfulnessService.createReflectionResponse(data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: mindfulnessKeys.reflections(),
      });
      queryClient.invalidateQueries({ queryKey: mindfulnessKeys.streak() });
    },
  });
};

export const useMindfulnessStreak = () => {
  return useQuery({
    queryKey: mindfulnessKeys.streak(),
    queryFn: () => mindfulnessService.getMindfulnessStreak(),
  });
};
