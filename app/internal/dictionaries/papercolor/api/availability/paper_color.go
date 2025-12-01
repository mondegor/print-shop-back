package availability

import (
	"context"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
)

type (
	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error)
	}
)
