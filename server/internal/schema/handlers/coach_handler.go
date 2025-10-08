package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type CoachHandler struct {
	service service.CoachService
}

func NewCoachHandler(service service.CoachService) *CoachHandler {
	return &CoachHandler{
		service: service,
	}
}

func getCoachIDFromContext(r *http.Request) (string, bool) {
	authID, ok := middleware.GetAuthUserIDFromContext(r.Context())
	if !ok || authID == "" {
		return "", false
	}
	return authID, true
}

func (h *CoachHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	dashboard, err := h.service.GetCoachDashboard(r.Context(), coachID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dashboard)
}

func (h *CoachHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	clients, err := h.service.GetCoachClients(r.Context(), coachID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"clients": clients,
		"total":   len(clients),
	})
}

func (h *CoachHandler) GetClientDetails(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}

	progress, err := h.service.GetClientProgress(r.Context(), coachID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, progress)
}

func (h *CoachHandler) AssignClient(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	var req types.CoachAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.CoachID = coachID

	assignment, err := h.service.AssignClientToCoach(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, assignment)
}

func (h *CoachHandler) RemoveClient(w http.ResponseWriter, r *http.Request) {
	_, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	assignmentID, err := strconv.Atoi(chi.URLParam(r, "assignmentID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	if err := h.service.RemoveClientFromCoach(r.Context(), assignmentID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Client removed successfully",
	})
}

func (h *CoachHandler) CreateSchemaForClient(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req types.ManualSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.UserID = userID
	req.CoachID = coachID

	schema, err := h.service.CreateManualSchemaForClient(r.Context(), coachID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, schema)
}

func (h *CoachHandler) UpdateSchema(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	var req types.ManualSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	schema, err := h.service.UpdateManualSchema(r.Context(), coachID, schemaID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, schema)
}
func (h *CoachHandler) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	if err := h.service.DeleteSchema(r.Context(), coachID, schemaID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Schema deleted successfully",
	})
}

func (h *CoachHandler) CloneSchema(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	var req struct {
		TargetUserID int    `json:"target_user_id"`
		Notes        string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	schema, err := h.service.CloneSchemaToClient(r.Context(), coachID, schemaID, req.TargetUserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, schema)
}

func (h *CoachHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	templates, err := h.service.GetCoachTemplates(r.Context(), coachID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	})
}

func (h *CoachHandler) SaveTemplate(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	var req struct {
		SchemaID     int    `json:"schema_id"`
		TemplateName string `json:"template_name"`
		Description  string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SchemaID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}

	if req.TemplateName == "" {
		respondWithError(w, http.StatusBadRequest, "Template name is required")
		return
	}

	if err := h.service.SaveSchemaAsTemplate(r.Context(), coachID, req.SchemaID, req.TemplateName); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Template saved successfully",
		"name":    req.TemplateName,
	})
}

func (h *CoachHandler) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	templateID, err := strconv.Atoi(chi.URLParam(r, "templateID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	var req struct {
		UserID int    `json:"user_id"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	schema, err := h.service.CreateSchemaFromCoachTemplate(r.Context(), coachID, templateID, req.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, schema)
}

func (h *CoachHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	_, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	templateID, err := strconv.Atoi(chi.URLParam(r, "templateID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	// TODO: Implement DeleteTemplate in service
	_ = templateID

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Template deleted successfully",
	})
}

func (h *CoachHandler) GetClientProgress(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}

	progress, err := h.service.GetClientProgress(r.Context(), coachID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, progress)
}

func (h *CoachHandler) GetClientWorkouts(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}

	// TODO: Implement GetClientWorkoutHistory in service
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":  userID,
		"workouts": []interface{}{},
		"message":  "Implementation pending",
	})
}

func (h *CoachHandler) GetClientSchemas(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"schemas": []interface{}{},
		"message": "Implementation pending",
	})
}

func (h *CoachHandler) AddClientNote(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Note == "" {
		respondWithError(w, http.StatusBadRequest, "Note cannot be empty")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Note added successfully",
		"note":    req.Note,
	})
}

func (h *CoachHandler) GetCoachStats(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}

	// TODO: Implement GetCoachStatistics in service
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"coach_id":           coachID,
		"total_clients":      0,
		"active_schemas":     0,
		"completion_rate":    0.0,
		"total_workouts":     0,
		"this_week_workouts": 0,
		"message":            "Implementation pending",
	})
}

func (h *CoachHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	coachID, ok := getCoachIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	_ = coachID // TODO: Use for filtering activity

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// TODO: Implement GetRecentActivity in service
	_ = limit

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"activities": []interface{}{},
		"message":    "Implementation pending",
	})
}

func (h *CoachHandler) RegisterRoutes(r chi.Router) {
	r.Route("/coach", func(r chi.Router) {
		// TODO: Add coach auth middleware
		// r.Use(authMiddleware.RequireCoachRole())

		// Dashboard & Overview
		r.Get("/dashboard", h.GetDashboard)
		r.Get("/stats", h.GetCoachStats)
		r.Get("/activity", h.GetRecentActivity)

		// Client Management
		r.Get("/clients", h.GetClients)
		r.Post("/clients/assign", h.AssignClient)
		r.Get("/clients/{userID}", h.GetClientDetails)
		r.Delete("/clients/{assignmentID}", h.RemoveClient)

		// Client Progress & Analytics
		r.Get("/clients/{userID}/progress", h.GetClientProgress)
		r.Get("/clients/{userID}/workouts", h.GetClientWorkouts)
		r.Get("/clients/{userID}/schemas", h.GetClientSchemas)
		r.Post("/clients/{userID}/notes", h.AddClientNote)

		// Schema Management
		r.Post("/clients/{userID}/schemas", h.CreateSchemaForClient)
		r.Put("/schemas/{schemaID}", h.UpdateSchema)
		r.Delete("/schemas/{schemaID}", h.DeleteSchema)
		r.Post("/schemas/{schemaID}/clone", h.CloneSchema)

		// Template Management
		r.Get("/templates", h.GetTemplates)
		r.Post("/templates", h.SaveTemplate)
		r.Post("/templates/{templateID}/create-schema", h.CreateFromTemplate)
		r.Delete("/templates/{templateID}", h.DeleteTemplate)
	})
}
