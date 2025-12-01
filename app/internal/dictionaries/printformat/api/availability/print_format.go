package availability

import (
	"context"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
)

type (
	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error)
	}
)
