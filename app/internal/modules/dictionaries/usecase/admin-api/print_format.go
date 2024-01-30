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
	PrintFormat struct {
		storage       PrintFormatStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewPrintFormat(
	storage PrintFormatStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *PrintFormat {
	return &PrintFormat{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *PrintFormat) GetList(ctx context.Context, params entity.PrintFormatParams) ([]entity.PrintFormat, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	if total < 1 {
		return []entity.PrintFormat{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	return items, total, nil
}

func (uc *PrintFormat) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.PrintFormat, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.PrintFormat{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePrintFormat, id)
	}

	return item, nil
}

func (uc *PrintFormat) Create(ctx context.Context, item *entity.PrintFormat) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *PrintFormat) Store(ctx context.Context, item *entity.PrintFormat) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePrintFormat, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *PrintFormat) ChangeStatus(ctx context.Context, item *entity.PrintFormat) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePrintFormat, item.ID)
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

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *PrintFormat) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePrintFormat, id)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *PrintFormat) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNamePrintFormat,
		data,
	)
}
