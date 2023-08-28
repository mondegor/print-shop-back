package mrredis

import (
    "context"
    "fmt"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "time"

    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v9"
    redislib "github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9
// go get github.com/go-redsync/redsync/v4

const ConnectionName = "redis"

type (
    Connection struct {
        conn redislib.UniversalClient
        redsync *redsync.Redsync
        logger mrapp.Logger
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
        logger: logger,
    }
}

func (c *Connection) Cli() redislib.UniversalClient {
    return c.conn
}

func (c *Connection) Connect(opt Options) error {
    if c.conn != nil {
        return mrerr.ErrStorageConnectionIsAlreadyCreated.New(ConnectionName)
    }

    conn := redislib.NewClient(getOptions(&opt))
    _, err := conn.Ping(context.Background()).Result()

    if err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, ConnectionName)
    }

    c.conn = conn

    pool := goredis.NewPool(conn)
    c.redsync = redsync.New(pool)

    return nil
}

func (c *Connection) Close() error {
    if c.conn == nil {
        return mrerr.ErrStorageConnectionIsNotOpened.New(ConnectionName)
    }

    conn := c.conn
    c.conn = nil
    err := conn.Close()

    if err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, ConnectionName)
    }

    return nil
}

func (c *Connection) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
    return c.redsync.NewMutex(name, options...)
}

func getOptions(o *Options) *redislib.Options {
    return &redislib.Options{
        Addr: fmt.Sprintf("%s:%s", o.Host, o.Port),
        Password: o.Password,
    }
}
