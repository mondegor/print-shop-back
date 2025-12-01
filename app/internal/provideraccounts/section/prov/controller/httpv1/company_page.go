package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/entity"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
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
	item, err := ht.useCase.GetItem(r.Context(), ht.parser.UserID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Store - comment method.
func (ht *CompanyPage) Store(w http.ResponseWriter, r *http.Request) error {
	req := StoreCompanyPageRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID:   ht.parser.UserID(r),
		RewriteName: req.RewriteName,
		PageTitle:   req.PageTitle,
		SiteURL:     req.SiteURL,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *CompanyPage) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	req := ChangePublicStatusRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := entity.CompanyPage{
		AccountID: ht.parser.UserID(r),
		Status:    req.Status,
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
