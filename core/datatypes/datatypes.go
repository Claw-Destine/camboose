package datatypes

type Project struct {
	ObjectId string `clover:"_id"`
	Name     string `clover:"name"`
	Recipe   string `clover:"recipe"`
}

type Version struct {
	ObjectId  string
	Name      string
	ProjectId string
}
