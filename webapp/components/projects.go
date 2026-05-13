package components

import (
	"log/slog"
	"net/http"
	"strings"

	dt "claw-destine.com/camboose/core/datatypes"

	pm "claw-destine.com/camboose/core/projects"
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
		p, err := ph.projectManager.GetProjectById(pid)
		if err != nil {
			slog.Error("Failed to fetch project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		}

		projectComponent(*p).Render(r.Context(), w)
	case "DELETE":
		err := ph.projectManager.DeleteProject(pid)
		if err != nil {
			slog.Error("Failed to delete project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		}
		projectsComponent(ph.projectManager.ListProjects(), dt.Project{}).Render(r.Context(), w)
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
		_, err := ph.projectManager.CreateProject(pname)
		if err != nil {
			slog.Error("Failed to create the project", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Info("Create project", "name", pname)
	case "GET":
		id := r.URL.Query().Get("project")
		if id != "" {
			pp, err := ph.projectManager.GetProjectById(id)
			p = *pp
			if err != nil {
				slog.Error("Requested unknown project", "id", id)
			}
		}
		print(r.Method)
	default:
		slog.Error("Unknown method", "method", r.Method)
	}
	projectsComponent(ph.projectManager.ListProjects(), p).Render(r.Context(), w)
}

func projectLink(p string) string {
	return "/components/project/" + p
}
