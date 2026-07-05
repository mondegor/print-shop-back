package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrevent"

	"print-shop-back/internal/warehousing/actiongroup/usr"
	"print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"print-shop-back/internal/warehousing/enum/locationkind"
	"print-shop-back/internal/warehousing/module"
)

type (
	// TransferStockContainer - comment struct.
	TransferStockContainer struct {
		storageStock         usr.StockStorage
		useCaseRefreshStores refreshStoresUseCase
		eventEmitter         mrevent.Emitter
		errorWrapper         errors.Wrapper
	}
)

// NewTransferStockContainer - создаёт объект TransferStockContainer.
func NewTransferStockContainer(
	storageStock usr.StockStorage,
	useCaseRefreshStores refreshStoresUseCase,
	eventEmitter mrevent.Emitter,
) *TransferStockContainer {
	return &TransferStockContainer{
		storageStock:         storageStock,
		useCaseRefreshStores: useCaseRefreshStores,
		eventEmitter:         eventEmitter,
		errorWrapper:         errors.NewServiceOperationFailedWrapper(),
	}
}

// Execute - comment method.
func (uc *TransferStockContainer) Execute(ctx context.Context, item dto.TransferStockContainers) error {
	if item.AccountID == uuid.Nil {
		return errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	if len(item.Stocks) == 0 {
		return nil
	}

	for i := range item.Stocks {
		if item.Stocks[i].StockID == 0 {
			return errors.ErrInternalIncorrectInputData.WithDetails("stockId is zero", "index", i)
		}

		if item.Stocks[i].Quantity == 0 {
			return errors.ErrInternalIncorrectInputData.WithDetails("itemQuantity is zero", "index", i)
		}
	}

	// TODO: транзакция

	locations := make([]uint64, 0, len(item.Stocks))

	for i := range item.Stocks {
		// TODO: заменить на FetchByIDs
		stock, err := uc.storageStock.FetchOne(ctx, item.AccountID, item.Stocks[i].StockID)
		if err != nil {
			return uc.errorWrapper.Wrap(err)
		}

		if locationkind.Is(stock.ContainerID, locationkind.Group) {
			// как списывать групповые контейнеры? варианты:
			// ++++1. не давать списывать непустые контейнеры
			// 2. списывать и привязывать находящиеся контейнеры к текущей локации
			// 3. списывать вместе с контейнерами, которые в нём находятся
			if err := uc.storageStock.CheckLocationAvailability(ctx, item.AccountID, stock.LocationID); err != nil {
				if errors.Is(err, module.ErrLocationIsOccupied) {
					return errors.ErrIncorrectInputData.New(err) // user error
				}

				return uc.errorWrapper.Wrap(err)
			}
		}

		if item.Stocks[i].Quantity > stock.ContainerQuantity {
			return errors.ErrInternalIncorrectInputData.WithDetails("itemQuantity > stock.Quantity", "index", i)
		}

		stock.ContainerQuantity -= item.Stocks[i].Quantity

		if stock.ContainerQuantity > 0 {
			if stock.ID, err = uc.storageStock.UpdateQuantity(ctx, stock.AccountID, stock.ID, stock.ContainerQuantity); err != nil {
				return uc.errorWrapper.Wrap(err)
			}
		} else {
			if err = uc.storageStock.Delete(ctx, stock.AccountID, stock.ID); err != nil {
				return uc.errorWrapper.Wrap(err)
			}
		}

		// TODO: добавить лог списания

		locations = append(locations, stock.LocationID)
	}

	// TODO: через событие сделать
	if err := uc.useCaseRefreshStores.Execute(ctx, locations); err != nil {
		return uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Transfer", "accountId", item.AccountID, "stocks", item.Stocks)

	return nil
}
