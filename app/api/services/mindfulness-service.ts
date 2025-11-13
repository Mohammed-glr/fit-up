import { httpClient } from '../client';
import type {
  MindfulnessSession,
  CreateMindfulnessSessionRequest,
  MindfulnessStats,
  BreathingExercise,
  CreateBreathingExerciseRequest,
  BreathingStats,
  GratitudeEntry,
  CreateGratitudeEntryRequest,
  ReflectionPrompt,
  ReflectionResponse,
  CreateReflectionResponseRequest,
  MindfulnessStreak,
} from '../../types/mindfulness';

const BASE_PATH = '/mindfulness';

export const createMindfulnessSession = async (
  data: CreateMindfulnessSessionRequest
): Promise<MindfulnessSession> => {
  const response = await httpClient.post(`${BASE_PATH}/sessions`, data);
  return response.data;
};

export const getMindfulnessSessions = async (
  limit = 20
): Promise<MindfulnessSession[]> => {
  const response = await httpClient.get(`${BASE_PATH}/sessions`, {
    params: { limit },
  });
  return response.data;
};

export const getMindfulnessStats = async (): Promise<MindfulnessStats> => {
  const response = await httpClient.get(`${BASE_PATH}/sessions/stats`);
  return response.data;
};

export const createBreathingExercise = async (
  data: CreateBreathingExerciseRequest
): Promise<BreathingExercise> => {
  const response = await httpClient.post(`${BASE_PATH}/breathing`, data);
  return response.data;
};

export const getBreathingExercises = async (
  limit = 20
): Promise<BreathingExercise[]> => {
  const response = await httpClient.get(`${BASE_PATH}/breathing`, {
    params: { limit },
  });
  return response.data;
};

export const getBreathingStats = async (): Promise<BreathingStats> => {
  const response = await httpClient.get(`${BASE_PATH}/breathing/stats`);
  return response.data;
};

export const createGratitudeEntry = async (
  data: CreateGratitudeEntryRequest
): Promise<GratitudeEntry> => {
  const response = await httpClient.post(`${BASE_PATH}/gratitude`, data);
  return response.data;
};

export const getGratitudeEntries = async (
  limit = 20
): Promise<GratitudeEntry[]> => {
  const response = await httpClient.get(`${BASE_PATH}/gratitude`, {
    params: { limit },
  });
  return response.data;
};

export const deleteGratitudeEntry = async (entryId: number): Promise<void> => {
  await httpClient.delete(`${BASE_PATH}/gratitude/${entryId}`);
};

export const getReflectionPrompts = async (
  category?: string
): Promise<ReflectionPrompt[]> => {
  const response = await httpClient.get(`${BASE_PATH}/reflections/prompts`, {
    params: category ? { category } : {},
  });
  return response.data;
};

export const createReflectionResponse = async (
  data: CreateReflectionResponseRequest
): Promise<ReflectionResponse> => {
  const response = await httpClient.post(
    `${BASE_PATH}/reflections/responses`,
    data
  );
  return response.data;
};

export const getReflectionResponses = async (
  limit = 20
): Promise<ReflectionResponse[]> => {
  const response = await httpClient.get(`${BASE_PATH}/reflections/responses`, {
    params: { limit },
  });
  return response.data;
};

export const getMindfulnessStreak = async (): Promise<MindfulnessStreak> => {
  const response = await httpClient.get(`${BASE_PATH}/streak`);
  return response.data;
};
