package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	SubmitForm struct {
		storage       SubmitFormStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewSubmitForm(
	storage SubmitFormStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *SubmitForm {
	return &SubmitForm{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if total < 1 {
		return []entity.SubmitForm{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return items, total, nil
}

func (uc *SubmitForm) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.SubmitForm, error) {
	if itemID < 1 {
		return entity.SubmitForm{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.SubmitForm{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	return item, nil
}

func (uc *SubmitForm) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.IsExists(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return usecase.FactoryErrSubmitFormNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return nil
}

func (uc *SubmitForm) Create(ctx context.Context, item entity.SubmitForm) (mrtype.KeyInt32, error) {
	if err := uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *SubmitForm) Store(ctx context.Context, item entity.SubmitForm) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *SubmitForm) ChangeStatus(ctx context.Context, item entity.SubmitForm) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, item.ID)
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

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *SubmitForm) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *SubmitForm) checkItem(ctx context.Context, item *entity.SubmitForm) error {
	if err := uc.checkParamName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *SubmitForm) checkParamName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIdByName(ctx, item.ParamName)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if item.ID != id {
		return usecase.FactoryErrSubmitFormParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *SubmitForm) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameSubmitForm,
		data,
	)
}
