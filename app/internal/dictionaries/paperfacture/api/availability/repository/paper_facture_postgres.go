package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
)

type (
	// PaperFacturePostgres - comment struct.
	PaperFacturePostgres struct {
		repoStatus db.FieldFetcher[uint64, mrenum.ItemStatus]
	}
)

// NewPaperFacturePostgres - создаёт объект PaperFacturePostgres.
func NewPaperFacturePostgres(client mrstorage.DBConnManager) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		repoStatus: db.NewFieldFetcher[uint64, mrenum.ItemStatus](
			client,
			module.DBTableNamePaperFactures,
			"facture_id",
			"facture_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
