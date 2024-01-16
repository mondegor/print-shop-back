package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CompanyPageService interface {
		GetItem(ctx context.Context, accountID mrtype.KeyString) (*entity.CompanyPage, error)
		Store(ctx context.Context, item *entity.CompanyPage) error
		ChangeStatus(ctx context.Context, item *entity.CompanyPage) error
	}

	CompanyPageStorage interface {
		LoadOne(ctx context.Context, row *entity.CompanyPage) error
		FetchStatus(ctx context.Context, row *entity.CompanyPage) (entity_shared.PublicStatus, error)
		InsertOrUpdate(ctx context.Context, row *entity.CompanyPage) error
		UpdateStatus(ctx context.Context, row *entity.CompanyPage) error
	}
)
