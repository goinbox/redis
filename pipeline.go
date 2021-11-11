package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/goinbox/golog"
)

type Pipeline struct {
	pipe redis.Pipeliner
	ctx  context.Context

	logger         golog.Logger
	logFieldKeyCmd string
}

func (p *Pipeline) Do(args ...interface{}) {
	cmd := redis.NewCmd(p.ctx, args...)

	p.log(cmd)

	_ = p.pipe.Process(p.ctx, cmd)
}

func (p *Pipeline) Exec() ([]*Reply, error) {
	if p.logger != nil {
		p.logger.Info("exec pipeline")
	}

	cmds, err := p.pipe.Exec(p.ctx)
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

func (p *Pipeline) Discard() error {
	if p.logger != nil {
		p.logger.Info("discard pipeline")
	}

	return p.pipe.Discard()
}

func (p *Pipeline) log(cmd redis.Cmder) {
	if p.logger == nil {
		return
	}

	p.logger.Info("pipeline cmd", &golog.Field{
		Key:   p.logFieldKeyCmd,
		Value: cmd.String(),
	})
}
