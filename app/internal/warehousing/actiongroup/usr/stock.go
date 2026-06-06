package usr

import (
	"context"

	"github.com/google/uuid"

	"print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// StockService - comment interface.
	StockService interface {
		GetList(ctx context.Context, params dto.StockParams) (items []entity.Stock, hasNext bool, err error)
	}

	// StockStorage - comment interface.
	StockStorage interface {
		FetchByCondition(ctx context.Context, params dto.StockParams) (rows []entity.Stock, hasNext bool, err error)
		FetchOne(ctx context.Context, accountID uuid.UUID, rowID uint64) (row entity.Stock, err error)
		CheckLocationAvailability(ctx context.Context, accountID uuid.UUID, locationID uint64) error
		InsertOrUpdate(ctx context.Context, row entity.Stock) (rowID uint64, err error)
		UpdateQuantity(ctx context.Context, accountID uuid.UUID, rowID uint64, quantity int) (newRowID uint64, err error)
		Delete(ctx context.Context, accountID uuid.UUID, rowID uint64) error
	}
)
