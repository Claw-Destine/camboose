package components

import (
	"io"
	"net/http"

	pm "claw-destine.com/camboose/core/controllers/projects"
	"claw-destine.com/camboose/core/controllers/specs"
)

func NewSpecsHandler(pm *pm.ProjectControler, sm *specs.SpecsController) SpecsCompHandler {
	return SpecsCompHandler{projectsCtl: pm, specsCtl: sm}
}

type SpecsCompHandler struct {
	projectsCtl *pm.ProjectControler
	specsCtl    *specs.SpecsController
}

func (sh SpecsCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setViewCookie(vRecipies, w)
	io.WriteString(w, "")
}

// 	var p *dt.Project
// 	var err error
// 	var si []dt.Version

// 	pid := r.URL.Query().Get("currentProject")

// 	if pid != "" && r.Method == http.MethodGet {
// 		p, err = sh.projectsCtl.GetProjectById(r.Context(), pid)
// 		if err != nil {
// 			slog.Error("Failed to fetch versions for project", "id", pid, "error", err)
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		si, err = sh.specsCtl.ListVersions(p.Id)
// 		statCounts, err := sh.specsCtl.VersionStatistics(si)
// 		if err != nil {
// 			slog.Error("Failed to fetch versions statistics for project", "id", pid, "error", err)
// 		}
// 		for idx, s := range si {
// 			si[idx].StoryStatusCounts = statCounts[s.Id]
// 		}
// 	}

// 	if strings.HasPrefix(r.URL.Path, "/components/specs") {
// 		switch r.Method {
// 		case "GET":
// 			setViewCookie(vSpecs, w)
// 			sh.displaySpecsPage(p, si, w, r)
// 		default:
// 			slog.Error("Unsupported method", "method", r.Method)
// 			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
// 		}
// 		// Order matters
// 	} else if strings.HasPrefix(r.URL.Path, "/components/versions") {
// 		switch r.Method {
// 		case "GET":
// 			sh.displayVersionList(si, w, r)
// 		default:
// 			slog.Error("Unsupported method", "method", r.Method)
// 			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
// 		}
// 	} else if strings.HasPrefix(r.URL.Path, "/components/version") {
// 		switch r.Method {
// 		case "POST":
// 			sh.createVersion(*p, r)
// 			http.Redirect(w, r, appendQueryParams("/components/versions", qkCurrProj, pid), http.StatusSeeOther)
// 		case http.MethodDelete:
// 			sh.deleteVersion(r)
// 			// Return empty string for swap
// 			io.WriteString(w, "")
// 		default:
// 			slog.Error("Unsupported method", "method", r.Method)
// 			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
// 		}
// 	} else {
// 		slog.Error("Wrong path", "path", r.URL.Path)
// 		http.Error(w, "Wrong url", http.StatusBadRequest)
// 	}

// }
// func (sh SpecsCompHandler) displaySpecsPage(p *dt.Project, si []dt.Version, w http.ResponseWriter, r *http.Request) {
// 	specsComponent(p, si).Render(r.Context(), w)
// }

// func (sh SpecsCompHandler) displayVersionList(si []dt.Version, w http.ResponseWriter, r *http.Request) {
// 	versionList(si).Render(r.Context(), w)
// }

// func (sh SpecsCompHandler) createVersion(project dt.Project, r *http.Request) {
// 	r.ParseForm()
// 	version_name := r.Form.Get("name")
// 	s := dt.Version{ProjectId: project.Id}
// 	s.Name = version_name
// 	sh.specsCtl.CreateVersion(s)
// }

// func (sh SpecsCompHandler) deleteVersion(r *http.Request) {
// 	urlPart := strings.Split(r.URL.Path, "/")
// 	vid := urlPart[len(urlPart)-1]
// 	sh.specsCtl.DeleteVersionById(vid)
// }
