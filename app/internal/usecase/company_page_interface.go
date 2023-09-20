package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CompanyPageService interface {
        GetList(ctx context.Context, listFilter *entity.CompanyPageListFilter) ([]entity.CompanyPage, error)
    }

    AccountCompanyPageService interface {
        GetItem(ctx context.Context, accountId mrentity.KeyString) (*entity.CompanyPage, error)
        Store(ctx context.Context, item *entity.CompanyPage) error
        ChangeStatus(ctx context.Context, item *entity.CompanyPage) error
    }

    PublicCompanyPageService interface {
        GetItemByName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error)
    }

    CompanyPageStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CompanyPageListFilter, rows *[]entity.CompanyPage) error
        LoadOne(ctx context.Context, row *entity.CompanyPage) error
        LoadOneByRewriteName(ctx context.Context, row *entity.CompanyPage) error
        FetchStatus(ctx context.Context, row *entity.CompanyPage) (entity.ResourceStatus, error)
        InsertOrUpdate(ctx context.Context, row *entity.CompanyPage) error
        UpdateStatus(ctx context.Context, row *entity.CompanyPage) error
    }
)
