package txns

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/sjain93/userservice/api/common"
)

var (
	once     sync.Once
	instance *txnService
)

const (
	csvTimeFormat = "2-Jan-06"
)

// service errors
var (
	ErrSvcTxnExists = errors.New("target transaction already exists")
)

type TxnServiceManager interface {
	SeedLocalTransactions() error
	SaveTransaction(tx Transaction) error
}

type txnService struct {
	txnRepo TxnRepoManager
}

func NewTxnService(r TxnRepoManager) TxnServiceManager {
	once.Do(func() {
		instance = &txnService{
			txnRepo: r,
		}
	})
	return instance
}

func (s *txnService) SaveTransaction(tx Transaction) error {
	balance := strconv.FormatInt(*tx.BalanceAmt, 10)
	tx.ID = common.GetMD5HashWithSum(
		tx.AccountNumber + balance + tx.CreatedAt.String(),
	)

	err := s.txnRepo.Create(&tx)
	if err != nil && errors.Is(err, common.ErrUniqueKeyViolated) {
		return ErrSvcTxnExists
	} else if err != nil {
		return err
	}

	return nil
}

func (s *txnService) SeedLocalTransactions() error {
	type cTransaction struct {
		AccountNumber      string `csv:"Account No"`
		Date               string `csv:"DATE"`
		TransactionDetails string `csv:"TRANSACTION DETAILS"`
		ValueDate          string `csv:"VALUE DATE"`
		WithdrawalAmt      string `csv:"WITHDRAWAL AMT"`
		DepositAmt         string `csv:"DEPOSIT AMT"`
		BalanceAmt         string `csv:"BALANCE AMT"`
	}

	in, err := os.Open("./api/txns/bank.csv")
	if err != nil {
		return err
	}
	defer in.Close()

	txns := []*cTransaction{}

	if err := gocsv.UnmarshalFile(in, &txns); err != nil {
		return err
	}
	for _, t := range txns {
		date, err := common.UnmarshalTime(t.Date, csvTimeFormat)
		if err != nil {
			return err
		}

		valdate, err := common.UnmarshalTime(t.ValueDate, csvTimeFormat)
		if err != nil {
			return err
		}

		balanceAmt, err := common.StringToCents(t.BalanceAmt)
		if err != nil {
			return err
		}

		txn := Transaction{
			AccountNumber:      t.AccountNumber,
			Date:               date,
			TransactionDetails: t.TransactionDetails,
			ValueDate:          valdate,
			BalanceAmt:         &balanceAmt,
			CreatedAt:          time.Now(),
		}
		if t.WithdrawalAmt != "" {
			wAmt, err := common.StringToCents(t.WithdrawalAmt)
			if err != nil {
				return err
			}
			txn.WithdrawalAmt = &wAmt
		}
		if t.DepositAmt != "" {
			dAmt, err := common.StringToCents(t.DepositAmt)
			if err != nil {
				return err
			}
			txn.DepositAmt = &dAmt
		}

		if err := s.SaveTransaction(txn); err != nil {
			return err
		}
	}

	return nil
}
