package rest

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/service/recipies"
)

type RestConfig struct {
	Port        string `env:"PORT" envDefault:"8080"`
	RecipiesCtr recipies.RecipeController
}

func StartRestService(cfg RestConfig) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.RecipiesCtr.ListRecipies())
	}
	http.HandleFunc("/api/recipies", handler)

	address := ":" + cfg.Port

	slog.Info("Starting HTTP server", "address", address)

	log.Fatal(http.ListenAndServe(address, nil))
}
