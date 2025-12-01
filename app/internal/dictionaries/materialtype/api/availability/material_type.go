package availability

import (
	"context"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
)

type (
	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error)
	}
)
