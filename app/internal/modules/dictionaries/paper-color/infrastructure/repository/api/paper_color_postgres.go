package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColorPostgres struct {
		client mrstorage.DBConn
	}
)

func NewPaperColorPostgres(
	client mrstorage.DBConn,
) *PaperColorPostgres {
	return &PaperColorPostgres{
		client: client,
	}
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperColorPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.PaperColorFetchStatusPostgres(ctx, re.client, rowID)
}
