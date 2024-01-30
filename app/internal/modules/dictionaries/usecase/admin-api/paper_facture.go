package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFacture struct {
		storage       PaperFactureStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewPaperFacture(
	storage PaperFactureStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperFacture {
	return &PaperFacture{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *PaperFacture) GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	if total < 1 {
		return []entity.PaperFacture{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	return items, total, nil
}

func (uc *PaperFacture) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.PaperFacture, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.PaperFacture{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, id)
	}

	return item, nil
}

func (uc *PaperFacture) Create(ctx context.Context, item *entity.PaperFacture) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *PaperFacture) Store(ctx context.Context, item *entity.PaperFacture) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *PaperFacture) ChangeStatus(ctx context.Context, item *entity.PaperFacture) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, item.ID)
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

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *PaperFacture) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaperFacture, id)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *PaperFacture) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNamePaperFacture,
		data,
	)
}
