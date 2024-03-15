package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAdapterUseCase interface {
		Get(ctx context.Context, filePath string) (mrtype.File, error)
	}
)
