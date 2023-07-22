package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
)

type (
    FormFieldItemService interface {
        GetList(ctx context.Context, listFilter *entity.FormFieldItemListFilter) ([]entity.FormFieldItem, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) (*entity.FormFieldItem, error)
        Create(ctx context.Context, item *entity.FormFieldItem) error
        Store(ctx context.Context, item *entity.FormFieldItem) error
        Remove(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error
        MoveAfterId(ctx context.Context, id mrentity.KeyInt32, afterId mrentity.KeyInt32, formId mrentity.KeyInt32) error
    }

    FormFieldItemStorage interface {
        GetMetaData(formId mrentity.KeyInt32) ItemMetaData
        LoadAll(ctx context.Context, listFilter *entity.FormFieldItemListFilter, rows *[]entity.FormFieldItem) error
        LoadOne(ctx context.Context, row *entity.FormFieldItem) error
        FetchIdByName(ctx context.Context, row *entity.FormFieldItem) (mrentity.KeyInt32, error)
        Insert(ctx context.Context, row *entity.FormFieldItem) error
        Update(ctx context.Context, row *entity.FormFieldItem) error
        Delete(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error
    }
)
