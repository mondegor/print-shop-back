package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/providers-api"
	usecase_shared "print-shop-back/internal/modules/provider-accounts/usecase/shared"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		imgBaseURL    mrlib.BuilderPath
		statusFlow    mrenum.StatusFlow
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
	imgBaseURL mrlib.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		imgBaseURL:    imgBaseURL,
		statusFlow:    enums.PublicStatusFlow,
	}
}

func (uc *CompanyPage) GetItem(ctx context.Context, accountID uuid.UUID) (entity.CompanyPage, error) {
	if accountID == uuid.Nil {
		return entity.CompanyPage{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, accountID)

	if err != nil {
		return entity.CompanyPage{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, accountID)
	}

	item.LogoURL = uc.imgBaseURL.FullPath(item.LogoURL)

	return item, nil
}

func (uc *CompanyPage) Store(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.checkItem(ctx, &item); err != nil {
		return err
	}

	item.Status = enums.PublicStatusDraft // only for insert

	if err := uc.storage.InsertOrUpdate(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"accountId": item.AccountID})

	return nil
}

func (uc *CompanyPage) ChangeStatus(ctx context.Context, item entity.CompanyPage) error {
	if item.AccountID == uuid.Nil {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"accountId": item.AccountID, "status": item.Status})

	return nil
}

func (uc *CompanyPage) checkItem(ctx context.Context, item *entity.CompanyPage) error {
	if err := uc.checkRewriteName(ctx, item); err != nil {
		return err
	}

	return nil
}

func (uc *CompanyPage) checkRewriteName(ctx context.Context, item *entity.CompanyPage) error {
	accountID, err := uc.storage.FetchAccountIdByRewriteName(ctx, item.RewriteName)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if item.AccountID != accountID {
		return usecase_shared.FactoryErrCompanyPageRewriteNameAlreadyExists.New(item.RewriteName)
	}

	return nil
}

func (uc *CompanyPage) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCompanyPage,
		data,
	)
}
