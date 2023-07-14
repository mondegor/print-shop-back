package usecase

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"
    "context"
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
        FetchStatus(ctx context.Context, row *entity.FormFieldTemplate) (entity.ItemStatus, error)
        Insert(ctx context.Context, row *entity.FormFieldTemplate) error
        Update(ctx context.Context, row *entity.FormFieldTemplate) error
        UpdateStatus(ctx context.Context, row *entity.FormFieldTemplate) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
