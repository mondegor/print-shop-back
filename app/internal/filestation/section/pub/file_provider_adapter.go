package pub

import (
	"context"

	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// FileProviderAdapterUseCase - comment interface.
	FileProviderAdapterUseCase interface {
		Get(ctx context.Context, filePath string) (mrtype.File, error)
	}
)
