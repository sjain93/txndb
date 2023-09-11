package main

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
