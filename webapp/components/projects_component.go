package components

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	dt "claw-destine.com/camboose/core/datatypes"

	pm "claw-destine.com/camboose/core/controllers/projects"
)

func NewProjectHandler(pm *pm.ProjectManager) ProjectHandler {
	return ProjectHandler{projectManager: pm}
}

type ProjectHandler struct {
	projectManager *pm.ProjectManager
}

func (ph ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	urlPart := strings.Split(r.URL.Path, "/")
	pid := urlPart[len(urlPart)-1]

	switch r.Method {
	case "GET":
		p, err := ph.projectManager.GetProjectById(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to fetch project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		}

		projectComponent(*p).Render(r.Context(), w)
	case "DELETE":
		err := ph.projectManager.DeleteProject(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to delete project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		}
		projects, err := ph.projectManager.ListProjects(nil)
		if err != nil {
			slog.Error("Failed to fetch projects", "path", r.URL.Path, "reason", err)
		}
		err = projectsComponent(projects, dt.Project{}).Render(r.Context(), w)
		if err != nil {
			slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
		}
	}
}

func NewProjectsHandler(pm *pm.ProjectManager) ProjectsHandler {
	return ProjectsHandler{projectManager: pm}
}

type ProjectsHandler struct {
	projectManager *pm.ProjectManager
}

func (ph ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p dt.Project
	switch r.Method {
	case "POST":
		r.ParseForm()
		pname := r.Form.Get("name")
		np, err := ph.projectManager.CreateProject(r.Context(), dt.Project{Base: dt.Base{Name: pname}})
		if err != nil {
			slog.Error("Failed to create the project", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Info("Create project", "name", pname, "id", np.Id)
		http.Redirect(w, r, fmt.Sprintf("/components/body?currentProject=%s", np.Id), http.StatusSeeOther)
		return

	case "GET":
		id := r.URL.Query().Get("project")
		if id != "" {
			pp, err := ph.projectManager.GetProjectById(r.Context(), id)
			p = *pp
			if err != nil {
				slog.Error("Requested unknown project", "id", id)
			}
		}
	default:
		slog.Error("Unknown method", "method", r.Method)
	}
	projects, err := ph.projectManager.ListProjects(nil)
	if err != nil {
		slog.Error("Failed to fetch projects", "path", r.URL.Path, "reason", err)
	}
	err = projectsComponent(projects, p).Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
	}
}

func projectLink(p string) string {
	return "/components/project/" + p
}
