package rest

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/service/recipies"
)

type RestConfig struct {
	Port          string `env:"PORT" envDefault:"8080"`
	BasicUser     string `env:"BASIC_USER"`
	BasicPassword string `env:"BASIC_PASSWORD"`
	RecipiesCtr   recipies.RecipeController
}

func newRecipiesHandler(cfg RestConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if cfg.BasicUser != "" && cfg.BasicPassword != "" {
			user, password, ok := r.BasicAuth()
			if !ok || user != cfg.BasicUser || password != cfg.BasicPassword {
				w.Header().Set("WWW-Authenticate", `Basic realm="camboose"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.RecipiesCtr.ListRecipies())
	}
}

func StartRestService(cfg RestConfig) {
	http.HandleFunc("/api/recipies", newRecipiesHandler(cfg))

	address := ":" + cfg.Port

	slog.Info("Starting HTTP server", "address", address)

	log.Fatal(http.ListenAndServe(address, nil))
}
