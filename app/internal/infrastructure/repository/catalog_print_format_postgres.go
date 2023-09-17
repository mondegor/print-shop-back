package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

type CatalogPrintFormat struct {
    client *mrpostgres.ConnAdapter
    builder squirrel.StatementBuilderType
}

func NewCatalogPrintFormat(client *mrpostgres.ConnAdapter,
                           queryBuilder squirrel.StatementBuilderType) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogPrintFormat) LoadAll(ctx context.Context, listFilter *entity.CatalogPrintFormatListFilter, rows *[]entity.CatalogPrintFormat) error {
    query := re.builder.
        Select(`
            format_id,
            tag_version,
            datetime_created,
            format_caption,
            format_length,
            format_width,
            format_status`).
        From("public.catalog_print_formats").
        Where(squirrel.NotEq{"format_status": mrcom.ItemStatusRemoved}).
        OrderBy("format_caption ASC, format_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"format_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

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
            return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
        }

        *rows = append(*rows, row)
    }

    if err = cursor.Err(); err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Caption, Length, Width, Status}
func (re *CatalogPrintFormat) LoadOne(ctx context.Context, row *entity.CatalogPrintFormat) error {
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

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrcom.ItemStatusRemoved,
    ).Scan(
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
func (re *CatalogPrintFormat) FetchStatus(ctx context.Context, row *entity.CatalogPrintFormat) (mrcom.ItemStatus, error) {
    sql := `
        SELECT format_status
        FROM
            public.catalog_print_formats
        WHERE format_id = $1 AND tag_version = $2 AND format_status <> $3;`

    var status mrcom.ItemStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// Insert
// uses: row{Caption, Length, Width, Status}
// modifies: row{Id}
func (re *CatalogPrintFormat) Insert(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        INSERT INTO public.catalog_print_formats
            (format_caption,
             format_length,
             format_width,
             format_status)
        VALUES
            ($1, $2, $3, $4)
        RETURNING format_id;`

    err := re.client.QueryRow(
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
func (re *CatalogPrintFormat) Update(ctx context.Context, row *entity.CatalogPrintFormat) error {
    filledFields, err := mrentity.FilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.catalog_print_formats").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"format_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"format_status": mrcom.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogPrintFormat) UpdateStatus(ctx context.Context, row *entity.CatalogPrintFormat) error {
    sql := `
        UPDATE public.catalog_print_formats
        SET
            tag_version = tag_version + 1,
            format_status = $4
        WHERE
            format_id = $1 AND tag_version = $2 AND format_status <> $3;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

func (re *CatalogPrintFormat) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_print_formats
        SET
            tag_version = tag_version + 1,
            format_status = $2
        WHERE
            format_id = $1 AND format_status <> $2;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}
