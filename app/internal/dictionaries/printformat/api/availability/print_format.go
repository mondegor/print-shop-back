package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
	}
)
