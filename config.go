package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultConnectTimeout = 10 * time.Second
	DefaultReadTimeout    = 10 * time.Second
	DefaultWriteTimeout   = 10 * time.Second

	DefaultLogFieldKeyCmd = "cmd"
)

type Config struct {
	*redis.Options

	LogFieldKeyAddr string
	LogFieldKeyCmd  string
}

func NewConfig(host, pass string, port int) *Config {
	return &Config{
		Options: &redis.Options{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Password:     pass,
			DB:           0,
			DialTimeout:  DefaultConnectTimeout,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		},

		LogFieldKeyCmd: DefaultLogFieldKeyCmd,
	}
}
