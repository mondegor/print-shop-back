package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

type (
    CompanyPage struct {
        client mrstorage.DbConn
        builder squirrel.StatementBuilderType
    }
)

func NewCompanyPage(client mrstorage.DbConn,
                    queryBuilder squirrel.StatementBuilderType) *CompanyPage {
    return &CompanyPage{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CompanyPage) LoadAll(ctx context.Context, listFilter *entity.CompanyPageListFilter, rows *[]entity.CompanyPage) error {
    query := re.builder.
        Select(`
            account_id,
            tag_version,
            datetime_updated,
            rewrite_name,
            page_head,
            logo_path,
            site_url,
            page_status`).
        From("public.accounts_companies_pages").
        OrderBy("page_head ASC, account_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"page_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        var row entity.CompanyPage

        err = cursor.Scan(
            &row.AccountId,
            &row.Version,
            &row.UpdateAt,
            &row.RewriteName,
            &row.PageHead,
            &row.LogoPath,
            &row.SiteUrl,
            &row.Status,
        )

        *rows = append(*rows, row)
    }

    return cursor.Err()
}

// LoadOne
// uses: row{AccountId}
// modifies: row{Version, UpdateAt, RewriteName, PageHead, LogoPath, SiteUrl, Status}
func (re *CompanyPage) LoadOne(ctx context.Context, row *entity.CompanyPage) error {
    sql := `
        SELECT
            tag_version,
            datetime_updated,
            rewrite_name,
            page_head,
            logo_path,
            site_url,
            page_status
        FROM
            public.accounts_companies_pages
        WHERE account_id = $1;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.AccountId,
    ).Scan(
        &row.Version,
        &row.UpdateAt,
        &row.RewriteName,
        &row.PageHead,
        &row.LogoPath,
        &row.SiteUrl,
        &row.Status,
    )
}

// LoadOneByRewriteName
// uses: row{RewriteName}
// modifies: row{PageHead, LogoPath, SiteUrl}
func (re *CompanyPage) LoadOneByRewriteName(ctx context.Context, row *entity.CompanyPage) error {
    sql := `
        SELECT
            page_head,
            logo_path,
            site_url
        FROM
            public.accounts_companies_pages
        WHERE rewrite_name = $1 AND page_status = $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.RewriteName,
        entity.ResourceStatusPublished,
    ).Scan(
        &row.PageHead,
        &row.LogoPath,
        &row.SiteUrl,
    )
}

// FetchStatus
// uses: row{AccountId, Version}
func (re *CompanyPage) FetchStatus(ctx context.Context, row *entity.CompanyPage) (entity.ResourceStatus, error) {
    sql := `
        SELECT page_status
        FROM
            public.accounts_companies_pages
        WHERE account_id = $1 AND tag_version = $2;`

    var status entity.ResourceStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.AccountId,
        row.Version,
    ).Scan(
        &status,
    )

    return status, err
}

// InsertOrUpdate
// uses: row{AccountId, Version, RewriteName, PageHead, SiteUrl, Status}
// WARNING: row.Status uses only for insert
func (re *CompanyPage) InsertOrUpdate(ctx context.Context, row *entity.CompanyPage) error {
    tx, err := re.client.Begin(ctx)

    if err != nil {
        return err // :TODO:
    }

    defer func() {
        var e error

        if err != nil {
            e = tx.Rollback(ctx)
        } else {
            e = tx.Commit(ctx)
        }

        if e != nil {
            mrctx.Logger(ctx).Err(e)
        }
    }()

    sql := `
        SELECT tag_version
        FROM public.accounts_companies_pages
        WHERE account_id = $1;`

    var version mrentity.Version

    err = tx.QueryRow(
        ctx,
        sql,
        row.AccountId,
    ).Scan(
        &version,
    )

    if err != nil {
        if !mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return err
        }

        sql = `
            INSERT INTO public.accounts_companies_pages
                (account_id,
                 rewrite_name,
                 page_head,
                 site_url,
                 page_status)
            VALUES
                ($1, $2, $3, $4, $5);`

        return tx.Exec(
            ctx,
            sql,
            row.AccountId,
            row.RewriteName,
            row.PageHead,
            row.SiteUrl,
            row.Status,
        )
    }

    //if row.Version != version {
    //   return fmt.Errorf("version %s <> %s", row.Version, version)
    //}

    sql = `
        UPDATE public.accounts_companies_pages
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            rewrite_name = $3,
            page_head = $4,
            site_url = $5
        WHERE account_id = $1 AND tag_version = $2;`

    return tx.Exec(
        ctx,
        sql,
        row.AccountId,
        row.Version,
        row.RewriteName,
        row.PageHead,
        row.SiteUrl,
    )
}

// UpdateStatus
// uses: row{AccountId, Version, Status}
func (re *CompanyPage) UpdateStatus(ctx context.Context, row *entity.CompanyPage) error {
    sql := `
        UPDATE public.accounts_companies_pages
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            page_status = $3
        WHERE
            account_id = $1 AND tag_version = $2;`

    return re.client.Exec(
        ctx,
        sql,
        row.AccountId,
        row.Version,
        row.Status,
    )
}
