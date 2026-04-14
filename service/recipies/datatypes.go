package recipies

type RecipeConfig struct {
	RecipePath string `env:"RECIPE_PATH"`
}

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

type ReqDescription struct {
	Me    ReqEntity
	Child *ReqEntity
}

type ProjectManagement struct {
	Relations []ReqDescription
}

type Recipe struct {
	Name        string
	Description string
}
