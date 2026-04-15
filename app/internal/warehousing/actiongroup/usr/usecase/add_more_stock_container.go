package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
)

type (
	// AddMoreStockContainer - comment struct.
	AddMoreStockContainer struct {
		storageContainer     usr.ContainerStorage
		storageStock         usr.StockStorage
		serviceStore         usr.StoreService
		useCaseRefreshStores refreshStoresUseCase
		eventEmitter         mrevent.Emitter
		errorWrapper         errors.Wrapper
		errorNotFoundWrapper errors.Wrapper
	}

	refreshStoresUseCase interface {
		Execute(ctx context.Context, storeIDs []uint64) error
	}
)

// NewAddMoreStockContainer - создаёт объект AddMoreStockContainer.
func NewAddMoreStockContainer(
	storageContainer usr.ContainerStorage,
	storageStock usr.StockStorage,
	serviceStore usr.StoreService,
	useCaseRefreshStores refreshStoresUseCase,
	eventEmitter mrevent.Emitter,
) *AddMoreStockContainer {
	return &AddMoreStockContainer{
		storageContainer:     storageContainer,
		storageStock:         storageStock,
		serviceStore:         serviceStore,
		useCaseRefreshStores: useCaseRefreshStores,
		eventEmitter:         eventEmitter,
		errorWrapper:         errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// Execute - comment method.
func (uc *AddMoreStockContainer) Execute(ctx context.Context, item dto.AddMoreStockContainer) (stockID uint64, err error) {
	if item.AccountID == uuid.Nil {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	if item.ContainerID == 0 {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("containerId is zero")
	}

	if item.LocationID == 0 {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("locationId is zero")
	}

	if item.Quantity == 0 {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("itemQuantity is zero")
	}

	if !locationkind.Is(item.ContainerID, locationkind.Container) {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("only containers can be add")
	}

	container, err := uc.storageContainer.FetchOne(ctx, item.AccountID, item.ContainerID)
	if err != nil {
		return 0, uc.errorNotFoundWrapper.Wrap(err)
	}

	if err = uc.serviceStore.CheckAvailability(ctx, item.AccountID, item.LocationID); err != nil {
		return 0, uc.errorWrapper.Wrap(err) // user  usecase error
	}

	stock := entity.Stock{
		AccountID:         item.AccountID,
		ContainerID:       item.ContainerID,
		LocationID:        item.LocationID,
		ContainerQuantity: item.Quantity,
		ContainerVolume:   container.Volume.Calc(),
	}

	// TODO: транзакция

	if stockID, err = uc.storageStock.InsertOrUpdate(ctx, stock); err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	// TODO: через событие сделать
	if err = uc.useCaseRefreshStores.Execute(ctx, []uint64{item.LocationID}); err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "AddMore", "accountId", item.AccountID, "stockId", stockID)

	return stockID, nil
}
