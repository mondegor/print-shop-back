package usecase

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
    CompanyPage struct {
        baseImageUrl string
        storage CompanyPageStorage
    }
)

func NewCompanyPage(baseImageUrl string,
                    storage CompanyPageStorage) *CompanyPage {
    return &CompanyPage{
        baseImageUrl: baseImageUrl,
        storage: storage,
    }
}

func (uc *CompanyPage) GetList(ctx context.Context, listFilter *entity.CompanyPageListFilter) ([]entity.CompanyPage, error) {
    items := make([]entity.CompanyPage, 0, 4)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCompanyPage)
    }

    for i, _ := range items {
        items[i].LogoPath = buildLogoUrl(uc.baseImageUrl, items[i].LogoPath)
    }

    return items, nil
}
