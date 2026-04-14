package projectmanagement

import "claw-destine.com/camboose/service/recipies"

type Project struct {
	Name   string
	Recipe string
}

type Version struct {
	VersionString string
	Status        recipies.ReqStatus
}

type Requirement struct {
	Ver    Version
	Parent *Requirement
	Type   recipies.ReqEntity
}

type Design struct {
	Requirement Requirement
	Parent      *Design
	Type        recipies.DesignEntity
}
