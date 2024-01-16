package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CompanyPageLogoPostgres struct {
		client mrstorage.DBConn
	}
)

func NewCompanyPageLogoPostgres(client mrstorage.DBConn) *CompanyPageLogoPostgres {
	return &CompanyPageLogoPostgres{
		client: client,
	}
}

func (re *CompanyPageLogoPostgres) FetchMeta(ctx context.Context, accountID mrtype.KeyString) (mrentity.ImageMeta, error) {
	sql := `
        SELECT
            logo_meta
        FROM
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        WHERE
            account_id = $1
        LIMIT 1;`

	var logoMeta mrentity.ImageMeta

	err := re.client.QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&logoMeta,
	)

	return logoMeta, err
}

func (re *CompanyPageLogoPostgres) UpdateMeta(ctx context.Context, accountID mrtype.KeyString, meta mrentity.ImageMeta) error {
	sql := `
        UPDATE
            ` + module.UnitCompanyPageDBSchema + `.companies_pages
        SET
            datetime_updated = NOW(),
            logo_meta = $2
        WHERE
            account_id = $1;`

	return re.client.Exec(
		ctx,
		sql,
		accountID,
		meta,
	)
}

func (re *CompanyPageLogoPostgres) DeleteMeta(ctx context.Context, accountID mrtype.KeyString) error {
	return re.UpdateMeta(ctx, accountID, mrentity.ImageMeta{})
}
