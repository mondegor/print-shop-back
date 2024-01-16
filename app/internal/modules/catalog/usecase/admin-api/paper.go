package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/usecase/shared"
	catalog "print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Paper struct {
		storage         PaperStorage
		paperColorAPI   catalog.PaperColorAPI
		paperFactureAPI catalog.PaperFactureAPI
		eventBox        mrcore.EventBox
		serviceHelper   *mrtool.ServiceHelper
		statusFlow      mrenum.StatusFlow
	}
)

func NewPaper(
	storage PaperStorage,
	paperColorAPI catalog.PaperColorAPI,
	paperFactureAPI catalog.PaperFactureAPI,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Paper {
	return &Paper{
		storage:         storage,
		paperColorAPI:   paperColorAPI,
		paperFactureAPI: paperFactureAPI,
		eventBox:        eventBox,
		serviceHelper:   serviceHelper,
		statusFlow:      mrenum.ItemStatusFlow,
	}
}

func (uc *Paper) GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	if total < 1 {
		return []entity.Paper{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, total, nil
}

func (uc *Paper) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Paper, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Paper{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, id)
	}

	return item, nil
}

func (uc *Paper) Create(ctx context.Context, item *entity.Paper) error {
	if err := uc.checkPaper(ctx, item); err != nil {
		return err
	}

	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.eventBoxEmitEntity(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *Paper) Store(ctx context.Context, item *entity.Paper) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, item.ID)
	}

	if err := uc.checkPaper(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Paper) ChangeStatus(ctx context.Context, item *entity.Paper) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, item.ID)
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

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Paper) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, id)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Paper) checkPaper(ctx context.Context, item *entity.Paper) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.ColorID > 0 {
		if err := uc.paperColorAPI.CheckingAvailability(ctx, item.ColorID); err != nil {
			return err
		}
	}

	if item.ID == 0 || item.FactureID > 0 {
		if err := uc.paperFactureAPI.CheckingAvailability(ctx, item.FactureID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *Paper) checkArticle(ctx context.Context, item *entity.Paper) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrPaperArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Paper) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNamePaper,
		eventName,
		data,
	)
}
