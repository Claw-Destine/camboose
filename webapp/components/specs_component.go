package components

import (
	"log/slog"
	"net/http"

	pm "claw-destine.com/camboose/core/controllers/projects"
	dt "claw-destine.com/camboose/core/datatypes"
)

func NewSpecsHandler(pm *pm.ProjectControler) SpecsCompHandler {
	return SpecsCompHandler{projectManager: pm}
}

type SpecsCompHandler struct {
	projectManager *pm.ProjectControler
}

func (sh SpecsCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p *dt.Project
	var err error
	pid := r.URL.Query().Get("currentProject")

	if pid != "" {
		p, err = sh.projectManager.GetProjectById(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to fetch project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	specsComponent(p).Render(r.Context(), w)
}
