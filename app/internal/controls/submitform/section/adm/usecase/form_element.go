package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrordering"
	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/controls/enum/elementdetailing"
)

type (
	// FormElement - comment struct.
	FormElement struct {
		storage            adm.FormElementStorage
		submitFormAPI      SubmitFormAPI
		elementTemplateAPI api.ElementTemplateHeader
		orderingAPI        mrordering.Mover
		eventEmitter       mrevent.Emitter
		errorWrapper       mrerr.UseCaseErrorWrapper
		logger             mrlog.Logger
	}

	// SubmitFormAPI - comment interface.
	SubmitFormAPI interface {
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error)
	}
)

// NewFormElement - создаёт объект FormElement.
func NewFormElement(
	storage adm.FormElementStorage,
	submitFormAPI SubmitFormAPI,
	elementTemplateAPI api.ElementTemplateHeader,
	orderingAPI mrordering.Mover,
	eventEmitter mrevent.Emitter,
	errorWrapper mrerr.UseCaseErrorWrapper,
	logger mrlog.Logger,
) *FormElement {
	return &FormElement{
		storage:            storage,
		submitFormAPI:      submitFormAPI,
		elementTemplateAPI: elementTemplateAPI,
		orderingAPI:        orderingAPI,
		eventEmitter:       mrevent.NewSourceEmitter(eventEmitter, entity.ModelNameFormElement),
		errorWrapper:       mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameFormElement),
		logger:             logger,
	}
}

// GetItem - comment method.
func (uc *FormElement) GetItem(ctx context.Context, itemID uint64) (entity.FormElement, error) {
	if itemID == 0 {
		return entity.FormElement{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.FormElement{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *FormElement) Create(ctx context.Context, item entity.FormElement) (itemID uint64, err error) {
	if err = uc.initItemBeforeCreate(ctx, &item); err != nil {
		return 0, err
	}

	if err = uc.checkForm(ctx, &item); err != nil {
		return 0, err
	}

	if err = uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrargs.Group{"id": itemID})

	if err = uc.orderingAPI.MoveToLast(ctx, itemID, uc.storage.NewCondition(item.FormID)); err != nil {
		uc.logger.Error(ctx, "Create", "error", err)
	}

	return itemID, nil
}

// Store - comment method.
func (uc *FormElement) Store(ctx context.Context, item entity.FormElement) error {
	if item.ID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mr.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if err := uc.storage.IsExist(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return mr.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrargs.Group{"id": item.ID, "ver": tagVersion})

	return nil
}

// Remove - comment method.
func (uc *FormElement) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	formID, err := uc.getFormID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.orderingAPI.Unlink(ctx, itemID, uc.storage.NewCondition(formID)); err != nil {
		return err
	}

	if err = uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrargs.Group{"id": itemID})

	return nil
}

// MoveAfterID - comment method.
func (uc *FormElement) MoveAfterID(ctx context.Context, itemID, afterID uint64) error {
	if itemID == 0 {
		return mr.ErrUseCaseEntityNotFound.New()
	}

	formID, err := uc.getFormID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.orderingAPI.MoveAfterID(ctx, itemID, afterID, uc.storage.NewCondition(formID)); err != nil {
		return err
	}

	uc.eventEmitter.Emit(ctx, "Move", mrargs.Group{"id": itemID, "afterId": afterID})

	return nil
}

func (uc *FormElement) initItemBeforeCreate(ctx context.Context, item *entity.FormElement) error {
	itemTemplate, err := uc.elementTemplateAPI.GetItemHeader(ctx, item.TemplateID)
	if err != nil {
		return err
	}

	if item.ParamName == "" {
		item.ParamName = itemTemplate.ParamName
	}

	if item.Caption == "" {
		item.Caption = itemTemplate.Caption
	}

	item.TemplateVersion = itemTemplate.TagVersion
	item.Detailing = itemTemplate.Detailing

	return nil
}

func (uc *FormElement) checkForm(ctx context.Context, item *entity.FormElement) error {
	if item.FormID == uuid.Nil {
		return module.ErrSubmitFormRequired.New()
	}

	form, err := uc.submitFormAPI.FetchOne(ctx, item.FormID)
	if err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return module.ErrSubmitFormNotFound.New(item.FormID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, "formId", item.FormID)
	}

	if form.Detailing != elementdetailing.Extended && form.Detailing != item.Detailing {
		return module.ErrFormElementDetailingNotAllowed.New(item.Detailing, form.Detailing)
	}

	if form.Status == itemstatus.Disabled {
		return module.ErrSubmitFormIsDisabled.New(item.FormID)
	}

	return nil
}

func (uc *FormElement) checkItem(ctx context.Context, item *entity.FormElement) error {
	return uc.checkParamName(ctx, item)
}

func (uc *FormElement) checkParamName(ctx context.Context, item *entity.FormElement) error {
	id, err := uc.storage.FetchIDByParamName(ctx, item.FormID, item.ParamName)
	if err != nil {
		if uc.errorWrapper.IsNotFoundOrNotAffectedError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err)
	}

	if item.ID != id {
		return module.ErrFormElementParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *FormElement) getFormID(ctx context.Context, itemID uint64) (uuid.UUID, error) {
	// TODO: можно оптимизировать загружая только FormID
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "itemId", itemID)
	}

	if item.FormID == uuid.Nil {
		return uuid.Nil, mr.ErrInternal.New().WithAttrs("model", entity.ModelNameFormElement, "formId", item.FormID)
	}

	return item.FormID, nil
}
