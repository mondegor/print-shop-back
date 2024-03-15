package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/admin-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/admin-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	companyPageListURL = "/v1/provider-accounts/companies-pages"
)

type (
	CompanyPage struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.CompanyPageUseCase
		listSorter mrview.ListSorter
	}
)

func NewCompanyPage(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.CompanyPageUseCase,
	listSorter mrview.ListSorter,
) *CompanyPage {
	return &CompanyPage{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, companyPageListURL, "", ht.GetList},
	}
}

func (ht *CompanyPage) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		CompanyPageListResponse{
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
