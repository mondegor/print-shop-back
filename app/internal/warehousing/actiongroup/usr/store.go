package usr

import (
	"context"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// StoreService - comment interface.
	StoreService interface {
		GetList(ctx context.Context, params dto.StoreParams) (items []entity.Store, hasNext bool, err error)
		CheckAvailability(ctx context.Context, accountID uuid.UUID, rowID uint64) error
	}

	// StoreStorage - comment interface.
	StoreStorage interface {
		FetchByCondition(ctx context.Context, params dto.StoreParams) (rows []entity.Store, hasNext bool, err error)
		FetchOne(ctx context.Context, accountID uuid.UUID, rowID uint64) (row entity.Store, err error)
		FetchState(ctx context.Context, accountID uuid.UUID, rowID uint64) (row entity.StoreState, err error)
	}
)
