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
	db.AutoMigrate(&Transaction{})
}

type User struct {
	ID            string `gorm:"primaryKey"`
	Username      string `gorm:"unique"`
	Email         string `gorm:"unique"`
	AccountNumber string `gorm:"unique"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Transaction struct {
	ID                 string `gorm:"primaryKey"`
	AccountNumber      string
	Date               time.Time
	TransactionDetails string
	ValueDate          time.Time
	WithdrawalAmt      *int
	DepositAmt         *int
	BalanceAmt         int
	CreatedAt          time.Time
}
