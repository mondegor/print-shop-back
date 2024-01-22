package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	"print-shop-back/internal/modules/provider-accounts/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/admin-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	companyPageURL = "/v1/provider-accounts/companies-pages"
)

type (
	CompanyPage struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.CompanyPageService
		listSorter mrview.ListSorter
	}
)

func NewCompanyPage(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.CompanyPageService,
	listSorter mrview.ListSorter,
) *CompanyPage {
	return &CompanyPage{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, companyPageURL, "", ht.GetList},
	}
}

func (ht *CompanyPage) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.CompanyPageListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *CompanyPage) listParams(r *http.Request) entity.CompanyPageParams {
	return entity.CompanyPageParams{
		Filter: entity.CompanyPageListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterPublicStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}
