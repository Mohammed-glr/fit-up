package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/services"
)

type WorkoutHandler struct {
	service service.WorkoutService
}

func NewWorkoutHandler(service service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		service: service,
	}
}

func (h *WorkoutHandler) GetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	workout, err := h.service.GetWorkoutByID(r.Context(), workoutID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) GetWorkoutWithExercises(w http.ResponseWriter, r *http.Request) {
	workoutID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	workout, err := h.service.GetWorkoutWithExercises(r.Context(), workoutID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) GetWorkoutsBySchemaID(w http.ResponseWriter, r *http.Request) {
	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	workouts, err := h.service.GetWorkoutsBySchemaID(r.Context(), schemaID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workouts)
}

func (h *WorkoutHandler) RegisterRoutes(r chi.Router) {
	r.Route("/workouts", func(r chi.Router) {
		r.Get("/{id}", h.GetWorkoutByID)
		r.Get("/{id}/full", h.GetWorkoutWithExercises)
	})

	r.Route("/schemas", func(r chi.Router) {
		r.Get("/{schemaID}/workouts", h.GetWorkoutsBySchemaID)
	})
}
