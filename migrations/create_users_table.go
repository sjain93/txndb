package migrations

import (
	"gorm.io/gorm"
)

/*
AutoMigration works for smaller scale projects and Proof of Concepts.
For larger more production serving APIs, versioned migrations are recommended.
*/
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
