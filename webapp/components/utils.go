package components

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	dt "claw-destine.com/camboose/core/datatypes"
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

const cvSpecs CambView = "specs"
const cvTasks CambView = "tasks"
const cvProjects CambView = "projects"
const cvRecipies CambView = "recipies"

func setViewCookie(view CambView, w http.ResponseWriter) {
	viewCookie := http.Cookie{
		Name:     "view",
		Value:    string(view),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
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
	"statsArray": func(rs map[dt.RequirementStatus]int) template.HTMLAttr {
		val, _ := json.Marshal(rs)
		return (template.HTMLAttr)("data-stats=" + string(val))
	},
}
