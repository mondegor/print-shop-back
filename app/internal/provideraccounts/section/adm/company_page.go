package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetList(ctx context.Context, params entity.CompanyPageParams) (items []entity.CompanyPage, countItems uint64, err error)
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		FetchWithTotal(ctx context.Context, params entity.CompanyPageParams) (rows []entity.CompanyPage, countRows uint64, err error)
	}
)
