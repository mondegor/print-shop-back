package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/box/usecase/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Box struct {
		storage       BoxStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewBox(
	storage BoxStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *Box {
	return &Box{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if total < 1 {
		return []entity.Box{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	return items, total, nil
}

func (uc *Box) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Box, error) {
	if itemID < 1 {
		return entity.Box{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.Box{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	return item, nil
}

func (uc *Box) Create(ctx context.Context, item entity.Box) (mrtype.KeyInt32, error) {
	if err := uc.checkArticle(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *Box) Store(ctx context.Context, item entity.Box) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
	}

	if err := uc.checkArticle(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *Box) ChangeStatus(ctx context.Context, item entity.Box) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

func (uc *Box) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Box) checkArticle(ctx context.Context, item *entity.Box) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrBoxArticleAlreadyExists.New(item.Article)
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
