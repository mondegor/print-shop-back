package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// PaperFactureFetchStatusPostgres - comment func.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func PaperFactureFetchStatusPostgres(ctx context.Context, client mrstorage.DBConnManager, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            facture_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePaperFactures + `
        WHERE
            facture_id = $1 AND deleted_at IS NULL
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
