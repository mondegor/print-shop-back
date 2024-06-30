package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/entity"
)

type (
	// BoxUseCase - comment interface.
	BoxUseCase interface {
		GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, error)
	}

	// BoxStorage - comment interface.
	BoxStorage interface {
		Fetch(ctx context.Context, params entity.BoxParams) ([]entity.Box, error)
	}
)
