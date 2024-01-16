package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAdapterService interface {
		Get(ctx context.Context, filePath string) (mrtype.File, error)
	}
)
