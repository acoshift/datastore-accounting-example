package share

import (
	"context"
	"log"

	"github.com/acoshift/ds"
)

// Validate checks accounts, and transactions
func Validate(client *ds.Client) {
	ctx := context.Background()
	log.Println("validate conflict: start")

	// 1. sum of all transactions must equal to zero
	{
		var trans []*Transaction
		err := client.Query(ctx, KindTransaction, &trans)
		if err != nil {
			log.Printf("query transactions error; %v", err)
			return
		}

		var sum int64
		for _, tran := range trans {
			sum += tran.Amount
		}
		if sum != 0 {
			log.Printf("found conflict: expected sum of all transactions to be 0; got %d", sum)
		} else {
			log.Println("no conflict in transactions: sum of all transactions equal to 0")
		}
	}

	// 2. sum of balance for all accounts must equal to zero
	{
		var accs []*Account
		err := client.Query(ctx, KindAccount, &accs)
		if err != nil {
			log.Printf("query accounts error; %v", err)
			return
		}

		var sum int64
		for _, acc := range accs {
			sum += acc.Balance
		}
		if sum != 0 {
			log.Printf("found conflict: expected sum of all account balances to be 0; got %d", sum)
		} else {
			log.Println("no conflict in accounts: sum of all accounts equal to 0")
		}
	}

	log.Println("validate conflict: done")
}
