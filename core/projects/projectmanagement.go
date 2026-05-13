package projects

import (
	dt "claw-destine.com/camboose/core/datatypes"
	ds "claw-destine.com/camboose/core/documentstore"
)

type ProjectManager struct {
	Db ds.ProjectStore
}

func (pm *ProjectManager) CreateProject(name string) (*dt.Project, error) {
	p, err := pm.Db.CreateProject(name)
	if err != nil {
		return nil, err
	}
	return p, err
}

func (pm *ProjectManager) GetProjectById(id string) (*dt.Project, error) {
	p, err := pm.Db.GetProject(id)
	if err != nil {
		return nil, err
	}
	return p, err
}

func (pm *ProjectManager) ListProjects() []dt.Project {
	return pm.Db.GetProjects()
}
