package components

import (
	"net/http"
)

func NewTasksHandler() TasksCompHandler {
	return TasksCompHandler{}
}

type TasksCompHandler struct {
}

func (ph TasksCompHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	setViewCookie(vTasks, w)
	tasksComponent().Render(r.Context(), w)
}
