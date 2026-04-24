package main

import (
	"net/http"

	"claw-destine.com/camboose/webapp/projects"
	"claw-destine.com/camboose/webapp/tasks"
	"github.com/a-h/templ"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/components/projects", templ.Handler(projects.Projects()))
	http.Handle("/components/tasks", templ.Handler(tasks.Tasks()))

	http.ListenAndServe(":3000", nil)
}
