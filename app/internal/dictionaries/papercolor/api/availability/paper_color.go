package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
	}
)
