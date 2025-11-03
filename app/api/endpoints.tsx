
export const API = {
    recipe: {
        system: {
            list: () => ({ method: 'GET', url: 'food-tracker/recipes/system' }),
            retrieve: (id: number) => ({ method: 'GET', url: `food-tracker/recipes/system/${id}` }),
            search: () => ({ method: 'GET', url: 'food-tracker/recipes/search' }),
            create: () => ({ method: 'POST', url: 'food-tracker/recipes/system' }),
            update: (id: number) => ({ method: 'PUT', url: `food-tracker/recipes/system/${id}` }),
            delete: (id: number) => ({ method: 'DELETE', url: `food-tracker/recipes/system/${id}` }),
        },
        user: {
            list: () => ({ method: 'GET', url: 'food-tracker/recipes/user' }),
            create: () => ({ method: 'POST', url: 'food-tracker/recipes/user' }),
            retrieve: (id: number) => ({ method: 'GET', url: `food-tracker/recipes/user/${id}` }),
            update: (id: number) => ({ method: 'PUT', url: `food-tracker/recipes/user/${id}` }),
            delete: (id: number) => ({ method: 'DELETE', url: `food-tracker/recipes/user/${id}` }),
            getFavorites: () => ({ method: 'GET', url: 'food-tracker/recipes/favorites' }),
            toggleFavorite: (recipeId: number) => ({ method: 'PATCH', url: `food-tracker/recipes/favorites/${recipeId}` }),
        },
        logs: {
            log: () => ({ method: 'POST', url: 'food-tracker/food-logs' }),
            logRecipe: () => ({ method: 'POST', url: 'food-tracker/food-logs/recipe' }),
            getlogsByDate: (date: string) => ({ method: 'GET', url: `food-tracker/food-logs/date/${date}` }),
            getLogsInRange: () => ({ method: 'GET', url: 'food-tracker/food-logs/range' }),
            getFoodLogEntry: (id: number) => ({ method: 'GET', url: `food-tracker/food-logs/${id}` }),
            updateFoodLogEntry: (id: number) => ({ method: 'PUT', url: `food-tracker/food-logs/${id}` }),
            deleteFoodLogEntry: (id: number) => ({ method: 'DELETE', url: `food-tracker/food-logs/${id}` }),
        },
        nutrition: {
            getDailySummary: (date: string) => ({ method: 'GET', url: `food-tracker/nutrition/daily/${date}` }),
            getWeeklySummary: () => ({ method: 'GET', url: `food-tracker/nutrition/weekly` }),
            getMonthlySummary: () => ({ method: 'GET', url: `food-tracker/nutrition/monthly` }),
            getGoals: () => ({ method: 'GET', url: `food-tracker/nutrition/goals` }),
            updateGoals: () => ({ method: 'POST', url: `food-tracker/nutrition/goals` }),
            compareWithGoals: (date: string) => ({ method: 'GET', url: `food-tracker/nutrition/comparison/${date}` }),
            getNutritionInsights: (date: string) => ({ method: 'GET', url: `food-tracker/nutrition/insights/${date}` }),
        }
    },
    schema: {
        exercises: {
            list:() => ({ method: 'GET', url: 'exercises/'}) ,
            listById: (id: number) => ({ method: 'GET', url: `exercises/${id}`}),
            filter: () => ({ method: 'POST', url: `exercises/filter`}),
            search: () => ({ method: 'GET', url: `exercises/search`}),
            listByMG: (muscleGroup: string) => ({ method: 'GET', url: `exercises/muscle-group/${muscleGroup}`}),
            getEquipment: (equipment: string) => ({ method: 'GET', url: `exercises/equipment/${equipment}`}),
            getRecommended: () => ({ method: 'GET', url: `exercises/recommended`}),
            getMostUsed: () => ({ method: 'GET', url : `exercises/most-used`}),
            getUsageStats: (id: number) => ({ method: 'GET', url: `exercises/${id}/usage-stats`}),
        },
        workouts: {
            retrieve: (id: number) => ({ method: 'GET', url: `workouts/${id}` }),
            getWorkoutExercises: (id: number) => ({ method: 'GET', url: `workouts/${id}/exercises` }),
        },
        plans: {
            create: () => ({ method: 'POST', url: 'plans/' }),
            getActivePlan: (userID: number) => ({ method: 'GET', url: `plans/users/${userID}/active` }),
            getPlanHistory: (userID: number) => ({ method: 'GET', url: `plans/users/${userID}/history` }),
            trackPerformance: (planID: number) => ({ method: 'POST', url: `plans/${planID}/performance` }),
            downloadPlanPDF: (planID: number) => ({ method: 'GET', url: `plans/${planID}/download` }),
            regeneratePlan: (planID: number) => ({ method: 'POST', url: `plans/${planID}/regenerate` }),
            getPlanEffectiveness: (planID: number) => ({ method: 'GET', url: `plans/${planID}/effectiveness` }),
            getAdaptationHistory: (userID: number) => ({ method: 'GET', url: `plans/adaptations/${userID}` }),
        },
        coach: {
            getDashboard: () => ({ method: 'GET', url: 'coach/dashboard' }),
            
            getClients: () => ({ method: 'GET', url: 'coach/clients' }),
            assignClient: () => ({ method: 'POST', url: 'coach/clients/assign' }),
            getClientDetails: (userID: number) => ({ method: 'GET', url: `coach/clients/${userID}` }),
            removeClient: (assignmentID: number) => ({ method: 'DELETE', url: `coach/clients/${assignmentID}` }),
            
            getClientProgress: (userID: number) => ({ method: 'GET', url: `coach/clients/${userID}/progress` }),
            getClientWorkouts: (userID: number) => ({ method: 'GET', url: `coach/clients/${userID}/workouts` }),
            getClientSchemas: (userID: number) => ({ method: 'GET', url: `coach/clients/${userID}/schemas` }),
            
            createSchemaForClient: (userID: number) => ({ method: 'POST', url: `coach/clients/${userID}/schemas` }),
            updateSchema: (schemaID: number) => ({ method: 'PUT', url: `coach/schemas/${schemaID}` }),
            deleteSchema: (schemaID: number) => ({ method: 'DELETE', url: `coach/schemas/${schemaID}` }),
            cloneSchema: (schemaID: number) => ({ method: 'POST', url: `coach/schemas/${schemaID}/clone` }),
            
            getTemplates: () => ({ method: 'GET', url: 'coach/templates' }),
            saveTemplate: () => ({ method: 'POST', url: 'coach/templates' }),
            createFromTemplate: (templateID: number) => ({ method: 'POST', url: `coach/templates/${templateID}/create-schema` }),
            deleteTemplate: (templateID: number) => ({ method: 'DELETE', url: `coach/templates/${templateID}` }),
        }
    },
    message: {
        conversations: {
            create: () => ({ method: 'POST', url: 'conversations/' }),
            list: () => ({ method: 'GET', url: 'conversations/' }),
            get: (conversation_id: number) => ({ method: 'GET', url: `conversations/${conversation_id}` }),
            getUnreadCount: (conversation_id: number) => ({ method: 'GET', url: `conversations/${conversation_id}/unread-count` }),
            getMessages: (conversation_id: number) => ({ method: 'GET', url: `conversations/${conversation_id}/messages` }),
            markAllAsRead: (conversation_id: number) => ({ method: 'POST', url: `conversations/${conversation_id}/messages/read-all` }),
        },
        messages: {
            send: () => ({ method: 'POST', url: `messages/` }),
            update: (message_id: number) => ({ method: 'PUT', url: `messages/${message_id}` }),
            delete: (message_id: number) => ({ method: 'DELETE', url: `messages/${message_id}` }),
            markAsRead: (message_id: number) => ({ method: 'POST', url: `messages/${message_id}/read` }),
        }
    },
    auth: {
    login: () => ({ method: 'POST', url: 'auth/login' }),
    logout: () => ({ method: 'POST', url: 'auth/logout' }),
    register: () => ({ method: 'POST', url: 'auth/register' }),
    refreshToken: () => ({ method: 'POST', url: 'auth/refresh-token' }),
    validateToken: () => ({ method: 'GET', url: 'auth/validate-token' }),
    forgetPassword: () => ({ method: 'POST', url: 'auth/forget-password' }),
    resetPassword: () => ({ method: 'POST', url: 'auth/reset-password' }),
    changePassword: () => ({ method: 'POST', url: 'auth/change-password' }),
    updateRole: () => ({ method: 'PUT', url: 'auth/update-role' }),
    updateProfile: () => ({ method: 'PUT', url: 'auth/profile' }),
    getProfile: (username: string) => ({ method: 'GET', url: `auth/${username}` }),
    oauthLogin: (provider: string) => ({ method: 'POST', url: `auth/oauth/${provider}` }),
    callbackOAuth: (provider: string) => ({ method: 'POST', url: `auth/oauth/callback/${provider}` }),
    linkOAuth: (provider: string) => ({ method: 'POST', url: `auth/oauth/link/${provider}` }),
    unlinkOAuth: (provider: string) => ({ method: 'DELETE', url: `auth/oauth/unlink/${provider}` }),
    getLinkedProviders: () => ({ method: 'GET', url: 'auth/oauth/linked-providers' }),
    verifyEmail: () => ({ method: 'POST', url: 'auth/verify-email' }),
    resendVerification: () => ({ method: 'POST', url: 'auth/verify-email/resend' }),
    },

}

