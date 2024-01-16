package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// LaminateTypeIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func LaminateTypeIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, id mrtype.KeyInt32) error {
	sql := `
        SELECT
            1
        FROM
            ` + module.UnitLaminateTypeDBSchema + `.laminate_types
        WHERE
            type_id = $1 AND type_status <> $2
        LIMIT 1;`

	return conn.QueryRow(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	).Scan(
		&id,
	)
}
