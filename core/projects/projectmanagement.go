package projects

import (
	dt "claw-destine.com/camboose/core/datatypes"
	ds "claw-destine.com/camboose/core/store"
)

type ProjectManager struct {
	Db ds.ProjectStore
}

func (pm *ProjectManager) CreateProject(name string) (*dt.Project, error) {
	return pm.Db.CreateProject(name)

}

func (pm *ProjectManager) GetProjectById(id string) (*dt.Project, error) {
	return pm.Db.GetProject(id)

}

func (pm *ProjectManager) DeleteProject(id string) error {
	return pm.Db.DeleteProject(id)

}

func (pm *ProjectManager) ListProjects() []dt.Project {
	return pm.Db.GetProjects()
}
