package usecase

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrentity"
    "context"
)

type (
    FormFieldItemService interface {
        GetList(ctx context.Context, listFilter *entity.FormFieldItemListFilter) ([]entity.FormFieldItem, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) (*entity.FormFieldItem, error)
        Create(ctx context.Context, item *entity.FormFieldItem) error
        Store(ctx context.Context, item *entity.FormFieldItem) error
        Remove(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error
    }

    FormFieldItemStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.FormFieldItemListFilter, rows *[]entity.FormFieldItem) error
        LoadOne(ctx context.Context, row *entity.FormFieldItem) error
        FetchIdByName(ctx context.Context, row *entity.FormFieldItem) (mrentity.KeyInt32, error)
        Insert(ctx context.Context, row *entity.FormFieldItem) error
        Update(ctx context.Context, row *entity.FormFieldItem) error
        Delete(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error
    }

    FormFieldItemOrdererService interface {
        InsertToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error
        InsertToLast(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveToLast(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveAfterId(ctx context.Context, nodeId mrentity.KeyInt32, afterNodeId mrentity.KeyInt32) error
        Unlink(ctx context.Context, nodeId mrentity.KeyInt32) error
    }

    FormFieldItemOrdererStorage interface {
        LoadNode(ctx context.Context, row *entity.Node) error
        LoadFirstNode(ctx context.Context, row *entity.Node) error
        LoadLastNode(ctx context.Context, row *entity.Node) error
        UpdateNode(ctx context.Context, row *entity.Node) error
        UpdateNodePrevId(ctx context.Context, id mrentity.KeyInt32, prevId mrentity.ZeronullInt32) error
        UpdateNodeNextId(ctx context.Context, id mrentity.KeyInt32, nextId mrentity.ZeronullInt32) error
        RecalcOrderField(ctx context.Context, minBorder mrentity.ZeronullInt64, step mrentity.ZeronullInt64) error
    }
)
