package availability

import (
	"context"

	"print-shop-back/internal/adapter/workflow"
)

type (
	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error)
	}
)
