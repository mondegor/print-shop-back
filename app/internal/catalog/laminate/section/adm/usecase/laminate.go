package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// Laminate - comment struct.
	Laminate struct {
		storage         adm.LaminateStorage
		materialTypeAPI api.MaterialTypeAvailability
		eventEmitter    mrsender.EventEmitter
		errorWrapper    mrcore.UsecaseErrorWrapper
		statusFlow      mrstatus.Flow
	}
)

// NewLaminate - создаёт объект NewLaminate.
func NewLaminate(
	storage adm.LaminateStorage,
	materialTypeAPI api.MaterialTypeAvailability,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
) *Laminate {
	return &Laminate{
		storage:         storage,
		materialTypeAPI: materialTypeAPI,
		eventEmitter:    eventEmitter,
		errorWrapper:    errorWrapper,
		statusFlow:      mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Laminate) GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	if total < 1 {
		return make([]entity.Laminate, 0), 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *Laminate) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Laminate, error) {
	if itemID < 1 {
		return entity.Laminate{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Laminate{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Laminate) Create(ctx context.Context, item entity.Laminate) (mrtype.KeyInt32, error) {
	if err := uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *Laminate) Store(ctx context.Context, item entity.Laminate) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Laminate) ChangeStatus(ctx context.Context, item entity.Laminate) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Laminate) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Laminate) checkItem(ctx context.Context, item *entity.Laminate) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.TypeID > 0 {
		if err := uc.materialTypeAPI.CheckingAvailability(ctx, item.TypeID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *Laminate) checkArticle(ctx context.Context, item *entity.Laminate) error {
	id, err := uc.storage.FetchIDByArticle(ctx, item.Article)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	if item.ID != id {
		return module.ErrLaminateArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Laminate) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameLaminate,
		data,
	)
}
