package components

import (
	"log/slog"
	"net/http"
)

func NewRecipuesHandler() RecipiesHandler {
	return RecipiesHandler{}
}

type RecipiesHandler struct {
}

func (ph RecipiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := recipiesComponent().Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
	}
}
