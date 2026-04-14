package datatypes

type ReqStatus string

const (
	ReqStatusNew       ReqStatus = "new"
	ReqStatusDefined   ReqStatus = "defined"
	ReqStatusDelivered ReqStatus = "delivered"
)

type ReqEntity string

const (
	ReqEntityEpic                ReqEntity = "epic"
	ReqEntityStory               ReqEntity = "story"
	ReqEntityAcceptanceCriterion ReqEntity = "acceptance_cryterion"
)

type DesignEntity string

const (
	DesignEntityView DesignEntity = "view"
)

type Project struct {
	Name   string
	Recipe string
}

type Version struct {
	VersionString string
	Status        ReqStatus
}

type Requirement struct {
	Ver    Version
	Parent *Requirement
	Type   ReqEntity
}

type Design struct {
	Requirement Requirement
	Parent      *Design
	Type        DesignEntity
}
