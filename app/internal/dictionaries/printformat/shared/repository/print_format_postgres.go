package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
)

// PrintFormatFetchStatusPostgres - comment func.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func PrintFormatFetchStatusPostgres(ctx context.Context, client mrstorage.DBConnManager, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            format_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        WHERE
            format_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}
