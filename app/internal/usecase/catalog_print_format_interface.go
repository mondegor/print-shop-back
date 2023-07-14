package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrentity"
    "context"
)

type (
    CatalogPrintFormatService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter) ([]entity.CatalogPrintFormat, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogPrintFormat, error)
        Create(ctx context.Context, item *entity.CatalogPrintFormat) error
        Store(ctx context.Context, item *entity.CatalogPrintFormat) error
        ChangeStatus(ctx context.Context, item *entity.CatalogPrintFormat) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogPrintFormatStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter, rows *[]entity.CatalogPrintFormat) error
        LoadOne(ctx context.Context, row *entity.CatalogPrintFormat) error
        FetchStatus(ctx context.Context, row *entity.CatalogPrintFormat) (entity.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogPrintFormat) error
        Update(ctx context.Context, row *entity.CatalogPrintFormat) error
        UpdateStatus(ctx context.Context, row *entity.CatalogPrintFormat) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
