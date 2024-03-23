package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/print-format"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PrintFormatFetchStatusPostgres
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func PrintFormatFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            format_status
        FROM
            ` + module.DBSchema + `.print_formats
        WHERE
            format_id = $1 AND format_status <> $2
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := conn.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}
