package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-facture"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PaperFactureIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func PaperFactureIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) error {
	sql := `
        SELECT
            1
        FROM
            ` + module.DBSchema + `.paper_factures
        WHERE
            facture_id = $1 AND facture_status <> $2
        LIMIT 1;`

	return conn.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}
