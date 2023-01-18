package config

import (
	"golang/golang-skeleton/entity"

	"gorm.io/gorm"
)

func migrationTable(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
}
