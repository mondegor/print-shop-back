package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
)

type (
	// PaperColorPostgres - comment struct.
	PaperColorPostgres struct {
		repoStatus db.FieldFetcher[uint64, itemstatus.Enum]
	}
)

// NewPaperColorPostgres - создаёт объект PaperColorPostgres.
func NewPaperColorPostgres(client mrstorage.DBConnManager) *PaperColorPostgres {
	return &PaperColorPostgres{
		repoStatus: db.NewFieldFetcher[uint64, itemstatus.Enum](
			client,
			module.DBTableNamePaperColors,
			"color_id",
			"color_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: itemstatus.Enum - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperColorPostgres) FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
