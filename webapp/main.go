package main

import (
	"log"
	"log/slog"
	"net/http"

	"claw-destine.com/camboose/core/documentstore/cloverstore"
	"claw-destine.com/camboose/core/projects"
	cmp "claw-destine.com/camboose/webapp/components"
	md "claw-destine.com/camboose/webapp/middleware"
	"github.com/a-h/templ"
)

func main() {
	// Init server and static dir
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Create controllers
	db, err := cloverstore.NewCloverStoreConnection("data/clover-store")
	if err != nil {
		log.Fatal(err)
	}
	projectManager := projects.ProjectManager{Db: db}

	// Set up routes and inject controlers
	mux.Handle("/components/projects", cmp.NewProjectsHandler(&projectManager))
	mux.Handle("/components/project/", cmp.NewProjectHandler(&projectManager))
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
