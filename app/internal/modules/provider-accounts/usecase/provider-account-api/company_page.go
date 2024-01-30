package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/go-webcore/mrtype"
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
		statusFlow:    entity_shared.PublicStatusFlow,
	}
}

func (uc *CompanyPage) GetItem(ctx context.Context, accountID mrtype.KeyString) (*entity.CompanyPage, error) {
	if accountID == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.CompanyPage{
		AccountID: accountID,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, accountID)
	}

	item.LogoURL = uc.imgBaseURL.FullPath(item.LogoURL)

	return item, nil
}

func (uc *CompanyPage) Store(ctx context.Context, item *entity.CompanyPage) error {
	if item.AccountID == "" {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item.Status = entity_shared.PublicStatusDraft // only for insert

	if err := uc.storage.InsertOrUpdate(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"accountId": item.AccountID})

	return nil
}

func (uc *CompanyPage) ChangeStatus(ctx context.Context, item *entity.CompanyPage) error {
	if item.AccountID == "" {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"accountId": item.AccountID, "status": item.Status})

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
