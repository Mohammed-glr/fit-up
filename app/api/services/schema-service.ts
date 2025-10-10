import { API } from '../endpoints';
import { executeAPI } from '../client';
import {
    FitnessLevel,
    FitnessGoal,
    EquipmentType,

    Exercise,
    ExerciseDetail,

    Workout,
    WorkoutExerciseDetail,
    WeeklySchemaWithWorkouts,

    PlanGenerationMetadata,
    CreatePlanRequest,
    GeneratedPlan,
    PlanPerformancePayload,

    PaginatedResponse,
} from '@/types/schema';

const exerciseService = {
    List: async (): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.list());
        return response.data as Exercise[];
    },

    ListById: async (id: number): Promise<ExerciseDetail> => {
        const response = await executeAPI(API.schema.exercises.listById(id));
        return response.data as ExerciseDetail;
    },

    Filter: async (data: { muscleGroups?: string[]; equipment?: string[]; exerciseTypes?: string[]; fitnessLevels?: string[]; page?: number; pageSize?: number; }): Promise<PaginatedResponse<Exercise>> => {
        const response = await executeAPI(API.schema.exercises.filter(), data);
        return response.data as PaginatedResponse<Exercise>;
    },

    Search: async (): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.search());
        return response.data as Exercise[];
    },

    ListByMuscleGroup: async (muscleGroup: string): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.listByMG(muscleGroup));
        return response.data as Exercise[];
    },

    GetEquipment: async (equipment: string): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.getEquipment(equipment));
        return response.data as Exercise[];
    },  

    GetRecommended: async (): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.getRecommended());
        return response.data as Exercise[];
    },

    GetMostUsed: async (): Promise<Exercise[]> => {
        const response = await executeAPI(API.schema.exercises.getMostUsed());
        return response.data as Exercise[];
    },
}

const workoutService = {
    Retrieve: async (id: number): Promise<Workout> => {
        const response = await executeAPI(API.schema.workouts.retrieve(id));
        return response.data as Workout;
    },

    GetWorkoutExercises: async (id: number): Promise<WorkoutExerciseDetail[]> => {
        const response = await executeAPI(API.schema.workouts.getWorkoutExercises(id));
        return response.data as WorkoutExerciseDetail[];
    },
}

const planService = {
    Create: async (data: CreatePlanRequest): Promise<GeneratedPlan> => {
        const response = await executeAPI(API.schema.plans.create(), data);
        return response.data as GeneratedPlan;
    },

    GetActivePlan: async (userID: number): Promise<WeeklySchemaWithWorkouts | null> => {
        const response = await executeAPI(API.schema.plans.getActivePlan(userID));
        return response.data as WeeklySchemaWithWorkouts | null;
    },

    GetPlanHistory: async (userID: number): Promise<GeneratedPlan[]> => {
        const response = await executeAPI(API.schema.plans.getPlanHistory(userID));
        return response.data as GeneratedPlan[];
    },

    TrackPerformance: async (planID: number, data: PlanPerformancePayload): Promise<void> => {
        await executeAPI(API.schema.plans.trackPerformance(planID), data);
    },

    DownloadPlanPDF: async (planID: number): Promise<Blob> => {
        const response = await executeAPI(API.schema.plans.downloadPlanPDF(planID), {}, { responseType: 'blob' });
        return response.data as Blob;
    },

    RegeneratePlan: async (userID: number, data: PlanGenerationMetadata): Promise<GeneratedPlan> => {
        const response = await executeAPI(API.schema.plans.create(), { user_id: userID, metadata: data });
        return response.data as GeneratedPlan;
    },
}

const coachService = {
    GetDashboard: async (): Promise<any> => {
        const response = await executeAPI(API.schema.coach.getDashboard());
        return response.data;
    },
    
    GetClients: async (): Promise<any[]> => {
        const response = await executeAPI(API.schema.coach.getClients());
        return response.data as any[];
    },

    GetClientDetails: async (userID: number): Promise<any> => {
        const response = await executeAPI(API.schema.coach.getClientDetails(userID));
        return response.data;
    },

    AssignClient: async (data: { user_id: number; coach_id: number; }): Promise<void> => {
        await executeAPI(API.schema.coach.assignClient(), data);
    },

    RemoveClient: async (userID: number): Promise<void> => {
        await executeAPI(API.schema.coach.removeClient(userID));
    },

    GetClientProgress: async (userID: number): Promise<any> => {
        const response = await executeAPI(API.schema.coach.getClientProgress(userID));
        return response.data;
    },

    GetClientWorkouts: async (userID: number): Promise<any[]> => {
        const response = await executeAPI(API.schema.coach.getClientWorkouts(userID));
        return response.data as any[];
    },

    GetClientSchemas: async (userID: number): Promise<any[]> => {
        const response = await executeAPI(API.schema.coach.getClientSchemas(userID));
        return response.data as any[];
    },

    CreateSchemaForClient: async (userID: number, data: { name: string; description?: string; weekly_frequency: number; focus_areas: string[]; equipment: EquipmentType[]; fitness_level: FitnessLevel; goals: FitnessGoal[]; time_per_workout: number; start_date: string; }): Promise<any> => {
        const response = await executeAPI(API.schema.coach.getClientSchemas(userID), data);
        return response.data;
    },

    DeleteSchema: async (schemaID: number): Promise<void> => {
        await executeAPI(API.schema.coach.updateSchema(schemaID));
    },

    UpdateSchema: async (schemaID: number, data: { name?: string; description?: string; weekly_frequency?: number; focus_areas?: string[]; equipment?: EquipmentType[]; fitness_level?: FitnessLevel; goals?: FitnessGoal[]; time_per_workout?: number; start_date?: string; }): Promise<any> => {
        const response = await executeAPI(API.schema.coach.updateSchema(schemaID), data);
        return response.data;
    },

    CloneSchema: async (schemaID: number, data: { new_start_date: string; }): Promise<any> => {
        const response = await executeAPI(API.schema.coach.updateSchema(schemaID), data);
        return response.data;
    },

    GetTemplate: async (): Promise<any> => {
        const response = await executeAPI(API.schema.coach.getTemplates());
        return response.data;
    },

    SaveTemplate: async (data: { name: string; description?: string; weekly_frequency: number; focus_areas: string[]; equipment: EquipmentType[]; fitness_level: FitnessLevel; goals: FitnessGoal[]; time_per_workout: number; }): Promise<any> => {
        const response = await executeAPI(API.schema.coach.saveTemplate(), data);
        return response.data;
    },

    CreateFromTemplate: async (templateID: number, data: { user_id: number; start_date: string; }): Promise<any> => {
        const response = await executeAPI(API.schema.coach.createFromTemplate(templateID), data);
        return response.data;
    },

    DeleteTemplate: async (templateID: number): Promise<void> => {
        await executeAPI(API.schema.coach.deleteTemplate(templateID));
    }, 
}

