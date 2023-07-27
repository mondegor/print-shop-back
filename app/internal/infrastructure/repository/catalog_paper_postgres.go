package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type CatalogPaper struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewCatalogPaper(client *mrpostgres.Connection,
                     queryBuilder squirrel.StatementBuilderType) *CatalogPaper {
    return &CatalogPaper{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogPaper) LoadAll(ctx context.Context, listFilter *entity.CatalogPaperListFilter, rows *[]entity.CatalogPaper) error {
    query := re.builder.
        Select(`
            paper_id,
            tag_version,
            datetime_created,
            paper_article,
            paper_caption,
            paper_length,
            paper_width,
            paper_density,
            color_id,
            facture_id,
            paper_thickness,
            paper_sides,
            paper_status`).
        From("public.catalog_papers").
        Where(squirrel.NotEq{"paper_status": entity.ItemStatusRemoved}).
        OrderBy("paper_caption ASC, paper_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"paper_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    for cursor.Next() {
        var row entity.CatalogPaper

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Article,
            &row.Caption,
            &row.Length,
            &row.Width,
            &row.Density,
            &row.ColorId,
            &row.FactureId,
            &row.Thickness,
            &row.Sides,
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
func (re *CatalogPaper) LoadOne(ctx context.Context, row *entity.CatalogPaper) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            paper_article,
            paper_caption,
            paper_length,
            paper_width,
            paper_density,
            color_id,
            facture_id,
            paper_thickness,
            paper_sides,
            paper_status
        FROM
            public.catalog_papers
        WHERE paper_id = $1 AND paper_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        entity.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Article,
        &row.Caption,
        &row.Length,
        &row.Width,
        &row.Density,
        &row.ColorId,
        &row.FactureId,
        &row.Thickness,
        &row.Sides,
        &row.Status,
    )
}

// FetchIdByArticle
// uses: row{Article}
func (re *CatalogPaper) FetchIdByArticle(ctx context.Context, row *entity.CatalogPaper) (mrentity.KeyInt32, error) {
    sql := `
        SELECT paper_id
        FROM
            public.catalog_papers
        WHERE paper_article = $1;`

    var id mrentity.KeyInt32

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Article,
    ).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogPaper) FetchStatus(ctx context.Context, row *entity.CatalogPaper) (entity.ItemStatus, error) {
    sql := `
        SELECT paper_status
        FROM
            public.catalog_papers
        WHERE paper_id = $1 AND tag_version = $2 AND paper_status <> $3;`

    var status entity.ItemStatus

    err := re.client.QueryRow(
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
func (re *CatalogPaper) Insert(ctx context.Context, row *entity.CatalogPaper) error {
    sql := `
        INSERT INTO public.catalog_papers
            (paper_article,
             paper_caption,
             paper_length,
             paper_width,
             paper_density,
             color_id,
             facture_id,
             paper_thickness,
             paper_sides,
             paper_status)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING paper_id;`

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Article,
        row.Caption,
        row.Length,
        row.Width,
        row.Density,
        row.ColorId,
        row.FactureId,
        row.Thickness,
        row.Sides,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Article, Caption, Length, Width, Density, ColorId, FactureId, Thickness, Sides, Status}
func (re *CatalogPaper) Update(ctx context.Context, row *entity.CatalogPaper) error {
    filledFields, err := mrentity.GetFilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.catalog_papers").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"paper_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"paper_status": entity.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogPaper) UpdateStatus(ctx context.Context, row *entity.CatalogPaper) error {
    sql := `
        UPDATE public.catalog_papers
        SET
            tag_version = tag_version + 1,
            paper_status = $4
        WHERE
            paper_id = $1 AND tag_version = $2 AND paper_status <> $3;`

    commandTag, err := re.client.Exec(
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

func (re *CatalogPaper) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_papers
        SET
            tag_version = tag_version + 1,
            paper_article = NULL,
            paper_status = $2
        WHERE
            paper_id = $1 AND paper_status <> $2;`

    commandTag, err := re.client.Exec(
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
