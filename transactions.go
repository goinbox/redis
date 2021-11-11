package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/goinbox/golog"
)

type Transactions struct {
	tx  redis.Pipeliner
	ctx context.Context

	logger         golog.Logger
	logFieldKeyCmd string
}

func (t *Transactions) Do(args ...interface{}) {
	cmd := redis.NewCmd(t.ctx, args...)

	t.log(cmd)

	_ = t.tx.Process(t.ctx, cmd)
}

func (t *Transactions) Exec() ([]*Reply, error) {
	if t.logger != nil {
		t.logger.Info("exec trans")
	}

	cmds, err := t.tx.Exec(t.ctx)
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

func (t *Transactions) Discard() error {
	if t.logger != nil {
		t.logger.Info("discard trans")
	}

	return t.tx.Discard()
}

func (t *Transactions) log(cmd redis.Cmder) {
	if t.logger == nil {
		return
	}

	t.logger.Info("trans cmd", &golog.Field{
		Key:   t.logFieldKeyCmd,
		Value: cmd.String(),
	})
}
