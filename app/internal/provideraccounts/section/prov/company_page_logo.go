package prov

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/mrentity"
	mrmodel "github.com/mondegor/go-core/mrmodel/media"
)

type (
	// CompanyPageLogoUseCase - comment interface.
	CompanyPageLogoUseCase interface {
		StoreFile(ctx context.Context, accountID uuid.UUID, image mrmodel.Image) error
		RemoveFile(ctx context.Context, accountID uuid.UUID) error
	}

	// CompanyPageLogoStorage - comment interface.
	CompanyPageLogoStorage interface {
		FetchMeta(ctx context.Context, accountID uuid.UUID) (mrentity.ImageMeta, error)
		UpdateMeta(ctx context.Context, accountID uuid.UUID, meta mrentity.ImageMeta) error
		DeleteMeta(ctx context.Context, accountID uuid.UUID) error
	}
)
