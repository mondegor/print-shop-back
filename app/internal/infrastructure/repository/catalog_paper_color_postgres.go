package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type CatalogPaperColor struct {
    client mrstorage.DbConn
    builder squirrel.StatementBuilderType
}

func NewCatalogPaperColor(client mrstorage.DbConn,
                          queryBuilder squirrel.StatementBuilderType) *CatalogPaperColor {
    return &CatalogPaperColor{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogPaperColor) LoadAll(ctx context.Context, listFilter *entity.CatalogPaperColorListFilter, rows *[]entity.CatalogPaperColor) error {
    query := re.builder.
        Select(`
            color_id,
            tag_version,
            datetime_created,
            color_caption,
            color_status`).
        From("public.catalog_paper_colors").
        Where(squirrel.NotEq{"color_status": mrcom.ItemStatusRemoved}).
        OrderBy("color_caption ASC, color_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"color_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        var row entity.CatalogPaperColor

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
            &row.Status,
        )

        *rows = append(*rows, row)
    }

    return cursor.Err()
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Caption, Status}
func (re *CatalogPaperColor) LoadOne(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            color_caption,
            color_status
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND color_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogPaperColor) FetchStatus(ctx context.Context, row *entity.CatalogPaperColor) (mrcom.ItemStatus, error) {
    sql := `
        SELECT color_status
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND tag_version = $2 AND color_status <> $3;`

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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CatalogPaperColor) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_paper_colors
        WHERE color_id = $1 AND color_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &id,
    )
}

// Insert
// uses: row{Caption, Status}
// modifies: row{Id}
func (re *CatalogPaperColor) Insert(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        INSERT INTO public.catalog_paper_colors
            (color_caption,
             color_status)
        VALUES
            ($1, $2)
        RETURNING color_id;`

    err := re.client.QueryRow(
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
func (re *CatalogPaperColor) Update(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_caption = $4
        WHERE color_id = $1 AND tag_version = $2 AND color_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Caption,
    )
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogPaperColor) UpdateStatus(ctx context.Context, row *entity.CatalogPaperColor) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_status = $4
        WHERE
            color_id = $1 AND tag_version = $2 AND color_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogPaperColor) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_paper_colors
        SET
            tag_version = tag_version + 1,
            color_status = $2
        WHERE
            color_id = $1 AND color_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )
}
