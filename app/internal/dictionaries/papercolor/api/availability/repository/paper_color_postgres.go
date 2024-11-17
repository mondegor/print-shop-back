package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
)

type (
	// PaperColorPostgres - comment struct.
	PaperColorPostgres struct {
		repoStatus db.FieldFetcher[uint64, mrenum.ItemStatus]
	}
)

// NewPaperColorPostgres - создаёт объект PaperColorPostgres.
func NewPaperColorPostgres(client mrstorage.DBConnManager) *PaperColorPostgres {
	return &PaperColorPostgres{
		repoStatus: db.NewFieldFetcher[uint64, mrenum.ItemStatus](
			client,
			module.DBTableNamePaperColors,
			"color_id",
			"color_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperColorPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
