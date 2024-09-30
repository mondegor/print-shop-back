package pub

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// LaminateUseCase - comment interface.
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error)
		GetTypeList(ctx context.Context) ([]mrtype.KeyInt32, error)
		GetThicknessList(ctx context.Context) ([]measure.Meter, error)
	}

	// LaminateStorage - comment interface.
	LaminateStorage interface {
		Fetch(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error)
		FetchTypeIDs(ctx context.Context) ([]mrtype.KeyInt32, error)
		FetchThicknesses(ctx context.Context) ([]measure.Meter, error)
	}
)
