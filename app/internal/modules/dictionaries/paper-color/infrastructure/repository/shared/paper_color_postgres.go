package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-color"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PaperColorFetchStatusPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func PaperColorFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            color_status
        FROM
            ` + module.DBSchema + `.paper_colors
        WHERE
            color_id = $1 AND color_status <> $2
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
