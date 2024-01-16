package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/admin-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	companyPageURL = "/v1/provider-accounts/companies-pages"
)

type (
	CompanyPage struct {
		section    mrcore.ClientSection
		service    usecase.CompanyPageService
		listSorter mrview.ListSorter
	}
)

func NewCompanyPage(
	section mrcore.ClientSection,
	service usecase.CompanyPageService,
	listSorter mrview.ListSorter,
) *CompanyPage {
	return &CompanyPage{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *CompanyPage) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCompanyPagePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(companyPageURL), moduleAccessFunc(ht.GetList()))
}

func (ht *CompanyPage) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.CompanyPageListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *CompanyPage) listParams(c mrcore.ClientContext) entity.CompanyPageParams {
	return entity.CompanyPageParams{
		Filter: entity.CompanyPageListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Statuses:   view_shared.ParseFilterPublicStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}
