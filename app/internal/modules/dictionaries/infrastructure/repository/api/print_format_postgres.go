package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PrintFormatPostgres struct {
		client mrstorage.DBConn
	}
)

func NewPrintFormatPostgres(
	client mrstorage.DBConn,
) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		client: client,
	}
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PrintFormatPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.PrintFormatIsExistsPostgres(ctx, re.client, id)
}
