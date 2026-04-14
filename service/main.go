package main

import (
	"log"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"claw-destine.com/camboose/service/recipies"
	"claw-destine.com/camboose/service/rest"
)

type appConfig struct {
	Rest     rest.RestConfig       `envPrefix:"HTTP_"`
	Recipies recipies.RecipeConfig `envPrefix:"RECIPIES_"`
}

type AppComponents struct {
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("Could not load .env.", "error", err)
	}

	var cfg appConfig
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	slog.Info("Recipe Path", "recipe-path", cfg.Recipies.RecipePath)

	rest.StartRestService(cfg.Rest)
}
