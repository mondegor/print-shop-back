package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/provider-accounts/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/provider-accounts/entity/providers-api"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/providers-api"
	usecase_shared "print-shop-back/internal/modules/provider-accounts/usecase/shared"

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
	CompanyPage struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.CompanyPageUseCase
	}
)

func NewCompanyPage(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.CompanyPageUseCase,
) *CompanyPage {
	return &CompanyPage{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, companyPageItemURL, "", ht.Get},
		{http.MethodPatch, companyPageItemURL, "", ht.Store},

		{http.MethodPatch, companyPageItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *CompanyPage) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), uuid.MustParse(tmpAccountID)) // :TODO:

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
		AccountID:   uuid.MustParse(tmpAccountID), // :TODO:
		RewriteName: request.RewriteName,
		PageTitle:   request.PageTitle,
		SiteURL:     request.SiteURL,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
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
		AccountID: uuid.MustParse(tmpAccountID),
		Status:    request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPage) wrapError(err error, r *http.Request) error {
	if usecase_shared.FactoryErrCompanyPageRewriteNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("rewriteName", err)
	}

	return err
}
