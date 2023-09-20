package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type CatalogLaminateType struct {
    client mrstorage.DbConn
    builder squirrel.StatementBuilderType
}

func NewCatalogLaminateType(client mrstorage.DbConn,
                            queryBuilder squirrel.StatementBuilderType) *CatalogLaminateType {
    return &CatalogLaminateType{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogLaminateType) LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateTypeListFilter, rows *[]entity.CatalogLaminateType) error {
    query := re.builder.
        Select(`
            type_id,
            tag_version,
            datetime_created,
            type_caption,
            type_status`).
        From("public.catalog_laminate_types").
        Where(squirrel.NotEq{"type_status": mrcom.ItemStatusRemoved}).
        OrderBy("type_caption ASC, type_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"type_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        var row entity.CatalogLaminateType

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
func (re *CatalogLaminateType) LoadOne(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            type_caption,
            type_status
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND type_status <> $2;`

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
func (re *CatalogLaminateType) FetchStatus(ctx context.Context, row *entity.CatalogLaminateType) (mrcom.ItemStatus, error) {
    sql := `
        SELECT type_status
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND tag_version = $2 AND type_status <> $3;`

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
func (re *CatalogLaminateType) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_laminate_types
        WHERE type_id = $1 AND type_status <> $2;`

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
func (re *CatalogLaminateType) Insert(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        INSERT INTO public.catalog_laminate_types
            (type_caption,
             type_status)
        VALUES
            ($1, $2)
        RETURNING type_id;`

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
func (re *CatalogLaminateType) Update(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_caption = $4
        WHERE type_id = $1 AND tag_version = $2 AND type_status <> $3;`

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
func (re *CatalogLaminateType) UpdateStatus(ctx context.Context, row *entity.CatalogLaminateType) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_status = $4
        WHERE
            type_id = $1 AND tag_version = $2 AND type_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogLaminateType) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_laminate_types
        SET
            tag_version = tag_version + 1,
            type_status = $2
        WHERE
            type_id = $1 AND type_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )
}
