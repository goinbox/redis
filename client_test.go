package redis

import (
	"encoding/json"
	"testing"

	"github.com/goinbox/golog"
)

func TestClientDo(t *testing.T) {
	client := getTestClient()

	key := "test"
	r := client.Do("get", key)
	t.Log("key not exist", r.Nil())

	r = client.Do("set", key, "test redis client")
	t.Log("set", r.Err)

	r = client.Do("get", key)
	t.Log("get", r.Value())

	client.Do("del", key)
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

	client := getTestClient()
	key := "test"
	r := client.Do("set", key, bs)
	t.Log("set", r.Err, bs)

	r = client.Do("get", key)
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

	reply := getTestClient().RunScript(src, keys, values...)
	t.Log(reply.Int())
}

func getTestClient() *Client {
	w, _ := golog.NewFileWriter("/dev/stdout", 0)
	logger := golog.NewSimpleLogger(w, golog.NewSimpleFormater())

	config := NewConfig("127.0.0.1", "123", 6379)

	return NewClient(config, logger)
}

func TestClose(t *testing.T) {
	RegisterDB("t1", NewConfig("127.0.0.1", "123", 6379))
	RegisterDB("t2", NewConfig("127.0.0.1", "123", 6379))
	c1, _ := NewClientFromPool("t1", &golog.NoopLogger{})
	c2, _ := NewClientFromPool("t2", &golog.NoopLogger{})

	err := c1.Do("get", "a").Err
	t.Log(err)
	err = c2.Do("get", "a").Err
	t.Log(err)

	_ = c1.Close()
	err = c2.Do("get", "a").Err
	t.Log(err)
}
