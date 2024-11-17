package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		repoStatus db.FieldFetcher[uint64, mrenum.ItemStatus]
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		repoStatus: db.NewFieldFetcher[uint64, mrenum.ItemStatus](
			client,
			module.DBTableNameMaterialTypes,
			"type_id",
			"type_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *MaterialTypePostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
