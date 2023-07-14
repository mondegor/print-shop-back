package repository

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/client/mrpostgres"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"

    "github.com/Masterminds/squirrel"
)

type CatalogPrintFormat struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogPrintFormat(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        client: client,
        builder: queryBuilder,
    }
}

func (f *CatalogPrintFormat) LoadAll(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter, rows *[]entity.CatalogPrintFormat) error {
    tbl := f.builder.
        Select(`
            format_id,
            tag_version,
            datetime_created,
            format_caption,
            format_length,
            format_width,
            format_status`).
        From("public.catalog_print_formats").
        Where(squirrel.NotEq{"format_status": entity.ItemStatusRemoved}).
        OrderBy("format_caption ASC, format_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"format_status": listFilter.Statuses})
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
        var row entity.CatalogPrintFormat

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
            &row.Length,
            &row.Width,
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
// modifies: row{Version, CreatedAt, Caption, Length, Width, Status}
func (f *CatalogPrintFormat) LoadOne(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            format_caption,
            format_length,
            format_width,
            format_status
        FROM
            public.catalog_print_formats
        WHERE format_id = $1 AND format_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Length,
        &row.Width,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (f *CatalogPrintFormat) FetchStatus(ctx context.Context, row *entity.CatalogPrintFormat) (entity.ItemStatus, error) {
    sql := `
        SELECT format_status
        FROM
            public.catalog_print_formats
        WHERE format_id = $1 AND tag_version = $2 AND format_status <> $3;`

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
// uses: row{Caption, Length, Width, Status}
// modifies: row{Id}
func (f *CatalogPrintFormat) Insert(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        INSERT INTO public.catalog_print_formats
            (format_caption,
             format_length,
             format_width,
             format_status)
        VALUES
            ($1, $2, $3, $4)
        RETURNING format_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Caption,
        row.Length,
        row.Width,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Caption, Length, Width, Status}
func (f *CatalogPrintFormat) Update(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        UPDATE public.catalog_print_formats
        SET
            tag_version = tag_version + 1,
            format_caption = $4,
            format_length = $5,
            format_width = $6
        WHERE format_id = $1 AND tag_version = $2 AND format_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Caption,
        row.Length,
        row.Width,
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
func (f *CatalogPrintFormat) UpdateStatus(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        UPDATE public.catalog_print_formats
        SET
            tag_version = tag_version + 1,
            format_status = $4
        WHERE
            format_id = $1 AND tag_version = $2 AND format_status <> $3;`

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

func (f *CatalogPrintFormat) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_print_formats
        SET
            tag_version = tag_version + 1,
            format_status = $2
        WHERE
            format_id = $1 AND format_status <> $2;`

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
