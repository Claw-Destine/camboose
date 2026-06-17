package components

import (
	"html/template"
	"log/slog"
	"net/http"
)

// qk - query param keys keys
const qkCurrProj = "currentProject"

func appendQueryParams(basePath string, params ...string) string {
	if len(params)%2 != 0 {
		slog.Warn("List of params for this method should be even")
	}
	ret := basePath
	separator := "?"
	i := 0
	for {
		if i+2 > len(params) {
			break
		}
		ret = ret + separator + params[i] + "=" + params[i+1]
		separator = "&"
		i = i + 2
	}
	return ret
}

// views
type CambView string

const vSpecs = "specs"
const vTasks = "tasks"
const vProjects = "projects"
const vRecipies = "recipies"

func setViewCookie(view CambView, w http.ResponseWriter) {
	viewCookie := http.Cookie{
		Name:     "view",
		Value:    string(view),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &viewCookie)
}

func redirectWithHX(w http.ResponseWriter, r *http.Request, target string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", target)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, target, http.StatusSeeOther)
}

var funcMap = template.FuncMap{
	"attr": func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	},
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
}
