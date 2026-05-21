package projects

import (
	c "context"
	"errors"
	"fmt"
	"log/slog"

	dt "claw-destine.com/camboose/core/datatypes"
	"gorm.io/gorm"
)

type ProjectControler struct {
	Db *gorm.DB
}

func (pm *ProjectControler) CreateProject(ctx c.Context, project dt.Project) (*dt.Project, error) {
	ret := &project
	err := gorm.G[dt.Project](pm.Db).Create(ctx, ret)
	if err == nil {
		slog.Info("Created a new project.", "id", ret.Id, "name", ret.Name)
	} else {
		slog.Warn("Failed to create a new project.", "name", ret.Name, "error", err)
	}

	return ret, err
}

func (pm *ProjectControler) GetProjectById(ctx c.Context, id string) (*dt.Project, error) {
	project, err := gorm.G[dt.Project](pm.Db).Where("id = ?", id).First(ctx)
	return &project, err

}

func (pm *ProjectControler) DeleteProject(ctx c.Context, id string) error {
	rows, err := gorm.G[dt.Project](pm.Db).Where("id = ?", id).Delete(ctx)
	if err == nil && rows > 0 {
		slog.Info("Deleted project with id.", "id", id, "rows affected", rows)
	} else {
		slog.Warn("Created a new project.", "id", id, "rows affected", rows, "error", err)
	}
	return err

}

func (pm *ProjectControler) UpdateProject(ctx c.Context, project dt.Project) error {
	if project.Id == "" {
		return errors.New("project id is required")
	}

	result := pm.Db.WithContext(ctx).
		Model(&dt.Project{}).
		Where("id = ?", project.Id).
		Updates(map[string]any{"recipe": project.Recipe})

	if result.Error != nil {
		slog.Warn("Failed to update project.", "id", project.Id, "recipe", project.Recipe, "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := fmt.Errorf("project with id %q not found", project.Id)
		slog.Warn("Failed to update project.", "id", project.Id, "recipe", project.Recipe, "error", err)
		return err
	}

	slog.Info("Updated project.", "id", project.Id, "recipe", project.Recipe, "rows affected", result.RowsAffected)
	return nil
}

type ListProjectsFilter struct {
	dt.Pagination
	dt.Ordering
}

func (pm *ProjectControler) ListProjects(filter *ListProjectsFilter) ([]dt.Project, error) {
	// todo - implement filter logic
	var projects []dt.Project
	result := pm.Db.Find(&projects)
	return projects, result.Error
}
