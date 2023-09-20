package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    AccountCompanyPageLogoService interface {
        Store(ctx context.Context, item *entity.CompanyPageLogoObject) error
        Remove(ctx context.Context, accountId mrentity.KeyString) error
    }

    AccountCompanyPageLogoStorage interface {
        Fetch(ctx context.Context, accountId mrentity.KeyString) (string, error)
        Update(ctx context.Context, accountId mrentity.KeyString, logoPath string) error
        Delete(ctx context.Context, accountId mrentity.KeyString) error
    }
)
