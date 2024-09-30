package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

type (
	// CompanyPagePostgres - comment struct.
	CompanyPagePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewCompanyPagePostgres - создаёт объект CompanyPagePostgres.
func NewCompanyPagePostgres(client mrstorage.DBConnManager) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client: client,
	}
}

// FetchOne - comment method.
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
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        WHERE
            account_id = $1
        LIMIT 1;`

	row := entity.CompanyPage{AccountID: accountID}

	err := re.client.Conn(ctx).QueryRow(
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

// FetchAccountIDByRewriteName - comment method.
func (re *CompanyPagePostgres) FetchAccountIDByRewriteName(ctx context.Context, rewriteName string) (uuid.UUID, error) {
	sql := `
        SELECT
            account_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        WHERE
            rewrite_name = $1
        LIMIT 1;`

	var accountID uuid.UUID

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rewriteName,
	).Scan(
		&accountID,
	)

	return accountID, err
}

// FetchStatus - comment method.
// result: enums.PublicStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *CompanyPagePostgres) FetchStatus(ctx context.Context, accountID uuid.UUID) (enum.PublicStatus, error) {
	sql := `
        SELECT
            page_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        WHERE
            account_id = $1
        LIMIT 1;`

	var status enum.PublicStatus

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&status,
	)

	return status, err
}

// InsertOrUpdate - comment method.
// WARNING: row.Status uses only for insert.
func (re *CompanyPagePostgres) InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error {
	return re.client.Do(ctx, func(ctx context.Context) error {
		conn := re.client.Conn(ctx)

		sql := `
	        UPDATE
    	        ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        	SET
            	updated_at = NOW(),
            	rewrite_name = $2,
            	page_title = $3,
            	site_url = $4
        	WHERE
            	account_id = $1;`

		err := conn.Exec(
			ctx,
			sql,
			row.AccountID,
			row.RewriteName,
			row.PageTitle,
			row.SiteURL,
		)
		// если сохранение удачное или если это системная ошибка
		if err == nil || !mrcore.ErrStorageRowsNotAffected.Is(err) {
			return err
		}

		sql = `
            INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
                (
                    account_id,
                    rewrite_name,
                    page_title,
                    site_url,
                    page_status
                )
            VALUES
                ($1, $2, $3, $4, $5);`

		return conn.Exec(
			ctx,
			sql,
			row.AccountID,
			row.RewriteName,
			row.PageTitle,
			row.SiteURL,
			row.Status,
		)
	})
}

// UpdateStatus - comment method.
func (re *CompanyPagePostgres) UpdateStatus(ctx context.Context, row entity.CompanyPage) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        SET
            updated_at = NOW(),
            page_status = $2
        WHERE
            account_id = $1;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		row.AccountID,
		row.Status,
	)
}
