package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// MaterialTypeFetchStatusPostgres - comment func.
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error.
func MaterialTypeFetchStatusPostgres(ctx context.Context, client mrstorage.DBConnManager, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            type_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
        WHERE
            type_id = $1 AND deleted_at IS NULL
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
