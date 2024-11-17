package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		repoStatus db.FieldFetcher[uint64, mrenum.ItemStatus]
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		repoStatus: db.NewFieldFetcher[uint64, mrenum.ItemStatus](
			client,
			module.DBTableNamePrintFormats,
			"format_id",
			"format_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
