package rest

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"claw-destine.com/camboose/service/datatypes"
	"claw-destine.com/camboose/service/graphdb"
	"claw-destine.com/camboose/service/recipies"
)

type RestConfig struct {
	Port          string `env:"PORT" envDefault:"8080"`
	BasicUser     string `env:"BASIC_USER"`
	BasicPassword string `env:"BASIC_PASSWORD"`
	RecipiesCtr   recipies.RecipeController
	ProjectsCtr   graphdb.ProjectController
}

func authorizeRequest(cfg RestConfig, w http.ResponseWriter, r *http.Request) bool {
	if cfg.BasicUser == "" || cfg.BasicPassword == "" {
		return true
	}

	user, password, ok := r.BasicAuth()
	if !ok || user != cfg.BasicUser || password != cfg.BasicPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="camboose"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func newRecipiesHandler(cfg RestConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !authorizeRequest(cfg, w, r) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.RecipiesCtr.ListRecipies())
	}
}

func newProjectsCollectionHandler(cfg RestConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authorizeRequest(cfg, w, r) {
			return
		}

		if cfg.ProjectsCtr == nil {
			http.Error(w, "project controller not configured", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodGet:
			projects, err := cfg.ProjectsCtr.ListProjects(r.Context())
			if err != nil {
				http.Error(w, "failed to list projects", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projects)
		case http.MethodPost:
			var payload datatypes.Project
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid json payload", http.StatusBadRequest)
				return
			}

			project, err := cfg.ProjectsCtr.CreateProject(r.Context(), payload)
			if err != nil {
				if errors.Is(err, graphdb.ErrProjectExists) {
					http.Error(w, "project already exists", http.StatusConflict)
					return
				}
				if errors.Is(err, graphdb.ErrProjectInvalid) {
					http.Error(w, "invalid project payload", http.StatusBadRequest)
					return
				}

				http.Error(w, "failed to create project", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(project)
		default:
			w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ", "))
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func newProjectsItemHandler(cfg RestConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authorizeRequest(cfg, w, r) {
			return
		}

		if cfg.ProjectsCtr == nil {
			http.Error(w, "project controller not configured", http.StatusInternalServerError)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/api/project/")
		name = strings.TrimSpace(name)
		if name == "" || strings.Contains(name, "/") {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			project, found, err := cfg.ProjectsCtr.GetProject(r.Context(), name)
			if err != nil {
				http.Error(w, "failed to get project", http.StatusInternalServerError)
				return
			}
			if !found {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(project)
		case http.MethodDelete:
			deleted, err := cfg.ProjectsCtr.DeleteProject(r.Context(), name)
			if err != nil {
				http.Error(w, "failed to delete project", http.StatusInternalServerError)
				return
			}
			if !deleted {
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodDelete}, ", "))
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func newMux(cfg RestConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/recipies", newRecipiesHandler(cfg))
	mux.HandleFunc("/api/project", newProjectsCollectionHandler(cfg))
	mux.HandleFunc("/api/project/", newProjectsItemHandler(cfg))

	return mux
}

func StartRestService(cfg RestConfig) {
	address := ":" + cfg.Port

	slog.Info("Starting HTTP server", "address", address)

	log.Fatal(http.ListenAndServe(address, newMux(cfg)))
}
