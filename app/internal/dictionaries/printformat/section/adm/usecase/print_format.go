package usecase

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrevent"
	"github.com/mondegor/go-core/mrworkflow/itemstatus"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/dictionaries/printformat/section/adm"
	"print-shop-back/internal/dictionaries/printformat/section/adm/entity"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage                     adm.PrintFormatStorage
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
		statusFlowMap               workflow.FlowMap[workflow.ItemStatus]
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(
	storage adm.PrintFormatStorage,
	eventEmitter mrevent.Emitter,
) *PrintFormat {
	return &PrintFormat{
		storage:                     storage,
		eventEmitter:                mrevent.EmitterWithSource(eventEmitter, entity.ModelNamePrintFormat),
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
		statusFlowMap:               itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *PrintFormat) GetList(ctx context.Context, params entity.PrintFormatParams) (items []entity.PrintFormat, countItems int, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.Wrap(err)
	}

	if countItems == 0 {
		return make([]entity.PrintFormat, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PrintFormat) GetItem(ctx context.Context, itemID uint64) (entity.PrintFormat, error) {
	if itemID == 0 {
		return entity.PrintFormat{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PrintFormat{}, uc.errorNotFoundWrapper.Wrap(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PrintFormat) Create(ctx context.Context, item entity.PrintFormat) (itemID uint64, err error) {
	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", "itemId", itemID)

	return itemID, err
}

// Save - comment method.
func (uc *PrintFormat) Save(ctx context.Context, item entity.PrintFormat) error {
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

	uc.eventEmitter.Emit(ctx, "Store", "itemId", item.ID, "tagVersion", tagVersion)

	return nil
}

// ChangeStatus - comment method.
func (uc *PrintFormat) ChangeStatus(ctx context.Context, item entity.PrintFormat) error {
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

	uc.eventEmitter.Emit(ctx, "ChangeStatus", "itemId", item.ID, "tagVersion", tagVersion, "status", item.Status)

	return nil
}

// Remove - comment method.
func (uc *PrintFormat) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return errors.ErrRecordNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", "itemId", itemID)

	return nil
}
