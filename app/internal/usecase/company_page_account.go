package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
    AccountCompanyPage struct {
        baseImageUrl string
        storage CompanyPageStorage
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
        statusFlow entity.ResourceStatusFlow
    }
)

func NewAccountCompanyPage(baseImageUrl string,
                           storage CompanyPageStorage,
                           eventBox mrcore.EventBox,
                           serviceHelper *mrtool.ServiceHelper) *AccountCompanyPage {
    return &AccountCompanyPage{
        baseImageUrl: baseImageUrl,
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: entity.ResourceStatusFlowDefault,
    }
}

func (uc *AccountCompanyPage) GetItem(ctx context.Context, accountId mrentity.KeyString) (*entity.CompanyPage, error) {
    if accountId == "" {
        return nil, mrcore.FactoryErrServiceEmptyInputData.New("accountId")
    }

    item := &entity.CompanyPage{AccountId: accountId}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCompanyPage)
    }

    item.LogoPath = buildLogoUrl(uc.baseImageUrl, item.LogoPath)

    return item, nil
}

func (uc *AccountCompanyPage) Store(ctx context.Context, item *entity.CompanyPage) error {
    if item.AccountId == "" {
        return mrcore.FactoryErrServiceEmptyInputData.New("item.AccountId")
    }

    item.Status = entity.ResourceStatusDraft // only for insert
    err := uc.storage.InsertOrUpdate(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCompanyPage)
    }

    uc.eventBox.Emit(
        "%s::Store: accountId=%s",
        entity.ModelNameCompanyPage,
        item.AccountId,
    )

    return nil
}

func (uc *AccountCompanyPage) ChangeStatus(ctx context.Context, item *entity.CompanyPage) error {
    if item.AccountId == "" || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.AccountId": item.AccountId, "Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCompanyPage)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCompanyPage, item.AccountId)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCompanyPage)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: accountId=%s, status=%s",
        entity.ModelNameCompanyPage,
        item.AccountId,
        item.Status,
    )

    return nil
}
