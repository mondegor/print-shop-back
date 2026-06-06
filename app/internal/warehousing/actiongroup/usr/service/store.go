package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"

	"print-shop-back/internal/warehousing/actiongroup/usr"
	"print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"print-shop-back/internal/warehousing/enum/activitystatus"
	"print-shop-back/internal/warehousing/enum/locationkind"
	"print-shop-back/internal/warehousing/module"
)

type (
	// Store - comment struct.
	Store struct {
		storageStore     usr.StoreStorage
		storageContainer usr.ContainerStorage
		errorWrapper     errors.Wrapper
	}
)

// NewStore - создаёт объект Container.
func NewStore(
	storageStore usr.StoreStorage,
	storageContainer usr.ContainerStorage,
) *Store {
	return &Store{
		storageStore:     storageStore,
		storageContainer: storageContainer,
		errorWrapper:     errors.NewServiceOperationFailedWrapper(),
	}
}

// GetList - comment method.
func (uc *Store) GetList(ctx context.Context, params dto.StoreParams) (items []entity.Store, hasNext bool, err error) {
	if params.AccountID == uuid.Nil {
		return nil, false, errors.ErrInternalIncorrectInputData.WithDetails("params.AccountId is empty")
	}

	items, hasNext, err = uc.storageStore.FetchByCondition(ctx, params)
	if err != nil {
		return nil, false, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Store, 0), false, nil
	}

	return items, hasNext, nil
}

// CheckAvailability - comment method.
func (uc *Store) CheckAvailability(ctx context.Context, accountID uuid.UUID, locationID uint64) error {
	if accountID == uuid.Nil {
		return errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	switch locationkind.ByID(locationID) {
	case locationkind.Group:
		if err := uc.storageContainer.IsExist(ctx, accountID, locationID); err != nil {
			if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
				return module.ErrContainerNotFound.New(locationID) // group container not found
			}

			return uc.errorWrapper.Wrap(err)
		}

		return nil
	case locationkind.Container:
		return errors.ErrIncorrectInputData.New("container cannot be used as location") // user error
	default:
		storeState, err := uc.storageStore.FetchState(ctx, accountID, locationID)
		if err != nil {
			if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
				return module.ErrStoreNotFound.New(locationID)
			}

			return uc.errorWrapper.Wrap(err)
		}

		if storeState.Status != activitystatus.Enabled {
			return errors.ErrIncorrectInputData.New("store is not enabled") // user error
		}

		return nil
	}
}
