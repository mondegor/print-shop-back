package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/laminate-type"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// LaminateTypeFetchStatusPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func LaminateTypeFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            type_status
        FROM
            ` + module.DBSchema + `.laminate_types
        WHERE
            type_id = $1 AND type_status <> $2
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
