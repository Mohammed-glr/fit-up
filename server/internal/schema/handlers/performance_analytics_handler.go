package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
)

type PerformanceAnalyticsHandler struct {
	service service.PerformanceAnalyticsService
}

func NewPerformanceAnalyticsHandler(service service.PerformanceAnalyticsService) *PerformanceAnalyticsHandler {
	return &PerformanceAnalyticsHandler{
		service: service,
	}
}

// CalculateStrengthProgression handles GET /api/v1/analytics/strength/{userID}/exercise/{exerciseID}
func (h *PerformanceAnalyticsHandler) CalculateStrengthProgression(w http.ResponseWriter, r *http.Request) {
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

	timeframe := 30 // Default 30 days
	if timeframeStr := r.URL.Query().Get("timeframe"); timeframeStr != "" {
		if tf, err := strconv.Atoi(timeframeStr); err == nil && tf > 0 {
			timeframe = tf
		}
	}

	progression, err := h.service.CalculateStrengthProgression(r.Context(), userID, exerciseID, timeframe)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, progression)
}

// DetectPerformancePlateau handles GET /api/v1/analytics/plateau/{userID}/exercise/{exerciseID}
func (h *PerformanceAnalyticsHandler) DetectPerformancePlateau(w http.ResponseWriter, r *http.Request) {
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

	plateau, err := h.service.DetectPerformancePlateau(r.Context(), userID, exerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plateau)
}

// PredictGoalAchievement handles GET /api/v1/analytics/goals/{userID}/{goalID}/prediction
func (h *PerformanceAnalyticsHandler) PredictGoalAchievement(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	goalID, err := strconv.Atoi(chi.URLParam(r, "goalID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	prediction, err := h.service.PredictGoalAchievement(r.Context(), userID, goalID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, prediction)
}

// CalculateTrainingVolume handles GET /api/v1/analytics/volume/{userID}
func (h *PerformanceAnalyticsHandler) CalculateTrainingVolume(w http.ResponseWriter, r *http.Request) {
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
		// Default to current week
		weekStart = time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	}

	volume, err := h.service.CalculateTrainingVolume(r.Context(), userID, weekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, volume)
}

// TrackIntensityProgression handles GET /api/v1/analytics/intensity/{userID}/exercise/{exerciseID}
func (h *PerformanceAnalyticsHandler) TrackIntensityProgression(w http.ResponseWriter, r *http.Request) {
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

	intensity, err := h.service.TrackIntensityProgression(r.Context(), userID, exerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, intensity)
}

// GetOptimalTrainingLoad handles GET /api/v1/analytics/optimal-load/{userID}
func (h *PerformanceAnalyticsHandler) GetOptimalTrainingLoad(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	load, err := h.service.GetOptimalTrainingLoad(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, load)
}

// RegisterRoutes registers all performance analytics routes
func (h *PerformanceAnalyticsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/strength/{userID}/exercise/{exerciseID}", h.CalculateStrengthProgression)
		r.Get("/plateau/{userID}/exercise/{exerciseID}", h.DetectPerformancePlateau)
		r.Get("/goals/{userID}/{goalID}/prediction", h.PredictGoalAchievement)
		r.Get("/volume/{userID}", h.CalculateTrainingVolume)
		r.Get("/intensity/{userID}/exercise/{exerciseID}", h.TrackIntensityProgression)
		r.Get("/optimal-load/{userID}", h.GetOptimalTrainingLoad)
	})
}
