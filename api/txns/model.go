package txns

import "time"

type Transaction struct {
	ID                 string `gorm:"primaryKey"`
	AccountNumber      string
	Date               time.Time
	TransactionDetails string
	ValueDate          time.Time
	WithdrawalAmt      *int64
	DepositAmt         *int64
	BalanceAmt         *int64
	CreatedAt          time.Time
}
