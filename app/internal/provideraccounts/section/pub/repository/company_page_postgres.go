package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/entity"

	"github.com/mondegor/go-storage/mrstorage"
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

// FetchByRewriteName - comment method.
func (re *CompanyPagePostgres) FetchByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error) {
	sql := `
        SELECT
            page_title,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        WHERE
            rewrite_name = $1 AND page_status IN ($2, $3)
        LIMIT 1;`

	var row entity.CompanyPage

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rewriteName,
		enum.PublicStatusPublished,
		enum.PublicStatusPublishedShared,
	).Scan(
		&row.PageTitle,
		&row.LogoURL,
		&row.SiteURL,
	)

	return row, err
}
