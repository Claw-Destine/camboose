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

type Recipe struct {
	Name        string
	Description string
}

type Project struct {
	Name   string `json:"name"`
	Recipe string `json:"recipe"`
}

type Version struct {
	Number int    `json:"number"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status"`
}
