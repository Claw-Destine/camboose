package components

import (
	"io"
	"net/http"
)

func NewTasksHandler() TasksCompHandler {
	return TasksCompHandler{}
}

type TasksCompHandler struct {
}

func (ph TasksCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	setViewCookie(cvTasks, w)
	io.WriteString(w, `<h1 class="title">Tasks</h1>`)
}
