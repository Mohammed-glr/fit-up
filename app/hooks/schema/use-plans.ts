import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { planService } from '@/api/services/schema-service';
import type { CreatePlanRequest, GeneratedPlan, PlanPerformancePayload, WeeklySchemaWithWorkouts } from '@/types/schema';
import { APIError } from '@/api/client';

const basePlanKey = ['plans'] as const;
const resolveUserCacheKey = (userId?: number | null) => (typeof userId === 'number' && userId > 0 ? userId : 'self');

export const planKeys = {
  all: basePlanKey,
  active: (userId?: number | null) => [...basePlanKey, 'active', resolveUserCacheKey(userId)] as const,
  history: (userId?: number | null) => [...basePlanKey, 'history', resolveUserCacheKey(userId)] as const,
  detail: (planId: number) => [...basePlanKey, 'detail', planId] as const,
};

export const useActivePlan = (userID?: number | null) => {
  const resolvedUserID = typeof userID === 'number' && userID > 0 ? userID : 0;

  return useQuery<WeeklySchemaWithWorkouts | null, APIError>({
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
  return useMutation<Blob, APIError, number>({
    mutationFn: (planID: number) => planService.DownloadPlanPDF(planID),
    onSuccess: (pdfBlob, planID) => {
      const url = window.URL.createObjectURL(pdfBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `workout_plan_${planID}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
  });
};


export const useRegeneratePlan = () => {
  const queryClient = useQueryClient();

  return useMutation<GeneratedPlan, APIError, CreatePlanRequest>({
    mutationFn: (data: CreatePlanRequest) => planService.RegeneratePlan(data.user_id, data.metadata),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: planKeys.active(variables.user_id) });
      queryClient.invalidateQueries({ queryKey: planKeys.history(variables.user_id) });
    },
  });
};
