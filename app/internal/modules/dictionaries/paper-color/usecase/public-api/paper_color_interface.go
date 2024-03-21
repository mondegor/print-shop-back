package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/public-api"
)

type (
	PaperColorUseCase interface {
		GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error)
	}

	PaperColorStorage interface {
		Fetch(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error)
	}
)
