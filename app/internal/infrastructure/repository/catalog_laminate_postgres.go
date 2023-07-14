package repository

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/client/mrpostgres"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"

    "github.com/Masterminds/squirrel"
)

type CatalogLaminate struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogLaminate(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *CatalogLaminate {
    return &CatalogLaminate{
        client: client,
        builder: queryBuilder,
    }
}

func (f *CatalogLaminate) LoadAll(ctx context.Context, listFilter *entity.CatalogLaminateListFilter, rows *[]entity.CatalogLaminate) error {
    tbl := f.builder.
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
        Where(squirrel.NotEq{"laminate_status": entity.ItemStatusRemoved}).
        OrderBy("laminate_caption ASC, laminate_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"laminate_status": listFilter.Statuses})
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
// modifies: row{Version, CreatedAt, Article, Caption, TypeId, Length, Weight, Thickness, Status}
func (f *CatalogLaminate) LoadOne(ctx context.Context, row *entity.CatalogLaminate) error {
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

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
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

// FetchIdByArticle
// uses: row{Article}
func (f *CatalogLaminate) FetchIdByArticle(ctx context.Context, row *entity.CatalogLaminate) (mrentity.KeyInt32, error) {
    sql := `
        SELECT laminate_id
        FROM
            public.catalog_laminates
        WHERE laminate_article = $1;`

    var id mrentity.KeyInt32

    err := f.client.QueryRow(ctx, sql, row.Article).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (f *CatalogLaminate) FetchStatus(ctx context.Context, row *entity.CatalogLaminate) (entity.ItemStatus, error) {
    sql := `
        SELECT laminate_status
        FROM
            public.catalog_laminates
        WHERE laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3;`

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
// uses: row{Article, Caption, TypeId, Length, Weight, Thickness, Status}
// modifies: row{Id}
func (f *CatalogLaminate) Insert(ctx context.Context, row *entity.CatalogLaminate) error {
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

    err := f.client.QueryRow(
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
func (f *CatalogLaminate) Update(ctx context.Context, row *entity.CatalogLaminate) error {
    sql := `
        UPDATE public.catalog_laminates
        SET
            tag_version = tag_version + 1,
            laminate_article = $4,
            laminate_caption = $5,
            type_id = $6,
            laminate_length = $7,
            laminate_weight = $8,
            laminate_thickness = $9
        WHERE laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Article,
        row.Caption,
        row.TypeId,
        row.Length,
        row.Weight,
        row.Thickness,
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
func (f *CatalogLaminate) UpdateStatus(ctx context.Context, row *entity.CatalogLaminate) error {
    sql := `
        UPDATE public.catalog_laminates
        SET
            tag_version = tag_version + 1,
            laminate_status = $4
        WHERE
            laminate_id = $1 AND tag_version = $2 AND laminate_status <> $3;`

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

func (f *CatalogLaminate) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_laminates
        SET
            tag_version = tag_version + 1,
            laminate_status = $2
        WHERE
            laminate_id = $1 AND laminate_status <> $2;`

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
