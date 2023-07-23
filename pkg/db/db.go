package db

import (
	"log"
	"sda-manager/pkg/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Migrate(db_name string) *gorm.DB {
	db := Get(db_name)
	db.AutoMigrate(models.Hymn{}, models.Verse{})
	return db
}

func Get(db_name string) *gorm.DB {
	db,err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	return db
}