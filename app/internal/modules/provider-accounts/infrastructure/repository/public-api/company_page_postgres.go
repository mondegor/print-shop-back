package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

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

func (re *CompanyPagePostgres) FetchByRewriteName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error) {
	sql := `
        SELECT
            page_head,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            rewrite_name = $1 AND page_status IN ($2, $3)
        LIMIT 1;`

	var row entity.CompanyPage

	err := re.client.QueryRow(
		ctx,
		sql,
		rewriteName,
		entity_shared.PublicStatusPublished,
		entity_shared.PublicStatusPublishedShared,
	).Scan(
		&row.PageHead,
		&row.LogoURL,
		&row.SiteURL,
	)

	if err != nil {
		return nil, err
	}

	return &row, err
}
