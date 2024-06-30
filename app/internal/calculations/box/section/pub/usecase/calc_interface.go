package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/entity"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CalcResultUseCase - comment interface.
	CalcResultUseCase interface {
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.CalcResult, error)
		Create(ctx context.Context, item entity.CalcResult) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.CalcResult) error
	}

	// CalcResultStorage - comment interface.
	CalcResultStorage interface {
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.CalcResult, error)
		Insert(ctx context.Context, row entity.CalcResult) (mrtype.KeyInt32, error)
	}
)
