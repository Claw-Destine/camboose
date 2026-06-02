package components

import (
	"io"
	"net/http"

	"claw-destine.com/camboose/core/controllers/projects"
)

func NewRecipuesHandler(recipeManager *projects.RecipeController) RecipiesCompHandler {
	return RecipiesCompHandler{recipeManager: recipeManager}
}

type RecipiesCompHandler struct {
	recipeManager *projects.RecipeController
}

func (ph RecipiesCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setViewCookie(vRecipies, w)
	io.WriteString(w, "recipies")
	// recipies, err := ph.recipeManager.ListRecipes()
	// if err != nil {
	// 	slog.Error("Failed to load recipies", "reason", err)
	// }
	// err = recipiesComponent(recipies).Render(r.Context(), w)
	// if err != nil {
	// 	slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
	// }
}
