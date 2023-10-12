package redis

import (
	"encoding/json"
	"testing"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
)

var ctx pcontext.Context
var client *Client

func init() {
	w, _ := golog.NewFileWriter("/dev/stdout", 0)
	logger := golog.NewSimpleLogger(w, golog.NewSimpleFormater())
	ctx = pcontext.NewSimpleContext(nil, logger)

	config := NewConfig("127.0.0.1", "123", 6379)
	client = NewClient(config)
}

func TestClientDo(t *testing.T) {
	key := "test"
	r := client.Do(ctx, "get", key)
	t.Log("key not exist", r.Nil())

	r = client.Do(ctx, "set", key, "test redis client")
	t.Log("set", r.Err)

	r = client.Do(ctx, "get", key)
	t.Log("get", r.Value())

	client.Do(ctx, "del", key)
}

func TestJson(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	bs, _ := json.Marshal(&person{
		Name: "zhangsan",
		Age:  10,
	})

	key := "test"
	r := client.Do(ctx, "set", key, bs)
	t.Log("set", r.Err, bs)

	r = client.Do(ctx, "get", key)
	t.Log(r.Value())
	p := new(person)
	s, _ := r.String()
	err := json.Unmarshal([]byte(s), p)
	t.Log("get", err, p)
}

func TestRunScript(t *testing.T) {
	src := `
local key = KEYS[1]
local change = ARGV[1]

local value = redis.call("GET", key)
if not value then
  value = 0
end

value = value + change
redis.call("SET", key, value)

return value
`
	keys := []string{"my_counter"}
	values := []interface{}{+1}

	reply := client.RunScript(ctx, src, keys, values...)
	t.Log(reply.Int())
}

func TestClose(t *testing.T) {
	RegisterDB("t1", NewConfig("127.0.0.1", "123", 6379))
	RegisterDB("t2", NewConfig("127.0.0.1", "123", 6379))
	c1, _ := NewClientFromPool("t1")
	c2, _ := NewClientFromPool("t2")

	err := c1.Do(ctx, "get", "a").Err
	t.Log(err)
	err = c2.Do(ctx, "get", "a").Err
	t.Log(err)

	_ = c1.Close(ctx)
	err = c2.Do(ctx, "get", "a").Err
	t.Log(err)
}
