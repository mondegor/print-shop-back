package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/dictionaries/materialtype/module"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		repoStatus db.FieldFetcher[uint64, workflow.ItemStatus]
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		repoStatus: db.NewFieldFetcher[uint64, workflow.ItemStatus](
			client,
			module.DBTableNameMaterialTypes,
			"type_id",
			"type_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: workflow.ItemStatus - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *MaterialTypePostgres) FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
