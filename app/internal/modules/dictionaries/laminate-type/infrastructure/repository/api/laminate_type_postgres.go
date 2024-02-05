package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/laminate-type/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *LaminateTypePostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.LaminateTypeIsExistsPostgres(ctx, re.client, id)
}
