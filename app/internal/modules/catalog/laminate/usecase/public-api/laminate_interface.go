package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/public-api"
)

type (
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error)
	}

	LaminateStorage interface {
		Fetch(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error)
	}
)
