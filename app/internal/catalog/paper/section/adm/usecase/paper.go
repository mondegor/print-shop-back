package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrstatus"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
)

type (
	// Paper - comment struct.
	Paper struct {
		storage                     adm.PaperStorage
		materialTypeAPI             api.MaterialTypeAvailability
		paperColorAPI               api.PaperColorAvailability
		paperFactureAPI             api.PaperFactureAvailability
		eventEmitter                mrevent.Emitter
		errorWrapper                errors.Wrapper
		errorNotFoundWrapper        errors.Wrapper
		errorVersionConflictWrapper errors.Wrapper
		statusFlowMap               mrstatus.FlowMap[itemstatus.Enum]
	}
)

// NewPaper - создаёт объект Paper.
func NewPaper(
	storage adm.PaperStorage,
	materialTypeAPI api.MaterialTypeAvailability,
	paperColorAPI api.PaperColorAvailability,
	paperFactureAPI api.PaperFactureAvailability,
	eventEmitter mrevent.Emitter,
) *Paper {
	return &Paper{
		storage:                     storage,
		materialTypeAPI:             materialTypeAPI,
		paperColorAPI:               paperColorAPI,
		paperFactureAPI:             paperFactureAPI,
		eventEmitter:                mrevent.EmitterWithSource(eventEmitter, entity.ModelNamePaper),
		errorWrapper:                errors.NewServiceOperationFailedWrapper(),
		errorNotFoundWrapper:        errors.NewServiceRecordNotFoundWrapper(),
		errorVersionConflictWrapper: errors.NewServiceRecordVersionConflictWrapper(),
		statusFlowMap:               itemstatus.NewFlowMap(),
	}
}

// GetList - comment method.
func (uc *Paper) GetList(ctx context.Context, params entity.PaperParams) (items []entity.Paper, countItems int, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.Wrap(err)
	}

	if countItems == 0 {
		return make([]entity.Paper, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *Paper) GetItem(ctx context.Context, itemID uint64) (entity.Paper, error) {
	if itemID == 0 {
		return entity.Paper{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Paper{}, uc.errorNotFoundWrapper.Wrap(err, "itemId", itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Paper) Create(ctx context.Context, item entity.Paper) (itemID uint64, err error) {
	if err = uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = itemstatus.Draft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Create", conv.Group{"id": itemID})

	return itemID, nil
}

// Save - comment method.
func (uc *Paper) Save(ctx context.Context, item entity.Paper) error {
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

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		return uc.errorVersionConflictWrapper.Wrap(err)
	}

	uc.eventEmitter.Emit(ctx, "Store", conv.Group{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Paper) ChangeStatus(ctx context.Context, item entity.Paper) error {
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

	uc.eventEmitter.Emit(ctx, "ChangeStatus", conv.Group{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Paper) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return errors.ErrRecordNotFound
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.Wrap(err, "itemId", itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", conv.Group{"id": itemID})

	return nil
}

func (uc *Paper) checkItem(ctx context.Context, item *entity.Paper) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.TypeID > 0 {
		if err := uc.materialTypeAPI.CheckAvailability(ctx, item.TypeID); err != nil {
			return err
		}
	}

	if item.ID == 0 || item.ColorID > 0 {
		if err := uc.paperColorAPI.CheckAvailability(ctx, item.ColorID); err != nil {
			return err
		}
	}

	if item.ID == 0 || item.FactureID > 0 {
		if err := uc.paperFactureAPI.CheckAvailability(ctx, item.FactureID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *Paper) checkArticle(ctx context.Context, item *entity.Paper) error {
	id, err := uc.storage.FetchIDByArticle(ctx, item.Article)
	if err != nil {
		if errors.Is(err, errors.ErrEventStorageNoRecordFound) {
			return nil
		}

		return uc.errorWrapper.Wrap(err)
	}

	if item.ID != id {
		return module.ErrPaperArticleAlreadyExists.New(item.Article)
	}

	return nil
}
