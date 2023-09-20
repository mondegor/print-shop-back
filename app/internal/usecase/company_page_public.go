package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
    PublicCompanyPage struct {
        baseImageUrl string
        storage CompanyPageStorage
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewPublicCompanyPage(baseImageUrl string,
                          storage CompanyPageStorage,
                          serviceHelper *mrtool.ServiceHelper) *PublicCompanyPage {
    return &PublicCompanyPage{
        baseImageUrl: baseImageUrl,
        storage: storage,
        serviceHelper: serviceHelper,
    }
}

func (uc *PublicCompanyPage) GetItemByName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error) {
    if rewriteName == "" {
        return nil, mrcore.FactoryErrServiceEmptyInputData.New("rewriteName")
    }

    item := &entity.CompanyPage{RewriteName: rewriteName}
    err := uc.storage.LoadOneByRewriteName(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCompanyPage)
    }

    item.LogoPath = buildLogoUrl(uc.baseImageUrl, item.LogoPath)

    return item, nil
}
