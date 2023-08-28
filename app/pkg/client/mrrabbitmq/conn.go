package mrrabbitmq

import (
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
)

// go get github.com/rabbitmq/amqp091-go@v1.8.1

const ConnectionName = "rabbitmq"

type (
    Connection struct {
        conn *amqp.Connection
        logger mrapp.Logger
    }

    Options struct {
        Host string
        Port string
        User string
        Password string
    }
)

func New(logger mrapp.Logger) *Connection {
    return &Connection{
        logger: logger,
    }
}

func (c *Connection) Cli() *amqp.Connection {
    return c.conn
}

func (c *Connection) Connect(opt Options) error {
    if c.conn != nil {
        return mrerr.ErrStorageConnectionIsAlreadyCreated.New(ConnectionName)
    }

    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", opt.User, opt.Password, opt.Host, opt.Port))

    if err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, ConnectionName)
    }

    c.conn = conn

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
