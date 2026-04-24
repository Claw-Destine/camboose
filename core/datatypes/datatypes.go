package datatypes

type Project struct {
	Id     string
	Name   string
	Recipe string
}

type Version struct {
	Id        string
	Name      string
	ProjectId string
}
