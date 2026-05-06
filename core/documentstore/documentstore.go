package documentstore

import dt "claw-destine.com/camboose/core/datatypes"

type StoreFactory interface {
	ProjectStore() *ProjectStore
	VersionStore()
	Close()
}

type ProjectStore interface {
	CreateProject(name string) *dt.Project
	UpdateProject(project dt.Project)
	GetProject(id string) *dt.Project
	GetProjects() []dt.Project
	DeleteProject() []dt.Project
}

type VersionStore interface {
	CreateVersion(name string, projectID string) *dt.Version
}
