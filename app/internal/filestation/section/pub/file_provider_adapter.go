package pub

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmodel"
)

type (
	// FileProviderAdapterUseCase - comment interface.
	FileProviderAdapterUseCase interface {
		Get(ctx context.Context, filePath string) (mrmodel.File, error)
	}
)
