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

type Recipe struct {
	Id           string
	Description  string
	Version      string
	SpecHierachy []string
	SpecSkill    string
}

type Base struct {
	Id          string `gorm:"primaryKey,type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"uniqueIndex:project_name_idx"`
	Description string
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.Id = uuid.NewString()
	return nil
}

type Project struct {
	Base
	Recipe string
}

type RequirementStatus string

const RS_Draft RequirementStatus = "draft"
const RS_New RequirementStatus = "new"
const RS_InProgress RequirementStatus = "progress"
const RS_InReview RequirementStatus = "review"
const RS_Done RequirementStatus = "done"

type Version struct {
	Base
	ProjectId string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status    RequirementStatus
}

type Story struct {
	Base
	VersionId string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status    RequirementStatus
}

type SpecItemType string

const SIT_Acceptance_Cryterion SpecItemType = "acc_crit"
const SIT_Wireframe SpecItemType = "wireframe"
const SIT_PixelPerfect SpecItemType = "pixelperf"
const SIT_ComplianceRefernce SpecItemType = "compliance"

type SpecItem struct {
	Base
	StoryId string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type    SpecItemType
	Status  RequirementStatus
}
