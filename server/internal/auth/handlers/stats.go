package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleGetUserStats(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	stats, err := h.store.GetUserStats(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to fetch user stats"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, stats)
}

func (h *AuthHandler) handleGetTodayWorkout(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	workout, err := h.store.GetTodayWorkout(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to fetch today's workout"))
		return
	}

	if workout == nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "No active workout plan found for today",
			"workout": nil,
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, workout)
}

func (h *AuthHandler) handleWorkoutCompletion(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	var completion types.WorkoutCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&completion); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	// Validate required fields
	// PlanID is optional as it's not used for coach schemas and not currently stored in progress_logs

	if completion.DurationSeconds <= 0 {
		utils.WriteError(w, http.StatusBadRequest, errors.New("duration_seconds must be greater than 0"))
		return
	}

	if len(completion.Exercises) == 0 {
		utils.WriteError(w, http.StatusBadRequest, errors.New("at least one exercise is required"))
		return
	}

	response, err := h.store.SaveWorkoutCompletion(r.Context(), userID, &completion)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to save workout completion"))
		return
	}

	// Check and award achievements after workout completion
	newlyEarned, err := h.store.CheckAndAwardAchievements(r.Context(), userID)
	if err != nil {
		// Log the error but don't fail the request - achievements are not critical
		// In production, you'd want proper logging here
	} else if len(newlyEarned) > 0 {
		// Add newly earned achievements to the response
		response.NewlyEarnedAchievements = newlyEarned
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) handleGetActivityFeed(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	// Default limit
	limit := 10

	// Get limit from query params if provided
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := json.Number(limitStr).Int64(); err == nil && parsedLimit > 0 {
			limit = int(parsedLimit)
		}
	}

	activities, err := h.store.GetActivityFeed(r.Context(), userID, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to fetch activity feed"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"activities": activities,
		"count":      len(activities),
	})
}

func (h *AuthHandler) handleGetWorkoutHistory(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	// Parse query parameters
	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	// Parse date filters
	var startDate, endDate *time.Time

	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}

	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &t
		}
	}

	history, err := h.store.GetWorkoutHistory(r.Context(), userID, startDate, endDate, page, pageSize)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to fetch workout history"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, history)
}

func (h *AuthHandler) handleGetExerciseProgress(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	// Get exercise name from query
	exerciseName := r.URL.Query().Get("exercise")
	if exerciseName == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("exercise name is required"))
		return
	}

	// Parse date filters
	var startDate, endDate *time.Time

	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}

	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &t
		}
	}

	progress, err := h.store.GetExerciseProgress(r.Context(), userID, exerciseName, startDate, endDate)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to fetch exercise progress"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, progress)
}
