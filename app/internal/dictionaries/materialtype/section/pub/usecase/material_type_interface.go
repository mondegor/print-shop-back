package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
)

type (
	// MaterialTypeUseCase - comment interface.
	MaterialTypeUseCase interface {
		GetList(ctx context.Context, params entity.MaterialTypeParams) ([]entity.MaterialType, error)
	}

	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		Fetch(ctx context.Context, params entity.MaterialTypeParams) ([]entity.MaterialType, error)
	}
)
