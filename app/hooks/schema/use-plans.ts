import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Alert, Platform } from 'react-native';
import * as FileSystem from 'expo-file-system';
import * as Sharing from 'expo-sharing';
import { planService } from '@/api/services/schema-service';
import { secureStorage } from '@/api/storage/secure-storage';
import { API } from '@/api/endpoints';
import { API_CONFIG } from '@/api/apiClient';
import type {
  CreatePlanRequest,
  GeneratedPlan,
  PlanAdaptation,
  PlanEffectivenessResponse,
  PlanPerformancePayload,
} from '@/types/schema';
import { APIError } from '@/api/client';

const basePlanKey = ['plans'] as const;
const resolveUserCacheKey = (userId?: number | null) => (typeof userId === 'number' && userId > 0 ? userId : 'self');

export const planKeys = {
  all: basePlanKey,
  active: (userId?: number | null) => [...basePlanKey, 'active', resolveUserCacheKey(userId)] as const,
  history: (userId?: number | null) => [...basePlanKey, 'history', resolveUserCacheKey(userId)] as const,
  detail: (planId: number) => [...basePlanKey, 'detail', planId] as const,
  effectiveness: (planId?: number | null) => [...basePlanKey, 'effectiveness', typeof planId === 'number' && planId > 0 ? planId : 'none'] as const,
  adaptations: (userId?: number | null) => [...basePlanKey, 'adaptations', resolveUserCacheKey(userId)] as const,
};

export const useActivePlan = (userID?: number | null) => {
  const resolvedUserID = typeof userID === 'number' && userID > 0 ? userID : 0;

  return useQuery<GeneratedPlan | null, APIError>({
    queryKey: planKeys.active(userID),
    queryFn: () => planService.GetActivePlan(resolvedUserID),
    enabled: userID !== undefined,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const usePlanHistory = (userID?: number | null) => {
  const resolvedUserID = typeof userID === 'number' && userID > 0 ? userID : 0;

  return useQuery<GeneratedPlan[], APIError>({
    queryKey: planKeys.history(userID),
    queryFn: () => planService.GetPlanHistory(resolvedUserID),
    enabled: userID !== undefined,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};


export const useCreatePlan = () => {
  const queryClient = useQueryClient();

  return useMutation<GeneratedPlan, APIError, CreatePlanRequest>({
    mutationFn: (data: CreatePlanRequest) => planService.Create(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: planKeys.active(variables.user_id) });
      queryClient.invalidateQueries({ queryKey: planKeys.history(variables.user_id) });
    },
  });
};


export const useTrackPlanPerformance = () => {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, { planID: number; data: PlanPerformancePayload }>({
    mutationFn: ({ planID, data }) => planService.TrackPerformance(planID, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: planKeys.detail(variables.planID) });
    },
  });
};


export const useDownloadPlanPDF = () => {
  return useMutation<void, APIError | Error, number>({
    mutationFn: async (planID: number) => {
      if (Platform.OS === 'web') {
        const { buffer, filename, mimeType } = await planService.DownloadPlanPDF(planID);
        const blob = new Blob([buffer], { type: mimeType });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
        return;
      }

      const token = await secureStorage.getToken('access_token');
      const endpoint = API.schema.plans.downloadPlanPDF(planID);
      const downloadUrl = `${API_CONFIG.BASE_URL}${endpoint.url}`;

      const baseDir = FileSystem.documentDirectory ?? FileSystem.cacheDirectory;
      if (!baseDir) {
        throw new Error('Storage directory unavailable');
      }

      const targetDir = `${baseDir}plans`;
      await FileSystem.makeDirectoryAsync(targetDir, { intermediates: true }).catch(() => undefined);

      const targetUri = `${targetDir}/workout_plan_${planID}_${Date.now()}.pdf`;
      const result = await FileSystem.downloadAsync(downloadUrl, targetUri, {
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      });

      if (await Sharing.isAvailableAsync()) {
        await Sharing.shareAsync(result.uri, {
          mimeType: 'application/pdf',
          dialogTitle: 'Workout Plan',
        });
      } else {
        Alert.alert('Download complete', `Saved to ${result.uri}`);
      }
    },
    onError: (error) => {
      let message: string | undefined;
      if (error instanceof Error) {
        message = error.message;
      } else if ((error as APIError)?.message) {
        message = (error as APIError).message;
      }
      Alert.alert('Download failed', message || 'Please try again later.');
    },
  });
};


export const useRegeneratePlan = () => {
  const queryClient = useQueryClient();

  return useMutation<{ message: string }, APIError, { planID: number; reason: string; userID?: number | null }>({
    mutationFn: ({ planID, reason }) => planService.RequestPlanRegeneration(planID, reason),
    onSuccess: (_, variables) => {
      if (typeof variables.userID === 'number') {
        queryClient.invalidateQueries({ queryKey: planKeys.history(variables.userID) });
        queryClient.invalidateQueries({ queryKey: planKeys.adaptations(variables.userID) });
      }
    },
  });
};


export const useDeletePlan = () => {
  const queryClient = useQueryClient();

  return useMutation<{ message: string }, APIError, { userID: number; planID: number }>({
    mutationFn: ({ userID, planID }) => planService.DeletePlan(userID, planID),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: planKeys.active(variables.userID) });
      queryClient.invalidateQueries({ queryKey: planKeys.history(variables.userID) });
    },
  });
};


export const usePlanEffectiveness = (planID?: number | null) => {
  return useQuery<PlanEffectivenessResponse, APIError>({
    queryKey: planKeys.effectiveness(planID ?? null),
    queryFn: () => planService.GetPlanEffectiveness(planID as number),
    enabled: typeof planID === 'number' && planID > 0,
    staleTime: 2 * 60 * 1000,
  });
};


export const useAdaptationHistory = (userID?: number | null) => {
  const resolvedUserID = typeof userID === 'number' && userID > 0 ? userID : 0;

  return useQuery<PlanAdaptation[], APIError>({
    queryKey: planKeys.adaptations(userID),
    queryFn: () => planService.GetAdaptationHistory(resolvedUserID),
    enabled: userID !== undefined,
    staleTime: 5 * 60 * 1000,
  });
};
