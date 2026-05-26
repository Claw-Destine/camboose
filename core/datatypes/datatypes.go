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

type QuerySettings struct {
	Limit       int
	Offset      int
	Ascending   bool
	OrderFields []string
}

func DefaultQuerySettings() QuerySettings {
	return QuerySettings{
		Limit:     -1,
		Offset:    -1,
		Ascending: false,
	}
}

type Recipe struct {
	Id           string
	Description  string
	Version      string
	SpecHierachy []string
	SpecSkill    string
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
	Description         string
	Recipe              string
	VersionStatusCounts map[RequirementStatus]int `gorm:"-"`
}

type RequirementStatus string

const RS_Draft RequirementStatus = "draft"
const RS_New RequirementStatus = "new"
const RS_InProgress RequirementStatus = "progress"
const RS_InReview RequirementStatus = "review"
const RS_Done RequirementStatus = "done"

var ALL_RS = [...]RequirementStatus{RS_Draft, RS_New, RS_InProgress, RS_InReview, RS_Done}

type Version struct {
	Base
	ProjectId         string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description       string
	Status            RequirementStatus
	StoryStatusCounts map[RequirementStatus]int `gorm:"-"`
}

type Story struct {
	Base
	VersionId   string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description string
	Status      RequirementStatus
}

type SpecItemType string

const SIT_Acceptance_Cryterion SpecItemType = "acc_crit"
const SIT_Wireframe SpecItemType = "wireframe"
const SIT_PixelPerfect SpecItemType = "pixelperf"
const SIT_ComplianceRefernce SpecItemType = "compliance"

type SpecItem struct {
	Base
	StoryId     string `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description string
	Type        SpecItemType
	Status      RequirementStatus
}
