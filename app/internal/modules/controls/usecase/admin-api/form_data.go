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
	FormData struct {
		storage       FormDataStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewFormData(
	storage FormDataStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *FormData {
	return &FormData{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *FormData) GetList(ctx context.Context, params entity.FormDataParams) ([]entity.FormData, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	if total < 1 {
		return []entity.FormData{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	return items, total, nil
}

func (uc *FormData) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.FormData, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.FormData{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormData, id)
	}

	return item, nil
}

func (uc *FormData) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return usecase.FactoryErrFormDataNotFound.New(id)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	return nil
}

func (uc *FormData) Create(ctx context.Context, item *entity.FormData) error {
	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *FormData) Store(ctx context.Context, item *entity.FormData) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormData, item.ID)
	}

	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *FormData) ChangeStatus(ctx context.Context, item *entity.FormData) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormData, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *FormData) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormData, id)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *FormData) checkItem(ctx context.Context, item *entity.FormData) error {
	if err := uc.checkParamName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *FormData) checkParamName(ctx context.Context, item *entity.FormData) error {
	id, err := uc.storage.FetchIdByName(ctx, item.ParamName)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormData)
	}

	if item.ID != id {
		return usecase.FactoryErrFormDataParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *FormData) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameFormData,
		data,
	)
}
