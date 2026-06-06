package repository

import (
	"context"

	"github.com/mondegor/go-sysmess/mrpostgres/db"
	"github.com/mondegor/go-sysmess/mrstorage"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/dictionaries/printformat/module"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		repoStatus db.FieldFetcher[uint64, workflow.ItemStatus]
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		repoStatus: db.NewFieldFetcher[uint64, workflow.ItemStatus](
			client,
			module.DBTableNamePrintFormats,
			"format_id",
			"format_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: workflow.ItemStatus - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
