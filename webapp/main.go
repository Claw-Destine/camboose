package main

import (
	"log"
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/core/controllers/projects"
	"claw-destine.com/camboose/core/database/postgres"
	dt "claw-destine.com/camboose/core/datatypes"
	cmp "claw-destine.com/camboose/webapp/components"
	md "claw-destine.com/camboose/webapp/middleware"
	"github.com/a-h/templ"
	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	cfg, err := env.ParseAs[dt.Config]()
	if err != nil {
		slog.Error("Failed to parse environmental variables")
		log.Panic("Exiting")
	}

	// Init server and static dir
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Create controllers
	db, err := postgres.ConnectToPostgres(cfg.PgConf)
	if err != nil {
		log.Fatal(err)
	}
	postgres.MigrateDatabase(db)

	projectManager := projects.ProjectManager{Db: db}
	recipiesManager := projects.RecipeManager{Conf: cfg}

	projectsHandler := cmp.NewProjectsHandler(&projectManager, &recipiesManager)

	// Set up routes and inject controlers
	mux.Handle("/components/body", cmp.NewBodyHandler(&projectManager))
	mux.Handle("/components/specs", cmp.NewSpecsHandler(&projectManager))
	mux.Handle("/components/tasks", templ.Handler(cmp.Tasks()))
	mux.Handle("/components/projects", projectsHandler)
	mux.Handle("/components/project/", projectsHandler)
	mux.Handle("/components/recipies", cmp.NewRecipuesHandler(&recipiesManager))

	// Serve
	server := &http.Server{
		Addr:    ":3000",
		Handler: md.WithRequestLogging(mux),
	}

	slog.Info("Starting HTTP server", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server stopped", "error", err)
	}

}
