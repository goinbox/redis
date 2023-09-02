package redis

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
	"github.com/redis/go-redis/v9"
)

type Pipeline struct {
	pipe redis.Pipeliner

	logFieldKeyCmd string
}

func (p *Pipeline) Do(ctx pcontext.Context, args ...interface{}) {
	cmd := redis.NewCmd(ctx, args...)

	p.log(ctx.Logger(), cmd)

	_ = p.pipe.Process(ctx, cmd)
}

func (p *Pipeline) Exec(ctx pcontext.Context) ([]*Reply, error) {
	ctx.Logger().Info("exec pipeline")

	cmds, err := p.pipe.Exec(ctx)
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

func (p *Pipeline) Discard(ctx pcontext.Context) {
	ctx.Logger().Info("discard pipeline")

	p.pipe.Discard()
}

func (p *Pipeline) log(logger golog.Logger, cmd redis.Cmder) {
	logger.Info("pipeline cmd", &golog.Field{
		Key:   p.logFieldKeyCmd,
		Value: cmd.String(),
	})
}
