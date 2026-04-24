package main

import (
	"net/http"

	"claw-destine.com/camboose/webapp/projects"
	"github.com/a-h/templ"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/components", templ.Handler(projects.Projects()))

	http.ListenAndServe(":3000", nil)
}
