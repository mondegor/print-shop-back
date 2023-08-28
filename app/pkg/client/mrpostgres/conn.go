package mrpostgres

import (
    "context"
    "fmt"
    "print-shop-back/pkg/mrerr"
    "time"

    "github.com/jackc/pgx/v5"
)

// go get -u github.com/jackc/pgx/v5
// go get -u github.com/Masterminds/squirrel

const ConnectionName = "postgres"
const ConnectionCloseTimeout = 10

type (
    Connection struct {
        conn *pgx.Conn
    }

    Options struct {
        Host string
        Port string
        Database string
        Username string
        Password string
        MaxPoolSize int32
        ConnAttempts int32
        ConnTimeout time.Duration
    }
)

func New() *Connection {
    return &Connection{}
}

func (c *Connection) Connect(ctx context.Context, opt Options) error {
    if c.conn != nil {
        return mrerr.ErrStorageConnectionIsAlreadyCreated.New(ConnectionName)
    }

    ctx, cancel := context.WithTimeout(ctx, opt.ConnTimeout * time.Second)
    defer cancel()

    var err error
    c.conn, err = pgx.Connect(ctx, getConnString(&opt))

    if err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, ConnectionName)
    }

    return nil
}

func (c *Connection) Close() error {
    if c.conn == nil {
        return mrerr.ErrStorageConnectionIsNotOpened.New(ConnectionName)
    }

    ctx, cancel := context.WithTimeout(context.Background(), ConnectionCloseTimeout * time.Second)
    defer cancel()

    conn := c.conn
    c.conn = nil
    err := conn.Close(ctx)

    if err != nil {
        return mrerr.ErrStorageConnectionFailed.Wrap(err, ConnectionName)
    }

    return nil
}

func getConnString(o *Options) string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        o.Username,
        o.Password,
        o.Host,
        o.Port,
        o.Database)
}
