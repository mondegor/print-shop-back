package repository_shared

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-facture"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PaperFactureFetchStatusPostgres
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func PaperFactureFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            facture_status
        FROM
            ` + module.DBSchema + `.paper_factures
        WHERE
            facture_id = $1 AND facture_status <> $2
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
