package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
)

type (
	// SubmitForm - comment struct.
	SubmitForm struct {
		storage        adm.SubmitFormStorage
		storageElement adm.FormElementStorage
		storageVersion adm.FormVersionStorage
		eventEmitter   mrevent.Emitter
		errorWrapper   errors.Wrapper
		statusFlowMap  mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewSubmitForm - создаёт объект SubmitForm.
func NewSubmitForm(
	storage adm.SubmitFormStorage,
	storageElement adm.FormElementStorage,
	storageVersion adm.FormVersionStorage,
	eventEmitter mrevent.Emitter,
) *SubmitForm {
	return &SubmitForm{
		storage:        storage,
		storageElement: storageElement,
		storageVersion: storageVersion,
		eventEmitter:   mrevent.EmitterWithSource(eventEmitter, entity.ModelNameSubmitForm),
		errorWrapper:   errors.NewUseCaseWrapper(),
		statusFlowMap:  itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *SubmitForm) GetList(ctx context.Context, params entity.SubmitFormParams) (items []entity.SubmitForm, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.Wrap(err)
	}

	if countItems == 0 {
		return make([]entity.SubmitForm, 0), 0, nil
	}

	return items, countItems, nil
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
func (uc *SubmitForm) Create(ctx context.Context, item entity.SubmitForm) (itemID uuid.UUID, err error) {
	if err = uc.checkItem(ctx, &item); err != nil {
		return uuid.Nil, err
	}

	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", conv.Group{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *SubmitForm) Store(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
		return errors.ErrUseCaseEntityNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrUseCaseEntityVersionConflict
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return errors.ErrUseCaseEntityVersionConflict.Wrap(err)
		}

		return uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", conv.Group{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *SubmitForm) ChangeStatus(ctx context.Context, item entity.SubmitForm) error {
	if item.ID == uuid.Nil {
		return errors.ErrUseCaseEntityNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrUseCaseEntityVersionConflict
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return errors.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return errors.ErrUseCaseEntityVersionConflict.Wrap(err)
		}

		return uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", conv.Group{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *SubmitForm) Remove(ctx context.Context, itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return errors.ErrUseCaseEntityNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", conv.Group{"id": itemID})

	return nil
}

// GetFormStatus - comment method.
func (uc *SubmitForm) GetFormStatus(ctx context.Context, formID uuid.UUID) (itemstatus.Enum, error) {
	if formID == uuid.Nil {
		return 0, module.ErrSubmitFormRequired
	}

	status, err := uc.storage.FetchStatus(ctx, formID)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return 0, module.ErrSubmitFormNotFound.Wrap(err, formID)
		}

		return 0, uc.errorWrapper.Wrap(err)
	}

	return status, nil
}

// GetFormWithElements - comment method.
func (uc *SubmitForm) GetFormWithElements(ctx context.Context, formID uuid.UUID) (entity.SubmitForm, error) {
	return uc.getForm(ctx, formID)
}

func (uc *SubmitForm) getForm(ctx context.Context, itemID uuid.UUID) (entity.SubmitForm, error) {
	if itemID == uuid.Nil {
		return entity.SubmitForm{}, errors.ErrUseCaseEntityNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.SubmitForm{}, uc.errorWrapper.Wrap(err, "itemId", itemID)
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
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return nil
		}

		return uc.errorWrapper.Wrap(err)
	}

	if item.ID != id {
		return module.ErrSubmitFormRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}

func (uc *SubmitForm) checkParamName(ctx context.Context, item *entity.SubmitForm) error {
	id, err := uc.storage.FetchIDByParamName(ctx, item.ParamName)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRowFound) {
			return nil
		}

		return uc.errorWrapper.Wrap(err)
	}

	if item.ID != id {
		return module.ErrSubmitFormParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *SubmitForm) setElements(ctx context.Context, form *entity.SubmitForm) error {
	elements, err := uc.storageElement.Fetch(ctx, form.ID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "formId", form.ID)
	}

	form.Elements = elements

	return nil
}

func (uc *SubmitForm) setVersions(ctx context.Context, form *entity.SubmitForm) error {
	versions, err := uc.storageVersion.Fetch(ctx, form.ID)
	if err != nil {
		return uc.errorWrapper.Wrap(err, "formId", form.ID)
	}

	form.Versions = versions

	return nil
}
