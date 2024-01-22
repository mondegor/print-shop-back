package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/provider-account-api"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

const (
	companyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
	CompanyPageLogo struct {
		// parser     mrserver.RequestParser
		sender  mrserver.ResponseSender
		service usecase.CompanyPageLogoService
	}
)

func NewCompanyPageLogo(
	// parser     mrserver.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.CompanyPageLogoService,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		// parser: parser,
		sender:  sender,
		service: service,
	}
}

func (ht *CompanyPageLogo) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodPut, companyPageItemLogoURL, "", ht.UploadLogo},
		{http.MethodDelete, companyPageItemLogoURL, "", ht.RemoveLogo},
	}
}

func (ht *CompanyPageLogo) UploadLogo(w http.ResponseWriter, r *http.Request) error {
	file, err := mrreq.File(r, module.ParamNameFileCompanyLogo)

	if err != nil {
		return err
	}

	defer file.Body.Close()

	if err = ht.service.StoreFile(r.Context(), tmpAccountID, file); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPageLogo) RemoveLogo(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.RemoveFile(r.Context(), tmpAccountID); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPageLogo) wrapError(err error, r *http.Request) error {
	return err
}
