package usecase

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"github.com/mondegor/print-shop-back/internal/warehousing/enum/locationkind"
)

const (
	maxInsertAttempts = 3
)

type (
	// CreateStockContainer - comment struct.
	CreateStockContainer struct {
		txManager               mrstorage.DBTxManager
		storageContainer        usr.ContainerStorage
		storageStock            usr.StockStorage
		serviceStore            usr.StoreService
		useCaseRefreshLocations refreshLocationsUseCase
		serviceSequence         accountSequenceService
		eventEmitter            mrevent.Emitter
		errorWrapper            errors.Wrapper
	}

	accountSequenceService interface {
		ContainerCode(ctx context.Context, accountID uuid.UUID, kind locationkind.Enum) (string, error)
	}

	refreshLocationsUseCase interface {
		Execute(ctx context.Context, locationIDs []uint64) error
	}
)

var errMarkerAlreadyExists = errors.New("marker already exists for code")

// NewCreateStockContainer - создаёт объект CreateStockContainer.
func NewCreateStockContainer(
	txManager mrstorage.DBTxManager,
	storageContainer usr.ContainerStorage,
	serviceSequence accountSequenceService,
	storageStock usr.StockStorage,
	serviceStore usr.StoreService,
	useCaseRefreshLocations refreshLocationsUseCase,
	eventEmitter mrevent.Emitter,
) *CreateStockContainer {
	return &CreateStockContainer{
		txManager:               txManager,
		storageContainer:        storageContainer,
		serviceSequence:         serviceSequence,
		storageStock:            storageStock,
		serviceStore:            serviceStore,
		useCaseRefreshLocations: useCaseRefreshLocations,
		eventEmitter:            eventEmitter,
		errorWrapper:            errors.NewServiceOperationFailedWrapper(),
	}
}

// Execute - comment method.
func (uc *CreateStockContainer) Execute(ctx context.Context, item dto.CreateStockContainer) (id dto.CreateStockContainerResult, err error) {
	if item.AccountID == uuid.Nil {
		return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	if item.LocationID == 0 {
		return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("locationId is zero")
	}

	if item.Quantity < 1 {
		item.Quantity = 1
	}

	switch item.Kind {
	case locationkind.Container:
	case locationkind.Group:
		if item.Quantity != 1 {
			return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("quantity for group must be only one")
		}

		// ограничение на вложение групповых контейнеров
		if !locationkind.Is(item.LocationID, locationkind.Store) {
			return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("G/ container add error") // user error
		}
	default:
		return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("unexpected container kind")
	}

	if item.Code == "" {
		if item.Code, err = uc.serviceSequence.ContainerCode(ctx, item.AccountID, item.Kind); err != nil {
			return dto.CreateStockContainerResult{}, uc.errorWrapper.Wrap(err) // system error
		}
	}

	c := entity.Container{
		Kind:      item.Kind,
		AccountID: item.AccountID,
		Code:      item.Code,
		Volume:    item.Volume,
		Tags:      item.Tags,
		Images:    item.Images,
	}

	s := entity.Stock{
		AccountID:         item.AccountID,
		LocationID:        item.LocationID,
		ContainerQuantity: item.Quantity,
		ContainerVolume:   item.Volume.Calc(),
	}

	attempts := maxInsertAttempts

	for {
		if err = uc.serviceStore.CheckAvailability(ctx, item.AccountID, item.LocationID); err != nil {
			return dto.CreateStockContainerResult{}, uc.errorWrapper.Wrap(err) // user error
		}

		c.Marker, err = uc.storageContainer.FetchMaxMarker(ctx, item.AccountID, item.Code)
		if err != nil && !errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return dto.CreateStockContainerResult{}, uc.errorWrapper.Wrap(err)
		}

		if c.Marker >= math.MaxInt16 {
			return dto.CreateStockContainerResult{}, errors.ErrInternalIncorrectInputData.WithDetails("marker has reached its maximum value") // user error
		}

		c.Marker++ // для нового контейнера всегда присваиваем новый маркер

		err = uc.txManager.Do(ctx, func(ctx context.Context) error {
			s.ContainerID, err = uc.storageContainer.Insert(ctx, c)
			if err != nil {
				if attempts < 1 || !errors.Is(err, errors.ErrInternalStorageDuplicateKeyViolation) {
					return uc.errorWrapper.Wrap(err)
				}

				attempts--

				return errMarkerAlreadyExists
			}

			s.ID, err = uc.storageStock.InsertOrUpdate(ctx, s)
			if err != nil {
				return uc.errorWrapper.Wrap(err)
			}

			uc.eventEmitter.Emit(
				ctx,
				"Create",
				conv.Group{
					"accountId":   item.AccountID,
					"itemCode":    item.Code,
					"containerId": s.ContainerID,
					"stockID":     s.ID,
					"locationID":  item.LocationID,
				},
			)

			// TODO: через событие сделать
			if err = uc.useCaseRefreshLocations.Execute(ctx, []uint64{item.LocationID}); err != nil {
				return uc.errorWrapper.Wrap(err)
			}

			return nil
		})
		if err != nil {
			if errors.Is(err, errMarkerAlreadyExists) {
				// мог кто-то опередить вставку с таким же маркером в результате гонки,
				// поэтому происходит ещё одна попытка добавления контейнера с новым маркером
				continue
			}

			return dto.CreateStockContainerResult{}, err
		}

		return dto.CreateStockContainerResult{
			ID:      s.ContainerID,
			Code:    c.Code,
			Marker:  c.Marker,
			StockID: s.ID,
		}, nil
	}
}
