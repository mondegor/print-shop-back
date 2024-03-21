package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/public-api"
)

type (
	LaminateTypeUseCase interface {
		GetList(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, error)
	}

	LaminateTypeStorage interface {
		Fetch(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, error)
	}
)
