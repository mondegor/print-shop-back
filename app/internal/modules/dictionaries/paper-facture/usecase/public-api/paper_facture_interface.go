package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/public-api"
)

type (
	PaperFactureUseCase interface {
		GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error)
	}

	PaperFactureStorage interface {
		Fetch(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error)
	}
)
