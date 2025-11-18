package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type SchemaHandler struct {
	repo repository.SchemaRepo
}

func NewSchemaHandler(repo repository.SchemaRepo) *SchemaHandler {
	return &SchemaHandler{
		repo: repo,
	}
}

// GetUserSchemas returns all schemas for a specific user
func (h *SchemaHandler) GetUserSchemas(w http.ResponseWriter, r *http.Request) {
	authUserID := chi.URLParam(r, "userID")
	if authUserID == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	pagination := types.PaginationParams{
		Page:     1,
		PageSize: 100,
		Limit:    100,
		Offset:   0,
	}

	result, err := h.repo.Schemas().GetWeeklySchemasByUserID(r.Context(), authUserID, pagination)
	if err != nil {
		log.Printf("[GetUserSchemas] Error fetching schemas for user %s: %v", authUserID, err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch schemas")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"schemas": result.Data,
		"total":   result.TotalCount,
	})
}

// GetSchemaWithWorkouts returns a specific schema with all its workouts
func (h *SchemaHandler) GetSchemaWithWorkouts(w http.ResponseWriter, r *http.Request) {
	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	schema, err := h.repo.Schemas().GetWeeklySchemaByID(r.Context(), schemaID)
	if err != nil {
		log.Printf("[GetSchemaWithWorkouts] Error fetching schema %d: %v", schemaID, err)
		respondWithError(w, http.StatusNotFound, "Schema not found")
		return
	}

	workouts, err := h.repo.Workouts().GetWorkoutsBySchemaID(r.Context(), schemaID)
	if err != nil {
		log.Printf("[GetSchemaWithWorkouts] Error fetching workouts for schema %d: %v", schemaID, err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch workouts")
		return
	}

	// Enrich each workout with exercises
	var enrichedWorkouts []types.WorkoutWithExercises
	for _, workout := range workouts {
		workoutWithExercises, err := h.repo.Workouts().GetWorkoutWithExercises(r.Context(), workout.WorkoutID)
		if err != nil {
			log.Printf("[GetSchemaWithWorkouts] Error fetching exercises for workout %d: %v", workout.WorkoutID, err)
			continue
		}
		enrichedWorkouts = append(enrichedWorkouts, *workoutWithExercises)
	}

	response := map[string]interface{}{
		"schema_id":  schema.SchemaID,
		"user_id":    schema.UserID,
		"week_start": schema.WeekStart,
		"active":     schema.Active,
		"workouts":   enrichedWorkouts,
	}

	respondWithJSON(w, http.StatusOK, response)
}
