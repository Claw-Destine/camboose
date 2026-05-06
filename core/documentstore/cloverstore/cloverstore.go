package cloverstore

import (
	dt "claw-destine.com/camboose/core/datatypes"
	ds "claw-destine.com/camboose/core/documentstore"
	c "github.com/ostafen/clover/v2"
)

type CloverStore struct {
	db *c.DB
}

func CreateCloverStore(dbname string) ds.ProjectStore {
	db, _ := c.Open(dbname)
	return &CloverStore{db: db}
}

func (c *CloverStore) CreateProject(name string) *dt.Project
func (c *CloverStore) UpdateProject(project dt.Project)
func (c *CloverStore) GetProject(id string) *dt.Project
func (c *CloverStore) GetProjects() []dt.Project
func (c *CloverStore) DeleteProject() []dt.Project
func (c *CloverStore) CreateVersion(name string, projectID string) *dt.Version
func (c *CloverStore) Close()
