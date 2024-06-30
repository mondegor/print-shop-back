package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FileProviderAdapterUseCase - comment interface.
	FileProviderAdapterUseCase interface {
		Get(ctx context.Context, filePath string) (mrtype.File, error)
	}
)
