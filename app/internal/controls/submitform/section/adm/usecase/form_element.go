package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FormElement - comment struct.
	FormElement struct {
		storage            FormElementStorage
		submitFormAPI      SubmitFormAPI
		elementTemplateAPI api.ElementTemplateHeader
		ordererAPI         mrsort.Orderer
		eventEmitter       mrsender.EventEmitter
		errorWrapper       mrcore.UsecaseErrorWrapper
	}

	// SubmitFormAPI - comment interface.
	SubmitFormAPI interface {
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error)
	}
)

// NewFormElement - создаёт объект FormElement.
func NewFormElement(
	storage FormElementStorage,
	submitFormAPI SubmitFormAPI,
	elementTemplateAPI api.ElementTemplateHeader,
	ordererAPI mrsort.Orderer,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
) *FormElement {
	return &FormElement{
		storage:            storage,
		submitFormAPI:      submitFormAPI,
		elementTemplateAPI: elementTemplateAPI,
		ordererAPI:         ordererAPI,
		eventEmitter:       eventEmitter,
		errorWrapper:       errorWrapper,
	}
}

// GetItem - comment method.
func (uc *FormElement) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.FormElement, error) {
	if itemID < 1 {
		return entity.FormElement{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.FormElement{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *FormElement) Create(ctx context.Context, item entity.FormElement) (mrtype.KeyInt32, error) {
	if err := uc.initItemBeforeCreate(ctx, &item); err != nil {
		return 0, err
	}

	if err := uc.checkForm(ctx, &item); err != nil {
		return 0, err
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	if err = uc.getOrdererAPI(item.FormID).MoveToLast(ctx, itemID); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}

	return itemID, nil
}

// Store - comment method.
func (uc *FormElement) Store(ctx context.Context, item entity.FormElement) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if err := uc.storage.IsExist(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// Remove - comment method.
func (uc *FormElement) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	formID, err := uc.getFormID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.getOrdererAPI(formID).Unlink(ctx, itemID); err != nil {
		return err
	}

	if err = uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

// MoveAfterID - comment method.
func (uc *FormElement) MoveAfterID(ctx context.Context, itemID, afterID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	formID, err := uc.getFormID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.getOrdererAPI(formID).MoveAfterID(ctx, itemID, afterID); err != nil {
		return err
	}

	uc.emitEvent(ctx, "Move", mrmsg.Data{"id": itemID, "afterId": afterID})

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
		if uc.errorWrapper.IsNotFoundError(err) {
			return module.ErrSubmitFormNotFound.New(item.FormID)
		}

		return uc.errorWrapper.WrapErrorEntityFailed(err, entity.ModelNameSubmitForm, item.FormID)
	}

	if form.Detailing != enum.ElementDetailingExtended && form.Detailing != item.Detailing {
		return module.ErrFormElementDetailingNotAllowed.New(item.Detailing, form.Detailing)
	}

	if form.Status == mrenum.ItemStatusDisabled {
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
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	if item.ID != id {
		return module.ErrFormElementParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *FormElement) getOrdererAPI(formID uuid.UUID) mrsort.Orderer {
	meta := uc.storage.NewOrderMeta(formID)

	return uc.ordererAPI.WithMetaData(meta)
}

func (uc *FormElement) getFormID(ctx context.Context, itemID mrtype.KeyInt32) (uuid.UUID, error) {
	// TODO: можно оптимизировать загружая только FormID
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	if item.FormID == uuid.Nil {
		return uuid.Nil, mrcore.ErrInternal.New().WithAttr(entity.ModelNameFormElement, mrmsg.Data{"formId": item.FormID})
	}

	return item.FormID, nil
}

func (uc *FormElement) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameFormElement,
		data,
	)
}
