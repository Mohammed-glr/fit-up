package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type PlanGenerationHandler struct {
	service service.PlanGenerationService
}

func NewPlanGenerationHandler(service service.PlanGenerationService) *PlanGenerationHandler {
	return &PlanGenerationHandler{
		service: service,
	}
}

func (h *PlanGenerationHandler) CreatePlanGeneration(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int                          `json:"user_id"`
		Metadata types.PlanGenerationMetadata `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	plan, err := h.service.CreatePlanGeneration(r.Context(), req.UserID, &req.Metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, plan)
}

func (h *PlanGenerationHandler) GetActivePlan(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := h.service.GetActivePlanForUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plan)
}

func (h *PlanGenerationHandler) GetPlanHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	plans, err := h.service.GetPlanGenerationHistory(r.Context(), userID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plans)
}

func (h *PlanGenerationHandler) TrackPlanPerformance(w http.ResponseWriter, r *http.Request) {
	planID, err := strconv.Atoi(chi.URLParam(r, "planID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var performance types.PlanPerformanceData
	if err := json.NewDecoder(r.Body).Decode(&performance); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.TrackPlanPerformance(r.Context(), planID, &performance); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Performance tracked successfully"})
}

func (h *PlanGenerationHandler) GetPlanEffectiveness(w http.ResponseWriter, r *http.Request) {
	planID, err := strconv.Atoi(chi.URLParam(r, "planID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	score, err := h.service.GetPlanEffectivenessScore(r.Context(), planID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"plan_id":             planID,
		"effectiveness_score": score,
	})
}

func (h *PlanGenerationHandler) MarkPlanForRegeneration(w http.ResponseWriter, r *http.Request) {
	planID, err := strconv.Atoi(chi.URLParam(r, "planID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.MarkPlanForRegeneration(r.Context(), planID, req.Reason); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Plan marked for regeneration"})
}

func (h *PlanGenerationHandler) GetAdaptationHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	adaptations, err := h.service.GetAdaptationHistory(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, adaptations)
}

func (h *PlanGenerationHandler) DownloadPlanPDF(w http.ResponseWriter, r *http.Request) {
	planID, err := strconv.Atoi(chi.URLParam(r, "planID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	// Generate PDF
	pdfBytes, err := h.service.ExportPlanToPDF(r.Context(), planID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate PDF: "+err.Error())
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=workout_plan_"+strconv.Itoa(planID)+".pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	// Write PDF bytes to response
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

func (h *PlanGenerationHandler) GetCurrentWeekSchema(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	schema, err := h.service.GetCurrentWeekSchema(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, schema)
}

func (h *PlanGenerationHandler) CreateWeeklySchemaFromTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     int       `json:"user_id"`
		TemplateID int       `json:"template_id"`
		WeekStart  time.Time `json:"week_start"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	schema, err := h.service.CreateWeeklySchemaFromTemplate(r.Context(), req.UserID, req.TemplateID, req.WeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, schema)
}

func (h *PlanGenerationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/plans", func(r chi.Router) {
		r.Post("/generate", h.CreatePlanGeneration)
		r.Get("/active/{userID}", h.GetActivePlan)
		r.Get("/history/{userID}", h.GetPlanHistory)
		r.Post("/{planID}/performance", h.TrackPlanPerformance)
		r.Get("/{planID}/effectiveness", h.GetPlanEffectiveness)
		r.Post("/{planID}/regenerate", h.MarkPlanForRegeneration)
		r.Get("/adaptations/{userID}", h.GetAdaptationHistory)
		r.Get("/{planID}/download", h.DownloadPlanPDF)
	})

	r.Route("/schemas", func(r chi.Router) {
		r.Get("/current/{userID}", h.GetCurrentWeekSchema)
		r.Post("/from-template", h.CreateWeeklySchemaFromTemplate)
	})
}
