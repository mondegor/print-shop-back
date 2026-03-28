package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/pub/entity"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		storage      pub.CompanyPageStorage
		errorWrapper errors.Wrapper
	}
)

// NewCompanyPage - создаёт объект CompanyPage.
func NewCompanyPage(
	storage pub.CompanyPageStorage,
) *CompanyPage {
	return &CompanyPage{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetItemByRewriteName - comment method.
func (uc *CompanyPage) GetItemByRewriteName(ctx context.Context, rewriteName string) (entity.CompanyPage, error) {
	if rewriteName == "" {
		return entity.CompanyPage{}, errors.ErrRecordNotFound
	}

	item, err := uc.storage.FetchByRewriteName(ctx, rewriteName)
	if err != nil {
		return entity.CompanyPage{}, uc.errorWrapper.Wrap(err, "rewriteName", rewriteName)
	}

	return item, nil
}
