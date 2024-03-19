package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/shared"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	SubmitForm struct {
		storage        SubmitFormStorage
		storageElement FormElementStorage
		eventEmitter   mrsender.EventEmitter
		usecaseHelper  *mrcore.UsecaseHelper
		statusFlow     mrenum.StatusFlow
	}
)

func NewSubmitForm(
	storage SubmitFormStorage,
	storageElement FormElementStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *SubmitForm {
	return &SubmitForm{
		storage:        storage,
		storageElement: storageElement,
		eventEmitter:   eventEmitter,
		usecaseHelper:  usecaseHelper,
		statusFlow:     mrenum.ItemStatusFlow,
	}
}

func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)
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

func (uc *SubmitForm) GetItem(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error) {
	if itemID == uuid.Nil {
		return entity.SubmitForm{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.SubmitForm{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	if err = uc.setElements(ctx, &item); err != nil {
		return entity.SubmitForm{}, err
	}

	return item, nil
}

func (uc *SubmitForm) Create(ctx context.Context, item entity.SubmitForm) (uuid.UUID, error) {
	if err := uc.checkItem(ctx, &item); err != nil {
		return uuid.Nil, err
	}

	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return uuid.Nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *SubmitForm) Store(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
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

	tagVersion, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *SubmitForm) ChangeStatus(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
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

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

func (uc *SubmitForm) Remove(ctx context.Context, itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *SubmitForm) checkItem(ctx context.Context, item *entity.SubmitForm) error {
	if err := uc.checkRewriteName(ctx, item); err != nil {
		return err
	}

	if err := uc.checkParamName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *SubmitForm) checkRewriteName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIdByRewriteName(ctx, item.ParamName)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if item.ID != id {
		return usecase.FactoryErrSubmitFormRewriteNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *SubmitForm) checkParamName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIdByParamName(ctx, item.ParamName)

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

func (uc *SubmitForm) setElements(ctx context.Context, form *entity.SubmitForm) error {
	elements, err := uc.storageElement.Fetch(ctx, form.ID)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, form.ID)
	}

	form.Elements = elements

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
