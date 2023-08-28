package mrpostgres

import (
    "print-shop-back/pkg/mrerr"

    "github.com/jackc/pgx/v5/pgconn"
)

func (c *Connection) wrapError(err error) error {
    _, ok := err.(*pgconn.PgError)

    if ok {
        // Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
        return mrerr.ErrStorageQueryFailed.Caller(2).Wrap(err)
    }

    if err.Error() == "no rows in result set" {
        return mrerr.ErrStorageNoRowFound.Caller(2).Wrap(err)
    }

    return mrerr.ErrInternal.Caller(2).Wrap(err)
}
