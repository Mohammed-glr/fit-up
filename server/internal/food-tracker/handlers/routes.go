package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	authRepo "github.com/tdmdh/fit-up-server/internal/auth/repository"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/services"
	"github.com/tdmdh/fit-up-server/internal/food-tracker/types"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type FoodTrackerHandler struct {
	authMiddleware *middleware.AuthMiddleware
	service        services.FoodTrackerService
}

func NewFoodTrackerHandler(
	service services.FoodTrackerService,
	schemaRepo repository.SchemaRepo,
	userStore authRepo.UserStore,
) *FoodTrackerHandler {
	return &FoodTrackerHandler{
		authMiddleware: middleware.NewAuthMiddleware(schemaRepo, userStore),
		service:        service,
	}
}

func (h *FoodTrackerHandler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get("/food-tracker/recipes/system", h.ListSystemRecipes)
		r.Get("/food-tracker/recipes/system/{id}", h.GetSystemRecipe)
		r.Get("/food-tracker/recipes/search", h.SearchRecipes)
	})

	router.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.RequireJWTAuth())
		r.Use(h.authMiddleware.RequireAdminRole())

		r.Post("/food-tracker/recipes/system", h.CreateSystemRecipe)
		r.Put("/food-tracker/recipes/system/{id}", h.UpdateSystemRecipe)
		r.Delete("/food-tracker/recipes/system/{id}", h.DeleteSystemRecipe)
	})

	router.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.RequireJWTAuth())

		r.Route("/food-tracker/recipes/user", func(r chi.Router) {
			r.Get("/", h.ListUserRecipes)
			r.Post("/", h.CreateUserRecipe)
			r.Get("/{id}", h.GetUserRecipe)
			r.Put("/{id}", h.UpdateUserRecipe)
			r.Delete("/{id}", h.DeleteUserRecipe)
		})

		r.Route("/food-tracker/recipes/favorites", func(r chi.Router) {
			r.Get("/", h.GetFavorites)
			r.Patch("/{recipeID}", h.ToggleFavorite)
		})

		r.Route("/food-tracker/food-logs", func(r chi.Router) {
			r.Post("/", h.LogFood)
			r.Post("/recipe", h.LogRecipe)
			r.Get("/date/{date}", h.GetLogsByDate)
			r.Get("/range", h.GetLogsByDateRange)
			r.Get("/{id}", h.GetFoodLogEntry)
			r.Put("/{id}", h.UpdateFoodLog)
			r.Delete("/{id}", h.DeleteFoodLog)
		})

		r.Route("/food-tracker/nutrition", func(r chi.Router) {
			r.Get("/daily/{date}", h.GetDailyNutrition)
			r.Get("/weekly", h.GetWeeklyNutrition)
			r.Get("/monthly", h.GetMonthlyNutrition)

			r.Route("/goals", func(r chi.Router) {
				r.Get("/", withContext(h.GetNutritionGoals))
				r.Post("/", h.CreateOrUpdateNutritionGoals)
			})

			r.Get("/comparison/{date}", withContext(h.GetNutritionComparison))
			r.Get("/insights/{date}", withContext(h.GetNutritionInsights))
		})
	})
}

func withContext(fn func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(r.Context(), w, r)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func getUserID(r *http.Request) string {
	if authID, ok := middleware.GetAuthUserIDFromContext(r.Context()); ok && authID != "" {
		return authID
	}

	userID := r.Context().Value("userID")
	if userID == nil {
		return ""
	}
	if uid, ok := userID.(string); ok {
		return uid
	}
	return ""
}

func parseRecipeFilters(r *http.Request) types.RecipeFilters {
	filters := types.RecipeFilters{
		Limit:  20,
		Offset: 0,
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filters.Offset = o
		}
	}

	if category := r.URL.Query().Get("category"); category != "" {
		cat := types.RecipeCategory(category)
		filters.Category = &cat
	}

	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		diff := types.RecipeDifficulty(difficulty)
		filters.Difficulty = &diff
	}

	if maxCals := r.URL.Query().Get("max_calories"); maxCals != "" {
		if mc, err := strconv.Atoi(maxCals); err == nil && mc > 0 {
			filters.MaxCalories = &mc
		}
	}

	if minProt := r.URL.Query().Get("min_protein"); minProt != "" {
		if mp, err := strconv.Atoi(minProt); err == nil && mp > 0 {
			filters.MinProtein = &mp
		}
	}

	if maxPrep := r.URL.Query().Get("max_prep_time"); maxPrep != "" {
		if mp, err := strconv.Atoi(maxPrep); err == nil && mp > 0 {
			filters.MaxPrepTime = &mp
		}
	}

	if search := r.URL.Query().Get("search"); search != "" {
		filters.SearchTerm = search
	}

	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}

	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder != "" {
		filters.SortOrder = sortOrder
	}

	if isFav := r.URL.Query().Get("favorites_only"); isFav != "" {
		fav := isFav == "true"
		filters.IsFavorite = &fav
	}

	return filters
}

func parseSearchQuery(r *http.Request) types.SearchQuery {
	query := types.SearchQuery{
		Limit:         20,
		Offset:        0,
		IncludeSystem: true,
		IncludeUser:   true,
	}

	if term := r.URL.Query().Get("term"); term != "" {
		query.Term = term
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			query.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			query.Offset = o
		}
	}

	if category := r.URL.Query().Get("category"); category != "" {
		cat := types.RecipeCategory(category)
		query.Category = &cat
	}

	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		diff := types.RecipeDifficulty(difficulty)
		query.Difficulty = &diff
	}

	if maxCals := r.URL.Query().Get("max_calories"); maxCals != "" {
		if mc, err := strconv.Atoi(maxCals); err == nil && mc > 0 {
			query.MaxCalories = &mc
		}
	}

	if minProt := r.URL.Query().Get("min_protein"); minProt != "" {
		if mp, err := strconv.Atoi(minProt); err == nil && mp > 0 {
			query.MinProtein = &mp
		}
	}

	if maxPrep := r.URL.Query().Get("max_prep_time"); maxPrep != "" {
		if mp, err := strconv.Atoi(maxPrep); err == nil && mp > 0 {
			query.MaxPrepTime = &mp
		}
	}

	if includeSystem := r.URL.Query().Get("include_system"); includeSystem != "" {
		query.IncludeSystem = includeSystem == "true"
	}

	if includeUser := r.URL.Query().Get("include_user"); includeUser != "" {
		query.IncludeUser = includeUser == "true"
	}

	if favOnly := r.URL.Query().Get("favorites_only"); favOnly != "" {
		query.FavoritesOnly = favOnly == "true"
	}

	return query
}
