package components

import (
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/core/controllers/projects"
)

func NewRecipuesHandler(recipeManager *projects.RecipeManager) RecipiesHandler {
	return RecipiesHandler{recipeManager: recipeManager}
}

type RecipiesHandler struct {
	recipeManager *projects.RecipeManager
}

func (ph RecipiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recipies, err := ph.recipeManager.ListRecipes()
	if err != nil {
		slog.Error("Failed to load recipies", "reason", err)
	}
	err = recipiesComponent(recipies).Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
	}
}
