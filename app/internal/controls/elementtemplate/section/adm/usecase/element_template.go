package usecase

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrworkflow/itemstatus"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/controls/elementtemplate/module"
	"print-shop-back/internal/controls/elementtemplate/section/adm"
	"print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplate - comment struct.
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		storage                     adm.ElementTemplateStorage
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
		statusFlowMap               workflow.FlowMap[workflow.ItemStatus]
	}
)

// NewElementTemplate - создаёт объект ElementTemplate.
func NewElementTemplate(
	storage adm.ElementTemplateStorage,
	eventEmitter mrevent.Emitter,
) *ElementTemplate {
	return &ElementTemplate{
		storage:                     storage,
		eventEmitter:                mrevent.EmitterWithSource(eventEmitter, entity.ModelNameElementTemplate),
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
		statusFlowMap:               itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *ElementTemplate) GetList(ctx context.Context, params entity.ElementTemplateParams) (items []entity.ElementTemplate, countItems int, err error) {
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
		return entity.ElementTemplate{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.ElementTemplate{}, uc.errorNotFoundWrapper.Wrap(err, "itemId", itemID)
	}

	// TODO: можно оптимизировать удалив body
	item.Body = nil

	return item, nil
}

// GetItemJson - comment method.
func (uc *ElementTemplate) GetItemJson(ctx context.Context, itemID uint64, pretty bool) ([]byte, error) {
	if itemID == 0 {
		return nil, errors.ErrRecordNotFound
	}

	// TODO: можно оптимизировать получая только body
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return nil, uc.errorNotFoundWrapper.Wrap(err, "itemId", itemID)
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

	uc.eventEmitter.Emit(ctx, "Create", "itemId", itemID)

	return itemID, nil
}

// Save - comment method.
func (uc *ElementTemplate) Save(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
		return errors.ErrRecordNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrRecordVersionConflict
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionConflict
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorNotFoundWrapper.Wrap(err, "itemId", item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", "itemId", item.ID, "tagVersion", tagVersion)

	return nil
}

// ChangeStatus - comment method.
func (uc *ElementTemplate) ChangeStatus(ctx context.Context, item entity.ElementTemplate) error {
	if item.ID == 0 {
		return errors.ErrRecordNotFound
	}

	if item.TagVersion == 0 {
		return errors.ErrRecordVersionConflict
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorNotFoundWrapper.Wrap(err, "itemId", item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlowMap.IsPossible(currentStatus, item.Status) {
		return errors.ErrSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "ChangeStatus", "itemId", item.ID, "tagVersion", tagVersion, "status", item.Status)

	return nil
}

// Remove - comment method.
func (uc *ElementTemplate) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return errors.ErrRecordNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", "itemId", itemID)

	return nil
}
