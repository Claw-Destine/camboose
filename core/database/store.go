package database

import (
	"fmt"

	dt "claw-destine.com/camboose/core/datatypes"
)

const COL_PROJECTS = "projects"
const COL_VERSION = "versions"
const COL_TASKS = "tasks"

type ErrorType string

const ENTITY_EXISTS ErrorType = "Entity exists"

type StoreError struct {
	What       ErrorType
	Collection string
}

func (e StoreError) Error() string {
	return fmt.Sprintf("Storage error: %s, Collection: %s", e.What, e.Collection)
}

type Store interface {
	ProjectStore
	VersionStore
	Close()
}

type ProjectStore interface {
	CreateProject(name string) (*dt.Project, error)
	UpdateProject(project dt.Project)
	GetProject(id string) (*dt.Project, error)
	GetProjects() []dt.Project
	DeleteProject(id string) error
}

type VersionStore interface {
	CreateVersion(name string, projectID string) *dt.Version
}
