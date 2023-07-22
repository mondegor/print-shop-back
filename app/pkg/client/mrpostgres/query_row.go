package mrpostgres

import "github.com/jackc/pgx/v5"

type QueryRow struct {
    conn *Connection
    row pgx.Row
    err error
}

func (qr QueryRow) Scan(dest ...any) error {
    if qr.err != nil {
        return qr.conn.wrapError(qr.err)
    }

    err := qr.row.Scan(dest...)

    if err != nil {
        return qr.conn.wrapError(err)
    }

    return nil
}
