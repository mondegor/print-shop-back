package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/paper/entity/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/paper/usecase/shared"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Paper struct {
		storage         PaperStorage
		paperColorAPI   dictionaries.PaperColorAPI
		paperFactureAPI dictionaries.PaperFactureAPI
		eventEmitter    mrsender.EventEmitter
		usecaseHelper   *mrcore.UsecaseHelper
		statusFlow      mrenum.StatusFlow
	}
)

func NewPaper(
	storage PaperStorage,
	paperColorAPI dictionaries.PaperColorAPI,
	paperFactureAPI dictionaries.PaperFactureAPI,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *Paper {
	return &Paper{
		storage:         storage,
		paperColorAPI:   paperColorAPI,
		paperFactureAPI: paperFactureAPI,
		eventEmitter:    eventEmitter,
		usecaseHelper:   usecaseHelper,
		statusFlow:      mrenum.ItemStatusFlow,
	}
}

func (uc *Paper) GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	if total < 1 {
		return []entity.Paper{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, total, nil
}

func (uc *Paper) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Paper, error) {
	if itemID < 1 {
		return entity.Paper{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.Paper{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, itemID)
	}

	return item, nil
}

func (uc *Paper) Create(ctx context.Context, item entity.Paper) (mrtype.KeyInt32, error) {
	if err := uc.checkItem(ctx, &item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *Paper) Store(ctx context.Context, item entity.Paper) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, item.ID)
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *Paper) ChangeStatus(ctx context.Context, item entity.Paper) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

func (uc *Paper) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNamePaper, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Paper) checkItem(ctx context.Context, item *entity.Paper) error {
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
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrPaperArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Paper) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNamePaper,
		data,
	)
}
