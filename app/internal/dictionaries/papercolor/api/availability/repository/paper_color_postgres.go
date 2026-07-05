package repository

import (
	"context"

	"github.com/mondegor/go-core/mrpostgres/db"
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/dictionaries/papercolor/module"
)

type (
	// PaperColorPostgres - comment struct.
	PaperColorPostgres struct {
		repoStatus db.FieldFetcher[uint64, workflow.ItemStatus]
	}
)

// NewPaperColorPostgres - создаёт объект PaperColorPostgres.
func NewPaperColorPostgres(client mrstorage.DBConnManager) *PaperColorPostgres {
	return &PaperColorPostgres{
		repoStatus: db.NewFieldFetcher[uint64, workflow.ItemStatus](
			client,
			module.DBTableNamePaperColors,
			"color_id",
			"color_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: workflow.ItemStatus - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *PaperColorPostgres) FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
