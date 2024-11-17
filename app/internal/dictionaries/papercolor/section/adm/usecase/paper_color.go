package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      adm.PaperColorStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(storage adm.PaperColorStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *PaperColor {
	return &PaperColor{
		storage:      storage,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNamePaperColor),
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *PaperColor) GetList(ctx context.Context, params entity.PaperColorParams) (items []entity.PaperColor, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	if countItems == 0 {
		return make([]entity.PaperColor, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PaperColor) GetItem(ctx context.Context, itemID uint64) (entity.PaperColor, error) {
	if itemID == 0 {
		return entity.PaperColor{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PaperColor{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PaperColor) Create(ctx context.Context, item entity.PaperColor) (itemID uint64, err error) {
	item.Status = mrenum.ItemStatusDraft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *PaperColor) Store(ctx context.Context, item entity.PaperColor) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *PaperColor) ChangeStatus(ctx context.Context, item entity.PaperColor) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *PaperColor) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}
