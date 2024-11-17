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

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      adm.PaperFactureStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(storage adm.PaperFactureStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNamePaperFacture),
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *PaperFacture) GetList(ctx context.Context, params entity.PaperFactureParams) (items []entity.PaperFacture, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	if countItems == 0 {
		return make([]entity.PaperFacture, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *PaperFacture) GetItem(ctx context.Context, itemID uint64) (entity.PaperFacture, error) {
	if itemID == 0 {
		return entity.PaperFacture{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.PaperFacture{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *PaperFacture) Create(ctx context.Context, item entity.PaperFacture) (itemID uint64, err error) {
	item.Status = mrenum.ItemStatusDraft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *PaperFacture) Store(ctx context.Context, item entity.PaperFacture) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *PaperFacture) ChangeStatus(ctx context.Context, item entity.PaperFacture) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *PaperFacture) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}
