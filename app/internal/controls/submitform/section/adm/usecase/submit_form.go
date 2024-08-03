package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct {
		storage        adm.SubmitFormStorage
		storageElement adm.FormElementStorage
		storageVersion adm.FormVersionStorage
		eventEmitter   mrsender.EventEmitter
		errorWrapper   mrcore.UsecaseErrorWrapper
		statusFlow     mrstatus.Flow
	}
)

// NewSubmitForm - создаёт объект SubmitForm.
func NewSubmitForm(
	storage adm.SubmitFormStorage,
	storageElement adm.FormElementStorage,
	storageVersion adm.FormVersionStorage,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
) *SubmitForm {
	return &SubmitForm{
		storage:        storage,
		storageElement: storageElement,
		storageVersion: storageVersion,
		eventEmitter:   eventEmitter,
		errorWrapper:   errorWrapper,
		statusFlow:     mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if total < 1 {
		return make([]entity.SubmitForm, 0), 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *SubmitForm) GetItem(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error) {
	item, err := uc.getForm(ctx, itemID)
	if err != nil {
		return entity.SubmitForm{}, err
	}

	if err = uc.setVersions(ctx, &item); err != nil {
		return entity.SubmitForm{}, err
	}

	return item, nil
}

// Create - comment method.
func (uc *SubmitForm) Create(ctx context.Context, item entity.SubmitForm) (uuid.UUID, error) {
	if err := uc.checkItem(ctx, &item); err != nil {
		return uuid.Nil, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *SubmitForm) Store(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *SubmitForm) ChangeStatus(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *SubmitForm) Remove(ctx context.Context, itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

// GetFormStatus - comment method.
func (uc *SubmitForm) GetFormStatus(ctx context.Context, formID uuid.UUID) (mrenum.ItemStatus, error) {
	if formID == uuid.Nil {
		return 0, module.ErrSubmitFormRequired.New()
	}

	status, err := uc.storage.FetchStatus(ctx, formID)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return 0, module.ErrSubmitFormNotFound.New(formID)
		}

		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	return status, nil
}

// GetFormWithElements - comment method.
func (uc *SubmitForm) GetFormWithElements(ctx context.Context, formID uuid.UUID) (entity.SubmitForm, error) {
	return uc.getForm(ctx, formID)
}

func (uc *SubmitForm) getForm(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error) {
	if itemID == uuid.Nil {
		return entity.SubmitForm{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.SubmitForm{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameSubmitForm, itemID)
	}

	if err = uc.setElements(ctx, &item); err != nil {
		return entity.SubmitForm{}, err
	}

	return item, nil
}

func (uc *SubmitForm) checkItem(ctx context.Context, item *entity.SubmitForm) error {
	if err := uc.checkRewriteName(ctx, item); err != nil {
		return err
	}

	return uc.checkParamName(ctx, item)
}

func (uc *SubmitForm) checkRewriteName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIDByRewriteName(ctx, item.RewriteName)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if item.ID != id {
		return module.ErrSubmitFormRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}

func (uc *SubmitForm) checkParamName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIDByParamName(ctx, item.ParamName)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameSubmitForm)
	}

	if item.ID != id {
		return module.ErrSubmitFormParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *SubmitForm) setElements(ctx context.Context, form *entity.SubmitForm) error {
	elements, err := uc.storageElement.Fetch(ctx, form.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, form.ID)
	}

	form.Elements = elements

	return nil
}

func (uc *SubmitForm) setVersions(ctx context.Context, form *entity.SubmitForm) error {
	versions, err := uc.storageVersion.Fetch(ctx, form.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormVersion, form.ID)
	}

	form.Versions = versions

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
