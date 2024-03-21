package repository

import (
	"context"
	module "print-shop-back/internal/modules/provider-accounts"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
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

func (re *CompanyPageLogoPostgres) FetchMeta(ctx context.Context, accountID uuid.UUID) (mrentity.ImageMeta, error) {
	sql := `
        SELECT
            logo_meta
        FROM
            ` + module.DBSchema + `.companies_pages
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

func (re *CompanyPageLogoPostgres) UpdateMeta(ctx context.Context, accountID uuid.UUID, meta mrentity.ImageMeta) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.companies_pages
        SET
            updated_at = NOW(),
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

func (re *CompanyPageLogoPostgres) DeleteMeta(ctx context.Context, accountID uuid.UUID) error {
	return re.UpdateMeta(ctx, accountID, mrentity.ImageMeta{})
}
