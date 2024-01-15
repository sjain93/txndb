package txns

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sjain93/userservice/api/common"
	"gorm.io/gorm"
)

type TxnRepoManager interface {
	Create(txn *Transaction) error
}

type txnRepository struct {
	DB *gorm.DB
}

func NewTxnRepository(db *gorm.DB) (TxnRepoManager, error) {
	if db != nil {
		return &txnRepository{
			DB: db,
		}, nil
	}

	return &txnRepository{}, common.ErrNoDatastore
}

func (r *txnRepository) Create(txn *Transaction) error {
	if r.DB != nil {
		if err := r.DB.Create(txn).Error; err != nil {
			// this is a GORM implementation detail
			var perr *pgconn.PgError
			if ok := errors.As(err, &perr); ok && perr.Code == common.UniqueViolationErr {
				return common.ErrUniqueKeyViolated
			} else {
				return err
			}
		}
	}
	return nil
}
