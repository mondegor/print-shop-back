package repository

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "context"

    "github.com/Masterminds/squirrel"
)

type CatalogLaminateType struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogLaminateType(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *CatalogLaminateType {
    return &CatalogLaminateType{
        client: client,
        builder: queryBuilder,
    }
}

func (f *CatalogLaminateType) LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter, rows *[]entity.CatalogLaminateType) error {
    tbl := f.builder.
        Select(`
            type_id,
            tag_version,
            datetime_created,
            type_caption,
            type_status`).
        From("public.catalog_laminate_types").
        Where(squirrel.NotEq{"type_status": entity.ItemStatusRemoved}).
        OrderBy("type_caption ASC, type_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"type_status": listFilter.Statuses})
    }

    sql, args, err := tbl.ToSql()

    if err != nil {
        return mrerr.ErrInternal.Wrap(err)
    }

    cursor, err := f.client.Query(ctx, sql, args...)

    if err != nil {
        return err
    }

    for cursor.Next() {
        var row entity.CatalogLaminateType

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
            &row.Status,
        )

        if err != nil {
            return mrerr.ErrStorageFetchDataFailed.Wrap(err)
        }

        *rows = append(*rows, row)
    }

    if err = cursor.Err(); err != nil {
        return mrerr.ErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Caption, Status}
func (f *CatalogLaminateType) LoadOne(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            type_caption,
            type_status
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND type_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (f *CatalogLaminateType) FetchStatus(ctx context.Context, row *entity.CatalogLaminateType) (entity.ItemStatus, error) {
    sql := `
        SELECT type_status
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND tag_version = $2 AND type_status <> $3;`

    var status entity.ItemStatus

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (f *CatalogLaminateType) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND type_status <> $2;`

    return f.client.QueryRow(ctx, sql, id, entity.ItemStatusRemoved).Scan(&id)
}

// Insert
// uses: row{Caption, Status}
// modifies: row{Id}
func (f *CatalogLaminateType) Insert(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        INSERT INTO public.catalog_laminate_types
            (type_caption,
             type_status)
        VALUES
            ($1, $2)
        RETURNING type_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Caption,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Caption, Status}
func (f *CatalogLaminateType) Update(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_caption = $4
        WHERE type_id = $1 AND tag_version = $2 AND type_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Caption,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (f *CatalogLaminateType) UpdateStatus(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_status = $4
        WHERE
            type_id = $1 AND tag_version = $2 AND type_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Status,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}

func (f *CatalogLaminateType) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_status = $2
        WHERE
            type_id = $1 AND type_status <> $2;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        id,
        entity.ItemStatusRemoved,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}
