package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type ExerciseHandler struct {
	service service.ExerciseService
}

func NewExerciseHandler(service service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		service: service,
	}
}

func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	exercise, err := h.service.GetExerciseByID(r.Context(), exerciseID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercise)
}

func (h *ExerciseHandler) ListExercises(w http.ResponseWriter, r *http.Request) {
	pagination := extractPaginationParams(r)

	exercises, err := h.service.ListExercises(r.Context(), pagination)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) FilterExercises(w http.ResponseWriter, r *http.Request) {
	filter := types.ExerciseFilter{}
	pagination := extractPaginationParams(r)

	if muscleGroup := r.URL.Query().Get("muscle_group"); muscleGroup != "" {
		filter.MuscleGroups = append(filter.MuscleGroups, muscleGroup)
	}
	if equipment := r.URL.Query().Get("equipment"); equipment != "" {
		filter.Equipment = append(filter.Equipment, types.EquipmentType(equipment))
	}
	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		level := types.FitnessLevel(difficulty)
		filter.Difficulty = &level
	}

	exercises, err := h.service.FilterExercises(r.Context(), filter, pagination)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) SearchExercises(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	pagination := extractPaginationParams(r)

	exercises, err := h.service.SearchExercises(r.Context(), query, pagination)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetExercisesByMuscleGroup(w http.ResponseWriter, r *http.Request) {
	muscleGroup := chi.URLParam(r, "group")

	exercises, err := h.service.GetExercisesByMuscleGroup(r.Context(), muscleGroup)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetExercisesByEquipment(w http.ResponseWriter, r *http.Request) {
	equipment := types.EquipmentType(chi.URLParam(r, "type"))

	exercises, err := h.service.GetExercisesByEquipment(r.Context(), equipment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetRecommendedExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	count := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 {
			count = c
		}
	}

	exercises, err := h.service.GetRecommendedExercises(r.Context(), userID, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetMostUsedExercises(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	exercises, err := h.service.GetMostUsedExercises(r.Context(), limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetExerciseUsageStats(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	stats, err := h.service.GetExerciseUsageStats(r.Context(), exerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *ExerciseHandler) RegisterRoutes(r chi.Router) {
	r.Route("/exercises", func(r chi.Router) {
		r.Get("/", h.ListExercises)
		r.Get("/filter", h.FilterExercises)
		r.Get("/search", h.SearchExercises)
		r.Get("/popular", h.GetMostUsedExercises)
		r.Get("/{id}", h.GetExerciseByID)
		r.Get("/{id}/stats", h.GetExerciseUsageStats)
		r.Get("/muscle-group/{group}", h.GetExercisesByMuscleGroup)
		r.Get("/equipment/{type}", h.GetExercisesByEquipment)
		r.Get("/recommended/{userID}", h.GetRecommendedExercises)
	})
}
