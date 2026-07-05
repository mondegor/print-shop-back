package usr

import (
	"context"

	"github.com/google/uuid"

	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// TerritoryStorage - comment interface.
	TerritoryStorage interface {
		FetchState(ctx context.Context, accountID uuid.UUID, rowID uint64) (row entity.TerritoryState, err error)
	}
)
