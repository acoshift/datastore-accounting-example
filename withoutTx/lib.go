package withoutTx

import (
	"context"
	"fmt"

	"github.com/acoshift/datastore-accounting-example/share"
	"github.com/acoshift/ds"
)

var (
	ctx    = context.Background()
	client *ds.Client
)

// SetClient sets client
func SetClient(c *ds.Client) {
	client = c
}

// Deposit deposits money to an account
func Deposit(accID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	var accs []*share.Account
	err := client.GetByNames(ctx, share.KindAccount, []string{accID, "Cash"}, &accs)
	if err != nil {
		return fmt.Errorf("get accounts error; %v", err)
	}
	accs[0].Balance += amount
	accs[1].Balance -= amount
	trans := []*share.Transaction{
		{Type: share.TransactionTypeDeposit, AccountID: accID, Amount: amount},
		{Type: share.TransactionTypeWithdraw, AccountID: "Cash", Amount: -amount},
	}
	for _, tran := range trans {
		tran.NewKey(share.KindTransaction)
	}
	err = client.SaveModels(ctx, "", []interface{}{accs[0], accs[1], trans[0], trans[1]})
	if err != nil {
		return fmt.Errorf("save models error; %v", err)
	}
	return nil
}

// Withdraw withdraws money from an account
func Withdraw(accID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	var accs []*share.Account
	err := client.GetByNames(ctx, share.KindAccount, []string{accID, "Cash"}, &accs)
	if err != nil {
		return fmt.Errorf("get accounts error; %v", err)
	}

	if accs[0].Balance < amount {
		return fmt.Errorf("balance of account %s is %d, not enough for withdraw %d", accID, accs[0].Balance, amount)
	}

	accs[0].Balance -= amount
	accs[1].Balance += amount
	trans := []*share.Transaction{
		{Type: share.TransactionTypeWithdraw, AccountID: accID, Amount: -amount},
		{Type: share.TransactionTypeDeposit, AccountID: "Cash", Amount: amount},
	}
	for _, tran := range trans {
		tran.NewKey(share.KindTransaction)
	}
	err = client.SaveModels(ctx, "", []interface{}{accs[0], accs[1], trans[0], trans[1]})
	if err != nil {
		return fmt.Errorf("save models error; %v", err)
	}
	return nil
}

// Transfer transfers money from an account to another account
func Transfer(from, to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	var accs []*share.Account
	err := client.GetByNames(ctx, share.KindAccount, []string{from, to}, &accs)
	if err != nil {
		return fmt.Errorf("get accounts error; %v", err)
	}

	if accs[0].Balance < amount {
		return fmt.Errorf("balance of account %s is %d, not enough for withdraw %d", from, accs[0].Balance, amount)
	}

	accs[0].Balance -= amount
	accs[1].Balance += amount
	trans := []*share.Transaction{
		{Type: share.TransactionTypeWithdraw, AccountID: from, Amount: -amount},
		{Type: share.TransactionTypeDeposit, AccountID: to, Amount: amount},
	}
	for _, tran := range trans {
		tran.NewKey(share.KindTransaction)
	}
	err = client.SaveModels(ctx, "", []interface{}{accs[0], accs[1], trans[0], trans[1]})
	if err != nil {
		return fmt.Errorf("save models error; %v", err)
	}
	return nil
}
