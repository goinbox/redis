package redis

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	pipe := client.Pipeline(ctx)
	key := "pipeline"

	pipe.Do(ctx, "set", key, "test pipeline")
	pipe.Do(ctx, "get", key)

	replies, err := pipe.Exec(ctx)
	t.Log(err)
	for _, reply := range replies {
		t.Log(reply.Value())
	}
}
