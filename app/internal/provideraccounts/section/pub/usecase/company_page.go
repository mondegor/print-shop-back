package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/entity"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		storage      pub.CompanyPageStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(
	storage pub.CompanyPageStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *CompanyPage {
	return &CompanyPage{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameCompanyPage),
	}
}

// GetItemByRewriteName - comment method.
func (uc *CompanyPage) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error) {
	if rewriteName == "" {
		return entity.CompanyPage{}, mr.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)
	if err != nil {
		return entity.CompanyPage{}, uc.errorWrapper.WrapErrorNotFoundOrFailed(err, "rewriteName", rewriteName)
	}

	return item, nil
}
