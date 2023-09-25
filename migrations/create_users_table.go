package migrations

import (
	"time"

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
	ID        string `gorm:"primaryKey"`
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
