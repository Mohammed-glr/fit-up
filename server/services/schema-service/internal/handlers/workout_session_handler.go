package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/service"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type WorkoutSessionHandler struct {
	service service.WorkoutSessionService
}

func NewWorkoutSessionHandler(service service.WorkoutSessionService) *WorkoutSessionHandler {
	return &WorkoutSessionHandler{
		service: service,
	}
}

func (h *WorkoutSessionHandler) StartWorkoutSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int `json:"user_id"`
		WorkoutID int `json:"workout_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.service.StartWorkoutSession(r.Context(), req.UserID, req.WorkoutID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, session)
}

func (h *WorkoutSessionHandler) CompleteWorkoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "sessionID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	var summary types.SessionSummary
	if err := json.NewDecoder(r.Body).Decode(&summary); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.service.CompleteWorkoutSession(r.Context(), sessionID, &summary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

// SkipWorkout handles POST /api/v1/sessions/skip
func (h *WorkoutSessionHandler) SkipWorkout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int    `json:"user_id"`
		WorkoutID int    `json:"workout_id"`
		Reason    string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	skipped, err := h.service.SkipWorkout(r.Context(), req.UserID, req.WorkoutID, req.Reason)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, skipped)
}

func (h *WorkoutSessionHandler) LogExercisePerformance(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "sessionID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	exerciseID, err := strconv.Atoi(chi.URLParam(r, "exerciseID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	var performance types.ExercisePerformance
	if err := json.NewDecoder(r.Body).Decode(&performance); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.LogExercisePerformance(r.Context(), sessionID, exerciseID, &performance); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Performance logged successfully"})
}

func (h *WorkoutSessionHandler) GetActiveSession(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	session, err := h.service.GetActiveSession(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

func (h *WorkoutSessionHandler) GetSessionHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	pagination := extractPaginationParams(r)

	sessions, err := h.service.GetSessionHistory(r.Context(), userID, pagination)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sessions)
}

func (h *WorkoutSessionHandler) GetSessionMetrics(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "sessionID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	metrics, err := h.service.GetSessionMetrics(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

func (h *WorkoutSessionHandler) GetWeeklySessionStats(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	weekStartStr := r.URL.Query().Get("week_start")
	var weekStart time.Time
	if weekStartStr != "" {
		parsed, err := time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid week_start format (use YYYY-MM-DD)")
			return
		}
		weekStart = parsed
	} else {
		weekStart = time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	}

	stats, err := h.service.GetWeeklySessionStats(r.Context(), userID, weekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *WorkoutSessionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/sessions", func(r chi.Router) {
		r.Post("/start", h.StartWorkoutSession)
		r.Post("/skip", h.SkipWorkout)
		r.Get("/active/{userID}", h.GetActiveSession)
		r.Get("/history/{userID}", h.GetSessionHistory)
		r.Get("/stats/{userID}", h.GetWeeklySessionStats)
		r.Post("/{sessionID}/complete", h.CompleteWorkoutSession)
		r.Get("/{sessionID}/metrics", h.GetSessionMetrics)
		r.Post("/{sessionID}/exercises/{exerciseID}/log", h.LogExercisePerformance)
	})
}
