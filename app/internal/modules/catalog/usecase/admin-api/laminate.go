package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/usecase/shared"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Laminate struct {
		storage         LaminateStorage
		laminateTypeAPI dictionaries.LaminateTypeAPI
		eventEmitter    mrsender.EventEmitter
		usecaseHelper   *mrcore.UsecaseHelper
		statusFlow      mrenum.StatusFlow
	}
)

func NewLaminate(
	storage LaminateStorage,
	laminateTypeAPI dictionaries.LaminateTypeAPI,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *Laminate {
	return &Laminate{
		storage:         storage,
		laminateTypeAPI: laminateTypeAPI,
		eventEmitter:    eventEmitter,
		usecaseHelper:   usecaseHelper,
		statusFlow:      mrenum.ItemStatusFlow,
	}
}

func (uc *Laminate) GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	if total < 1 {
		return []entity.Laminate{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, total, nil
}

func (uc *Laminate) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Laminate, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Laminate{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, id)
	}

	return item, nil
}

func (uc *Laminate) Create(ctx context.Context, item *entity.Laminate) error {
	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *Laminate) Store(ctx context.Context, item *entity.Laminate) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, item.ID)
	}

	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Laminate) ChangeStatus(ctx context.Context, item *entity.Laminate) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Laminate) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameLaminate, id)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Laminate) checkItem(ctx context.Context, item *entity.Laminate) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.TypeID > 0 {
		if err := uc.laminateTypeAPI.CheckingAvailability(ctx, item.TypeID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *Laminate) checkArticle(ctx context.Context, item *entity.Laminate) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrLaminateArticleAlreadyExists.New(item.Article)
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
