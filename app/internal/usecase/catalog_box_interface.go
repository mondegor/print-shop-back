package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrentity"
    "context"
)

type (
    CatalogBoxService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogBoxListFilter) ([]entity.CatalogBox, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogBox, error)
        Create(ctx context.Context, item *entity.CatalogBox) error
        Store(ctx context.Context, item *entity.CatalogBox) error
        ChangeStatus(ctx context.Context, item *entity.CatalogBox) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogBoxStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogBoxListFilter, rows *[]entity.CatalogBox) error
        LoadOne(ctx context.Context, row *entity.CatalogBox) error
        FetchIdByArticle(ctx context.Context, row *entity.CatalogBox) (mrentity.KeyInt32, error)
        FetchStatus(ctx context.Context, row *entity.CatalogBox) (entity.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogBox) error
        Update(ctx context.Context, row *entity.CatalogBox) error
        UpdateStatus(ctx context.Context, row *entity.CatalogBox) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)