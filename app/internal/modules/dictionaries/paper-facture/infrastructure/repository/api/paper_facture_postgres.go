package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-facture/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperFacturePostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	return repository_shared.PaperFactureIsExistsPostgres(ctx, re.client, rowID)
}
