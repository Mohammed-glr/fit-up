package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
)

// System Recipe Handlers

func (h *FoodTrackerHandler) GetSystemRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	recipe, err := h.service.Recipes().GetSystemRecipe(r.Context(), id)
	if err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recipe)
}

func (h *FoodTrackerHandler) ListSystemRecipes(w http.ResponseWriter, r *http.Request) {
	filters := parseRecipeFilters(r)

	recipes, err := h.service.Recipes().ListSystemRecipes(r.Context(), filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"recipes": recipes,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	})
}

func (h *FoodTrackerHandler) CreateSystemRecipe(w http.ResponseWriter, r *http.Request) {
	var req types.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.Recipes().CreateSystemRecipe(r.Context(), &req)
	if err != nil {
		if err == types.ErrInvalidRequest {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, recipe)
}

func (h *FoodTrackerHandler) UpdateSystemRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	var req types.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.Recipes().UpdateSystemRecipe(r.Context(), id, &req)
	if err != nil {
		if err == types.ErrInvalidRequest || err == types.ErrInvalidID {
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

	respondWithJSON(w, http.StatusOK, recipe)
}

func (h *FoodTrackerHandler) DeleteSystemRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	if err := h.service.Recipes().DeleteSystemRecipe(r.Context(), id); err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Recipe deleted successfully"})
}

func (h *FoodTrackerHandler) GetUserRecipe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	recipe, err := h.service.Recipes().GetUserRecipe(r.Context(), id, userID)
	if err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recipe)
}

func (h *FoodTrackerHandler) ListUserRecipes(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filters := parseRecipeFilters(r)

	recipes, err := h.service.Recipes().ListUserRecipes(r.Context(), userID, filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"recipes": recipes,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	})
}

func (h *FoodTrackerHandler) CreateUserRecipe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req types.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.Recipes().CreateUserRecipe(r.Context(), userID, &req)
	if err != nil {
		if err == types.ErrInvalidRequest {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, recipe)
}

func (h *FoodTrackerHandler) UpdateUserRecipe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	var req types.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.Recipes().UpdateUserRecipe(r.Context(), id, userID, &req)
	if err != nil {
		if err == types.ErrInvalidRequest || err == types.ErrInvalidID {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		if err == types.ErrUnauthorized {
			respondWithError(w, http.StatusForbidden, "You don't have permission to update this recipe")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recipe)
}

func (h *FoodTrackerHandler) DeleteUserRecipe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	if err := h.service.Recipes().DeleteUserRecipe(r.Context(), id, userID); err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		if err == types.ErrUnauthorized {
			respondWithError(w, http.StatusForbidden, "You don't have permission to delete this recipe")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Recipe deleted successfully"})
}

func (h *FoodTrackerHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	query := parseSearchQuery(r)

	recipes, err := h.service.Recipes().SearchRecipes(r.Context(), userID, query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"recipes": recipes,
		"limit":   query.Limit,
		"offset":  query.Offset,
	})
}

func (h *FoodTrackerHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	recipeIDStr := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	if err := h.service.Recipes().ToggleFavorite(r.Context(), userID, recipeID); err != nil {
		if err == types.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Favorite status toggled successfully"})
}

func (h *FoodTrackerHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	favorites, err := h.service.Recipes().GetFavorites(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"favorites": favorites,
	})
}
