package main

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"claw-destine.com/camboose/service/recipies"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("Could not load .env. If the env variables are not provided by another mean it may result in a crash")
	}

	// parse
	var cfg recipies.RecipeConfig
	err = env.Parse(&cfg)

	// parse with generics
	cfg, err = env.ParseAs[recipies.RecipeConfig]()

	slog.Info("Recipe Path", "recipe-path", cfg.RecipePath)
}
