package redis

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
	"github.com/redis/go-redis/v9"
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
	ctx.Logger().Info("exec trans")

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

func (t *Transactions) Discard(ctx pcontext.Context) {
	ctx.Logger().Info("discard trans")

	t.tx.Discard()
}

func (t *Transactions) log(logger golog.Logger, cmd redis.Cmder) {
	logger.Info("trans cmd", &golog.Field{
		Key:   t.logFieldKeyCmd,
		Value: cmd.String(),
	})
}
