package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"
)

type (
	CompanyPageUseCase interface {
		GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error)
	}

	CompanyPageStorage interface {
		FetchByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error)
	}
)
