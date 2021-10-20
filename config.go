package redis

import (
	"fmt"
	"time"
)

const (
	DefaultConnectTimeout = 10 * time.Second
	DefaultReadTimeout    = 10 * time.Second
	DefaultWriteTimeout   = 10 * time.Second

	DefaultLogFieldKeyAddr = "redis"
	DefaultLogFieldKeyCmd  = "cmd"
)

type Config struct {
	Addr string
	Pass string

	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	TimeoutAutoReconnect bool

	LogFieldKeyAddr string
	LogFieldKeyCmd  string
}

func NewConfig(host, pass string, port int) *Config {
	return &Config{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Pass: pass,

		ConnectTimeout: DefaultConnectTimeout,
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,

		TimeoutAutoReconnect: true,

		LogFieldKeyAddr: DefaultLogFieldKeyAddr,
		LogFieldKeyCmd:  DefaultLogFieldKeyCmd,
	}
}
