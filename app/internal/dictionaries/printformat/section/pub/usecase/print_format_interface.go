package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/entity"
)

type (
	// PrintFormatUseCase - comment interface.
	PrintFormatUseCase interface {
		GetList(ctx context.Context, params entity.PrintFormatParams) ([]entity.PrintFormat, error)
	}

	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		Fetch(ctx context.Context, params entity.PrintFormatParams) ([]entity.PrintFormat, error)
	}
)
