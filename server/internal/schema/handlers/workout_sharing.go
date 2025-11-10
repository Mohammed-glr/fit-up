package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type WorkoutSharingHandler struct {
	store *repository.Store
}

func NewWorkoutSharingHandler(store *repository.Store) *WorkoutSharingHandler {
	return &WorkoutSharingHandler{
		store: store,
	}
}

// handleGetWorkoutShareSummary returns a complete workout summary for sharing
func (h *WorkoutSharingHandler) handleGetWorkoutShareSummary(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		respondWithError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	sessionIDStr := chi.URLParam(r, "sessionId")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid session ID")
		return
	}

	summary, err := h.store.WorkoutSharing().GetWorkoutShareSummary(r.Context(), sessionID, userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

// handleShareWorkout processes workout sharing requests
func (h *WorkoutSharingHandler) handleShareWorkout(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		respondWithError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	var req types.ShareWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get workout summary
	summary, err := h.store.WorkoutSharing().GetWorkoutShareSummary(r.Context(), req.SessionID, userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Generate share text based on type
	var response types.ShareWorkoutResponse
	response.Success = true

	switch req.ShareType {
	case "text":
		response.ShareText = generateShareText(summary)
		response.Message = "Workout summary copied to clipboard"

	case "image":
		// Image generation would be handled on frontend
		response.Message = "Workout summary ready for image export"

	case "coach":
		// Sharing with coach would create a message
		response.Message = "Workout shared with your coach"

	case "social":
		response.ShareText = generateShareText(summary)
		response.Message = "Workout ready to share"

	default:
		respondWithError(w, http.StatusBadRequest, "invalid share type")
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// generateShareText creates a formatted text summary of the workout
func generateShareText(summary *types.WorkoutShareSummary) string {
	text := fmt.Sprintf("💪 %s\n\n", summary.WorkoutTitle)
	text += fmt.Sprintf("⏱️ Duration: %d minutes\n", summary.DurationMinutes)
	text += fmt.Sprintf("🏋️ Exercises: %d\n", summary.TotalExercises)
	text += fmt.Sprintf("📊 Total Sets: %d\n", summary.TotalSets)
	text += fmt.Sprintf("🔢 Total Reps: %d\n", summary.TotalReps)
	text += fmt.Sprintf("💯 Total Volume: %.1f lbs\n\n", summary.TotalVolumeLbs)

	if summary.PRsAchieved > 0 {
		text += fmt.Sprintf("🎉 %d PRs Achieved!\n\n", summary.PRsAchieved)
	}

	text += "Exercises:\n"
	for _, ex := range summary.Exercises {
		text += fmt.Sprintf("• %s: %d sets", ex.ExerciseName, ex.SetsCompleted)
		if ex.BestSet != nil {
			text += fmt.Sprintf(" (best: %.1f lbs × %d reps)", ex.BestSet.Weight, ex.BestSet.Reps)
		}
		if ex.PRAchieved {
			text += " 🔥 PR!"
		}
		text += "\n"
	}

	text += "\n#FitUp #WorkoutComplete"
	return text
}

// RegisterWorkoutSharingRoutes registers all workout sharing routes
func (h *WorkoutSharingHandler) RegisterRoutes(r chi.Router, authMW *middleware.AuthMiddleware) {
	r.Route("/workout-sessions", func(r chi.Router) {
		r.Use(authMW.RequireJWTAuth())

		r.Get("/{sessionId}/share-summary", h.handleGetWorkoutShareSummary)
		r.Post("/share", h.handleShareWorkout)
	})
}
