package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

// ItemOrderer -
type ItemOrderer struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
    meta usecase.ItemMetaData
}

// NewItemOrderer -
func NewItemOrderer(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *ItemOrderer {
    return &ItemOrderer{
        client: client,
        builder: queryBuilder,
    }
}

// WithMetaData -
func (f *ItemOrderer) WithMetaData(meta usecase.ItemMetaData) usecase.ItemOrdererStorage {
    return &ItemOrderer{
        client: f.client,
        builder: f.builder,
        meta: meta,
    }
}

// LoadNode -
func (f *ItemOrderer) LoadNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Select(`
            prev_field_id,
            next_field_id,
            order_field`).
        From(f.meta.TableInfo().Name).
        Where(squirrel.Eq{f.meta.TableInfo().PrimaryKey: row.Id})

    query = f.meta.PrepareSelect(query)

    return f.client.SqQueryRow(ctx, query).Scan(&row.PrevId, &row.NextId, &row.OrderField)
}

// LoadFirstNode -
func (f *ItemOrderer) LoadFirstNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Select(`MIN(order_field)`).
        From(f.meta.TableInfo().Name)

    query = f.meta.PrepareSelect(query)

    err := f.client.SqQueryRow(ctx, query).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = f.loadNodeByOrderField(ctx, row)

    if err != nil {
        return err
    }

    if row.PrevId > 0 {
        return mrerr.ErrStorageFetchedInvalidData.New(mrerr.Arg{"row.Id": row.Id, "row.PrevId": row.PrevId})
    }

    return nil
}

// LoadLastNode -
func (f *ItemOrderer) LoadLastNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Select(`MAX(order_field)`).
        From(f.meta.TableInfo().Name)

    query = f.meta.PrepareSelect(query)

    err := f.client.SqQueryRow(ctx, query).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = f.loadNodeByOrderField(ctx, row)

    if err != nil {
        return err
    }

    if row.NextId > 0 {
        return mrerr.ErrStorageFetchedInvalidData.New(mrerr.Arg{"row.Id": row.Id, "row.NextId": row.NextId})
    }

    return nil
}

// UpdateNode -
func (f *ItemOrderer) UpdateNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Update(f.meta.TableInfo().Name).
        SetMap(map[string]any{
            "prev_field_id": row.PrevId,
            "next_field_id": row.NextId,
            "order_field": row.OrderField,
        }).
        Where(squirrel.Eq{f.meta.TableInfo().PrimaryKey: row.Id})

    query = f.meta.PrepareUpdate(query)

    err := f.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{f.meta.TableInfo().PrimaryKey: row.Id})
    }

    return err
}

// UpdateNodePrevId -
func (f *ItemOrderer) UpdateNodePrevId(ctx context.Context, id mrentity.KeyInt32, prevId mrentity.ZeronullInt32) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Update(f.meta.TableInfo().Name).
        Set("prev_field_id", prevId).
        Where(squirrel.Eq{f.meta.TableInfo().PrimaryKey: id})

    query = f.meta.PrepareUpdate(query)

    err := f.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{f.meta.TableInfo().PrimaryKey: id})
    }

    return nil
}

// UpdateNodeNextId -
func (f *ItemOrderer) UpdateNodeNextId(ctx context.Context, id mrentity.KeyInt32, nextId mrentity.ZeronullInt32) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Update(f.meta.TableInfo().Name).
        Set("next_field_id", nextId).
        Where(squirrel.Eq{f.meta.TableInfo().PrimaryKey: id})

    query = f.meta.PrepareUpdate(query)

    err := f.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{f.meta.TableInfo().PrimaryKey: id})
    }

    return nil
}

// RecalcOrderField -
func (f *ItemOrderer) RecalcOrderField(ctx context.Context, minBorder mrentity.Int64, step mrentity.Int64) error {
    if f.meta != nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := f.builder.
        Update(f.meta.TableInfo().Name).
        Set("order_field", squirrel.Expr("order_field + ?", step)).
        Where(squirrel.Gt{"order_field": minBorder})

    query = f.meta.PrepareUpdate(query)

    return f.client.SqUpdate(ctx, query)
}

func (f *ItemOrderer) loadNodeByOrderField(ctx context.Context, row *entity.ItemOrdererNode) error {
    query := f.builder.
        Select(f.meta.TableInfo().PrimaryKey, `
            prev_field_id,
            next_field_id`).
        From(f.meta.TableInfo().Name).
        Where(squirrel.Eq{"order_field": row.OrderField}).
        Suffix("FETCH FIRST 1 ROWS ONLY")

    query = f.meta.PrepareSelect(query)

    return f.client.SqQueryRow(ctx, query).Scan(&row.Id, &row.PrevId, &row.NextId)
}
