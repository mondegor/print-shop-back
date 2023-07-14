package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "context"
)

type (
    FormDataService interface {
        GetList(ctx context.Context, listFilter *entity.FormDataListFilter) ([]entity.FormData, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormData, error)
        CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error
        Create(ctx context.Context, item *entity.FormData) error
        Store(ctx context.Context, item *entity.FormData) error
        ChangeStatus(ctx context.Context, item *entity.FormData) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    FormDataStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.FormDataListFilter, rows *[]entity.FormData) error
        LoadOne(ctx context.Context, row *entity.FormData) error
        FetchStatus(ctx context.Context, row *entity.FormData) (entity.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.FormData) error
        Update(ctx context.Context, row *entity.FormData) error
        UpdateStatus(ctx context.Context, row *entity.FormData) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
