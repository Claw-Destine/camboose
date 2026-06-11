package components

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/core/controllers/projects"
	dt "claw-destine.com/camboose/core/datatypes"
)

func NewRecipuesHandler(recipeManager *projects.RecipeController) RecipiesCompHandler {
	rh := RecipiesCompHandler{recipeManager: recipeManager}
	tpl := `<h1 class="title">Recipies</h1>
	{{range .Recipies}}<p>{{.Id}}, {{.Description}}</p>{{end}}`
	t, err := template.New("projects").Funcs(funcMap).Parse(tpl)
	if err != nil {
		slog.Error("Cannot parse template", "err", err)
		log.Panic("exiting")
	}
	rh.recipiesTpl = t
	return rh
}

type recipiesData struct {
	Recipies []dt.Recipe
}

type RecipiesCompHandler struct {
	recipeManager *projects.RecipeController
	recipiesTpl   *template.Template
}

func (ph RecipiesCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setViewCookie(vRecipies, w)
	ph.displayRecipiesView(w, r)
}

func (ph RecipiesCompHandler) displayRecipiesView(w http.ResponseWriter, r *http.Request) {
	recipies, err := ph.recipeManager.ListRecipes()
	if err != nil {
		slog.Error("Failed to load recipies", "reason", err)
	}
	data := recipiesData{Recipies: recipies}
	if err := ph.recipiesTpl.Execute(w, data); err != nil {
		slog.Error("Cannot render recipies component", "err", err)
		http.Error(w, "failed to recipies body", http.StatusInternalServerError)
	}
}
