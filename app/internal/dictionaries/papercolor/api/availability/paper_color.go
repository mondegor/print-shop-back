package availability

import (
	"context"

	"print-shop-back/internal/adapter/workflow"
)

type (
	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error)
	}
)
