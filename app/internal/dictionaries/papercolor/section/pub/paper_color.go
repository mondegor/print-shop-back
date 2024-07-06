package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
)

type (
	// PaperColorUseCase - comment interface.
	PaperColorUseCase interface {
		GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error)
	}

	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		Fetch(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error)
	}
)
