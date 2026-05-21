package specs

import "gorm.io/gorm"

type SpecsController struct {
	Db *gorm.DB
}
