package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/usecase/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Box struct {
		storage       BoxStorage
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewBox(
	storage BoxStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Box {
	return &Box{
		storage:       storage,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if total < 1 {
		return []entity.Box{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	return items, total, nil
}

func (uc *Box) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Box, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Box{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, id)
	}

	return item, nil
}

func (uc *Box) Create(ctx context.Context, item *entity.Box) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventBoxEmitEntity(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *Box) Store(ctx context.Context, item *entity.Box) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
	}

	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Box) ChangeStatus(ctx context.Context, item *entity.Box) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Box) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameBox, id)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Box) checkArticle(ctx context.Context, item *entity.Box) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrBoxArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Box) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameBox,
		eventName,
		data,
	)
}
