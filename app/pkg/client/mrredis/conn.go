package mrredis

import (
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis"
)

// go get -u github.com/go-redis/redis

type (
	Connection struct {
        conn *redis.Client
        Logger mrapp.Logger
    }

    Options struct {
        Host string
        Port string
        Password string
        ConnTimeout time.Duration
    }
)

func New(logger mrapp.Logger) *Connection {
	return &Connection{
        Logger: logger,
    }
}

func (c *Connection) Cli() *redis.Client {
    return c.conn
}

func (c *Connection) Connect(opt Options) error {
	if c.conn != nil {
		return mrerr.ErrStorageConnectionIsAlreadyCreated.New("redis")
	}

	conn := redis.NewClient(getOptions(&opt))

    if _, err := conn.Ping().Result(); err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, "redis")
    }

    c.conn = conn

    return nil
}

func (c *Connection) Close(ctx context.Context) error {
	if c.conn == nil {
        panic("connection had not opened")
	}

	conn := c.conn
	c.conn = nil

    if err := conn.Close(); err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, "redis")
    }

	return nil
}

func getOptions(o *Options) *redis.Options {
    return &redis.Options{
        Addr: fmt.Sprintf("%s:%s", o.Host, o.Port),
        Password: o.Password,
    }
}
