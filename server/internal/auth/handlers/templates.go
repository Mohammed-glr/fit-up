package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

func (h *AuthHandler) handleGetUserTemplates(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	page := 1
	pageSize := 20

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	templates, err := h.store.GetUserTemplates(r.Context(), userID, page, pageSize)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, templates)
}

func (h *AuthHandler) handleGetPublicTemplates(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 20

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	templates, err := h.store.GetPublicTemplates(r.Context(), page, pageSize)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, templates)
}

func (h *AuthHandler) handleGetTemplateByID(w http.ResponseWriter, r *http.Request) {
	templateIDStr := chi.URLParam(r, "templateId")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	template, err := h.store.GetTemplateByID(r.Context(), templateID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, template)
}

func (h *AuthHandler) handleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	var req types.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding template request: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	log.Printf("Received template request - Name: %s, IsPublic: %v, ExerciseCount: %d", req.Name, req.IsPublic, len(req.Exercises))
	if len(req.Exercises) > 0 {
		log.Printf("First exercise: %+v", req.Exercises[0])
	}

	if req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name is required"))
		return
	}

	if len(req.Exercises) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one exercise is required"))
		return
	}

	template, err := h.store.CreateTemplate(r.Context(), userID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, template)
}

func (h *AuthHandler) handleUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	templateIDStr := chi.URLParam(r, "templateId")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req types.UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	template, err := h.store.UpdateTemplate(r.Context(), templateID, userID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, template)
}

func (h *AuthHandler) handleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	templateIDStr := chi.URLParam(r, "templateId")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.DeleteTemplate(r.Context(), templateID, userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}

func (h *AuthHandler) RegisterTemplateRoutes(r chi.Router, authMW *middleware.AuthMiddleware) {
	r.Route("/templates", func(r chi.Router) {
		r.Get("/public", h.handleGetPublicTemplates)

		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireJWTAuth())

			r.Get("/", h.handleGetUserTemplates)
			r.Post("/", h.handleCreateTemplate)
			r.Put("/{templateId}", h.handleUpdateTemplate)
			r.Delete("/{templateId}", h.handleDeleteTemplate)
		})
	})
}
