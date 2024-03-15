package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplate struct {
		storage       ElementTemplateStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewElementTemplate(
	storage ElementTemplateStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *ElementTemplate {
	return &ElementTemplate{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *ElementTemplate) GetList(ctx context.Context, params entity.ElementTemplateParams) ([]entity.ElementTemplate, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	if total < 1 {
		return []entity.ElementTemplate{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	return items, total, nil
}

func (uc *ElementTemplate) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.ElementTemplate, error) {
	if itemID < 1 {
		return entity.ElementTemplate{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.ElementTemplate{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, itemID)
	}

	return item, nil
}

func (uc *ElementTemplate) Create(ctx context.Context, item entity.ElementTemplate) (mrtype.KeyInt32, error) {
	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *ElementTemplate) Store(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *ElementTemplate) ChangeStatus(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, item.ID)
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

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *ElementTemplate) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *ElementTemplate) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameElementTemplate,
		data,
	)
}
