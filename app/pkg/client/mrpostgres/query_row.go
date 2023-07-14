package mrpostgres

import "github.com/jackc/pgx/v5"

type QueryRow struct {
    conn *Connection
    row pgx.Row
}

func wrapQueryRow(conn *Connection, row pgx.Row) QueryRow {
    return QueryRow{
        conn: conn,
        row: row,
    }
}

func (qr QueryRow) Scan(dest ...any) error {
    err := qr.row.Scan(dest...)

    if err != nil {
        return qr.conn.wrapError(err)
    }

    return nil
}
