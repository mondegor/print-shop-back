package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
	}
)
