package datatypes

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit  int
	Offset int
}

type OrderField string

type Ordering struct {
	Field     OrderField
	Ascending bool
}

type Base struct {
	Id        string `gorm:"primaryKey,type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"uniqueIndex:project_name_idx"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.Id = uuid.NewString()
	return nil
}

type Project struct {
	Base
	Recipe string
}

type SpecType int

const (
	Version            SpecType = iota
	Story              SpecType = 1 << (4 * iota)
	AcceptanceCryteria SpecType = 1 << (4 * iota)
	Leaf               SpecType = 1 << (4 * iota)
)

type SpecItem struct {
	Base
	ProjectId string     `gorm:"constraint:OnUpdate:CqASCADE,OnDelete:CASCADE;"`
	ParentId  *string    `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Children  []SpecItem `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      SpecType
}

type Recipe struct {
	Id           string
	Description  string
	Version      string
	SpecHierachy []string
	SpecSkill    string
}
