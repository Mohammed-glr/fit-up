package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type FitnessProfileHandler struct {
	service service.FitnessProfileService
}

func NewFitnessProfileHandler(service service.FitnessProfileService) *FitnessProfileHandler {
	return &FitnessProfileHandler{
		service: service,
	}
}

func (h *FitnessProfileHandler) CreateFitnessAssessment(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var assessment types.FitnessAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.CreateFitnessAssessment(r.Context(), userID, &assessment)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (h *FitnessProfileHandler) GetUserFitnessProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	profile, err := h.service.GetUserFitnessProfile(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

func (h *FitnessProfileHandler) UpdateFitnessLevel(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Level types.FitnessLevel `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateFitnessLevel(r.Context(), userID, req.Level); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Fitness level updated successfully"})
}

func (h *FitnessProfileHandler) UpdateFitnessGoals(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var goals []types.FitnessGoalTarget
	if err := json.NewDecoder(r.Body).Decode(&goals); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateFitnessGoals(r.Context(), userID, goals); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Fitness goals updated successfully"})
}

func (h *FitnessProfileHandler) EstimateOneRepMax(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	exerciseID, err := strconv.Atoi(chi.URLParam(r, "exerciseID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	var performance types.PerformanceData
	if err := json.NewDecoder(r.Body).Decode(&performance); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	estimate, err := h.service.EstimateOneRepMax(r.Context(), userID, exerciseID, &performance)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, estimate)
}

func (h *FitnessProfileHandler) GetOneRepMaxHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	exerciseID, err := strconv.Atoi(chi.URLParam(r, "exerciseID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	history, err := h.service.GetOneRepMaxHistory(r.Context(), userID, exerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, history)
}

func (h *FitnessProfileHandler) CreateMovementAssessment(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var assessment types.MovementAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.CreateMovementAssessment(r.Context(), userID, &assessment)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (h *FitnessProfileHandler) GetMovementLimitations(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limitations, err := h.service.GetMovementLimitations(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, limitations)
}

func (h *FitnessProfileHandler) CreateWorkoutProfile(w http.ResponseWriter, r *http.Request) {
	authUserID := r.Header.Get("X-User-ID")
	if authUserID == "" {
		respondWithError(w, http.StatusUnauthorized, "User ID not found in request")
		return
	}

	var profile types.WorkoutProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.CreateWorkoutProfile(r.Context(), authUserID, &profile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (h *FitnessProfileHandler) GetWorkoutProfile(w http.ResponseWriter, r *http.Request) {
	authUserID := r.Header.Get("X-User-ID")
	if authUserID == "" {
		respondWithError(w, http.StatusUnauthorized, "User ID not found in request")
		return
	}

	profile, err := h.service.GetWorkoutProfileByAuthID(r.Context(), authUserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

func (h *FitnessProfileHandler) CreateFitnessGoal(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var goal types.FitnessGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.CreateFitnessGoal(r.Context(), userID, &goal)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (h *FitnessProfileHandler) GetActiveGoals(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	goals, err := h.service.GetActiveGoals(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goals)
}

func (h *FitnessProfileHandler) RegisterRoutes(r chi.Router) {
	r.Route("/fitness-profiles", func(r chi.Router) {
		r.Get("/{userID}", h.GetUserFitnessProfile)
		r.Post("/{userID}/assessment", h.CreateFitnessAssessment)
		r.Put("/{userID}/level", h.UpdateFitnessLevel)
		r.Put("/{userID}/goals", h.UpdateFitnessGoals)
		r.Post("/{userID}/goals", h.CreateFitnessGoal)
		r.Get("/{userID}/goals/active", h.GetActiveGoals)
		r.Post("/{userID}/1rm/exercise/{exerciseID}", h.EstimateOneRepMax)
		r.Get("/{userID}/1rm/exercise/{exerciseID}/history", h.GetOneRepMaxHistory)
		r.Post("/{userID}/movement-assessment", h.CreateMovementAssessment)
		r.Get("/{userID}/limitations", h.GetMovementLimitations)
	})

	r.Route("/workout-profiles", func(r chi.Router) {
		r.Post("/", h.CreateWorkoutProfile)
		r.Get("/", h.GetWorkoutProfile)
	})
}
