package usecase

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/mrentity"

    "github.com/Masterminds/squirrel"
)

type (
    ItemOrdererComponent interface {
        WithMetaData(meta ItemMetaData) ItemOrdererComponent
        InsertToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error
        InsertToLast(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveToLast(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveToFirst(ctx context.Context, nodeId mrentity.KeyInt32) error
        MoveAfterId(ctx context.Context, nodeId mrentity.KeyInt32, afterNodeId mrentity.KeyInt32) error
        Unlink(ctx context.Context, nodeId mrentity.KeyInt32) error
    }

    ItemOrdererStorage interface {
        WithMetaData(meta ItemMetaData) ItemOrdererStorage
        LoadNode(ctx context.Context, row *entity.ItemOrdererNode) error
        LoadFirstNode(ctx context.Context, row *entity.ItemOrdererNode) error
        LoadLastNode(ctx context.Context, row *entity.ItemOrdererNode) error
        UpdateNode(ctx context.Context, row *entity.ItemOrdererNode) error
        UpdateNodePrevId(ctx context.Context, id mrentity.KeyInt32, prevId mrentity.ZeronullInt32) error
        UpdateNodeNextId(ctx context.Context, id mrentity.KeyInt32, nextId mrentity.ZeronullInt32) error
        RecalcOrderField(ctx context.Context, minBorder mrentity.Int64, step mrentity.Int64) error
    }

    ItemMetaData interface {
        TableInfo() *entity.TableInfo
        Select(queryBuilder squirrel.SelectBuilder) squirrel.SelectBuilder
        Update(queryBuilder squirrel.UpdateBuilder) squirrel.UpdateBuilder
        Delete(queryBuilder squirrel.DeleteBuilder) squirrel.DeleteBuilder
    }
)
