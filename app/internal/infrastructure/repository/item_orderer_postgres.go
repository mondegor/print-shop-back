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
func NewItemOrderer(client *mrpostgres.Connection,
                    queryBuilder squirrel.StatementBuilderType) *ItemOrderer {
    return &ItemOrderer{
        client: client,
        builder: queryBuilder,
    }
}

// WithMetaData -
func (re *ItemOrderer) WithMetaData(meta usecase.ItemMetaData) usecase.ItemOrdererStorage {
    return &ItemOrderer{
        client:  re.client,
        builder: re.builder,
        meta:    meta,
    }
}

// LoadNode -
func (re *ItemOrderer) LoadNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Select(`
            prev_field_id,
            next_field_id,
            order_field`).
        From(re.meta.TableInfo().Name).
        Where(squirrel.Eq{re.meta.TableInfo().PrimaryKey: row.Id})

    query = re.meta.PrepareSelect(query)

    return re.client.SqQueryRow(ctx, query).Scan(&row.PrevId, &row.NextId, &row.OrderField)
}

// LoadFirstNode -
func (re *ItemOrderer) LoadFirstNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Select(`MIN(order_field)`).
        From(re.meta.TableInfo().Name)

    query = re.meta.PrepareSelect(query)

    err := re.client.SqQueryRow(ctx, query).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = re.loadNodeByOrderField(ctx, row)

    if err != nil {
        return err
    }

    if row.PrevId > 0 {
        return mrerr.ErrStorageFetchedInvalidData.New(mrerr.Arg{"row.Id": row.Id, "row.PrevId": row.PrevId})
    }

    return nil
}

// LoadLastNode -
func (re *ItemOrderer) LoadLastNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Select(`MAX(order_field)`).
        From(re.meta.TableInfo().Name)

    query = re.meta.PrepareSelect(query)

    err := re.client.SqQueryRow(ctx, query).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = re.loadNodeByOrderField(ctx, row)

    if err != nil {
        return err
    }

    if row.NextId > 0 {
        return mrerr.ErrStorageFetchedInvalidData.New(mrerr.Arg{"row.Id": row.Id, "row.NextId": row.NextId})
    }

    return nil
}

// UpdateNode -
func (re *ItemOrderer) UpdateNode(ctx context.Context, row *entity.ItemOrdererNode) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Update(re.meta.TableInfo().Name).
        SetMap(map[string]any{
            "prev_field_id": row.PrevId,
            "next_field_id": row.NextId,
            "order_field": row.OrderField,
        }).
        Where(squirrel.Eq{re.meta.TableInfo().PrimaryKey: row.Id})

    query = re.meta.PrepareUpdate(query)

    err := re.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{re.meta.TableInfo().PrimaryKey: row.Id})
    }

    return err
}

// UpdateNodePrevId -
func (re *ItemOrderer) UpdateNodePrevId(ctx context.Context, id mrentity.KeyInt32, prevId mrentity.ZeronullInt32) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Update(re.meta.TableInfo().Name).
        Set("prev_field_id", prevId).
        Where(squirrel.Eq{re.meta.TableInfo().PrimaryKey: id})

    query = re.meta.PrepareUpdate(query)

    err := re.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{re.meta.TableInfo().PrimaryKey: id})
    }

    return nil
}

// UpdateNodeNextId -
func (re *ItemOrderer) UpdateNodeNextId(ctx context.Context, id mrentity.KeyInt32, nextId mrentity.ZeronullInt32) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Update(re.meta.TableInfo().Name).
        Set("next_field_id", nextId).
        Where(squirrel.Eq{re.meta.TableInfo().PrimaryKey: id})

    query = re.meta.PrepareUpdate(query)

    err := re.client.SqUpdate(ctx, query)

    if err != nil {
        return mrerr.ErrDataContainer.Wrap(err, mrerr.Arg{re.meta.TableInfo().PrimaryKey: id})
    }

    return nil
}

// RecalcOrderField -
func (re *ItemOrderer) RecalcOrderField(ctx context.Context, minBorder mrentity.Int64, step mrentity.Int64) error {
    if re.meta == nil {
        return mrerr.ErrInternalNilPointer.New()
    }

    query := re.builder.
        Update(re.meta.TableInfo().Name).
        Set("order_field", squirrel.Expr("order_field + ?", step)).
        Where(squirrel.Gt{"order_field": minBorder})

    query = re.meta.PrepareUpdate(query)

    return re.client.SqUpdate(ctx, query)
}

func (re *ItemOrderer) loadNodeByOrderField(ctx context.Context, row *entity.ItemOrdererNode) error {
    query := re.builder.
        Select(re.meta.TableInfo().PrimaryKey, `
            prev_field_id,
            next_field_id`).
        From(re.meta.TableInfo().Name).
        Where(squirrel.Eq{"order_field": row.OrderField}).
        Suffix("FETCH FIRST 1 ROWS ONLY")

    query = re.meta.PrepareSelect(query)

    return re.client.SqQueryRow(ctx, query).Scan(&row.Id, &row.PrevId, &row.NextId)
}
