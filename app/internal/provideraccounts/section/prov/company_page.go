package prov

import (
	"context"

	"github.com/google/uuid"

	"print-shop-back/internal/provideraccounts/section/prov/entity"
	"print-shop-back/pkg/provideraccounts/enum/publicstatus"
)

type (
	// CompanyPageUseCase - comment interface.
	CompanyPageUseCase interface {
		GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		Save(ctx context.Context, item entity.CompanyPage) error
		ChangeStatus(ctx context.Context, item entity.CompanyPage) error
	}

	// CompanyPageStorage - comment interface.
	CompanyPageStorage interface {
		FetchOne(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error)
		FetchAccountIDByRewriteName(ctx context.Context, rewriteName string) (rowID uuid.UUID, err error)
		FetchStatus(ctx context.Context, accountID uuid.UUID) (publicstatus.Enum, error)
		InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error
		UpdateStatus(ctx context.Context, row entity.CompanyPage) error
	}
)
