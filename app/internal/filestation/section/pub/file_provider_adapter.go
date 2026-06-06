package pub

import (
	"context"

	mrmodel "github.com/mondegor/go-sysmess/mrmodel/media"
)

type (
	// FileProviderAdapterUseCase - comment interface.
	FileProviderAdapterUseCase interface {
		Get(ctx context.Context, filePath string) (mrmodel.File, error)
	}
)
