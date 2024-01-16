package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		serviceHelper *mrtool.ServiceHelper
		imgBaseURL    mrcore.BuilderPath
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	serviceHelper *mrtool.ServiceHelper,
	imgBaseURL mrcore.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		serviceHelper: serviceHelper,
		imgBaseURL:    imgBaseURL,
	}
}

func (uc *CompanyPage) GetList(ctx context.Context, params entity.CompanyPageParams) ([]entity.CompanyPage, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if total < 1 {
		return []entity.CompanyPage{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	for i := range items {
		items[i].LogoURL = uc.imgBaseURL.FullPath(items[i].LogoURL)
	}

	return items, total, nil
}
