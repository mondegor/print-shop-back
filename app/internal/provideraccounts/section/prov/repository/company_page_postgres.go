package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

type (
	// CompanyPagePostgres - comment struct.
	CompanyPagePostgres struct {
		client            mrstorage.DBConnManager
		repoByRewriteName db.FieldFetcher[string, uuid.UUID]
		repoStatus        db.FieldUpdater[uuid.UUID, enum.PublicStatus]
	}
)

// NewCompanyPagePostgres - создаёт объект CompanyPagePostgres.
func NewCompanyPagePostgres(client mrstorage.DBConnManager) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client: client,
		repoByRewriteName: db.NewFieldFetcher[string, uuid.UUID](
			client,
			module.DBTableNameCompaniesPages,
			"rewrite_name",
			"account_id",
			module.DBFieldWithoutDeletedAt,
		),
		repoStatus: db.NewFieldUpdater[uuid.UUID, enum.PublicStatus](
			client,
			module.DBTableNameCompaniesPages,
			"account_id",
			"page_status",
			module.DBFieldWithoutDeletedAt,
		),
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
            ` + module.DBTableNameCompaniesPages + `
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
func (re *CompanyPagePostgres) FetchAccountIDByRewriteName(ctx context.Context, rewriteName string) (rowID uuid.UUID, err error) {
	return re.repoByRewriteName.Fetch(ctx, rewriteName)
}

// FetchStatus - comment method.
// result: enums.PublicStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *CompanyPagePostgres) FetchStatus(ctx context.Context, accountID uuid.UUID) (enum.PublicStatus, error) {
	return re.repoStatus.Fetch(ctx, accountID)
}

// InsertOrUpdate - comment method.
// WARNING: row.Status uses only for insert.
func (re *CompanyPagePostgres) InsertOrUpdate(ctx context.Context, row entity.CompanyPage) error {
	return re.client.Do(ctx, func(ctx context.Context) error {
		conn := re.client.Conn(ctx)

		sql := `
	        UPDATE
    	        ` + module.DBTableNameCompaniesPages + `
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
            INSERT INTO ` + module.DBTableNameCompaniesPages + `
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
	return re.repoStatus.Update(ctx, row.AccountID, row.Status)
}
