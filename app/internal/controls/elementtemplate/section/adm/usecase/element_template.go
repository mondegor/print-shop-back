package usecase

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplate - comment struct.
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage      adm.ElementTemplateStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage adm.ElementTemplateStorage,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UseCaseErrorWrapper,
) *ElementTemplate {
	return &ElementTemplate{
		storage:      storage,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameElementTemplate),
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *ElementTemplate) GetList(ctx context.Context, params entity.ElementTemplateParams) (items []entity.ElementTemplate, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	if countItems == 0 {
		return make([]entity.ElementTemplate, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *ElementTemplate) GetItem(ctx context.Context, itemID uint64) (entity.ElementTemplate, error) {
	if itemID == 0 {
		return entity.ElementTemplate{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.ElementTemplate{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, itemID)
	}

	// TODO: можно оптимизировать удалив body
	item.Body = nil

	return item, nil
}

// GetItemJson - comment method.
func (uc *ElementTemplate) GetItemJson(ctx context.Context, itemID uint64, pretty bool) ([]byte, error) {
	if itemID == 0 {
		return nil, mrcore.ErrUseCaseEntityNotFound.New()
	}

	// TODO: можно оптимизировать получая только body
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, itemID)
	}

	if pretty {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, item.Body, "", module.JsonPrettyIndent); err != nil {
			return nil, uc.errorWrapper.WrapErrorEntityFailed(err, entity.ModelNameElementTemplate, itemID)
		}

		return prettyJSON.Bytes(), nil
	}

	return item.Body, nil
}

// Create - comment method.
func (uc *ElementTemplate) Create(ctx context.Context, item entity.ElementTemplate) (itemID uint64, err error) {
	item.Status = mrenum.ItemStatusDraft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *ElementTemplate) Store(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.eventEmitter.Emit(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *ElementTemplate) ChangeStatus(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameElementTemplate)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *ElementTemplate) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameElementTemplate, itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}
