package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
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

func (re *CompanyPagePostgres) LoadOne(ctx context.Context, row *entity.CompanyPage) error {
	sql := `
        SELECT
            datetime_updated,
            rewrite_name,
            page_head,
            COALESCE(logo_meta ->> 'path', '') as logo_url,
            site_url,
            page_status
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.AccountID,
	).Scan(
		&row.UpdatedAt,
		&row.RewriteName,
		&row.PageHead,
		&row.LogoURL,
		&row.SiteURL,
		&row.Status,
	)
}

func (re *CompanyPagePostgres) FetchStatus(ctx context.Context, row *entity.CompanyPage) (entity_shared.PublicStatus, error) {
	sql := `
        SELECT
            page_status
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	var status entity_shared.PublicStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		row.AccountID,
	).Scan(
		&status,
	)

	return status, err
}

// InsertOrUpdate
// WARNING: row.Status uses only for insert
func (re *CompanyPagePostgres) InsertOrUpdate(ctx context.Context, row *entity.CompanyPage) error {
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
        SELECT
            1
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	var value int

	err = tx.QueryRow(
		ctx,
		sql,
		row.AccountID,
	).Scan(
		&value,
	)

	if err != nil {
		if !mrcore.FactoryErrStorageNoRowFound.Is(err) {
			return err
		}

		sql = `
            INSERT INTO ` + module.UnitCompanyPageDBSchema + `.companies_pages
                (
                    account_id,
                    rewrite_name,
                    page_head,
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
			row.PageHead,
			row.SiteURL,
			row.Status,
		)

		if err != nil {
			return err
		}
	}

	sql = `
        UPDATE
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        SET
            datetime_updated = NOW(),
            rewrite_name = $2,
            page_head = $3,
            site_url = $4
        WHERE
            account_id = $1;`

	// WARNING: err is used in defer func()
	err = tx.Exec(
		ctx,
		sql,
		row.AccountID,
		row.RewriteName,
		row.PageHead,
		row.SiteURL,
	)

	return err
}

func (re *CompanyPagePostgres) UpdateStatus(ctx context.Context, row *entity.CompanyPage) error {
	sql := `
        UPDATE
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        SET
            datetime_updated = NOW(),
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
