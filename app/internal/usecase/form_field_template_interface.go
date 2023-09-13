package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    FormFieldTemplateService interface {
        GetList(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter) ([]entity.FormFieldTemplate, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.FormFieldTemplate, error)
        Create(ctx context.Context, item *entity.FormFieldTemplate) error
        Store(ctx context.Context, item *entity.FormFieldTemplate) error
        ChangeStatus(ctx context.Context, item *entity.FormFieldTemplate) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    FormFieldTemplateStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter, rows *[]entity.FormFieldTemplate) error
        LoadOne(ctx context.Context, row *entity.FormFieldTemplate) error
        FetchStatus(ctx context.Context, row *entity.FormFieldTemplate) (mrcom.ItemStatus, error)
        Insert(ctx context.Context, row *entity.FormFieldTemplate) error
        Update(ctx context.Context, row *entity.FormFieldTemplate) error
        UpdateStatus(ctx context.Context, row *entity.FormFieldTemplate) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
