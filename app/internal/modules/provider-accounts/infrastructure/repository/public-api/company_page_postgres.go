package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/mondegor/go-storage/mrstorage"
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

func (re *CompanyPagePostgres) FetchByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error) {
	sql := `
        SELECT
            page_title,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url
        FROM
            ` + module.DBSchema + `.companies_pages
        WHERE
            rewrite_name = $1 AND page_status IN ($2, $3)
        LIMIT 1;`

	var row entity.CompanyPage

	err := re.client.QueryRow(
		ctx,
		sql,
		rewriteName,
		enums.PublicStatusPublished,
		enums.PublicStatusPublishedShared,
	).Scan(
		&row.PageTitle,
		&row.LogoURL,
		&row.SiteURL,
	)

	return row, err
}
