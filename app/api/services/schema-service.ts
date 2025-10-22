import { API } from '../endpoints';
import { executeAPI } from '../client';
import type {
    FitnessLevel,
    FitnessGoal,
    EquipmentType,

    Exercise,
    ExerciseDetail,

    Workout,
    WorkoutExerciseDetail,
    WeeklySchemaWithWorkouts,
    WeeklySchemaExtended,
    ManualSchemaRequest,
    CoachDashboard,
    ClientSummary,
    WorkoutTemplate,

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
    GetDashboard: async (): Promise<CoachDashboard> => {
        const response = await executeAPI(API.schema.coach.getDashboard());
        return response.data as CoachDashboard;
    },
    
    GetClients: async (): Promise<{ clients: ClientSummary[]; total: number }> => {
        const response = await executeAPI(API.schema.coach.getClients());
        return response.data as { clients: ClientSummary[]; total: number };
    },

    GetClientDetails: async (userID: number): Promise<ClientSummary> => {
        const response = await executeAPI(API.schema.coach.getClientDetails(userID));
        return response.data as ClientSummary;
    },

    AssignClient: async (data: { user_id: string; notes?: string }): Promise<{ assignment_id: number; message: string }> => {
        const response = await executeAPI(API.schema.coach.assignClient(), data);
        return response.data as { assignment_id: number; message: string };
    },

    RemoveClient: async (assignmentID: number): Promise<{ message: string }> => {
        const response = await executeAPI(API.schema.coach.removeClient(assignmentID));
        return response.data as { message: string };
    },

    GetClientProgress: async (userID: number): Promise<{
        user_id: number;
        total_workouts: number;
        current_streak: number;
        last_workout?: string | null;
        personal_bests: unknown[];
    }> => {
        const response = await executeAPI(API.schema.coach.getClientProgress(userID));
        return response.data as {
            user_id: number;
            total_workouts: number;
            current_streak: number;
            last_workout?: string | null;
            personal_bests: unknown[];
        };
    },

    GetClientWorkouts: async (userID: number): Promise<{ workouts: Workout[] }> => {
        const response = await executeAPI(API.schema.coach.getClientWorkouts(userID));
        return response.data as { workouts: Workout[] };
    },

    GetClientSchemas: async (userID: number): Promise<{ schemas: WeeklySchemaExtended[] }> => {
        const response = await executeAPI(API.schema.coach.getClientSchemas(userID));
        return response.data as { schemas: WeeklySchemaExtended[] };
    },

    CreateSchemaForClient: async (userID: number, schema: ManualSchemaRequest): Promise<WeeklySchemaExtended> => {
        const response = await executeAPI(API.schema.coach.createSchemaForClient(userID), schema);
        return response.data as WeeklySchemaExtended;
    },

    UpdateSchema: async (schemaID: number, schema: ManualSchemaRequest): Promise<WeeklySchemaExtended> => {
        const response = await executeAPI(API.schema.coach.updateSchema(schemaID), schema);
        return response.data as WeeklySchemaExtended;
    },

    DeleteSchema: async (schemaID: number): Promise<{ message: string }> => {
        const response = await executeAPI(API.schema.coach.deleteSchema(schemaID));
        return response.data as { message: string };
    },

    CloneSchema: async (schemaID: number, targetUserID: number): Promise<WeeklySchemaExtended> => {
        const response = await executeAPI(API.schema.coach.cloneSchema(schemaID), { target_user_id: targetUserID });
        return response.data as WeeklySchemaExtended;
    },

    GetTemplates: async (): Promise<{ templates: WorkoutTemplate[]; total: number }> => {
        const response = await executeAPI(API.schema.coach.getTemplates());
        return response.data as { templates: WorkoutTemplate[]; total: number };
    },

    SaveTemplate: async (schemaID: number, templateName: string, description?: string): Promise<{ message: string; template_id: number }> => {
        const response = await executeAPI(API.schema.coach.saveTemplate(), {
            schema_id: schemaID,
            template_name: templateName,
            description,
        });
        return response.data as { message: string; template_id: number };
    },

    CreateFromTemplate: async (templateID: number, userID: number): Promise<WeeklySchemaExtended> => {
        const response = await executeAPI(API.schema.coach.createFromTemplate(templateID), { user_id: userID });
        return response.data as WeeklySchemaExtended;
    },

    DeleteTemplate: async (templateID: number): Promise<{ message: string }> => {
        const response = await executeAPI(API.schema.coach.deleteTemplate(templateID));
        return response.data as { message: string };
    }, 
}

export { exerciseService, workoutService, planService, coachService };

