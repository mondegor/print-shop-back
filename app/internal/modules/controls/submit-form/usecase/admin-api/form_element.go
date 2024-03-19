package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/shared"
	"print-shop-back/pkg/modules/controls"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElement struct {
		storage            FormElementStorage
		submitFormAPI      SubmitFormAPI
		elementTemplateAPI controls.ElementTemplateAPI
		ordererAPI         mrorderer.API
		eventEmitter       mrsender.EventEmitter
		usecaseHelper      *mrcore.UsecaseHelper
	}

	SubmitFormAPI interface {
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.SubmitForm, error)
	}
)

func NewFormElement(
	storage FormElementStorage,
	submitFormAPI SubmitFormAPI,
	elementTemplateAPI controls.ElementTemplateAPI,
	ordererAPI mrorderer.API,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *FormElement {
	return &FormElement{
		storage:            storage,
		submitFormAPI:      submitFormAPI,
		elementTemplateAPI: elementTemplateAPI,
		ordererAPI:         ordererAPI,
		eventEmitter:       eventEmitter,
		usecaseHelper:      usecaseHelper,
	}
}

func (uc *FormElement) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.FormElement, error) {
	if itemID < 1 {
		return entity.FormElement{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.FormElement{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	return item, nil
}

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
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	meta := uc.storage.NewOrderMeta(item.FormID)

	if err = uc.ordererAPI.WithMetaData(meta).MoveToLast(ctx, itemID); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}

	return itemID, nil
}

func (uc *FormElement) Store(ctx context.Context, item entity.FormElement) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *FormElement) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	ordererAPI, err := uc.getOrdererAPI(ctx, itemID)

	if err != nil {
		return err
	}

	if err = ordererAPI.Unlink(ctx, itemID); err != nil {
		return err
	}

	if err = uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *FormElement) MoveAfterID(ctx context.Context, itemID mrtype.KeyInt32, afterID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	ordererAPI, err := uc.getOrdererAPI(ctx, itemID)

	if err != nil {
		return err
	}

	if err = ordererAPI.MoveAfterID(ctx, itemID, afterID); err != nil {
		return err
	}

	uc.emitEvent(ctx, "Move", mrmsg.Data{"id": itemID, "afterId": afterID})

	return nil
}

func (uc *FormElement) initItemBeforeCreate(ctx context.Context, item *entity.FormElement) error {
	itemTemplate, err := uc.elementTemplateAPI.GetItemHead(ctx, item.TemplateID)

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
		return usecase.FactoryErrSubmitFormRequired.New()
	}

	form, err := uc.submitFormAPI.FetchOne(ctx, item.FormID)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return usecase.FactoryErrSubmitFormNotFound.New(item.FormID)
		}

		return uc.usecaseHelper.WrapErrorEntityFailed(err, entity.ModelNameSubmitForm, item.FormID)
	}

	if form.Detailing != enums.ElementDetailingExtended && form.Detailing != item.Detailing {
		return usecase.FactoryErrFormElementDetailingNotAllowed.New(item.Detailing, form.Detailing)
	}

	if form.Status == mrenum.ItemStatusDisabled {
		return usecase.FactoryErrSubmitFormIsDisabled.New(item.FormID)
	}

	return nil
}

func (uc *FormElement) checkItem(ctx context.Context, item *entity.FormElement) error {
	if err := uc.checkParamName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *FormElement) checkParamName(ctx context.Context, item *entity.FormElement) error {
	id, err := uc.storage.FetchIdByParamName(ctx, item.FormID, item.ParamName)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	if item.ID != id {
		return usecase.FactoryErrFormElementParamNameAlreadyExists.New(item.ParamName)
	}

	return nil
}

func (uc *FormElement) getOrdererAPI(ctx context.Context, itemID mrtype.KeyInt32) (mrorderer.API, error) {
	// :TODO: можно оптимизировать загружая только FormID
	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, itemID)
	}

	if item.FormID == uuid.Nil {
		return nil, mrcore.FactoryErrInternal.WithAttr(entity.ModelNameFormElement, mrmsg.Data{"formId": item.FormID}).New()
	}

	meta := uc.storage.NewOrderMeta(item.FormID)

	return uc.ordererAPI.WithMetaData(meta), nil
}

func (uc *FormElement) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameFormElement,
		data,
	)
}
