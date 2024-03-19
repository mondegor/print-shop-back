package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		usecaseHelper *mrcore.UsecaseHelper
		imgBaseURL    mrlib.BuilderPath
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	usecaseHelper *mrcore.UsecaseHelper,
	imgBaseURL mrlib.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		usecaseHelper: usecaseHelper,
		imgBaseURL:    imgBaseURL,
	}
}

func (uc *CompanyPage) GetList(ctx context.Context, params entity.CompanyPageParams) ([]entity.CompanyPage, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	if total < 1 {
		return []entity.CompanyPage{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCompanyPage)
	}

	for i := range items {
		items[i].LogoURL = uc.imgBaseURL.FullPath(items[i].LogoURL)
	}

	return items, total, nil
}
