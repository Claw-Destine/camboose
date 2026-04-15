package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"claw-destine.com/camboose/service/graphdb"
	"claw-destine.com/camboose/service/recipies"
	"claw-destine.com/camboose/service/rest"
)

type appConfig struct {
	Rest     rest.RestConfig       `envPrefix:"HTTP_"`
	GraphDB  graphdb.GraphDBConfig `envPrefix:"GRAPHDB_"`
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

	cfg.Rest.RecipiesCtr = recipies.NewRecipeController(cfg.Recipies)
	projectController, err := graphdb.NewNeo4jProjectController(context.Background(), cfg.GraphDB)
	if err != nil {
		log.Fatal(err)
	}
	defer projectController.Close(context.Background())
	cfg.Rest.ProjectsCtr = projectController

	rest.StartRestService(cfg.Rest)
}
