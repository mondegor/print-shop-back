package repository

import (
    "context"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type (
    CompanyPageLogo struct {
        client mrstorage.DbConn
        builder squirrel.StatementBuilderType
    }
)

func NewCompanyPageLogo(client mrstorage.DbConn,
                        queryBuilder squirrel.StatementBuilderType) *CompanyPageLogo {
    return &CompanyPageLogo{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CompanyPageLogo) Fetch(ctx context.Context, accountId mrentity.KeyString) (string, error) {
    sql := `
        SELECT logo_path
        FROM
            public.accounts_companies_pages
        WHERE account_id = $1;`

    var logoPath string

    err := re.client.QueryRow(
        ctx,
        sql,
        accountId,
    ).Scan(
        &logoPath,
    )

    return logoPath, err
}

func (re *CompanyPageLogo) Update(ctx context.Context, accountId mrentity.KeyString, logoPath string) error {
    sql := `
        UPDATE public.accounts_companies_pages
        SET
            datetime_updated = NOW(),
            logo_path = $2
        WHERE account_id = $1;`

    return re.client.Exec(
        ctx,
        sql,
        accountId,
        logoPath,
    )
}

func (re *CompanyPageLogo) Delete(ctx context.Context, accountId mrentity.KeyString) error {
    return re.Update(ctx, accountId, "")
}
