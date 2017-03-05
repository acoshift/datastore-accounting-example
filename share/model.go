package share

import "github.com/acoshift/ds"

// Model Kinds
const (
	KindAccount     = "Account"
	KindTransaction = "Transaction"
)

// Account model
type Account struct {
	ds.StringIDModel
	ds.StampModel
	Balance int64
}

// Transaction model
type Transaction struct {
	ds.Model
	ds.StampModel
	AccountID string
	Type      TransactionType
	Amount    int64
}

// TransactionType type
type TransactionType int

// TransactionType values
const (
	_ TransactionType = iota
	TransactionTypeDeposit
	TransactionTypeWithdraw
)

var mapTransactionTypeString = map[TransactionType]string{
	TransactionTypeDeposit:  "deposit",
	TransactionTypeWithdraw: "withdraw",
}

func (t TransactionType) String() string {
	return mapTransactionTypeString[t]
}
