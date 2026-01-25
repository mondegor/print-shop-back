package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		repoStatus db.FieldFetcher[uint64, itemstatus.Enum]
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		repoStatus: db.NewFieldFetcher[uint64, itemstatus.Enum](
			client,
			module.DBTableNamePrintFormats,
			"format_id",
			"format_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: itemstatus.Enum - exists, errors.ErrEventStorageNoRowFound - not exists, error - query error.
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
