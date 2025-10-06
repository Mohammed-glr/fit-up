package handlers

import (
	"github.com/go-chi/chi/v5"
	authRepo "github.com/tdmdh/fit-up-server/internal/auth/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type SchemaRoutes struct {
	authMiddleware        *middleware.AuthMiddleware
	exerciseHandler       *ExerciseHandler
	workoutHandler        *WorkoutHandler
	workoutSessionHandler *WorkoutSessionHandler
	fitnessProfileHandler *FitnessProfileHandler
	planGenerationHandler *PlanGenerationHandler
	coachHandler          *CoachHandler
}

func NewSchemaRoutes(
	schemaRepo repository.SchemaRepo,
	userStore authRepo.UserStore,
	exerciseService service.ExerciseService,
	workoutService service.WorkoutService,
	workoutSessionService service.WorkoutSessionService,
	fitnessProfileService service.FitnessProfileService,
	planGenerationService service.PlanGenerationService,
	coachService service.CoachService,
) *SchemaRoutes {
	return &SchemaRoutes{
		authMiddleware:        middleware.NewAuthMiddleware(schemaRepo, userStore),
		exerciseHandler:       NewExerciseHandler(exerciseService),
		workoutHandler:        NewWorkoutHandler(workoutService),
		workoutSessionHandler: NewWorkoutSessionHandler(workoutSessionService),
		fitnessProfileHandler: NewFitnessProfileHandler(fitnessProfileService),
		planGenerationHandler: NewPlanGenerationHandler(planGenerationService),
		coachHandler:          NewCoachHandler(coachService),
	}
}

