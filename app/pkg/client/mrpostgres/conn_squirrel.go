package mrpostgres

import (
    "context"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

func (c *Connection) SqUpdate(ctx context.Context, query squirrel.UpdateBuilder) (pgconn.CommandTag, error) {
    sql, args, err := query.ToSql()

    if err != nil {
        return pgconn.CommandTag{}, mrerr.ErrInternal.Caller(1).Wrap(err)
    }

    c.debugQuery(ctx, sql)

    commandTag, err := c.conn.Exec(ctx, sql, args...)

    if err != nil {
        return commandTag, c.wrapError(err)
    }

    return commandTag, nil
}

func (c *Connection) SqQuery(ctx context.Context, query squirrel.SelectBuilder) (pgx.Rows, error) {
    sql, args, err := query.ToSql()

    if err != nil {
        return nil, mrerr.ErrInternal.Caller(1).Wrap(err)
    }

    c.debugQuery(ctx, sql)

    rows, err := c.conn.Query(ctx, sql, args...)

    if err != nil {
        return nil, c.wrapError(err)
    }

    return rows, nil
}

func (c *Connection) SqQueryRow(ctx context.Context, query squirrel.SelectBuilder) QueryRow {
    sql, args, err := query.ToSql()

    if err != nil {
        return QueryRow{err: err}
    }

    c.debugQuery(ctx, sql)

    return QueryRow{
        conn: c,
        row: c.conn.QueryRow(ctx, sql, args...),
    }
}
