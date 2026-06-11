package components

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"

	pm "claw-destine.com/camboose/core/controllers/projects"
	dt "claw-destine.com/camboose/core/datatypes"
)

type bodyData struct {
	Project  *dt.Project
	View     string
	Projects []dt.Project
}

func NewBodyHandler(pm *pm.ProjectControler) BodyCompHandler {
	bh := BodyCompHandler{projectManager: pm}
	tpl := `<camb-body {{if .Project}}project-id="{{.Project.Id}}"{{end}} view="{{.View}}" id="main-body">
	{{range .Projects}}<a slot="project-dropdown" class="dropdown-item{{if and $.Project (eq $.Project.Id .Id)}} is-active{{end}}"
		hx-get="/components/body?currentProject={{.Id}}"
		hx-swap="outerHTML" hx-target="#main-body"> {{.Name}}
	</a>
	{{end}}</camb-body>
	`

	t, err := template.New("main-body").Parse(tpl)
	if err != nil {
		slog.Error("Cannot parse template", "err", err)
		log.Panic("exiting")
	}
	bh.templ = t
	return bh
}

type BodyCompHandler struct {
	projectManager *pm.ProjectControler
	templ          *template.Template
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
		var projectCookie http.Cookie
		if err != nil {
			slog.Info("Failed to load current project", "id", currentProjectId)
			projectCookie = http.Cookie{
				Name:   "project",
				MaxAge: -1,
			}
		} else {
			currentProject = cp
			projectCookie = http.Cookie{
				Name:     "project",
				Value:    currentProjectId,
				SameSite: http.SameSiteLaxMode,
			}

		}
		http.SetCookie(w, &projectCookie)
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

	data := bodyData{
		Project:  currentProject,
		View:     currentView,
		Projects: lastProjects,
	}

	if err := ph.templ.Execute(w, data); err != nil {
		slog.Error("Cannot render body component", "err", err)
		http.Error(w, "failed to render body", http.StatusInternalServerError)
	}
}
