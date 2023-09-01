package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/goinbox/pcontext"

	"github.com/goinbox/golog"
)

type Transactions struct {
	tx redis.Pipeliner

	logFieldKeyCmd string
}

func (t *Transactions) Do(ctx pcontext.Context, args ...interface{}) {
	cmd := redis.NewCmd(ctx, args...)

	t.log(ctx.Logger(), cmd)

	_ = t.tx.Process(ctx, cmd)
}

func (t *Transactions) Exec(ctx pcontext.Context) ([]*Reply, error) {
	logger := ctx.Logger()
	if logger != nil {
		logger.Info("exec trans")
	}

	cmds, err := t.tx.Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*Reply, len(cmds))
	for i, cmd := range cmds {
		result[i] = &Reply{
			cmd: cmd.(*redis.Cmd),
			Err: nil,
		}
	}

	return result, nil
}

func (t *Transactions) Discard(ctx pcontext.Context) error {
	logger := ctx.Logger()
	if logger != nil {
		logger.Info("discard trans")
	}

	return t.tx.Discard()
}

func (t *Transactions) log(logger golog.Logger, cmd redis.Cmder) {
	if logger == nil {
		return
	}

	logger.Info("trans cmd", &golog.Field{
		Key:   t.logFieldKeyCmd,
		Value: cmd.String(),
	})
}
