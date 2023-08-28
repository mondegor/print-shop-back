package mrpostgres

import (
    "context"

    "github.com/jackc/pgx/v5"
)

func (c *Connection) Begin(ctx context.Context) (pgx.Tx, error) {
    tx, err := c.conn.Begin(ctx)

    if err != nil {
        return nil, c.wrapError(err)
    }

    return tx, nil
}
