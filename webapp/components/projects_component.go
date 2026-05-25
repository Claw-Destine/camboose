package components

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	dt "claw-destine.com/camboose/core/datatypes"

	pm "claw-destine.com/camboose/core/controllers/projects"
)

func NewProjectsHandler(pm *pm.ProjectControler, rm *pm.RecipeController) ProjectsCompHandler {
	return ProjectsCompHandler{projectManager: pm, recipeManager: rm}
}

type ProjectsCompHandler struct {
	projectManager *pm.ProjectControler
	recipeManager  *pm.RecipeController
}

type projectView int

const (
	allProjectsView   = iota
	singleProjectView = iota
)

func (ph ProjectsCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Order matters
	if strings.HasPrefix(r.URL.Path, "/components/projects") {
		ph.displayProjectView(allProjectsView, w, r)
	} else if strings.HasPrefix(r.URL.Path, "/components/project/") {
		switch r.Method {
		case "GET":
			ph.displayProjectView(singleProjectView, w, r)
		case "POST":
			ph.createProject(w, r)
		case "PUT":
			ph.updateProject(w, r)
			ph.displayProjectView(singleProjectView, w, r)
		case "DELETE":
			ph.deleteProject(w, r)
			ph.displayProjectView(allProjectsView, w, r)
		default:
			slog.Error("Unknown method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
	} else {
		slog.Error("Wrong path", "path", r.URL.Path)
		http.Error(w, "Wrong url", http.StatusBadRequest)
	}
}

func (ph ProjectsCompHandler) updateProject(_ http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	pid := urlPart[len(urlPart)-1]
	project := dt.Project{Base: dt.Base{Id: pid}}

	r.ParseForm()
	recipe := r.Form.Get("recipe")
	if recipe != "" {
		project.Recipe = recipe
	}

	ph.projectManager.UpdateProject(r.Context(), project)
}

func (ph ProjectsCompHandler) deleteProject(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	pid := urlPart[len(urlPart)-1]
	err := ph.projectManager.DeleteProject(r.Context(), pid)
	if err != nil {
		slog.Error("Failed to delete project", "id", pid, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
}

func (ph ProjectsCompHandler) createProject(w http.ResponseWriter, r *http.Request) {
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
}

func (ph ProjectsCompHandler) displayProjectView(view projectView, w http.ResponseWriter, r *http.Request) {
	projects, err := ph.projectManager.ListProjects(nil)
	if err != nil {
		slog.Error("Failed to fetch projects", "path", r.URL.Path, "reason", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pid string
	var p *dt.Project

	switch view {
	case singleProjectView:
		urlPart := strings.Split(r.URL.Path, "/")
		pid = urlPart[len(urlPart)-1]

	case allProjectsView:
		pid = r.URL.Query().Get("currentProject")
	}

	if pid != "" {
		p, err = ph.projectManager.GetProjectById(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to fetch project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats, err := ph.projectManager.ProjectStatistics([]string{pid})
		if err != nil {
			if err != nil {
				slog.Error("Failed to stats for project", "id", pid, "error", err)
			}
		}
		p.VersionStatusCounts = stats[pid]
	}

	recipies, err := ph.recipeManager.ListRecipes()
	if err != nil {
		slog.Error("Failed to fetch recipies", "path", r.URL.Path, "reason", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch view {
	case singleProjectView:
		err = projectComponent(p, recipies).Render(r.Context(), w)
	case allProjectsView:
		err = projectsComponent(projects, p, recipies).Render(r.Context(), w)
	}
	if err != nil {
		slog.Error("Failed to render component", "path", r.URL.Path, "reason", err)
	}
}

func projectLink(p string) string {
	return "/components/project/" + p
}
