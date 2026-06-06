package availability

import (
	"context"

	"print-shop-back/internal/adapter/workflow"
)

type (
	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error)
	}
)
