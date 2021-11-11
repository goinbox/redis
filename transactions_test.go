package redis

import (
	"testing"
)

func TestTransactions(t *testing.T) {
	tx := getTestClient().Transactions()
	key := "trans"

	tx.Do("set", key, "test trans")
	tx.Do("get", key)

	replies, err := tx.Exec()
	t.Log(err)
	for _, reply := range replies {
		t.Log(reply.Value())
	}
}
