package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
)

type (
	// CompanyPageLogoPostgres - comment struct.
	CompanyPageLogoPostgres struct {
		repoMeta db.FieldUpdater[uuid.UUID, mrentity.ImageMeta]
	}
)

// NewCompanyPageLogoPostgres - создаёт объект CompanyPageLogoPostgres.
func NewCompanyPageLogoPostgres(client mrstorage.DBConnManager) *CompanyPageLogoPostgres {
	return &CompanyPageLogoPostgres{
		repoMeta: db.NewFieldUpdater[uuid.UUID, mrentity.ImageMeta](
			client,
			module.DBTableNameCompaniesPages,
			"account_id",
			"logo_meta",
			module.DBFieldWithoutDeletedAt,
		),
	}
}

// FetchMeta - comment method.
func (re *CompanyPageLogoPostgres) FetchMeta(ctx context.Context, accountID uuid.UUID) (mrentity.ImageMeta, error) {
	return re.repoMeta.Fetch(ctx, accountID)
}

// UpdateMeta - comment method.
func (re *CompanyPageLogoPostgres) UpdateMeta(ctx context.Context, accountID uuid.UUID, meta mrentity.ImageMeta) error {
	return re.repoMeta.Update(ctx, accountID, meta)
}

// DeleteMeta - comment method.
func (re *CompanyPageLogoPostgres) DeleteMeta(ctx context.Context, accountID uuid.UUID) error {
	return re.repoMeta.Update(ctx, accountID, mrentity.ImageMeta{})
}
