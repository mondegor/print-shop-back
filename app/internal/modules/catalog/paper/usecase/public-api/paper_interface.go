package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/paper/entity/public-api"
)

type (
	PaperUseCase interface {
		GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
	}

	PaperStorage interface {
		Fetch(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
	}
)
