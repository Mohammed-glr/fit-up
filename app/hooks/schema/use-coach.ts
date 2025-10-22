import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { coachService } from '@/api/services/schema-service';
import type { ManualSchemaRequest, ClientSummary, WeeklySchemaExtended, CoachDashboard, WorkoutTemplate } from '@/types/schema';
import { APIError } from '@/api/client';

export const coachKeys = {
  all: ['coach'] as const,
  dashboard: () => [...coachKeys.all, 'dashboard'] as const,
  clients: () => [...coachKeys.all, 'clients'] as const,
  client: (id: number) => [...coachKeys.all, 'client', id] as const,
  clientProgress: (id: number) => [...coachKeys.client(id), 'progress'] as const,
  clientWorkouts: (id: number) => [...coachKeys.client(id), 'workouts'] as const,
  clientSchemas: (id: number) => [...coachKeys.client(id), 'schemas'] as const,
  templates: () => [...coachKeys.all, 'templates'] as const,
};

export const useCoachDashboard = () => {
  return useQuery<CoachDashboard, APIError>({
    queryKey: coachKeys.dashboard(),
    queryFn: () => coachService.GetDashboard(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCoachClients = () => {
  return useQuery<{ clients: ClientSummary[]; total: number }, APIError>({
    queryKey: coachKeys.clients(),
    queryFn: () => coachService.GetClients(),
    staleTime: 5 * 60 * 1000,
  });
};

export const useClientDetails = (userID: number) => {
  return useQuery<ClientSummary, APIError>({
    queryKey: coachKeys.client(userID),
    queryFn: () => coachService.GetClientDetails(userID),
    enabled: !!userID,
    staleTime: 5 * 60 * 1000,
  });
};

export const useClientProgress = (userID: number) => {
  return useQuery({
    queryKey: coachKeys.clientProgress(userID),
    queryFn: () => coachService.GetClientProgress(userID),
    enabled: !!userID,
    staleTime: 5 * 60 * 1000,
  });
};

export const useClientWorkouts = (userID: number) => {
  return useQuery({
    queryKey: coachKeys.clientWorkouts(userID),
    queryFn: () => coachService.GetClientWorkouts(userID),
    enabled: !!userID,
    staleTime: 5 * 60 * 1000,
  });
};

export const useClientSchemas = (userID: number) => {
  return useQuery<{ schemas: WeeklySchemaExtended[] }, APIError>({
    queryKey: coachKeys.clientSchemas(userID),
    queryFn: () => coachService.GetClientSchemas(userID),
    enabled: !!userID,
    staleTime: 5 * 60 * 1000,
  });
};

export const useAssignClient = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { user_id: string; notes?: string }) => coachService.AssignClient(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clients() });
      queryClient.invalidateQueries({ queryKey: coachKeys.dashboard() });
    },
  });
};

export const useRemoveClient = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (assignmentID: number) => coachService.RemoveClient(assignmentID),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clients() });
      queryClient.invalidateQueries({ queryKey: coachKeys.dashboard() });
    },
  });
};

export const useCreateSchema = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userID, schema }: { userID: number; schema: ManualSchemaRequest }) =>
      coachService.CreateSchemaForClient(userID, schema),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clientSchemas(variables.userID) });
      queryClient.invalidateQueries({ queryKey: coachKeys.dashboard() });
    },
  });
};

export const useUpdateSchema = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schemaID, schema }: { schemaID: number; schema: ManualSchemaRequest }) =>
      coachService.UpdateSchema(schemaID, schema),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clientSchemas(data.user_id) });
    },
  });
};

export const useDeleteSchema = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (schemaID: number) => coachService.DeleteSchema(schemaID),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: coachKeys.all });
    },
  });
};

export const useCloneSchema = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schemaID, targetUserID }: { schemaID: number; targetUserID: number }) =>
      coachService.CloneSchema(schemaID, targetUserID),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clientSchemas(variables.targetUserID) });
    },
  });
};

export const useCoachTemplates = () => {
  return useQuery<{ templates: WorkoutTemplate[]; total: number }, APIError>({
    queryKey: coachKeys.templates(),
    queryFn: () => coachService.GetTemplates(),
    staleTime: 10 * 60 * 1000,
  });
};

export const useSaveTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schemaID, templateName, description }: { schemaID: number; templateName: string; description?: string }) =>
      coachService.SaveTemplate(schemaID, templateName, description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: coachKeys.templates() });
    },
  });
};

export const useCreateFromTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ templateID, userID }: { templateID: number; userID: number }) =>
      coachService.CreateFromTemplate(templateID, userID),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: coachKeys.clientSchemas(variables.userID) });
    },
  });
};

export const useDeleteTemplate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (templateID: number) => coachService.DeleteTemplate(templateID),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: coachKeys.templates() });
    },
  });
};
