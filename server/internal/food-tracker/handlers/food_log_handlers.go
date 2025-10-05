package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (h *FoodTrackerHandler) LogFood(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req types.CreateFoodLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	logEntry, err := h.service.FoodLogs().LogFood(r.Context(), userID, &req)
	if err != nil {
		if err == types.ErrInvalidRequest || err == types.ErrInvalidMealType || err == types.ErrNurtritionValues {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, logEntry)
}

func (h *FoodTrackerHandler) LogRecipe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		RecipeID       int            `json:"recipe_id"`
		IsSystemRecipe bool           `json:"is_system_recipe"`
		Date           string         `json:"date"`
		MealType       types.MealType `json:"meal_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	logEntry, err := h.service.FoodLogs().LogRecipe(r.Context(), userID, req.RecipeID, req.IsSystemRecipe, req.Date, req.MealType)
	if err != nil {
		if err == types.ErrInvalidID || err == types.ErrInvalidRequest {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, logEntry)
}

func (h *FoodTrackerHandler) GetFoodLogEntry(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Note: The repository method GetFoodLogEntryByID is not exposed in the service interface
	// This would need to be added to the FoodLogService interface if needed
	// For now, we'll return an error indicating this endpoint is not yet implemented
	respondWithError(w, http.StatusNotImplemented, "Get single log entry not yet implemented")
}

func (h *FoodTrackerHandler) UpdateFoodLog(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid log entry ID")
		return
	}

	var req types.CreateFoodLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	logEntry, err := h.service.FoodLogs().UpdateLog(r.Context(), id, userID, &req)
	if err != nil {
		if err == types.ErrInvalidRequest || err == types.ErrInvalidMealType || err == types.ErrNurtritionValues {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Log entry not found")
			return
		}
		if err == types.ErrUnauthorized {
			respondWithError(w, http.StatusForbidden, "You don't have permission to update this log entry")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logEntry)
}

func (h *FoodTrackerHandler) DeleteFoodLog(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid log entry ID")
		return
	}

	if err := h.service.FoodLogs().DeleteLog(r.Context(), id, userID); err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Log entry not found")
			return
		}
		if err == types.ErrUnauthorized {
			respondWithError(w, http.StatusForbidden, "You don't have permission to delete this log entry")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Log entry deleted successfully"})
}

func (h *FoodTrackerHandler) GetLogsByDate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	date := chi.URLParam(r, "date")
	if date == "" {
		respondWithError(w, http.StatusBadRequest, "Date is required")
		return
	}

	logs, err := h.service.FoodLogs().GetLogsByDate(r.Context(), userID, date)
	if err != nil {
		if err == types.ErrInvalidRequest {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"date": date,
		"logs": logs,
	})
}

func (h *FoodTrackerHandler) GetLogsByDateRange(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	startDate := chi.URLParam(r, "startDate")
	endDate := chi.URLParam(r, "endDate")

	if startDate == "" || endDate == "" {
		respondWithError(w, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	logs, err := h.service.FoodLogs().GetLogsByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		if err == types.ErrInvalidRequest {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
		"logs":       logs,
	})
}
