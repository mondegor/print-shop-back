package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"

	"print-shop-back/internal/warehousing/actiongroup/usr"
	"print-shop-back/internal/warehousing/actiongroup/usr/dto"
	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"print-shop-back/internal/warehousing/enum/locationkind"
	"print-shop-back/internal/warehousing/util"
)

type (
	// Container - comment struct.
	Container struct {
		storageContainer            usr.ContainerStorage
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
	}
)

// NewContainer - создаёт объект Container.
func NewContainer(
	storageContainer usr.ContainerStorage,
	eventEmitter mrevent.Emitter,
) *Container {
	return &Container{
		storageContainer:            storageContainer,
		eventEmitter:                eventEmitter,
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
	}
}

// GetList - comment method.
func (uc *Container) GetList(ctx context.Context, params dto.ContainerParams) (items []entity.Container, hasNext bool, err error) {
	if params.AccountID == uuid.Nil {
		return nil, false, errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	items, hasNext, err = uc.storageContainer.FetchByCondition(ctx, params)
	if err != nil {
		return nil, false, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Container, 0), false, nil
	}

	return items, hasNext, nil
}

// SaveTags - comment method.
func (uc *Container) SaveTags(ctx context.Context, item entity.UpdateContainerTags) (tagVersion uint32, err error) {
	if item.AccountID == uuid.Nil {
		return 0, errors.ErrInternalIncorrectInputData.WithDetails("accountId is empty")
	}

	if item.TagVersion == 0 {
		return 0, errors.ErrRecordVersionConflict
	}

	if !locationkind.Is(item.ID, locationkind.Container) {
		return 0, errors.ErrIncorrectInputData.New("tags can only be updated for container") // user error
	}

	item.Tags = util.PrepareTags(item.Tags)

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionConflict
	if err = uc.storageContainer.IsExist(ctx, item.AccountID, item.ID); err != nil {
		return 0, uc.errorNotFoundWrapper.Wrap(err)
	}

	tagVersion, err = uc.storageContainer.UpdateTags(ctx, item)
	if err != nil {
		return 0, uc.errorVersionConflictWrapper.Wrap(err)
	}

	return tagVersion, nil
}
