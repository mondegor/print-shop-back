package repository_api

import (
	"context"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PaperColorPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.PaperColorIsExistsPostgres(ctx, re.client, id)
}
