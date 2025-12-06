package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage       adm.PaperColorStorage
		eventEmitter  mrevent.Emitter
		errorWrapper  mrerr.UseCaseErrorWrapper
		statusFlowMap mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage adm.PaperColorStorage,
	eventEmitter mrevent.Emitter,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *PaperColor {
	return &PaperColor{
		storage:       storage,
		eventEmitter:  mrevent.NewSourceEmitter(eventEmitter, entity.ModelNamePaperColor),
		errorWrapper:  mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNamePaperColor),
		statusFlowMap: itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *PaperColor) GetList(ctx context.Context, params entity.PaperColorParams) (items []entity.PaperColor, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err)
	}

	if countItems == 0 {
		return make([]entity.PaperColor, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PaperColor) GetItem(ctx context.Context, itemID uint64) (entity.PaperColor, error) {
	if itemID == 0 {
		return entity.PaperColor{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PaperColor{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PaperColor) Create(ctx context.Context, item entity.PaperColor) (itemID uint64, err error) {
	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrargs.Group{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *PaperColor) Store(ctx context.Context, item entity.PaperColor) error {
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
func (uc *PaperColor) ChangeStatus(ctx context.Context, item entity.PaperColor) error {
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
func (uc *PaperColor) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrargs.Group{"id": itemID})

	return nil
}
