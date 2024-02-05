package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/provider-account-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	tmpAccountID = "b8a75cd5-ddd1-46ab-be84-406e669cbfa9"

	companyPageItemURL             = "/v1/account/company-page"
	companyPageItemChangeStatusURL = "/v1/account/company-page/status"
)

type (
	CompanyPage struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		service usecase.CompanyPageService
	}
)

func NewCompanyPage(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.CompanyPageService,
) *CompanyPage {
	return &CompanyPage{
		parser:  parser,
		sender:  sender,
		service: service,
	}
}

func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, companyPageItemURL, "", ht.Get},
		{http.MethodPut, companyPageItemURL, "", ht.Store},

		{http.MethodPut, companyPageItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *CompanyPage) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), tmpAccountID)

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *CompanyPage) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreCompanyPageRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID:   tmpAccountID,
		RewriteName: request.RewriteName,
		PageHead:    request.PageHead,
		SiteURL:     request.SiteURL,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPage) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangePublicStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID: tmpAccountID,
		Status:    request.Status,
	}

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPage) wrapError(err error, r *http.Request) error {
	return err
}
