package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	CompanyPageUseCase interface {
		GetList(ctx context.Context, params entity.CompanyPageParams) ([]entity.CompanyPage, int64, error)
	}

	CompanyPageStorage interface {
		NewFetchParams(params entity.CompanyPageParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CompanyPage, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
	}
)
