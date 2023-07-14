package mrpostgres

import (
    "print-shop-back/pkg/mrerr"
    "errors"

    "github.com/jackc/pgx/v5/pgconn"
)

func (c *Connection) wrapError(err error) error {
    var pgErr *pgconn.PgError

    if errors.As(err, &pgErr) {
        // Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
        return mrerr.ErrStorageQueryFailed.Wrap(err)
    }

    if err.Error() == "no rows in result set" {
        return mrerr.ErrStorageNoRowFound.Wrap(err)
    }

    return mrerr.ErrInternal.Wrap(err)
}
