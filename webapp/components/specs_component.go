package components

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	pm "claw-destine.com/camboose/core/controllers/projects"
	"claw-destine.com/camboose/core/controllers/specs"
	dt "claw-destine.com/camboose/core/datatypes"
)

func NewSpecsHandler(pm *pm.ProjectControler, sm *specs.SpecsController) SpecsCompHandler {
	return SpecsCompHandler{projectsCtl: pm, specsCtl: sm}
}

type SpecsCompHandler struct {
	projectsCtl *pm.ProjectControler
	specsCtl    *specs.SpecsController
}

func (sh SpecsCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var p *dt.Project
	var err error
	var si []dt.SpecItem

	pid := r.URL.Query().Get("currentProject")

	if pid != "" {
		p, err = sh.projectsCtl.GetProjectById(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to fetch project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		si, err = sh.specsCtl.ListSpecItems(p.Id)
	}

	if strings.HasPrefix(r.URL.Path, "/components/specs") {
		switch r.Method {
		case "GET":
			sh.displaySpecsPage(p, si, w, r)
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
		// Order matters
	} else if strings.HasPrefix(r.URL.Path, "/components/versions") {
		switch r.Method {
		case "GET":
			sh.displayVersionList(si, w, r)
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}

	} else if strings.HasPrefix(r.URL.Path, "/components/version") {
		switch r.Method {
		case "POST":
			sh.createVersion(*p, r)
			http.Redirect(w, r, appendQueryParams("/components/versions", paramCurrentProjec, pid), http.StatusSeeOther)
		case http.MethodDelete:
			sh.deleteVersion(r)
			io.WriteString(w, " \n")
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
	} else {
		slog.Error("Wrong path", "path", r.URL.Path)
		http.Error(w, "Wrong url", http.StatusBadRequest)
	}

}
func (sh SpecsCompHandler) displaySpecsPage(p *dt.Project, si []dt.SpecItem, w http.ResponseWriter, r *http.Request) {
	specsComponent(p, si).Render(r.Context(), w)
}

func (sh SpecsCompHandler) displayVersionList(si []dt.SpecItem, w http.ResponseWriter, r *http.Request) {
	versionList(si).Render(r.Context(), w)
}

func (sh SpecsCompHandler) createVersion(project dt.Project, r *http.Request) {
	r.ParseForm()
	version_name := r.Form.Get("name")
	s := dt.SpecItem{Type: dt.Version, ProjectId: project.Id}
	s.Name = version_name
	sh.specsCtl.CreateSpecItem(s)
}

func (sh SpecsCompHandler) deleteVersion(r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	vid := urlPart[len(urlPart)-1]
	sh.specsCtl.DeleteSpecItemById(vid)
}
