package adm

import (
	"context"

	"print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetList(ctx context.Context, params entity.CompanyPageParams) (items []entity.CompanyPage, countItems int, err error)
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		FetchWithTotal(ctx context.Context, params entity.CompanyPageParams) (rows []entity.CompanyPage, countRows int, err error)
	}
)
