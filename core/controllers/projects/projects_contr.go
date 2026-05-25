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

func (pm *ProjectControler) ProjectStatistics(projects any) (map[string]map[dt.RequirementStatus]int, error) {
	var projectIDs []string
	switch pp := projects.(type) {
	case []string:
		projectIDs = pp
	case []dt.Project:
		projectIDs = make([]string, len(pp))
		for _, p := range pp {
			projectIDs = append(projectIDs, p.Id)
		}
	default:
		slog.Error("ProjectStatistics method accepts only []string and []Project")
		return nil, &dt.WrongTypeError{What: "ProjectStatistics method accepts only []string and []Project"}
	}

	type versionCountRow struct {
		ProjectId string
		Status    dt.RequirementStatus
		Count     int
	}

	var versionCounts []versionCountRow

	if err := pm.Db.Model(&dt.Version{}).
		Select("project_id, status, COUNT(*) as count").
		Where("project_id IN ?", projectIDs).
		Group("project_id, status").
		Scan(&versionCounts).Error; err != nil {
		return nil, err
	}

	countsByProject := make(map[string]map[dt.RequirementStatus]int, len(projectIDs))
	for _, row := range versionCounts {
		if _, ok := countsByProject[row.ProjectId]; !ok {
			countsByProject[row.ProjectId] = make(map[dt.RequirementStatus]int)
		}
		countsByProject[row.ProjectId][row.Status] = row.Count
	}
	return countsByProject, nil
}

func (pm *ProjectControler) ListProjects(filter *dt.QuerySettings) ([]dt.Project, error) {
	var projects []dt.Project
	result := pm.Db.Find(&projects)
	if result.Error != nil {
		return projects, result.Error
	}

	if len(projects) == 0 {
		return projects, nil
	}

	return projects, nil
}
