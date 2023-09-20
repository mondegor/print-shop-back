package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type CatalogBox struct {
    client mrstorage.DbConn
    builder squirrel.StatementBuilderType
}

func NewCatalogBox(client mrstorage.DbConn,
                   queryBuilder squirrel.StatementBuilderType) *CatalogBox {
    return &CatalogBox{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogBox) LoadAll(ctx context.Context, listFilter *entity.CatalogBoxListFilter, rows *[]entity.CatalogBox) error {
    query := re.builder.
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
        Where(squirrel.NotEq{"box_status": mrcom.ItemStatusRemoved}).
        OrderBy("box_caption ASC, box_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"box_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

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

        *rows = append(*rows, row)
    }

    return cursor.Err()
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Article, Caption, Length, Width, Depth, Status}
func (re *CatalogBox) LoadOne(ctx context.Context, row *entity.CatalogBox) error {
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

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrcom.ItemStatusRemoved,
    ).Scan(
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

func (re *CatalogBox) FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error) {
    sql := `
        SELECT box_id
        FROM
            public.catalog_boxes
        WHERE box_article = $1;`

    var id mrentity.KeyInt32

    err := re.client.QueryRow(
        ctx,
        sql,
        article,
    ).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogBox) FetchStatus(ctx context.Context, row *entity.CatalogBox) (mrcom.ItemStatus, error) {
    sql := `
        SELECT box_status
        FROM
            public.catalog_boxes
        WHERE box_id = $1 AND tag_version = $2 AND box_status <> $3;`

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
// uses: row{Article, Caption, Length, Width, Depth, Status}
// modifies: row{Id}
func (re *CatalogBox) Insert(ctx context.Context, row *entity.CatalogBox) error {
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

    err := re.client.QueryRow(
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
func (re *CatalogBox) Update(ctx context.Context, row *entity.CatalogBox) error {
    filledFields, err := mrentity.FilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.catalog_boxes").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"box_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"box_status": mrcom.ItemStatusRemoved})

    return re.client.SqExec(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogBox) UpdateStatus(ctx context.Context, row *entity.CatalogBox) error {
    sql := `
        UPDATE public.catalog_boxes
        SET
            tag_version = tag_version + 1,
            box_status = $4
        WHERE
            box_id = $1 AND tag_version = $2 AND box_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogBox) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_boxes
        SET
            tag_version = tag_version + 1,
            box_article = NULL,
            box_status = $2
        WHERE
            box_id = $1 AND box_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )
}
