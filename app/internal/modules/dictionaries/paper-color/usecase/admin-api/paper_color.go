package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColor struct {
		storage       PaperColorStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewPaperColor(
	storage PaperColorStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperColor {
	return &PaperColor{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *PaperColor) GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	if total < 1 {
		return []entity.PaperColor{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	return items, total, nil
}

func (uc *PaperColor) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PaperColor, error) {
	if itemID < 1 {
		return entity.PaperColor{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.PaperColor{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, itemID)
	}

	return item, nil
}

func (uc *PaperColor) Create(ctx context.Context, item entity.PaperColor) (mrtype.KeyInt32, error) {
	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *PaperColor) Store(ctx context.Context, item entity.PaperColor) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *PaperColor) ChangeStatus(ctx context.Context, item entity.PaperColor) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *PaperColor) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperColor, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *PaperColor) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNamePaperColor,
		data,
	)
}
