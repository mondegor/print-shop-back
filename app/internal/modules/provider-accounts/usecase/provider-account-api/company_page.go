package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
		imgBaseURL    mrcore.BuilderPath
		statusFlow    mrenum.StatusFlow
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
	imgBaseURL mrcore.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
		imgBaseURL:    imgBaseURL,
		statusFlow:    entity_shared.PublicStatusFlow,
	}
}

func (uc *CompanyPage) GetItem(ctx context.Context, accountID mrtype.KeyString) (*entity.CompanyPage, error) {
	if accountID == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.CompanyPage{AccountID: accountID}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, accountID)
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
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"accountId": item.AccountID})

	return nil
}

func (uc *CompanyPage) ChangeStatus(ctx context.Context, item *entity.CompanyPage) error {
	if item.AccountID == "" {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	if err = uc.storage.UpdateStatus(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, item.AccountID)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"accountId": item.AccountID, "status": item.Status})

	return nil
}

func (uc *CompanyPage) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameCompanyPage,
		eventName,
		data,
	)
}
