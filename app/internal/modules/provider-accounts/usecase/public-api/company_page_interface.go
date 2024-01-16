package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"
)

type (
	CompanyPageService interface {
		GetItemByName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error)
	}

	CompanyPageStorage interface {
		FetchByRewriteName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error)
	}
)
