package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		repoStatus db.FieldFetcher[uint64, itemstatus.Enum]
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		repoStatus: db.NewFieldFetcher[uint64, itemstatus.Enum](
			client,
			module.DBTableNameMaterialTypes,
			"type_id",
			"type_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: itemstatus.Enum - exists, errors.ErrEventStorageNoRecordFound - not exists, error - query error.
func (re *MaterialTypePostgres) FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
