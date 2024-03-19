package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
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

func (re *CompanyPagePostgres) FetchOne(ctx context.Context, accountID mrtype.KeyString) (entity.CompanyPage, error) {
	sql := `
        SELECT
            updated_at,
            rewrite_name,
            page_head,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url,
            page_status
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	row := entity.CompanyPage{AccountID: accountID}

	err := re.client.QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&row.UpdatedAt,
		&row.RewriteName,
		&row.PageHead,
		&row.LogoURL,
		&row.SiteURL,
		&row.Status,
	)

	return row, err
}

func (re *CompanyPagePostgres) FetchStatus(ctx context.Context, row entity.CompanyPage) (enums.PublicStatus, error) {
	sql := `
        SELECT
            page_status
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	var status enums.PublicStatus

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
            updated_at = NOW(),
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

func (re *CompanyPagePostgres) UpdateStatus(ctx context.Context, row entity.CompanyPage) error {
	sql := `
        UPDATE
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
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