func (sr *SchemaRoutes) RegisterRoutes(r chi.Router) {
	r.Route("/exercises", func(r chi.Router) {
		r.Get("/", sr.exerciseHandler.ListExercises)
		r.Get("/{id}", sr.exerciseHandler.GetExerciseByID)
		r.Post("/filter", sr.exerciseHandler.FilterExercises)
		r.Get("/search", sr.exerciseHandler.SearchExercises)
		r.Get("/muscle-group/{muscleGroup}", sr.exerciseHandler.GetExercisesByMuscleGroup)
		r.Get("/equipment/{equipment}", sr.exerciseHandler.GetExercisesByEquipment)
		r.Get("/recommended", sr.exerciseHandler.GetRecommendedExercises)
		r.Get("/most-used", sr.exerciseHandler.GetMostUsedExercises)
		r.Get("/{id}/usage-stats", sr.exerciseHandler.GetExerciseUsageStats)
	})

	r.Group(func(r chi.Router) {
		r.Use(sr.authMiddleware.RequireJWTAuth())

		r.Route("/workouts", func(r chi.Router) {
			r.Get("/{id}", sr.workoutHandler.GetWorkoutByID)
			r.Get("/{id}/exercises", sr.workoutHandler.GetWorkoutWithExercises)
		})

		r.Route("/workout-sessions", func(r chi.Router) {
			r.Post("/start", sr.workoutSessionHandler.StartWorkoutSession)
			r.Post("/{sessionID}/complete", sr.workoutSessionHandler.CompleteWorkoutSession)
			r.Post("/{sessionID}/skip", sr.workoutSessionHandler.SkipWorkout)
			r.Post("/{sessionID}/log-exercise", sr.workoutSessionHandler.LogExercisePerformance)
			r.Get("/users/{userID}/active", sr.workoutSessionHandler.GetActiveSession)
			r.Get("/users/{userID}/history", sr.workoutSessionHandler.GetSessionHistory)
			r.Get("/users/{userID}/metrics", sr.workoutSessionHandler.GetSessionMetrics)
			r.Get("/users/{userID}/weekly-stats", sr.workoutSessionHandler.GetWeeklySessionStats)
		})

		r.Route("/fitness-profile", func(r chi.Router) {
			r.Post("/users/{userID}/assessment", sr.fitnessProfileHandler.CreateFitnessAssessment)
			r.Get("/users/{userID}", sr.fitnessProfileHandler.GetUserFitnessProfile)
			r.Put("/users/{userID}/fitness-level", sr.fitnessProfileHandler.UpdateFitnessLevel)
			r.Put("/users/{userID}/goals", sr.fitnessProfileHandler.UpdateFitnessGoals)
			r.Post("/users/{userID}/1rm-estimate", sr.fitnessProfileHandler.EstimateOneRepMax)
			r.Get("/users/{userID}/1rm-history", sr.fitnessProfileHandler.GetOneRepMaxHistory)
			r.Post("/users/{userID}/movement-assessment", sr.fitnessProfileHandler.CreateMovementAssessment)
			r.Get("/users/{userID}/movement-limitations", sr.fitnessProfileHandler.GetMovementLimitations)
			r.Post("/users/{userID}/workout-profile", sr.fitnessProfileHandler.CreateWorkoutProfile)
			r.Get("/users/{userID}/workout-profile", sr.fitnessProfileHandler.GetWorkoutProfile)
			r.Post("/users/{userID}/fitness-goals", sr.fitnessProfileHandler.CreateFitnessGoal)
			r.Get("/users/{userID}/active-goals", sr.fitnessProfileHandler.GetActiveGoals)
		})

		r.Route("/plans", func(r chi.Router) {
			r.Post("/", sr.planGenerationHandler.CreatePlanGeneration)
			r.Get("/users/{userID}/active", sr.planGenerationHandler.GetActivePlan)
			r.Get("/users/{userID}/history", sr.planGenerationHandler.GetPlanHistory)
			r.Post("/{planID}/performance", sr.planGenerationHandler.TrackPlanPerformance)
			r.Get("/{planID}/download", sr.planGenerationHandler.DownloadPlanPDF)
			r.Post("/{planID}/regenerate", sr.planGenerationHandler.MarkPlanForRegeneration)
		})

		r.Route("/coach", func(r chi.Router) {
			r.Use(sr.authMiddleware.RequireCoachRole())

			r.Get("/dashboard", sr.coachHandler.GetDashboard)
			r.Get("/stats", sr.coachHandler.GetCoachStats)
			r.Get("/activity", sr.coachHandler.GetRecentActivity)

			r.Get("/clients", sr.coachHandler.GetClients)
			r.Post("/clients/assign", sr.coachHandler.AssignClient)
			r.Get("/clients/{userID}", sr.coachHandler.GetClientDetails)
			r.Delete("/clients/{assignmentID}", sr.coachHandler.RemoveClient)

			r.Get("/clients/{userID}/progress", sr.coachHandler.GetClientProgress)
			r.Get("/clients/{userID}/workouts", sr.coachHandler.GetClientWorkouts)
			r.Get("/clients/{userID}/schemas", sr.coachHandler.GetClientSchemas)
			r.Post("/clients/{userID}/notes", sr.coachHandler.AddClientNote)

			r.Post("/clients/{userID}/schemas", sr.coachHandler.CreateSchemaForClient)
			r.Put("/schemas/{schemaID}", sr.coachHandler.UpdateSchema)
			r.Delete("/schemas/{schemaID}", sr.coachHandler.DeleteSchema)
			r.Post("/schemas/{schemaID}/clone", sr.coachHandler.CloneSchema)

			r.Get("/templates", sr.coachHandler.GetTemplates)
			r.Post("/templates", sr.coachHandler.SaveTemplate)
			r.Post("/templates/{templateID}/create-schema", sr.coachHandler.CreateFromTemplate)
			r.Delete("/templates/{templateID}", sr.coachHandler.DeleteTemplate)
		})

		r.Group(func(r chi.Router) {
			r.Use(sr.authMiddleware.RequireAdminRole())

			// Admin-specific routes can be added here
			// For example: user management, system settings, etc.
		})
	})
}
