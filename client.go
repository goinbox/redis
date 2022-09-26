package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/goinbox/golog"
)

type dbItem struct {
	config *Config
	db     *redis.Client
}

var dbPool = map[string]*dbItem{}

func RegisterDB(key string, config *Config) {
	db := redis.NewClient(config.Options)

	dbPool[key] = &dbItem{
		config: config,
		db:     db,
	}
}

type Client struct {
	db  *redis.Client
	ctx context.Context

	config *Config
	logger golog.Logger
}

func NewClientFromPool(key string, logger golog.Logger) (*Client, error) {
	item, ok := dbPool[key]
	if !ok {
		return nil, errors.New("DB " + key + " not exist")
	}

	return newClient(item.db, item.config, logger), nil
}

func NewClient(config *Config, logger golog.Logger) *Client {
	return newClient(redis.NewClient(config.Options), config, logger)
}

func newClient(db *redis.Client, config *Config, logger golog.Logger) *Client {
	c := &Client{
		db:     db,
		ctx:    context.Background(),
		config: config,
	}

	if logger != nil {
		c.logger = logger.With(&golog.Field{
			Key:   config.LogFieldKeyAddr,
			Value: config.Addr,
		})
	}

	return c
}

func (c *Client) SetLogger(logger golog.Logger) *Client {
	if logger != nil {
		c.logger = logger.With(&golog.Field{
			Key:   c.config.LogFieldKeyAddr,
			Value: c.config.Addr,
		})
	}

	return c
}

func (c *Client) Do(args ...interface{}) *Reply {
	cmd := redis.NewCmd(c.ctx, args...)

	c.log(cmd)

	err := c.db.Process(c.ctx, cmd)

	return &Reply{
		cmd: cmd,
		Err: err,
	}
}

func (c *Client) Pipeline() *Pipeline {
	if c.logger != nil {
		c.logger.Info("start pipeline")
	}

	return &Pipeline{
		pipe: c.db.Pipeline(),
		ctx:  c.ctx,

		logger:         c.logger,
		logFieldKeyCmd: c.config.LogFieldKeyCmd,
	}
}

func (c *Client) Transactions() *Transactions {
	if c.logger != nil {
		c.logger.Info("start trans")
	}

	return &Transactions{
		tx:  c.db.TxPipeline(),
		ctx: c.ctx,

		logger:         c.logger,
		logFieldKeyCmd: c.config.LogFieldKeyCmd,
	}
}

func (c *Client) RunScript(src string, keys []string, args ...interface{}) *Reply {
	script := redis.NewScript(src)

	if c.logger != nil {
		c.logger.Info("run script", []*golog.Field{
			{
				Key:   "src",
				Value: []byte(src),
			},
			{
				Key:   "keys",
				Value: keys,
			},
			{
				Key:   "args",
				Value: args,
			},
		}...)
	}

	cmd := script.Run(c.ctx, c.db, keys, args...)

	return &Reply{
		cmd: cmd,
		Err: nil,
	}
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) log(cmd redis.Cmder) {
	if c.logger == nil {
		return
	}

	c.logger.Info("run cmd", &golog.Field{
		Key:   c.config.LogFieldKeyCmd,
		Value: cmd.String(),
	})
}
