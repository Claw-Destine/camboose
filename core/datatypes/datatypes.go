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

type Base struct {
	gorm.Model
	Id        string `gorm:"primaryKey,type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.Id = uuid.NewString()
	return nil
}

type Project struct {
	Base
	Name   string
	Recipe string
}

type Version struct {
	Base
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ProjectId string
	Project   Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
