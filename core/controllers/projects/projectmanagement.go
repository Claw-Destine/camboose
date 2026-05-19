package projects

import (
	c "context"
	"log/slog"

	dt "claw-destine.com/camboose/core/datatypes"
	"gorm.io/gorm"
)

type ProjectManager struct {
	Db *gorm.DB
}

func (pm *ProjectManager) CreateProject(ctx c.Context, project dt.Project) (*dt.Project, error) {
	ret := &project
	err := gorm.G[dt.Project](pm.Db).Create(ctx, ret)
	if err == nil {
		slog.Info("Created a new project.", "id", ret.Id, "name", ret.Name)
	} else {
		slog.Warn("Failed to create a new project.", "name", ret.Name, "error", err)
	}

	return ret, err
}

func (pm *ProjectManager) GetProjectById(ctx c.Context, id string) (*dt.Project, error) {
	project, err := gorm.G[dt.Project](pm.Db).Where("id = ?", id).First(ctx)
	return &project, err

}

func (pm *ProjectManager) DeleteProject(ctx c.Context, id string) error {
	rows, err := gorm.G[dt.Project](pm.Db).Where("id = ?", id).Delete(ctx)
	if err == nil && rows > 0 {
		slog.Info("Deleted project with id.", "id", id, "rows affected", rows)
	} else {
		slog.Warn("Created a new project.", "id", id, "rows affected", rows, "error", err)
	}
	return err

}

func (pm *ProjectManager) ListProjects() ([]dt.Project, error) {
	var projects []dt.Project
	result := pm.Db.Find(&projects)
	return projects, result.Error
}
