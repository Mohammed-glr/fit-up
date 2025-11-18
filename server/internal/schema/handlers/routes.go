package handlers

import (
	"fmt"
	"net/http"

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
	schemaHandler         *SchemaHandler
	planGenerationHandler *PlanGenerationHandler
	coachHandler          *CoachHandler
	invitationHandler     *InvitationHandler
	workoutSharingHandler *WorkoutSharingHandler
}

func NewSchemaRoutes(
	schemaRepo repository.SchemaRepo,
	userStore authRepo.UserStore,
	exerciseService service.ExerciseService,
	workoutService service.WorkoutService,
	planGenerationService service.PlanGenerationService,
	coachService service.CoachService,
	invitationService service.InvitationService,
) *SchemaRoutes {
	store, ok := schemaRepo.(*repository.Store)
	if !ok {
		panic("schemaRepo must be *repository.Store")
	}

	return &SchemaRoutes{
		authMiddleware:        middleware.NewAuthMiddleware(schemaRepo, userStore),
		exerciseHandler:       NewExerciseHandler(exerciseService),
		workoutHandler:        NewWorkoutHandler(workoutService),
		schemaHandler:         NewSchemaHandler(schemaRepo),
		planGenerationHandler: NewPlanGenerationHandler(planGenerationService),
		coachHandler:          NewCoachHandler(coachService),
		invitationHandler:     NewInvitationHandler(invitationService),
		workoutSharingHandler: NewWorkoutSharingHandler(store),
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

		r.Route("/schemas", func(r chi.Router) {
			r.Get("/user/{userID}", sr.schemaHandler.GetUserSchemas)
			r.Get("/{schemaID}/workouts", sr.schemaHandler.GetSchemaWithWorkouts)
		})

		r.Route("/workouts", func(r chi.Router) {
			r.Get("/{id}", sr.workoutHandler.GetWorkoutByID)
			r.Get("/{id}/exercises", sr.workoutHandler.GetWorkoutWithExercises)
		})

		r.Route("/plans", func(r chi.Router) {
			r.Post("/", sr.planGenerationHandler.CreatePlanGeneration)
			r.Get("/users/{userID}/active", sr.planGenerationHandler.GetActivePlan)
			r.Get("/users/{userID}/history", sr.planGenerationHandler.GetPlanHistory)
			r.Get("/adaptations/{userID}", sr.planGenerationHandler.GetAdaptationHistory)
			r.Delete("/users/{userID}/{planID}", sr.planGenerationHandler.DeletePlan)
			r.Post("/{planID}/performance", sr.planGenerationHandler.TrackPlanPerformance)
			r.Get("/{planID}/effectiveness", sr.planGenerationHandler.GetPlanEffectiveness)
			r.Get("/{planID}/download", sr.planGenerationHandler.DownloadPlanPDF)
			r.Post("/{planID}/regenerate", sr.planGenerationHandler.MarkPlanForRegeneration)
		})

		// Workout Sharing Routes
		r.Route("/workout-sessions", func(r chi.Router) {
			fmt.Println("✅ Registering workout sharing routes")
			if sr.workoutSharingHandler == nil {
				fmt.Println("❌ ERROR: workoutSharingHandler is nil!")
			} else {
				fmt.Println("✅ workoutSharingHandler is initialized")
			}

			// Test endpoint to verify routing works
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message":"workout-sessions routes are working"}`))
			})

			r.Get("/{sessionId}/share-summary", sr.workoutSharingHandler.HandleGetWorkoutShareSummary)
			r.Post("/share", sr.workoutSharingHandler.HandleShareWorkout)
		})

		r.Route("/coach", func(r chi.Router) {
			r.Use(sr.authMiddleware.RequireCoachRole())

			r.Get("/dashboard", sr.coachHandler.GetDashboard)

			r.Get("/clients", sr.coachHandler.GetClients)
			r.Get("/search-users", sr.coachHandler.SearchUsers)
			r.Post("/clients/assign", sr.coachHandler.AssignClient)
			r.Get("/clients/{userID}", sr.coachHandler.GetClientDetails)
			r.Delete("/clients/{assignmentID}", sr.coachHandler.RemoveClient)

			r.Get("/clients/{userID}/progress", sr.coachHandler.GetClientProgress)
			r.Get("/clients/{userID}/workouts", sr.coachHandler.GetClientWorkouts)
			r.Get("/clients/{userID}/schemas", sr.coachHandler.GetClientSchemas)

			r.Post("/clients/{userID}/schemas", sr.coachHandler.CreateSchemaForClient)
			r.Put("/schemas/{schemaID}", sr.coachHandler.UpdateSchema)
			r.Delete("/schemas/{schemaID}", sr.coachHandler.DeleteSchema)
			r.Post("/schemas/{schemaID}/clone", sr.coachHandler.CloneSchema)

			r.Get("/templates", sr.coachHandler.GetTemplates)
			r.Post("/templates", sr.coachHandler.SaveTemplate)
			r.Post("/templates/{templateID}/create-schema", sr.coachHandler.CreateFromTemplate)
			r.Delete("/templates/{templateID}", sr.coachHandler.DeleteTemplate)

			// Invitation routes
			r.Post("/invitations", sr.invitationHandler.CreateInvitation)
			r.Get("/invitations", sr.invitationHandler.GetInvitations)
			r.Post("/invitations/{id}/resend", sr.invitationHandler.ResendInvitation)
			r.Delete("/invitations/{id}", sr.invitationHandler.CancelInvitation)
		})

		// Public invitation acceptance route (requires auth but not coach role)
		r.Post("/invitations/accept", sr.invitationHandler.AcceptInvitation)

		r.Group(func(r chi.Router) {
			r.Use(sr.authMiddleware.RequireAdminRole())

			// Admin-specific routes can be added here
			// For example: user management, system settings, etc.
		})
	})
}
