package specs

import (
	dt "claw-destine.com/camboose/core/datatypes"
	"gorm.io/gorm"
)

type SpecsController struct {
	Db *gorm.DB
}

func (sc *SpecsController) ListSpecItems(projectId string) ([]dt.SpecItem, error) {
	var items []dt.SpecItem
	err := sc.Db.Where("project_id = ?", projectId).Find(&items).Error
	return items, err
}

func (sc *SpecsController) GetSpecItemById(id string) (*dt.SpecItem, error) {
	var item dt.SpecItem
	err := sc.Db.First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (sc *SpecsController) CreateSpecItem(si dt.SpecItem) (*dt.SpecItem, error) {
	if si.Id == "" {
		si.Id = ""
	}
	err := sc.Db.Create(&si).Error
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func (sc *SpecsController) UpdateSpecItem(si dt.SpecItem) (*dt.SpecItem, error) {
	err := sc.Db.Save(&si).Error
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func (sc *SpecsController) DeleteSpecItemById(id string) error {
	return sc.Db.Delete(&dt.SpecItem{}, "id = ?", id).Error
}
