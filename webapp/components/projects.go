package components

import (
	"log/slog"
	"net/http"

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
		}
		slog.Info("Create project", "name", pname)
	case "GET":
		pname := r.URL.Query().Get("project")
		p = dt.Project{ObjectId: "123", Name: pname}
		print(r.Method)
	}
	projectsComponent(ph.projectManager.ListProjects(), p).Render(r.Context(), w)
}

func projectLink(p string) string {
	return "/components/project/" + p
}
