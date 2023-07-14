package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "context"
)

type (
    CatalogPaperFactureService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogPaperFactureListFilter) ([]entity.CatalogPaperFacture, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPaperFacture, error)
        Create(ctx context.Context, item *entity.CatalogPaperFacture) error
        Store(ctx context.Context, item *entity.CatalogPaperFacture) error
        ChangeStatus(ctx context.Context, item *entity.CatalogPaperFacture) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogPaperFactureStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogPaperFactureListFilter, rows *[]entity.CatalogPaperFacture) error
        LoadOne(ctx context.Context, row *entity.CatalogPaperFacture) error
        FetchStatus(ctx context.Context, row *entity.CatalogPaperFacture) (entity.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogPaperFacture) error
        Update(ctx context.Context, row *entity.CatalogPaperFacture) error
        UpdateStatus(ctx context.Context, row *entity.CatalogPaperFacture) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
