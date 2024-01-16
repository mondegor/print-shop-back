package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/controller/http_v1/provider-account-api/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	tmpAccountID = "b8a75cd5-ddd1-46ab-be84-406e669cbfa9"

	companyPageItemURL         = "/v1/account/company-page"
	companyPageChangeStatusURL = "/v1/account/company-page/status"
)

type (
	CompanyPage struct {
		section mrcore.ClientSection
		service usecase.CompanyPageService
	}
)

func NewCompanyPage(
	section mrcore.ClientSection,
	service usecase.CompanyPageService,
) *CompanyPage {
	return &CompanyPage{
		section: section,
		service: service,
	}
}

func (ht *CompanyPage) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCompanyPagePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(companyPageItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(companyPageItemURL), moduleAccessFunc(ht.Store()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(companyPageChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *CompanyPage) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), tmpAccountID)

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *CompanyPage) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreCompanyPageRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.CompanyPage{
			AccountID:   tmpAccountID,
			RewriteName: request.RewriteName,
			PageHead:    request.PageHead,
			SiteURL:     request.SiteURL,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CompanyPage) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangePublicStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.CompanyPage{
			AccountID: tmpAccountID,
			Status:    request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CompanyPage) wrapError(err error, c mrcore.ClientContext) error {
	return err
}
