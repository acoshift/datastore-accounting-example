package share

import (
	"context"
	"log"

	"github.com/acoshift/ds"
)

// Cleanup removes all account and transaction entities
func Cleanup(client *ds.Client) {
	log.Println("Running Cleanup")
	ctx := context.Background()
	keys, _ := client.QueryKeys(ctx, KindAccount)
	client.DeleteMulti(ctx, keys)
	keys, _ = client.QueryKeys(ctx, KindTransaction)
	client.DeleteMulti(ctx, keys)
}
