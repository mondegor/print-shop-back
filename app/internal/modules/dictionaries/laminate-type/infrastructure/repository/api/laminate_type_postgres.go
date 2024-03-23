package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/laminate-type/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateTypePostgres struct {
		client mrstorage.DBConn
	}
)

func NewLaminateTypePostgres(
	client mrstorage.DBConn,
) *LaminateTypePostgres {
	return &LaminateTypePostgres{
		client: client,
	}
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *LaminateTypePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.LaminateTypeFetchStatusPostgres(ctx, re.client, rowID)
}
