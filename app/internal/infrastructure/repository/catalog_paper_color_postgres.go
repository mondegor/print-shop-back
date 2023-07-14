package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type CatalogPaperColor struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogPaperColor(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *CatalogPaperColor {
    return &CatalogPaperColor{
        client: client,
        builder: queryBuilder,
    }
}

func (f *CatalogPaperColor) LoadAll(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter, rows *[]entity.CatalogPaperColor) error {
    tbl := f.builder.
        Select(`
            color_id,
            tag_version,
            datetime_created,
            color_caption,
            color_status`).
        From("public.catalog_paper_colors").
        Where(squirrel.NotEq{"color_status": entity.ItemStatusRemoved}).
        OrderBy("color_caption ASC, color_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"color_status": listFilter.Statuses})
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
        var row entity.CatalogPaperColor

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
func (f *CatalogPaperColor) LoadOne(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            color_caption,
            color_status
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND color_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (f *CatalogPaperColor) FetchStatus(ctx context.Context, row *entity.CatalogPaperColor) (entity.ItemStatus, error) {
    sql := `
        SELECT color_status
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND tag_version = $2 AND color_status <> $3;`

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
func (f *CatalogPaperColor) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND color_status <> $2;`

    return f.client.QueryRow(ctx, sql, id, entity.ItemStatusRemoved).Scan(&id)
}

// Insert
// uses: row{Caption, Status}
// modifies: row{Id}
func (f *CatalogPaperColor) Insert(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        INSERT INTO public.catalog_paper_colors
            (color_caption,
             color_status)
        VALUES
            ($1, $2)
        RETURNING color_id;`

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
func (f *CatalogPaperColor) Update(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_caption = $4
        WHERE color_id = $1 AND tag_version = $2 AND color_status <> $3;`

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
func (f *CatalogPaperColor) UpdateStatus(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_status = $4
        WHERE
            color_id = $1 AND tag_version = $2 AND color_status <> $3;`

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

func (f *CatalogPaperColor) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_status = $2
        WHERE
            color_id = $1 AND color_status <> $2;`

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
