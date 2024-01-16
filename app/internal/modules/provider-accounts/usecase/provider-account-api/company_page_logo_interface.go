package usecase

import (
	"context"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CompanyPageLogoService interface {
		StoreFile(ctx context.Context, accountID mrtype.KeyString, image mrtype.File) error
		RemoveFile(ctx context.Context, accountID mrtype.KeyString) error
	}

	CompanyPageLogoStorage interface {
		FetchMeta(ctx context.Context, accountID mrtype.KeyString) (mrentity.ImageMeta, error)
		UpdateMeta(ctx context.Context, accountID mrtype.KeyString, meta mrentity.ImageMeta) error
		DeleteMeta(ctx context.Context, accountID mrtype.KeyString) error
	}
)
