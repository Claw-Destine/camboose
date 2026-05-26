package postgres

import (
	"fmt"

	dt "claw-destine.com/camboose/core/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgres(pgconf dt.PostgresConfig) (*gorm.DB, error) {

	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgconf.Host, pgconf.User, pgconf.Password, pgconf.Db, pgconf.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&dt.Project{}, &dt.Version{}, &dt.Story{})
}
