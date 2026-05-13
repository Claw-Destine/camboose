package cloverstore

import (
	"log/slog"
	"os"

	dt "claw-destine.com/camboose/core/datatypes"
	ds "claw-destine.com/camboose/core/documentstore"
	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

type CloverStore struct {
	db *c.DB
}

func NewCloverStoreConnection(dbname string) (ds.Store, error) {
	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		err := os.MkdirAll(dbname, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	db, err := c.Open(dbname)
	if err != nil {
		return nil, err
	}

	if err := db.CreateCollection(ds.COL_PROJECTS); (err != nil) && (err != c.ErrCollectionExist) {
		return nil, err
	}

	return &CloverStore{db: db}, nil
}

func (cs *CloverStore) Close() {
	cs.db.Close()
}

func (cs *CloverStore) CreateProject(name string) (*dt.Project, error) {
	doc := d.NewDocument()
	doc.Set("name", name)
	id, err := cs.db.InsertOne(ds.COL_PROJECTS, doc)
	if err != nil {
		return nil, err
	}
	p := &dt.Project{ObjectId: id, Name: name}
	return p, nil
}
func (cs *CloverStore) UpdateProject(project dt.Project) {

}
func (cs *CloverStore) GetProject(id string) *dt.Project {
	return &dt.Project{}
}
func (cs *CloverStore) GetProjects() []dt.Project {
	docs, _ := cs.db.FindAll(q.NewQuery(ds.COL_PROJECTS))

	projects := []dt.Project{}
	for _, doc := range docs {
		p := &dt.Project{}
		err := doc.Unmarshal(p)
		if err != nil {
			slog.Error("Failed to unmarshal document", "Error", err)
		}
		projects = append(projects, *p)
	}
	return projects
}
func (c *CloverStore) DeleteProject() []dt.Project {
	return []dt.Project{}
}
func (c *CloverStore) CreateVersion(name string, projectID string) *dt.Version {
	return &dt.Version{}
}
