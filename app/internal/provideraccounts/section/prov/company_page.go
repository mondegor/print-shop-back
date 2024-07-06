package prov

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"

	"github.com/google/uuid"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		Store(ctx context.Context, item entity.CompanyPage) error
		ChangeStatus(ctx context.Context, item entity.CompanyPage) error
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		FetchOne(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		FetchAccountIDByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error)
		FetchStatus(ctx context.Context, accountID uuid.UUID) (enum.PublicStatus, error)
		InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error
		UpdateStatus(ctx context.Context, row entity.CompanyPage) error
	}
)
