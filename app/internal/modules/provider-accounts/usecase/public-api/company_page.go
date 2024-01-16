package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/provider-accounts/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	CompanyPage struct {
		storage       CompanyPageStorage
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewCompanyPage(
	storage CompanyPageStorage,
	serviceHelper *mrtool.ServiceHelper,
) *CompanyPage {
	return &CompanyPage{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *CompanyPage) GetItemByName(ctx context.Context, rewriteName string) (*entity.CompanyPage, error) {
	if rewriteName == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)

	if err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCompanyPage, rewriteName)
	}

	return item, nil
}
