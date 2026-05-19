package components

import (
	"net/http"

	pm "claw-destine.com/camboose/core/controllers/projects"
)

func NewSpecsHandler(pm *pm.ProjectManager) SpecsHandler {
	return SpecsHandler{projectManager: pm}
}

type SpecsHandler struct {
	projectManager *pm.ProjectManager
}

func (sh SpecsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	specsComponent().Render(r.Context(), w)
}
