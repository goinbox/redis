package redis

import (
	"fmt"
	"io"

	"github.com/garyburd/redigo/redis"

	"github.com/goinbox/golog"
)

type cmdArgs struct {
	cmd  string
	args []interface{}
}

type Client struct {
	config *Config
	logger golog.Logger

	conn      redis.Conn
	connected bool

	pipeCmds []*cmdArgs
}

func NewClient(config *Config, logger golog.Logger) *Client {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	c := &Client{
		config: config,
		logger: logger.With(&golog.Field{
			Key:   config.LogFieldKeyAddr,
			Value: config.Addr,
		}),

		pipeCmds: []*cmdArgs{},
	}

	return c
}

func (c *Client) SetLogger(logger golog.Logger) *Client {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}
	c.logger = logger

	return c
}

func (c *Client) Connected() bool {
	return c.connected
}

func (c *Client) Free() {
	if c.conn != nil {
		_ = c.conn.Close()
	}

	c.connected = false
}

func (c *Client) Connect() error {
	options := []redis.DialOption{
		redis.DialConnectTimeout(c.config.ConnectTimeout),
		redis.DialReadTimeout(c.config.ReadTimeout),
		redis.DialWriteTimeout(c.config.WriteTimeout),
	}

	conn, err := redis.Dial("tcp", c.config.Addr, options...)
	if err != nil {
		return err
	}

	_, err = conn.Do("auth", c.config.Pass)
	if err != nil {
		return err
	}

	c.conn = conn
	c.connected = true

	return nil
}

func (c *Client) Do(cmd string, args ...interface{}) *Reply {
	if !c.connected {
		if err := c.Connect(); err != nil {
			return NewReply(nil, err)
		}
	}

	c.log(cmd, args...)

	return c.do(cmd, args...)
}

func (c *Client) DoWithoutLog(cmd string, args ...interface{}) *Reply {
	if !c.connected {
		if err := c.Connect(); err != nil {
			return NewReply(nil, err)
		}
	}

	return c.do(cmd, args...)
}

func (c *Client) do(cmd string, args ...interface{}) *Reply {
	defer func() {
		c.pipeCmds = []*cmdArgs{}
	}()

	for _, ca := range c.pipeCmds {
		if err := c.conn.Send(ca.cmd, ca.args...); err != nil {
			return NewReply(nil, err)
		}
	}

	reply, err := c.conn.Do(cmd, args...)
	if err != nil {
		if err != io.EOF {
			return NewReply(nil, err)
		}
		if !c.config.TimeoutAutoReconnect {
			return NewReply(nil, err)
		}
		err = c.reconnect()
		if err != nil {
			return NewReply(nil, err)
		}

		for _, ca := range c.pipeCmds {
			if err = c.conn.Send(ca.cmd, ca.args...); err != nil {
				return NewReply(nil, err)
			}
		}
		reply, err = c.conn.Do(cmd, args...)
		if err != nil {
			return NewReply(nil, err)
		}
	}

	return NewReply(reply, err)
}

func (c *Client) Send(cmd string, args ...interface{}) {
	c.log(cmd, args...)
	c.pipeCmds = append(c.pipeCmds, &cmdArgs{cmd, args})
}

func (c *Client) ExecPipelining() ([]*Reply, []int) {
	if !c.connected {
		if err := c.Connect(); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	defer func() {
		c.pipeCmds = []*cmdArgs{}
	}()

	for _, ca := range c.pipeCmds {
		if err := c.conn.Send(ca.cmd, ca.args...); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}
	if err := c.conn.Flush(); err != nil {
		return []*Reply{NewReply(nil, err)}, []int{0}
	}

	reply, err := c.conn.Receive()
	if err != nil {
		if err != io.EOF {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		if !c.config.TimeoutAutoReconnect {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		err = c.reconnect()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}

		for _, ca := range c.pipeCmds {
			if err = c.conn.Send(ca.cmd, ca.args...); err != nil {
				return []*Reply{NewReply(nil, err)}, []int{0}
			}
		}

		if err = c.conn.Flush(); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
		reply, err = c.conn.Receive()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	replies := make([]*Reply, len(c.pipeCmds))
	var errIndexes []int

	replies[0] = NewReply(reply, nil)
	for i := 1; i < len(c.pipeCmds); i++ {
		reply, err := c.conn.Receive()
		replies[i] = NewReply(reply, err)
		if err != nil {
			errIndexes = append(errIndexes, i)
		}
	}

	return replies, errIndexes
}

func (c *Client) BeginTrans() {
	c.Send("multi")
}

func (c *Client) DiscardTrans() error {
	return c.Do("discard").Err
}

func (c *Client) ExecTrans() ([]*Reply, error) {
	reply := c.Do("exec")
	values, err := redis.Values(reply.reply, reply.Err)
	if err != nil {
		return nil, err
	}

	replies := make([]*Reply, len(values))
	for i, value := range values {
		replies[i] = NewReply(value, nil)
	}

	return replies, nil
}

func (c *Client) log(cmd string, args ...interface{}) {
	if len(cmd) == 0 {
		return
	}

	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	c.logger.Info("run cmd", &golog.Field{
		Key:   c.config.LogFieldKeyCmd,
		Value: cmd,
	})
}

func (c *Client) reconnect() error {
	c.Free()

	return c.Connect()
}
