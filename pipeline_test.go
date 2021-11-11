package redis

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	pipe := getTestClient().Pipeline()
	key := "pipeline"

	pipe.Do("set", key, "test pipeline")
	pipe.Do("get", key)

	replies, err := pipe.Exec()
	t.Log(err)
	for _, reply := range replies {
		t.Log(reply.Value())
	}
}
