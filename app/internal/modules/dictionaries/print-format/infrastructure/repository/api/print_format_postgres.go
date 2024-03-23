package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/print-format/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
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

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.PrintFormatFetchStatusPostgres(ctx, re.client, rowID)
}
