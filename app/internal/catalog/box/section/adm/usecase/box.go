package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// Box - comment struct.
	Box struct {
		storage      adm.BoxStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewBox - создаёт объект Box.
func NewBox(storage adm.BoxStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *Box {
	return &Box{
		storage:      storage,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if total < 1 {
		return nil, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *Box) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Box, error) {
	if itemID < 1 {
		return entity.Box{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Box{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Box) Create(ctx context.Context, item entity.Box) (mrtype.KeyInt32, error) {
	if err := uc.checkArticle(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *Box) Store(ctx context.Context, item entity.Box) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
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

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Box) ChangeStatus(ctx context.Context, item entity.Box) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
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

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Box) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

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

func (uc *Box) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameBox,
		data,
	)
}
