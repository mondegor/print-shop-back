package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrpath"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		storage      adm.CompanyPageStorage
		errorWrapper mrcore.UsecaseErrorWrapper
		imgBaseURL   mrpath.PathBuilder
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(storage adm.CompanyPageStorage, errorWrapper mrcore.UsecaseErrorWrapper, imgBaseURL mrpath.PathBuilder) *CompanyPage {
	return &CompanyPage{
		storage:      storage,
		errorWrapper: errorWrapper,
		imgBaseURL:   imgBaseURL,
	}
}

// GetList - comment method.
func (uc *CompanyPage) GetList(ctx context.Context, params entity.CompanyPageParams) ([]entity.CompanyPage, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if total < 1 {
		return nil, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	for i := range items {
		uc.prepareItem(&items[i])
	}

	return items, total, nil
}

func (uc *CompanyPage) prepareItem(item *entity.CompanyPage) {
	item.LogoURL = uc.imgBaseURL.BuildPath(item.LogoURL)
}
