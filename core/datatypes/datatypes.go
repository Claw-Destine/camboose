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
	Name      string
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
	ProjectId string
	Project   Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId  *string
	Parent    *SpecItem  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Children  []SpecItem `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      SpecType
}
