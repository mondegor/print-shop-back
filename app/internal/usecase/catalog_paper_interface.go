package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "context"
)

type (
    CatalogPaperService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogPaperListFilter) ([]entity.CatalogPaper, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaper, error)
        Create(ctx context.Context, item *entity.CatalogPaper) error
        Store(ctx context.Context, item *entity.CatalogPaper) error
        ChangeStatus(ctx context.Context, item *entity.CatalogPaper) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogPaperStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogPaperListFilter, rows *[]entity.CatalogPaper) error
        LoadOne(ctx context.Context, row *entity.CatalogPaper) error
        FetchIdByArticle(ctx context.Context, row *entity.CatalogPaper) (mrentity.KeyInt32, error)
        FetchStatus(ctx context.Context, row *entity.CatalogPaper) (entity.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogPaper) error
        Update(ctx context.Context, row *entity.CatalogPaper) error
        UpdateStatus(ctx context.Context, row *entity.CatalogPaper) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
