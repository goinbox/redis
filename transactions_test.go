package redis

import (
	"testing"
)

func TestTransactions(t *testing.T) {
	tx := client.Transactions(ctx)
	key := "trans"

	tx.Do(ctx, "set", key, "test trans")
	tx.Do(ctx, "get", key)

	replies, err := tx.Exec(ctx)
	t.Log(err)
	for _, reply := range replies {
		t.Log(reply.Value())
	}
}
