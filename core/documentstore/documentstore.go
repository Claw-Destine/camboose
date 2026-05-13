package documentstore

import dt "claw-destine.com/camboose/core/datatypes"

const COL_PROJECTS = "projects"
const COL_VERSION = "versions"
const COL_TASKS = "tasks"

type Store interface {
	ProjectStore
	VersionStore
	Close()
}

type ProjectStore interface {
	CreateProject(name string) (*dt.Project, error)
	UpdateProject(project dt.Project)
	GetProject(id string) *dt.Project
	GetProjects() []dt.Project
	DeleteProject() []dt.Project
}

type VersionStore interface {
	CreateVersion(name string, projectID string) *dt.Version
}
