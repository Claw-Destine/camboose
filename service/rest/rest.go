package rest

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"
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

		path := strings.TrimPrefix(r.URL.Path, "/api/project/")
		path = strings.Trim(path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		segments := strings.Split(path, "/")
		name := strings.TrimSpace(segments[0])
		if name == "" {
			http.NotFound(w, r)
			return
		}

		if len(segments) == 1 {
			handleProjectItemRequest(cfg, w, r, name)
			return
		}

		if len(segments) == 2 && segments[1] == "version" {
			handleVersionCollectionRequest(cfg, w, r, name)
			return
		}

		if len(segments) == 3 && segments[1] == "version" {
			number, err := strconv.Atoi(segments[2])
			if err != nil {
				http.Error(w, "invalid version number", http.StatusBadRequest)
				return
			}

			handleVersionItemRequest(cfg, w, r, name, number)
			return
		}

		http.NotFound(w, r)
	}
}

func handleProjectItemRequest(cfg RestConfig, w http.ResponseWriter, r *http.Request, name string) {
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

func handleVersionCollectionRequest(cfg RestConfig, w http.ResponseWriter, r *http.Request, projectName string) {
	switch r.Method {
	case http.MethodGet:
		versions, err := cfg.ProjectsCtr.ListVersions(r.Context(), projectName)
		if err != nil {
			if errors.Is(err, graphdb.ErrProjectNotFound) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(err, graphdb.ErrVersionInvalid) {
				http.Error(w, "invalid version payload", http.StatusBadRequest)
				return
			}

			http.Error(w, "failed to list versions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versions)
	case http.MethodPost:
		var payload datatypes.Version
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}

		version, err := cfg.ProjectsCtr.CreateVersion(r.Context(), projectName, payload)
		if err != nil {
			if errors.Is(err, graphdb.ErrProjectNotFound) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(err, graphdb.ErrVersionExists) {
				http.Error(w, "version already exists", http.StatusConflict)
				return
			}
			if errors.Is(err, graphdb.ErrVersionInvalid) {
				http.Error(w, "invalid version payload", http.StatusBadRequest)
				return
			}

			http.Error(w, "failed to create version", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(version)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ", "))
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleVersionItemRequest(cfg RestConfig, w http.ResponseWriter, r *http.Request, projectName string, number int) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	version, found, err := cfg.ProjectsCtr.GetVersion(r.Context(), projectName, number)
	if err != nil {
		if errors.Is(err, graphdb.ErrProjectNotFound) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, graphdb.ErrVersionInvalid) {
			http.Error(w, "invalid version number", http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to get version", http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
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
