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
	postgres.MigrateColumns(db, &dt.Project{}, &dt.Version{})

	projectManager := projects.ProjectManager{Db: db}

	// Set up routes and inject controlers
	mux.Handle("/components/projects", cmp.NewProjectsHandler(&projectManager))
	mux.Handle("/components/project/", cmp.NewProjectHandler(&projectManager))
	mux.Handle("/components/recipies", cmp.NewRecipuesHandler())
	mux.Handle("/components/tasks", templ.Handler(cmp.Tasks()))

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
