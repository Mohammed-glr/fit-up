package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

func (h *FoodTrackerHandler) GetDailyNutrition(w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.service.FoodLogs().GetDailyNutrition(r.Context(), userID, date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

func (h *FoodTrackerHandler) GetWeeklyNutrition(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	startDate := chi.URLParam(r, "startDate")
	if startDate == "" {
		respondWithError(w, http.StatusBadRequest, "Start date is required")
		return
	}

	summaries, err := h.service.FoodLogs().GetWeeklyNutrition(r.Context(), userID, startDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"start_date": startDate,
		"summaries":  summaries,
	})
}

func (h *FoodTrackerHandler) GetMonthlyNutrition(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid year")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid month")
		return
	}

	summaries, err := h.service.FoodLogs().GetMonthlyNutrition(r.Context(), userID, year, month)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"year":      year,
		"month":     month,
		"summaries": summaries,
	})
}

func (h *FoodTrackerHandler) GetNutritionGoals(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	goals, err := h.service.Nutrition().GetNutritionGoals(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goals)
}

func (h *FoodTrackerHandler) GetNutritionComparison(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.service.FoodLogs().GetDailyNutrition(r.Context(), userID, date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	goals, err := h.service.Nutrition().GetNutritionGoals(ctx,userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	comparison := h.service.Nutrition().CompareToGoals(summary, goals)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"date":       date,
		"summary":    summary,
		"goals":      goals,
		"comparison": comparison,
	})
}

func (h *FoodTrackerHandler) GetNutritionInsights(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.service.FoodLogs().GetDailyNutrition(r.Context(), userID, date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	goals, err := h.service.Nutrition().GetNutritionGoals(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	macroDistribution := h.service.Nutrition().GetMacroDistribution(
		summary.TotalCalories,
		summary.TotalProtein,
		summary.TotalCarbs,
		summary.TotalFat,
	)

	comparison := h.service.Nutrition().CompareToGoals(summary, goals)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"date":               date,
		"summary":            summary,
		"goals":              goals,
		"comparison":         comparison,
		"macro_distribution": macroDistribution,
	})
}

func (h *FoodTrackerHandler) CreateOrUpdateNutritionGoals(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var goals types.NutritionGoals
	if err := json.NewDecoder(r.Body).Decode(&goals); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	goals.UserID = userID

	if err := h.service.Nutrition().ValidateNutritionGoals(&goals); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Nutrition().SetNutritionGoals(r.Context(), &goals); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goals)
}