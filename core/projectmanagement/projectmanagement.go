package projectmanagement

import (
	dt "claw-destine.com/camboose/core/datatypes"
	ds "claw-destine.com/camboose/core/documentstore"
)

type ProjectManager struct {
	Db *ds.Store
}

func (pm *ProjectManager) CreateProject(name string) *dt.Project {
	return (*pm.Db).CreateProject(name)
}

func ListProjects() []string {
	mock := [4]string{"p1", "p2", "p3"}
	return mock[1:3]
}
