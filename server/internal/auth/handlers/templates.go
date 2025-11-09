package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

func (h *AuthHandler) handleGetUserTemplates(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

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
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req types.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
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
	userID := r.Context().Value(middleware.UserIDKey).(string)

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
	userID := r.Context().Value(middleware.UserIDKey).(string)

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
