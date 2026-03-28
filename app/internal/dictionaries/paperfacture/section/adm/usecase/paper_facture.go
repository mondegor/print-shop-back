package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage                     adm.PaperFactureStorage
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
		statusFlowMap               mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	storage adm.PaperFactureStorage,
	eventEmitter mrevent.Emitter,
) *PaperFacture {
	return &PaperFacture{
		storage:                     storage,
		eventEmitter:                mrevent.EmitterWithSource(eventEmitter, entity.ModelNamePaperFacture),
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
		statusFlowMap:               itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *PaperFacture) GetList(ctx context.Context, params entity.PaperFactureParams) (items []entity.PaperFacture, countItems int, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.Wrap(err)
	}

	if countItems == 0 {
		return make([]entity.PaperFacture, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PaperFacture) GetItem(ctx context.Context, itemID uint64) (entity.PaperFacture, error) {
	if itemID == 0 {
		return entity.PaperFacture{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PaperFacture{}, uc.errorNotFoundWrapper.Wrap(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PaperFacture) Create(ctx context.Context, item entity.PaperFacture) (itemID uint64, err error) {
	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", conv.Group{"id": itemID})

	return itemID, nil
}

// Save - comment method.
func (uc *PaperFacture) Save(ctx context.Context, item entity.PaperFacture) error {
	if item.ID == 0 {
		return errors.ErrRecordNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrRecordVersionConflict
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionConflict
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorNotFoundWrapper.Wrap(err, "itemId", item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", conv.Group{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *PaperFacture) ChangeStatus(ctx context.Context, item entity.PaperFacture) error {
	if item.ID == 0 {
		return errors.ErrRecordNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrRecordVersionConflict
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorNotFoundWrapper.Wrap(err, "itemId", item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return errors.ErrSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", conv.Group{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *PaperFacture) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return errors.ErrRecordNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", conv.Group{"id": itemID})

	return nil
}
