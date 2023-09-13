package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogLaminateService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogLaminateListFilter) ([]entity.CatalogLaminate, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminate, error)
        Create(ctx context.Context, item *entity.CatalogLaminate) error
        Store(ctx context.Context, item *entity.CatalogLaminate) error
        ChangeStatus(ctx context.Context, item *entity.CatalogLaminate) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogLaminateStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateListFilter, rows *[]entity.CatalogLaminate) error
        LoadOne(ctx context.Context, row *entity.CatalogLaminate) error
        FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error)
        FetchStatus(ctx context.Context, row *entity.CatalogLaminate) (mrcom.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogLaminate) error
        Update(ctx context.Context, row *entity.CatalogLaminate) error
        UpdateStatus(ctx context.Context, row *entity.CatalogLaminate) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
