package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
)

type (
	// MaterialTypeUseCase - comment interface.
	MaterialTypeUseCase interface {
		GetList(ctx context.Context, lz mrcore.Localizer, params entity.MaterialTypeParams) ([]entity.MaterialType, error)
	}

	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		Fetch(ctx context.Context, params entity.MaterialTypeParams) ([]entity.MaterialType, error)
	}
)
