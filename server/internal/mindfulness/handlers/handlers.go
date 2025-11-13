package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/mindfulness/repository"
	"github.com/tdmdh/fit-up-server/internal/mindfulness/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type MindfulnessHandler struct {
	repo repository.MindfulnessRepo
}

func NewMindfulnessHandler(repo repository.MindfulnessRepo) *MindfulnessHandler {
	return &MindfulnessHandler{repo: repo}
}

func (h *MindfulnessHandler) CreateMindfulnessSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req types.CreateMindfulnessSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.repo.CreateMindfulnessSession(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, session)
}

func (h *MindfulnessHandler) GetMindfulnessSessions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	sessions, err := h.repo.GetMindfulnessSessions(r.Context(), userID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sessions)
}

func (h *MindfulnessHandler) GetMindfulnessStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	stats, err := h.repo.GetMindfulnessStats(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *MindfulnessHandler) CreateBreathingExercise(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req types.CreateBreathingExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	exercise, err := h.repo.CreateBreathingExercise(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, exercise)
}

func (h *MindfulnessHandler) GetBreathingExercises(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	exercises, err := h.repo.GetBreathingExercises(r.Context(), userID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func (h *MindfulnessHandler) GetBreathingStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	stats, err := h.repo.GetBreathingStats(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *MindfulnessHandler) CreateGratitudeEntry(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req types.CreateGratitudeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	entry, err := h.repo.CreateGratitudeEntry(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, entry)
}

func (h *MindfulnessHandler) GetGratitudeEntries(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	entries, err := h.repo.GetGratitudeEntries(r.Context(), userID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entries)
}

func (h *MindfulnessHandler) DeleteGratitudeEntry(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	entryIDStr := chi.URLParam(r, "entryId")
	entryID, err := strconv.Atoi(entryIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid entry ID")
		return
	}

	if err := h.repo.DeleteGratitudeEntry(r.Context(), userID, entryID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MindfulnessHandler) GetReflectionPrompts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	prompts, err := h.repo.GetReflectionPrompts(r.Context(), categoryPtr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, prompts)
}

func (h *MindfulnessHandler) CreateReflectionResponse(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req types.CreateReflectionResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.repo.CreateReflectionResponse(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *MindfulnessHandler) GetReflectionResponses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	responses, err := h.repo.GetReflectionResponses(r.Context(), userID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, responses)
}

func (h *MindfulnessHandler) GetStreak(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	streak, err := h.repo.GetOrCreateStreak(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, streak)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *MindfulnessHandler) RegisterRoutes(r chi.Router, authMW *middleware.AuthMiddleware) {
	r.Route("/mindfulness", func(r chi.Router) {
		r.Use(authMW.RequireJWTAuth())

		r.Post("/sessions", h.CreateMindfulnessSession)
		r.Get("/sessions", h.GetMindfulnessSessions)
		r.Get("/sessions/stats", h.GetMindfulnessStats)

		r.Post("/breathing", h.CreateBreathingExercise)
		r.Get("/breathing", h.GetBreathingExercises)
		r.Get("/breathing/stats", h.GetBreathingStats)

		r.Post("/gratitude", h.CreateGratitudeEntry)
		r.Get("/gratitude", h.GetGratitudeEntries)
		r.Delete("/gratitude/{entryId}", h.DeleteGratitudeEntry)

		r.Get("/reflections/prompts", h.GetReflectionPrompts)
		r.Post("/reflections/responses", h.CreateReflectionResponse)
		r.Get("/reflections/responses", h.GetReflectionResponses)

		r.Get("/streak", h.GetStreak)
	})
}
