package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *CompanyPage) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error) {
	if rewriteName == "" {
		return entity.CompanyPage{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)

	if err != nil {
		return entity.CompanyPage{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, rewriteName)
	}

	return item, nil
}
