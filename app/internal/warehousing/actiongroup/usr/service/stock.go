package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"

	"print-shop-back/internal/warehousing/actiongroup/usr"
	"print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
)

type (
	// Stock - comment struct.
	Stock struct {
		storageStock usr.StockStorage
		errorWrapper errors.Wrapper
	}
)

// NewStock - создаёт объект Stock.
func NewStock(
	storageStock usr.StockStorage,
) *Stock {
	return &Stock{
		storageStock: storageStock,
		errorWrapper: errors.NewServiceOperationFailedWrapper(),
	}
}

// GetList - comment method.
func (uc *Stock) GetList(ctx context.Context, params dto.StockParams) (items []entity.Stock, hasNext bool, err error) {
	if params.AccountID == uuid.Nil {
		return nil, false, errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	items, hasNext, err = uc.storageStock.FetchByCondition(ctx, params)
	if err != nil {
		return nil, false, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Stock, 0), false, nil
	}

	return items, hasNext, nil
}
