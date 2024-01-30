package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	usecase_api "print-shop-back/internal/modules/controls/usecase/api"
	usecase "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElement struct {
		storage            FormElementStorage
		elementTemplateAPI usecase_api.ElementTemplateAPI
		ordererAPI         mrorderer.API
		eventEmitter       mrsender.EventEmitter
		usecaseHelper      *mrcore.UsecaseHelper
	}
)

func NewFormElement(
	storage FormElementStorage,
	elementTemplateAPI usecase_api.ElementTemplateAPI,
	ordererAPI mrorderer.API,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *FormElement {
	return &FormElement{
		storage:            storage,
		elementTemplateAPI: elementTemplateAPI,
		ordererAPI:         ordererAPI,
		eventEmitter:       eventEmitter,
		usecaseHelper:      usecaseHelper,
	}
}

func (uc *FormElement) GetList(ctx context.Context, params entity.FormElementParams) ([]entity.FormElement, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	if total < 1 {
		return []entity.FormElement{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	return items, total, nil
}

func (uc *FormElement) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.FormElement, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.FormElement{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, id)
	}

	return item, nil
}

func (uc *FormElement) Create(ctx context.Context, item *entity.FormElement) error {
	itemTemplate, err := uc.elementTemplateAPI.GetHead(ctx, item.TemplateID)

	if err != nil {
		return err
	}

	if item.ParamName == "" {
		item.ParamName = itemTemplate.ParamName
	}

	if item.Caption == "" {
		item.Caption = itemTemplate.Caption
	}

	if err = uc.checkItem(ctx, item); err != nil {
		return err
	}

	if err = uc.storage.Insert(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": item.ID})

	meta := uc.storage.GetMetaData(item.FormID)
	ordererAPI := uc.ordererAPI.WithMetaData(meta)

	if err = ordererAPI.MoveToLast(ctx, item.ID); err != nil {
		mrlog.Ctx(ctx).Error().Caller().Err(err)
	}

	return nil
}

func (uc *FormElement) Store(ctx context.Context, item *entity.FormElement) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, item.ID)
	}

	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameFormElement)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *FormElement) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, id)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *FormElement) MoveAfterID(ctx context.Context, id mrtype.KeyInt32, afterID mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := entity.FormElement{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, &item); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameFormElement, id)
	}

	if item.FormID < 1 {
		return mrcore.FactoryErrInternal.WithAttr(entity.ModelNameFormElement, mrmsg.Data{"formId": item.FormID}).New()
	}

	meta := uc.storage.GetMetaData(item.FormID)
	ordererAPI := uc.ordererAPI.WithMetaData(meta)

	if err := ordererAPI.MoveAfterID(ctx, id, afterID); err != nil {
		return err
	}

	uc.emitEvent(ctx, "Move", mrmsg.Data{"id": id, "afterId": afterID})

	return nil
}

func (uc *FormElement) checkItem(ctx context.Context, item *entity.FormElement) error {
	if err := uc.checkParamName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *FormElement) checkParamName(ctx context.Context, item *entity.FormElement) error {
	id, err := uc.storage.FetchIdByName(ctx, item.FormID, item.ParamName)

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

func (uc *FormElement) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameFormElement,
		data,
	)
}
