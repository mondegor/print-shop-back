package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/entity"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error)
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		FetchByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error)
	}
)
