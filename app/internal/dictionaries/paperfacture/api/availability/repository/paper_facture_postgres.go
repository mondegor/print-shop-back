package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
)

type (
	// PaperFacturePostgres - comment struct.
	PaperFacturePostgres struct {
		repoStatus db.FieldFetcher[uint64, itemstatus.Enum]
	}
)

// NewPaperFacturePostgres - создаёт объект PaperFacturePostgres.
func NewPaperFacturePostgres(client mrstorage.DBConnManager) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		repoStatus: db.NewFieldFetcher[uint64, itemstatus.Enum](
			client,
			module.DBTableNamePaperFactures,
			"facture_id",
			"facture_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: itemstatus.Enum - exists, errors.ErrEventStorageNoRowFound - not exists, error - query error.
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
