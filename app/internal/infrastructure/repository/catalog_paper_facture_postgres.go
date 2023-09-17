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

type CatalogPaperFacture struct {
    client *mrpostgres.ConnAdapter
    builder squirrel.StatementBuilderType
}

func NewCatalogPaperFacture(client *mrpostgres.ConnAdapter,
                            queryBuilder squirrel.StatementBuilderType) *CatalogPaperFacture {
    return &CatalogPaperFacture{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogPaperFacture) LoadAll(ctx context.Context, listFilter *entity.CatalogPaperFactureListFilter, rows *[]entity.CatalogPaperFacture) error {
    query := re.builder.
        Select(`
            facture_id,
            tag_version,
            datetime_created,
            facture_caption,
            facture_status`).
        From("public.catalog_paper_factures").
        Where(squirrel.NotEq{"facture_status": mrcom.ItemStatusRemoved}).
        OrderBy("facture_caption ASC, facture_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"facture_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        var row entity.CatalogPaperFacture

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
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
// modifies: row{Version, CreatedAt, Caption, Status}
func (re *CatalogPaperFacture) LoadOne(ctx context.Context, row *entity.CatalogPaperFacture) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            facture_caption,
            facture_status
        FROM
            public.catalog_paper_factures
        WHERE facture_id = $1 AND facture_status <> $2;`

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
func (re *CatalogPaperFacture) FetchStatus(ctx context.Context, row *entity.CatalogPaperFacture) (mrcom.ItemStatus, error) {
    sql := `
        SELECT facture_status
        FROM
            public.catalog_paper_factures
        WHERE facture_id = $1 AND tag_version = $2 AND facture_status <> $3;`

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
func (re *CatalogPaperFacture) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_paper_factures
        WHERE facture_id = $1 AND facture_status <> $2;`

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
func (re *CatalogPaperFacture) Insert(ctx context.Context, row *entity.CatalogPaperFacture) error {
    sql := `
        INSERT INTO public.catalog_paper_factures
            (facture_caption,
             facture_status)
        VALUES
            ($1, $2)
        RETURNING facture_id;`

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
func (re *CatalogPaperFacture) Update(ctx context.Context, row *entity.CatalogPaperFacture) error {
    sql := `
        UPDATE public.catalog_paper_factures
        SET
            tag_version = tag_version + 1,
            facture_caption = $4
        WHERE facture_id = $1 AND tag_version = $2 AND facture_status <> $3;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Caption,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogPaperFacture) UpdateStatus(ctx context.Context, row *entity.CatalogPaperFacture) error {
    sql := `
        UPDATE public.catalog_paper_factures
        SET
            tag_version = tag_version + 1,
            facture_status = $4
        WHERE
            facture_id = $1 AND tag_version = $2 AND facture_status <> $3;`

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

func (re *CatalogPaperFacture) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_paper_factures
        SET
            tag_version = tag_version + 1,
            facture_status = $2
        WHERE
            facture_id = $1 AND facture_status <> $2;`

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
