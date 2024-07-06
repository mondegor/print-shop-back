package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	tmpAccountID = "b8a75cd5-ddd1-46ab-be84-406e669cbfa9"

	companyPageItemURL             = "/v1/account/company-page"
	companyPageItemChangeStatusURL = "/v1/account/company-page/status"
)

type (
	// CompanyPage - comment struct.
	CompanyPage struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase prov.CompanyPageUseCase
	}
)

// NewCompanyPage - создаёт контроллер CompanyPage.
func NewCompanyPage(parser validate.RequestParser, sender mrserver.ResponseSender, useCase prov.CompanyPageUseCase) *CompanyPage {
	return &CompanyPage{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера CompanyPage.
func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: companyPageItemURL, Func: ht.Get},
		{Method: http.MethodPatch, URL: companyPageItemURL, Func: ht.Store},

		{Method: http.MethodPatch, URL: companyPageItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// Get - comment method.
func (ht *CompanyPage) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), uuid.MustParse(tmpAccountID)) // TODO:
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Store - comment method.
func (ht *CompanyPage) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreCompanyPageRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID:   uuid.MustParse(tmpAccountID), // TODO:
		RewriteName: request.RewriteName,
		PageTitle:   request.PageTitle,
		SiteURL:     request.SiteURL,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *CompanyPage) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangePublicStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID: uuid.MustParse(tmpAccountID),
		Status:    request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPage) wrapError(err error, _ *http.Request) error {
	if module.ErrCompanyPageRewriteNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("rewriteName", err)
	}

	return err
}
