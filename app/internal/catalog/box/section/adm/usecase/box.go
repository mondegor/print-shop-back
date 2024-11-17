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

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
)

type (
	// Box - comment struct.
	Box struct {
		storage      adm.BoxStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewBox - создаёт объект Box.
func NewBox(storage adm.BoxStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *Box {
	return &Box{
		storage:      storage,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameBox),
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) (items []entity.Box, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if countItems == 0 {
		return make([]entity.Box, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *Box) GetItem(ctx context.Context, itemID uint64) (entity.Box, error) {
	if itemID == 0 {
		return entity.Box{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Box{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Box) Create(ctx context.Context, item entity.Box) (itemID uint64, err error) {
	if err = uc.checkArticle(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *Box) Store(ctx context.Context, item entity.Box) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
	}

	if err := uc.checkArticle(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Box) ChangeStatus(ctx context.Context, item entity.Box) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Box) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Box) checkArticle(ctx context.Context, item *entity.Box) error {
	id, err := uc.storage.FetchIDByArticle(ctx, item.Article)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if item.ID != id {
		return module.ErrBoxArticleAlreadyExists.New(item.Article)
	}

	return nil
}
