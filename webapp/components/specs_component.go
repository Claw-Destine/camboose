package components

import (
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"

	pm "claw-destine.com/camboose/core/controllers/projects"
	"claw-destine.com/camboose/core/controllers/specs"
	dt "claw-destine.com/camboose/core/datatypes"
)

func NewSpecsHandler(pm *pm.ProjectControler, sm *specs.SpecsController) SpecsCompHandler {
	tpl := `<camb-specs data-project={{.Project.Id}}>
{{template "versions-list" .Versions}}
</camb-specs>

{{define "versions-list"}}
{{range .}}
<version-item slot="version-list" data-id={{.Id}} data-name={{.Name}} 
{{if .Description}}{{$attr1 := print "data-desc=" .Description}}{{$attr1 | attr}}{{end}}
data-status={{.Status}}>
{{range $key,$val := .StoryStatusCounts}}<div slot="vi-story-status" class="level-item has-text-centered">
<div><p class="heading">{{$key}}</p><p class="has-text-weight-semibold is-size-4">{{$val}}</p></div></div>{{end}}
</version-item>{{end}}
{{end}}`

	t, err := template.New("specs-view").Funcs(funcMap).Parse(tpl)
	if err != nil {
		slog.Error("Cannot parse template", "err", err)
		log.Panic("exiting")
	}
	return SpecsCompHandler{projectsCtl: pm, specsCtl: sm, specsViewTmpl: t}
}

type SpecsCompHandler struct {
	projectsCtl   *pm.ProjectControler
	specsCtl      *specs.SpecsController
	specsViewTmpl *template.Template
}

type specsViewData struct {
	Project  *dt.Project
	Versions []dt.Version
}

func (sh SpecsCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p *dt.Project
	var err error
	var vs []dt.Version

	pid := r.URL.Query().Get("currentProject")

	if pid != "" && r.Method == http.MethodGet {
		p, err = sh.projectsCtl.GetProjectById(r.Context(), pid)
		if err != nil {
			slog.Error("Failed to fetch versions for project", "id", pid, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vs, err = sh.specsCtl.ListVersions(p.Id)
		statCounts, err := sh.specsCtl.VersionStatistics(vs)
		if err != nil {
			slog.Error("Failed to fetch versions statistics for project", "id", pid, "error", err)
		}
		for idx, s := range vs {
			vs[idx].StoryStatusCounts = make(map[dt.RequirementStatus]int)
			for _, rs := range dt.ALL_RS {
				vs[idx].StoryStatusCounts[rs] = statCounts[s.Id][rs]
			}
		}
	}

	if strings.HasPrefix(r.URL.Path, "/components/specs") {
		switch r.Method {
		case "GET":
			setViewCookie(vSpecs, w)
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
		// Order matters
	} else if strings.HasPrefix(r.URL.Path, "/components/versions") {
		switch r.Method {
		case "GET":
			if err := sh.specsViewTmpl.ExecuteTemplate(w, "versions-list", vs); err != nil {
				slog.Error("Cannot render versions-list", "err", err)
				http.Error(w, "failed to render versions", http.StatusInternalServerError)
			}
			return
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
	} else if strings.HasPrefix(r.URL.Path, "/components/version") {
		switch r.Method {
		case "POST":
			sh.createVersion(r)
			http.Redirect(w, r, appendQueryParams("/components/versions", qkCurrProj, pid), http.StatusSeeOther)
		case http.MethodDelete:
			sh.deleteVersion(r)
			// Return empty string for swap
			io.WriteString(w, "")
			return
		default:
			slog.Error("Unsupported method", "method", r.Method)
			http.Error(w, "Wrong url", http.StatusMethodNotAllowed)
		}
	} else {
		slog.Error("Wrong path", "path", r.URL.Path)
		http.Error(w, "Wrong url", http.StatusBadRequest)
	}

	data := specsViewData{
		Project:  p,
		Versions: vs,
	}

	if err := sh.specsViewTmpl.Execute(w, data); err != nil {
		slog.Error("Cannot render specs component", "err", err)
		http.Error(w, "failed to render specs", http.StatusInternalServerError)
	}
}

func (sh SpecsCompHandler) createVersion(r *http.Request) {
	r.ParseForm()
	version_name := r.Form.Get("version_name")
	pid := r.Form.Get("pid")
	s := dt.Version{ProjectId: pid}
	s.Name = version_name
	sh.specsCtl.CreateVersion(s)
}

func (sh SpecsCompHandler) deleteVersion(r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	vid := urlPart[len(urlPart)-1]
	sh.specsCtl.DeleteVersionById(vid)
}
