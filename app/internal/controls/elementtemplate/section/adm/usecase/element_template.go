package usecase

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplate - comment struct.
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage       adm.ElementTemplateStorage
		eventEmitter  mrevent.Emitter
		errorWrapper  errors.Wrapper
		statusFlowMap mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage adm.ElementTemplateStorage,
	eventEmitter mrevent.Emitter,
) *ElementTemplate {
	return &ElementTemplate{
		storage:       storage,
		eventEmitter:  mrevent.EmitterWithSource(eventEmitter, entity.ModelNameElementTemplate),
		errorWrapper:  errors.NewUseCaseWrapper(),
		statusFlowMap: itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *ElementTemplate) GetList(ctx context.Context, params entity.ElementTemplateParams) (items []entity.ElementTemplate, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.Wrap(err)
	}

	if countItems == 0 {
		return make([]entity.ElementTemplate, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *ElementTemplate) GetItem(ctx context.Context, itemID uint64) (entity.ElementTemplate, error) {
	if itemID == 0 {
		return entity.ElementTemplate{}, errors.ErrUseCaseEntityNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.ElementTemplate{}, uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	// TODO: можно оптимизировать удалив body
	item.Body = nil

	return item, nil
}

// GetItemJson - comment method.
func (uc *ElementTemplate) GetItemJson(ctx context.Context, itemID uint64, pretty bool) ([]byte, error) {
	if itemID == 0 {
		return nil, errors.ErrUseCaseEntityNotFound
	}

	// TODO: можно оптимизировать получая только body
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	if pretty {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, item.Body, "", module.JsonPrettyIndent); err != nil {
			return nil, uc.errorWrapper.Wrap(err, "itemId", itemID)
		}

		return prettyJSON.Bytes(), nil
	}

	return item.Body, nil
}

// Create - comment method.
func (uc *ElementTemplate) Create(ctx context.Context, item entity.ElementTemplate) (itemID uint64, err error) {
	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", conv.Group{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *ElementTemplate) Store(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
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
func (uc *ElementTemplate) ChangeStatus(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
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
func (uc *ElementTemplate) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return errors.ErrUseCaseEntityNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", conv.Group{"id": itemID})

	return nil
}
