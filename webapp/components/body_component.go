package components

import (
	"log/slog"
	"net/http"

	pm "claw-destine.com/camboose/core/controllers/projects"
	dt "claw-destine.com/camboose/core/datatypes"
)

func NewBodyHandler(pm *pm.ProjectControler) BodyCompHandler {
	return BodyCompHandler{projectManager: pm}
}

type BodyCompHandler struct {
	projectManager *pm.ProjectControler
}

func (ph BodyCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var currentProject *dt.Project

	currentProjectId := r.URL.Query().Get("currentProject")

	if currentProjectId == "" {
		projectCookie, err := r.Cookie("project")
		if err == nil {
			currentProjectId = projectCookie.Value
		}
	}

	if currentProjectId != "" {
		cp, err := ph.projectManager.GetProjectById(r.Context(), currentProjectId)
		if err != nil {
			slog.Error("Failed to load current project", "id", currentProjectId)
		} else {
			currentProject = cp
			projectCookie := http.Cookie{
				Name:     "project",
				Value:    currentProjectId,
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, &projectCookie)
		}
	}
	var currentView string
	viewCookie, err := r.Cookie("view")
	if err == nil {
		currentView = viewCookie.Value
	} else {
		currentView = "specs"
	}

	// First let's get last projects
	filter := dt.QuerySettings{
		Limit:       5,
		Offset:      0,
		OrderFields: []string{"updated_at"},
		Ascending:   false,
	}
	lastProjects, error := ph.projectManager.ListProjects(&filter)

	if error != nil {
		// Nothing much we can do on error, let's just report it
		slog.Error("Cannot load last projects")
	}

	bodyComponent(currentProject, lastProjects, currentView).Render(r.Context(), w)
}
