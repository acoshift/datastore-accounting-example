package withTx

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/datastore-accounting-example/share"
	"github.com/acoshift/ds"
)

var (
	ctx    = context.Background()
	client *ds.Client
)

const maxAttempts = 10

// SetClient sets client
func SetClient(c *ds.Client) {
	client = c
}

// Deposit deposits money to an account
func Deposit(accID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	keyCash := datastore.NameKey(share.KindAccount, "Cash", nil)
	keyAcc := datastore.NameKey(share.KindAccount, accID, nil)

	trans := []*share.Transaction{
		{Type: share.TransactionTypeDeposit, AccountID: accID, Amount: amount},
		{Type: share.TransactionTypeWithdraw, AccountID: "Cash", Amount: -amount},
	}

	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		accs := make([]*share.Account, 2)
		err := tx.GetMulti([]*datastore.Key{keyAcc, keyCash}, accs)
		if err != nil {
			return fmt.Errorf("get accounts error; %v", err)
		}
		accs[0].Balance += amount
		accs[1].Balance -= amount
		accs[0].Stamp()
		accs[1].Stamp()
		for _, tran := range trans {
			tran.Stamp()
		}
		_, err = tx.PutMulti([]*datastore.Key{
			keyAcc,
			keyCash,
			datastore.IncompleteKey(share.KindTransaction, nil),
			datastore.IncompleteKey(share.KindTransaction, nil),
		}, []interface{}{accs[0], accs[1], trans[0], trans[1]})
		if err != nil {
			return fmt.Errorf("save models error; %v", err)
		}
		return nil
	}, datastore.MaxAttempts(maxAttempts))
	if err != nil {
		return fmt.Errorf("transaction error; %v", err)
	}
	return nil
}

// Withdraw withdraws money from an account
func Withdraw(accID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	keyCash := datastore.NameKey(share.KindAccount, "Cash", nil)
	keyAcc := datastore.NameKey(share.KindAccount, accID, nil)

	trans := []*share.Transaction{
		{Type: share.TransactionTypeWithdraw, AccountID: accID, Amount: -amount},
		{Type: share.TransactionTypeDeposit, AccountID: "Cash", Amount: amount},
	}

	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		accs := make([]*share.Account, 2)
		err := tx.GetMulti([]*datastore.Key{keyAcc, keyCash}, accs)
		if err != nil {
			return fmt.Errorf("get accounts error; %v", err)
		}
		if accs[0].Balance < amount {
			return fmt.Errorf("balance of account %s is %d, not enough for withdraw %d", accID, accs[0].Balance, amount)
		}
		accs[0].Balance -= amount
		accs[1].Balance += amount
		accs[0].Stamp()
		accs[1].Stamp()
		for _, tran := range trans {
			tran.Stamp()
		}
		_, err = tx.PutMulti([]*datastore.Key{
			keyAcc,
			keyCash,
			datastore.IncompleteKey(share.KindTransaction, nil),
			datastore.IncompleteKey(share.KindTransaction, nil),
		}, []interface{}{accs[0], accs[1], trans[0], trans[1]})
		if err != nil {
			return fmt.Errorf("save models error; %v", err)
		}
		return nil
	}, datastore.MaxAttempts(maxAttempts))
	if err != nil {
		return fmt.Errorf("transaction error; %v", err)
	}
	return nil
}

// Transfer transfers money from an account to another account
func Transfer(from, to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount; %d", amount)
	}

	keyFrom := datastore.NameKey(share.KindAccount, from, nil)
	keyTo := datastore.NameKey(share.KindAccount, to, nil)

	trans := []*share.Transaction{
		{Type: share.TransactionTypeWithdraw, AccountID: from, Amount: -amount},
		{Type: share.TransactionTypeDeposit, AccountID: to, Amount: amount},
	}

	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		accs := make([]*share.Account, 2)
		err := tx.GetMulti([]*datastore.Key{keyFrom, keyTo}, accs)
		if err != nil {
			return fmt.Errorf("get accounts error; %v", err)
		}
		if accs[0].Balance < amount {
			return fmt.Errorf("balance of account %s is %d, not enough for withdraw %d", from, accs[0].Balance, amount)
		}
		accs[0].Balance -= amount
		accs[1].Balance += amount
		accs[0].Stamp()
		accs[1].Stamp()
		for _, tran := range trans {
			tran.Stamp()
		}
		_, err = tx.PutMulti([]*datastore.Key{
			keyFrom,
			keyTo,
			datastore.IncompleteKey(share.KindTransaction, nil),
			datastore.IncompleteKey(share.KindTransaction, nil),
		}, []interface{}{accs[0], accs[1], trans[0], trans[1]})
		if err != nil {
			return fmt.Errorf("save models error; %v", err)
		}
		return nil
	}, datastore.MaxAttempts(maxAttempts))
	if err != nil {
		return fmt.Errorf("transaction error; %v", err)
	}
	return nil
}
