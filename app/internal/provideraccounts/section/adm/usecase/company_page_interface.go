package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetList(ctx context.Context, params entity.CompanyPageParams) ([]entity.CompanyPage, int64, error)
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		NewSelectParams(params entity.CompanyPageParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.CompanyPage, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
	}
)
