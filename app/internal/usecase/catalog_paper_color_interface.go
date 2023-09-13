package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogPaperColorService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter) ([]entity.CatalogPaperColor, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperColor, error)
        Create(ctx context.Context, item *entity.CatalogPaperColor) error
        Store(ctx context.Context, item *entity.CatalogPaperColor) error
        ChangeStatus(ctx context.Context, item *entity.CatalogPaperColor) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogPaperColorStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter, rows *[]entity.CatalogPaperColor) error
        LoadOne(ctx context.Context, row *entity.CatalogPaperColor) error
        FetchStatus(ctx context.Context, row *entity.CatalogPaperColor) (mrcom.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogPaperColor) error
        Update(ctx context.Context, row *entity.CatalogPaperColor) error
        UpdateStatus(ctx context.Context, row *entity.CatalogPaperColor) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
