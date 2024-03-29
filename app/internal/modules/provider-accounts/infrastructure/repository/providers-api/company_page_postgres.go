package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/providers-api"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	CompanyPagePostgres struct {
		client mrstorage.DBConn
	}
)

func NewCompanyPagePostgres(
	client mrstorage.DBConn,
) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client: client,
	}
}

func (re *CompanyPagePostgres) FetchOne(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error) {
	sql := `
        SELECT
            rewrite_name,
            page_title,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url,
            page_status,
			created_at,
            updated_at
        FROM
            ` + module.DBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	row := entity.CompanyPage{AccountID: accountID}

	err := re.client.QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&row.RewriteName,
		&row.PageTitle,
		&row.LogoURL,
		&row.SiteURL,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (re *CompanyPagePostgres) FetchAccountIdByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error) {
	sql := `
        SELECT
            account_id
        FROM
            ` + module.DBSchema + `.companies_pages
        WHERE
            rewrite_name = $1
        LIMIT 1;`

	var accountID uuid.UUID

	err := re.client.QueryRow(
		ctx,
		sql,
		rewriteName,
	).Scan(
		&accountID,
	)

	return accountID, err
}

// FetchStatus
// result: enums.PublicStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CompanyPagePostgres) FetchStatus(ctx context.Context, accountID uuid.UUID) (enums.PublicStatus, error) {
	sql := `
        SELECT
            page_status
        FROM
            ` + module.DBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	var status enums.PublicStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&status,
	)

	return status, err
}

// InsertOrUpdate
// WARNING: row.Status uses only for insert
func (re *CompanyPagePostgres) InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error {
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
			mrlog.Ctx(ctx).Error().Err(e).Msg("defer func()")
		}
	}()

	sql := `
        SELECT
            account_id
        FROM
            ` + module.DBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	err = tx.QueryRow(
		ctx,
		sql,
		row.AccountID,
	).Scan(
		&row.AccountID,
	)

	if err == nil {
		sql = `
        UPDATE
            ` + module.DBSchema + `.companies_pages
        SET
            updated_at = NOW(),
            rewrite_name = $2,
            page_title = $3,
            site_url = $4
        WHERE
            account_id = $1;`

		// WARNING: err is used in defer func()
		err = tx.Exec(
			ctx,
			sql,
			row.AccountID,
			row.RewriteName,
			row.PageTitle,
			row.SiteURL,
		)

		return err
	}

	if mrcore.FactoryErrStorageNoRowFound.Is(err) {
		sql = `
            INSERT INTO ` + module.DBSchema + `.companies_pages
                (
                    account_id,
                    rewrite_name,
                    page_title,
                    site_url,
                    page_status
                )
            VALUES
                ($1, $2, $3, $4, $5);`

		// WARNING: err is used in defer func()
		err = tx.Exec(
			ctx,
			sql,
			row.AccountID,
			row.RewriteName,
			row.PageTitle,
			row.SiteURL,
			row.Status,
		)
	}

	return err
}

func (re *CompanyPagePostgres) UpdateStatus(ctx context.Context, row entity.CompanyPage) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.companies_pages
        SET
            updated_at = NOW(),
            page_status = $2
        WHERE
            account_id = $1;`

	return re.client.Exec(
		ctx,
		sql,
		row.AccountID,
		row.Status,
	)
}
