package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-facture/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFacturePostgres struct {
		client mrstorage.DBConn
	}
)

func NewPaperFacturePostgres(
	client mrstorage.DBConn,
) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client: client,
	}
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperFacturePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.PaperFactureFetchStatusPostgres(ctx, re.client, rowID)
}
