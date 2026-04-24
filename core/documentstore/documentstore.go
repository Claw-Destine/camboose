package documentstore

import dt "claw-destine.com/camboose/core/datatypes"

type Store interface {
	Close()
	CreateProject(name string) *dt.Project
	UpdateProject(project dt.Project)
	GetProject(id string) *dt.Project
	GetProjects() []dt.Project
	DeleteProject() []dt.Project
	CreateVersion(name string, projectID string) *dt.Version
}
