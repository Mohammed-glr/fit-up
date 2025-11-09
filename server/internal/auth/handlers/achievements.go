package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tdmdh/fit-up-server/shared/middleware"
)

func (h *AuthHandler) handleGetAchievements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	achievements, err := h.store.GetUserAchievements(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch achievements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(achievements)
}

func (h *AuthHandler) handleGetAchievementStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	stats, err := h.store.GetAchievementStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch achievement stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AuthHandler) handleCheckAchievements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	newlyEarned, err := h.store.CheckAndAwardAchievements(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to check achievements", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"newly_earned": newlyEarned,
		"count":        len(newlyEarned),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
