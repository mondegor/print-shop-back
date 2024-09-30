package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
)

type (
	// CompanyPageLogoPostgres - comment struct.
	CompanyPageLogoPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewCompanyPageLogoPostgres - создаёт объект CompanyPageLogoPostgres.
func NewCompanyPageLogoPostgres(client mrstorage.DBConnManager) *CompanyPageLogoPostgres {
	return &CompanyPageLogoPostgres{
		client: client,
	}
}

// FetchMeta - comment method.
func (re *CompanyPageLogoPostgres) FetchMeta(ctx context.Context, accountID uuid.UUID) (mrentity.ImageMeta, error) {
	sql := `
        SELECT
            logo_meta
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        WHERE
            account_id = $1
        LIMIT 1;`

	var logoMeta mrentity.ImageMeta

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&logoMeta,
	)

	return logoMeta, err
}

// UpdateMeta - comment method.
func (re *CompanyPageLogoPostgres) UpdateMeta(ctx context.Context, accountID uuid.UUID, meta mrentity.ImageMeta) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameCompaniesPages + `
        SET
            updated_at = NOW(),
            logo_meta = $2
        WHERE
            account_id = $1;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		accountID,
		meta,
	)
}

// DeleteMeta - comment method.
func (re *CompanyPageLogoPostgres) DeleteMeta(ctx context.Context, accountID uuid.UUID) error {
	return re.UpdateMeta(ctx, accountID, mrentity.ImageMeta{})
}
