package pub

import (
	"context"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/entity"
)

type (
	// QueryResultUseCase - comment interface.
	QueryResultUseCase interface {
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.QueryHistoryItem, error)
		Create(ctx context.Context, item entity.QueryHistoryItem) (uuid.UUID, error)
	}

	// QueryResultStorage - comment interface.
	QueryResultStorage interface {
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.QueryHistoryItem, error)
		Insert(ctx context.Context, row entity.QueryHistoryItem) (id uuid.UUID, err error)
		UpdateQuantity(ctx context.Context, rowID uuid.UUID) error
	}
)
