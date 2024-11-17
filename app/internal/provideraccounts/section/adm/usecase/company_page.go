package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrpath"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		storage      adm.CompanyPageStorage
		errorWrapper mrcore.UseCaseErrorWrapper
		imgBaseURL   mrpath.PathBuilder
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(storage adm.CompanyPageStorage, errorWrapper mrcore.UseCaseErrorWrapper, imgBaseURL mrpath.PathBuilder) *CompanyPage {
	return &CompanyPage{
		storage:      storage,
		errorWrapper: errorWrapper,
		imgBaseURL:   imgBaseURL,
	}
}

// GetList - comment method.
func (uc *CompanyPage) GetList(ctx context.Context, params entity.CompanyPageParams) (items []entity.CompanyPage, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if countItems == 0 {
		return make([]entity.CompanyPage, 0), 0, nil
	}

	for i := range items {
		uc.prepareItem(&items[i])
	}

	return items, countItems, nil
}

func (uc *CompanyPage) prepareItem(item *entity.CompanyPage) {
	item.LogoURL = uc.imgBaseURL.BuildPath(item.LogoURL)
}
