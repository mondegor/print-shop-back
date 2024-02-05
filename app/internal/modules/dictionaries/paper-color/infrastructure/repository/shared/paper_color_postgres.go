package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-color"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PaperColorIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func PaperColorIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, id mrtype.KeyInt32) error {
	sql := `
        SELECT
            1
        FROM
            ` + module.DBSchema + `.paper_colors
        WHERE
            color_id = $1 AND color_status <> $2
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
