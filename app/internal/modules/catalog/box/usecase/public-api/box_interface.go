package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/box/entity/public-api"
)

type (
	BoxUseCase interface {
		GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, error)
	}

	BoxStorage interface {
		Fetch(ctx context.Context, params entity.BoxParams) ([]entity.Box, error)
	}
)
