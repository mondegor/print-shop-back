package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
)

type (
	// PaperFactureUseCase - comment interface.
	PaperFactureUseCase interface {
		GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperFactureParams) ([]entity.PaperFacture, error)
	}

	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		Fetch(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error)
	}
)
