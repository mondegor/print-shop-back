package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	companyPageItemURL = "/v1/company/:rewriteName"
)

type (
	CompanyPage struct {
		section    mrcore.ClientSection
		service    usecase.CompanyPageService
		imgBaseURL mrcore.BuilderPath
	}
)

func NewCompanyPage(
	section mrcore.ClientSection,
	service usecase.CompanyPageService,
	imgBaseURL mrcore.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		section:    section,
		service:    service,
		imgBaseURL: imgBaseURL,
	}
}

func (ht *CompanyPage) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCompanyPagePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(companyPageItemURL), moduleAccessFunc(ht.Get()))
}

func (ht *CompanyPage) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItemByName(c.Context(), c.ParamFromPath("rewriteName"))

		if err != nil {
			return err
		}

		item.LogoURL = ht.imgBaseURL.FullPath(item.LogoURL)

		return c.SendResponse(http.StatusOK, item)
	}
}
