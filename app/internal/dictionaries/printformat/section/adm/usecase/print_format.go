package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage       adm.PrintFormatStorage
		eventEmitter  mrevent.Emitter
		errorWrapper  mrerr.UseCaseErrorWrapper
		statusFlowMap mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(
	storage adm.PrintFormatStorage,
	eventEmitter mrevent.Emitter,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *PrintFormat {
	return &PrintFormat{
		storage:       storage,
		eventEmitter:  mrevent.NewSourceEmitter(eventEmitter, entity.ModelNamePrintFormat),
		errorWrapper:  mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNamePrintFormat),
		statusFlowMap: itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *PrintFormat) GetList(ctx context.Context, params entity.PrintFormatParams) (items []entity.PrintFormat, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err)
	}

	if countItems == 0 {
		return make([]entity.PrintFormat, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PrintFormat) GetItem(ctx context.Context, itemID uint64) (entity.PrintFormat, error) {
	if itemID == 0 {
		return entity.PrintFormat{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PrintFormat{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PrintFormat) Create(ctx context.Context, item entity.PrintFormat) (itemID uint64, err error) {
	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrargs.Group{"id": itemID})

	return itemID, err
}

// Store - comment method.
func (uc *PrintFormat) Store(ctx context.Context, item entity.PrintFormat) error {
	if item.ID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mr.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mr.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrargs.Group{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *PrintFormat) ChangeStatus(ctx context.Context, item entity.PrintFormat) error {
	if item.ID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mr.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return mr.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mr.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrargs.Group{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *PrintFormat) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrargs.Group{"id": itemID})

	return nil
}
