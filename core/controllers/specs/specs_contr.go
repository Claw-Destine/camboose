package specs

import (
	"log/slog"

	dt "claw-destine.com/camboose/core/datatypes"
	"gorm.io/gorm"
)

type SpecsController struct {
	Db *gorm.DB
}

func (sc *SpecsController) ListVersions(projectId string) ([]dt.Version, error) {
	var items []dt.Version
	err := sc.Db.Where("project_id = ?", projectId).Find(&items).Error
	return items, err
}

func (sc *SpecsController) GetVersionById(id string) (*dt.Version, error) {
	var item dt.Version
	err := sc.Db.First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (sc *SpecsController) CreateVersion(si dt.Version) (*dt.Version, error) {
	if si.Id == "" {
		si.Id = ""
	}
	err := sc.Db.Create(&si).Error
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func (sc *SpecsController) UpdateVersion(si dt.Version) (*dt.Version, error) {
	err := sc.Db.Save(&si).Error
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func (sc *SpecsController) DeleteVersionById(id string) error {
	return sc.Db.Delete(&dt.Version{}, "id = ?", id).Error
}

func (pm *SpecsController) VersionStatistics(versions any) (map[string]map[dt.RequirementStatus]int, error) {
	var projectIDs []string
	switch pp := versions.(type) {
	case []string:
		projectIDs = pp
	case []dt.Version:
		projectIDs = make([]string, len(pp))
		for _, p := range pp {
			projectIDs = append(projectIDs, p.Id)
		}
	default:
		slog.Error("VersionStatistics method accepts only []string and []Project")
		return nil, &dt.WrongTypeError{What: "VersionStatistics method accepts only []string and []Project"}
	}

	type storyCountRow struct {
		ProjectId string
		Status    dt.RequirementStatus
		Count     int
	}

	var storyCounts []storyCountRow

	if err := pm.Db.Model(&dt.Story{}).
		Select("version_id, status, COUNT(*) as count").
		Where("version_id IN ?", projectIDs).
		Group("version_id, status").
		Scan(&storyCounts).Error; err != nil {
		return nil, err
	}

	countsByVersion := make(map[string]map[dt.RequirementStatus]int, len(projectIDs))
	for _, row := range storyCounts {
		if _, ok := countsByVersion[row.ProjectId]; !ok {
			countsByVersion[row.ProjectId] = make(map[dt.RequirementStatus]int)
		}
		countsByVersion[row.ProjectId][row.Status] = row.Count
	}
	return countsByVersion, nil
}
