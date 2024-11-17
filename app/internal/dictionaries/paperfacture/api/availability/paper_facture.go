package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
	}
)
