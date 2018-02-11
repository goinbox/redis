package redis

import (
	"github.com/goinbox/pool"
)

type PConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)
}

type Pool struct {
	pl *pool.Pool

	config *PConfig
}

func NewPool(config *PConfig) *Pool {
	p := &Pool{
		config: config,
	}

	config.NewConnFunc = p.newConn
	config.KeepAliveFunc = keepAlive

	p.pl = pool.NewPool(&p.config.Config)

	return p
}

func (p *Pool) Get() (*Client, error) {
	conn, err := p.pl.Get()
	if err != nil {
		return nil, err
	}

	return conn.(*Client), nil
}

func (p *Pool) Put(client *Client) error {
	return p.pl.Put(client)
}

func (p *Pool) newConn() (pool.IConn, error) {
	return p.config.NewClientFunc()
}

func keepAlive(conn pool.IConn) error {
	client := conn.(*Client)

	return client.Do("ping").Err
}
