package redis

import (
	"errors"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
	"github.com/redis/go-redis/v9"
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
	db *redis.Client

	config *Config
}

func NewClientFromPool(key string) (*Client, error) {
	item, ok := dbPool[key]
	if !ok {
		return nil, errors.New("DB " + key + " not exist")
	}

	return newClient(item.db, item.config), nil
}

func NewClient(config *Config) *Client {
	return newClient(redis.NewClient(config.Options), config)
}

func newClient(db *redis.Client, config *Config) *Client {
	return &Client{
		db:     db,
		config: config,
	}
}

func (c *Client) Do(ctx pcontext.Context, args ...interface{}) *Reply {
	cmd := redis.NewCmd(ctx, args...)

	c.log(ctx.Logger(), cmd)

	err := c.db.Process(ctx, cmd)

	return &Reply{
		cmd: cmd,
		Err: err,
	}
}

func (c *Client) Pipeline(ctx pcontext.Context) *Pipeline {
	logger := ctx.Logger()
	if logger != nil {
		logger.Info("start pipeline")
	}

	return &Pipeline{
		pipe: c.db.Pipeline(),

		logFieldKeyCmd: c.config.LogFieldKeyCmd,
	}
}

func (c *Client) Transactions(ctx pcontext.Context) *Transactions {
	logger := ctx.Logger()
	if logger != nil {
		logger.Info("start trans")
	}

	return &Transactions{
		tx: c.db.TxPipeline(),

		logFieldKeyCmd: c.config.LogFieldKeyCmd,
	}
}

func (c *Client) RunScript(ctx pcontext.Context, src string, keys []string, args ...interface{}) *Reply {
	script := redis.NewScript(src)

	logger := ctx.Logger()
	if logger != nil {
		logger.Info("run script", []*golog.Field{
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

	cmd := script.Run(ctx, c.db, keys, args...)

	return &Reply{
		cmd: cmd,
		Err: nil,
	}
}

func (c *Client) Close(ctx pcontext.Context) error {
	logger := ctx.Logger()
	if logger != nil {
		logger.Info("close db")
	}

	return c.db.Close()
}

func (c *Client) log(logger golog.Logger, cmd redis.Cmder) {
	if logger == nil {
		return
	}

	logger.Info("run cmd", &golog.Field{
		Key:   c.config.LogFieldKeyCmd,
		Value: cmd.String(),
	})
}
