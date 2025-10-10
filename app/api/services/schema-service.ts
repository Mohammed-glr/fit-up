import { API } from '../endpoints';
import { executeAPI } from '../client';

const schemaService = {
    ListExercises:  () => executeAPI(API.schema.exercises.list()),
    GetExerciseById: (id: number) => executeAPI(API.schema.exercises.listById(id)),
    FilterExercises: (data: any) => executeAPI(API.schema.exercises.filter(), data),
    SearchExercises: () => executeAPI(API.schema.exercises.search()),
    ListExercisesByMuscleGroup: (muscleGroup: string) => executeAPI(API.schema.exercises.listByMG(muscleGroup)),
    GetExercisesByEquipment: (equipment: string) => executeAPI(API.schema.exercises.getEquipment(equipment)),
    GetRecommendedExercises: () => executeAPI(API.schema.exercises.getRecommended()),
    GetMostUsedExercises: () => executeAPI(API.schema.exercises.getMostUsed()),
    GetExerciseUsageStats: (id: number) => executeAPI(API.schema.exercises.getUsageStats(id)),

    RetrieveWorkout: (id: number) => executeAPI(API.schema.workouts.retrieve(id)),
    GetWorkoutExercises: (id: number) => executeAPI(API.schema.workouts.getWorkoutExercises(id)),

    CreatePlan: (data: any) => executeAPI(API.schema.plans.create(), data),
    GetActivePlan: (userID: number) => executeAPI(API.schema.plans.getActivePlan(userID)),
    GetPlanHistory: (userID: number) => executeAPI(API.schema.plans.getPlanHistory(userID)),
    TrackPlanPerformance: (planID: number, data: any) => executeAPI(API.schema.plans.trackPerformance(planID), data),
    DownloadPlanPDF: (planID: number) => executeAPI(API.schema.plans.downloadPlanPDF(planID)),
    RegeneratePlan: (planID: number) => executeAPI(API.schema.plans.regeneratePlan(planID)),

    GetCoachDashboard: () => executeAPI(API.schema.coach.getDashboard()),
    GetCoachClients: () => executeAPI(API.schema.coach.getClients()),
    AssignClientToCoach: () => executeAPI(API.schema.coach.assignClient()),
    GetClientDetails: (userID: number) => executeAPI(API.schema.coach.getClientDetails(userID)),
    RemoveClientFromCoach: (assignmentID: number) => executeAPI(API.schema.coach.removeClient(assignmentID)),
    GetClientProgress: (userID: number) => executeAPI(API.schema.coach.getClientProgress(userID)),
    GetClientWorkouts: (userID: number) => executeAPI(API.schema.coach.getClientWorkouts(userID)),
    GetClientSchemas: (userID: number) => executeAPI(API.schema.coach.getClientSchemas(userID)),
    CreateSchemaForClient: (userID: number) => executeAPI(API.schema.coach.createSchemaForClient(userID)),
    UpdateSchema: (schemaID: number, data: any) => executeAPI(API.schema.coach.updateSchema(schemaID), data),
    DeleteSchema: (schemaID: number) => executeAPI(API.schema.coach.deleteSchema(schemaID)),
    CloneSchema: (schemaID: number) => executeAPI(API.schema.coach.cloneSchema(schemaID)),
    GetTemplates: () => executeAPI(API.schema.coach.getTemplates()),
    SaveTemplate: (data: any) => executeAPI(API.schema.coach.saveTemplate(), data),
    CreateFromTemplate: (templateID: number) => executeAPI(API.schema.coach.createFromTemplate(templateID)),
    DeleteTemplate: (templateID: number) => executeAPI(API.schema.coach.deleteTemplate(templateID)),
};

export default schemaService;