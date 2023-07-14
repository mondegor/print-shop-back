package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type CatalogBox struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogBox(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *CatalogBox {
    return &CatalogBox{
        client: client,
        builder: queryBuilder,
    }
}

func (f *CatalogBox) LoadAll(ctx context.Context, listFilter *entity.CatalogBoxListFilter, rows *[]entity.CatalogBox) error {
    tbl := f.builder.
        Select(`
            box_id,
            tag_version,
            datetime_created,
            box_article,
            box_caption,
            box_length,
            box_width,
            box_depth,
            box_status`).
        From("public.catalog_boxes").
        Where(squirrel.NotEq{"box_status": entity.ItemStatusRemoved}).
        OrderBy("box_caption ASC, box_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"box_status": listFilter.Statuses})
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
        var row entity.CatalogBox

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Article,
            &row.Caption,
            &row.Length,
            &row.Width,
            &row.Depth,
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
// modifies: row{Version, CreatedAt, Article, Caption, Length, Width, Depth, Status}
func (f *CatalogBox) LoadOne(ctx context.Context, row *entity.CatalogBox) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            box_article,
            box_caption,
            box_length,
            box_width,
            box_depth,
            box_status
        FROM
            public.catalog_boxes
        WHERE box_id = $1 AND box_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Article,
        &row.Caption,
        &row.Length,
        &row.Width,
        &row.Depth,
        &row.Status,
    )
}

// FetchIdByArticle
// uses: row{Article}
func (f *CatalogBox) FetchIdByArticle(ctx context.Context, row *entity.CatalogBox) (mrentity.KeyInt32, error) {
    sql := `
        SELECT box_id
        FROM
            public.catalog_boxes
        WHERE box_article = $1;`

    var id mrentity.KeyInt32

    err := f.client.QueryRow(ctx, sql, row.Article).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (f *CatalogBox) FetchStatus(ctx context.Context, row *entity.CatalogBox) (entity.ItemStatus, error) {
    sql := `
        SELECT box_status
        FROM
            public.catalog_boxes
        WHERE box_id = $1 AND tag_version = $2 AND box_status <> $3;`

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

// Insert
// uses: row{Article, Caption, Length, Width, Depth, Status}
// modifies: row{Id}
func (f *CatalogBox) Insert(ctx context.Context, row *entity.CatalogBox) error {
    sql := `
        INSERT INTO public.catalog_boxes
            (box_article,
             box_caption,
             box_length,
             box_width,
             box_depth,
             box_status)
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING box_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Article,
        row.Caption,
        row.Length,
        row.Width,
        row.Depth,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Article, Caption, Length, Width, Depth, Status}
func (f *CatalogBox) Update(ctx context.Context, row *entity.CatalogBox) error {
    sql := `
        UPDATE public.catalog_boxes
        SET
            tag_version = tag_version + 1,
            box_article = $4,
            box_caption = $5,
            box_length = $6,
            box_width = $7,
            box_depth = $8
        WHERE box_id = $1 AND tag_version = $2 AND box_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Article,
        row.Caption,
        row.Length,
        row.Width,
        row.Depth,
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
func (f *CatalogBox) UpdateStatus(ctx context.Context, row *entity.CatalogBox) error {
    sql := `
        UPDATE public.catalog_boxes
        SET
            tag_version = tag_version + 1,
            box_status = $4
        WHERE
            box_id = $1 AND tag_version = $2 AND box_status <> $3;`

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

func (f *CatalogBox) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_boxes
        SET
            tag_version = tag_version + 1,
            box_status = $2
        WHERE
            box_id = $1 AND box_status <> $2;`

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
