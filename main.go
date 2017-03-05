package main

import (
	"context"
	"log"
	"sync"

	"github.com/acoshift/datastore-accounting-example/share"
	"github.com/acoshift/datastore-accounting-example/withTx"
	"github.com/acoshift/datastore-accounting-example/withoutTx"
	"github.com/acoshift/ds"
)

const projectID = "acoshift-test"

func main() {
	ctx := context.Background()
	client, err := ds.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("new datastore client error; %v", err)
	}

	withoutTx.SetClient(client)
	withTx.SetClient(client)

	setupInitial := func() {
		cash := &share.Account{Balance: -20000}
		cash.SetNameID(share.KindAccount, "Cash")
		acc1 := &share.Account{Balance: 10000}
		acc1.SetNameID(share.KindAccount, "Acc1")
		acc2 := &share.Account{Balance: 10000}
		acc2.SetNameID(share.KindAccount, "Acc2")
		err = client.SaveModels(ctx, "", []interface{}{cash, acc1, acc2})
		if err != nil {
			log.Fatalf("save initial value error; %v", err)
		}
	}

	run := func(actions []func() error) {
		wg := &sync.WaitGroup{}
		wg.Add(len(actions))
		for i, action := range actions {
			i, action := i, action
			go func() {
				defer wg.Done()
				err := action()
				if err != nil {
					log.Printf("Execute action %d error; %v", i, err)
					return
				}
				log.Printf("Execute action %d completed", i)
			}()
		}
		wg.Wait()
	}

	// withoutTx
	{
		log.Println("Start withoutTx")
		setupInitial()
		actions := []func() error{
			func() error { return withoutTx.Deposit("Acc1", 1000) },
			func() error { return withoutTx.Deposit("Acc1", 2000) },
			func() error { return withoutTx.Deposit("Acc2", 5000) },
			func() error { return withoutTx.Withdraw("Acc1", 1000) },
			func() error { return withoutTx.Withdraw("Acc2", 1000) },
			func() error { return withoutTx.Transfer("Acc1", "Acc2", 2000) },
			func() error { return withoutTx.Transfer("Acc2", "Acc1", 4000) },
			func() error { return withoutTx.Deposit("Acc1", 10000) },
			func() error { return withoutTx.Withdraw("Acc1", 9000) },
		}
		run(actions)
		share.Validate(client)
		share.Cleanup(client)
		log.Println("Done withoutTx")
	}

	log.Println("-----------------------------------")

	// withTx
	{
		log.Println("Start withTx")
		setupInitial()
		actions := []func() error{
			func() error { return withTx.Deposit("Acc1", 1000) },
			func() error { return withTx.Deposit("Acc1", 2000) },
			func() error { return withTx.Deposit("Acc2", 5000) },
			func() error { return withTx.Withdraw("Acc1", 1000) },
			func() error { return withTx.Withdraw("Acc2", 1000) },
			func() error { return withTx.Transfer("Acc1", "Acc2", 2000) },
			func() error { return withTx.Transfer("Acc2", "Acc1", 4000) },
			func() error { return withTx.Deposit("Acc1", 10000) },
			func() error { return withTx.Withdraw("Acc1", 9000) },
		}
		run(actions)
		share.Validate(client)
		share.Cleanup(client)
		log.Println("Done withTx")
	}
}
