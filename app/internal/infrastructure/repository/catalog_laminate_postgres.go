package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type CatalogLaminate struct {
    client mrstorage.DbConn
    builder squirrel.StatementBuilderType
}

func NewCatalogLaminate(client mrstorage.DbConn,
                        queryBuilder squirrel.StatementBuilderType) *CatalogLaminate {
    return &CatalogLaminate{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogLaminate) LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateListFilter, rows *[]entity.CatalogLaminate) error {
    query := re.builder.
        Select(`
            laminate_id,
            tag_version,
            datetime_created,
            laminate_article,
            laminate_caption,
            type_id,
            laminate_length,
            laminate_weight,
            laminate_thickness,
            laminate_status`).
        From("public.catalog_laminates").
        Where(squirrel.NotEq{"laminate_status": mrcom.ItemStatusRemoved}).
        OrderBy("laminate_caption ASC, laminate_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"laminate_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        var row entity.CatalogLaminate

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Article,
            &row.Caption,
            &row.TypeId,
            &row.Length,
            &row.Weight,
            &row.Thickness,
            &row.Status,
        )

        *rows = append(*rows, row)
    }

    return cursor.Err()
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Article, Caption, TypeId, Length, Weight, Thickness, Status}
func (re *CatalogLaminate) LoadOne(ctx context.Context, row *entity.CatalogLaminate) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            laminate_article,
            laminate_caption,
            type_id,
            laminate_length,
            laminate_weight,
            laminate_thickness,
            laminate_status
        FROM
            public.catalog_laminates
        WHERE laminate_id = $1 AND laminate_status <> $2;`

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
        &row.TypeId,
        &row.Length,
        &row.Weight,
        &row.Thickness,
        &row.Status,
    )
}

func (re *CatalogLaminate) FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error) {
    sql := `
        SELECT laminate_id
        FROM
            public.catalog_laminates
        WHERE laminate_article = $1;`

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
func (re *CatalogLaminate) FetchStatus(ctx context.Context, row *entity.CatalogLaminate) (mrcom.ItemStatus, error) {
    sql := `
        SELECT laminate_status
        FROM
            public.catalog_laminates
        WHERE laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3;`

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
// uses: row{Article, Caption, TypeId, Length, Weight, Thickness, Status}
// modifies: row{Id}
func (re *CatalogLaminate) Insert(ctx context.Context, row *entity.CatalogLaminate) error {
    sql := `
        INSERT INTO public.catalog_laminates
            (laminate_article,
             laminate_caption,
             type_id,
             laminate_length,
             laminate_weight,
             laminate_thickness,
             laminate_status)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7)
        RETURNING laminate_id;`

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Article,
        row.Caption,
        row.TypeId,
        row.Length,
        row.Weight,
        row.Thickness,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Article, Caption, TypeId, Length, Weight, Thickness, Status}
func (re *CatalogLaminate) Update(ctx context.Context, row *entity.CatalogLaminate) error {
    filledFields, err := mrentity.FilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.catalog_laminates").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"laminate_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"laminate_status": mrcom.ItemStatusRemoved})

    return re.client.SqExec(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogLaminate) UpdateStatus(ctx context.Context, row *entity.CatalogLaminate) error {
    sql := `
        UPDATE public.catalog_laminates
        SET
            tag_version = tag_version + 1,
            laminate_status = $4
        WHERE
            laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogLaminate) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_laminates
        SET
            tag_version = tag_version + 1,
            laminate_article = NULL,
            laminate_status = $2
        WHERE
            laminate_id = $1 AND laminate_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )
}
