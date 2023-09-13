package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogLaminateTypeService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter) ([]entity.CatalogLaminateType, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogLaminateType, error)
        Create(ctx context.Context, item *entity.CatalogLaminateType) error
        Store(ctx context.Context, item *entity.CatalogLaminateType) error
        ChangeStatus(ctx context.Context, item *entity.CatalogLaminateType) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogLaminateTypeStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter, rows *[]entity.CatalogLaminateType) error
        LoadOne(ctx context.Context, row *entity.CatalogLaminateType) error
        FetchStatus(ctx context.Context, row *entity.CatalogLaminateType) (mrcom.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogLaminateType) error
        Update(ctx context.Context, row *entity.CatalogLaminateType) error
        UpdateStatus(ctx context.Context, row *entity.CatalogLaminateType) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
