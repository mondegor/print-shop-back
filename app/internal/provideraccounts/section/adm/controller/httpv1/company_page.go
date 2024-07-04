package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/shared/validate"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	companyPageListURL = "/v1/prov/companies-pages"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		parser     validate.RequestProviderAccountsParser
		sender     mrserver.ResponseSender
		useCase    usecase.CompanyPageUseCase
		listSorter mrview.ListSorter
	}
)

// NewCompanyPage - создаёт контроллер CompanyPage.
func NewCompanyPage(
	parser validate.RequestProviderAccountsParser,
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

// Handlers - возвращает обработчики контроллера CompanyPage.
func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: companyPageListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
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
