package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/providers-api"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/google/uuid"
)

type (
	CompanyPageUseCase interface {
		GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		Store(ctx context.Context, item entity.CompanyPage) error
		ChangeStatus(ctx context.Context, item entity.CompanyPage) error
	}

	CompanyPageStorage interface {
		FetchOne(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		FetchAccountIdByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error)
		FetchStatus(ctx context.Context, row entity.CompanyPage) (enums.PublicStatus, error)
		InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error
		UpdateStatus(ctx context.Context, row entity.CompanyPage) error
	}
)
