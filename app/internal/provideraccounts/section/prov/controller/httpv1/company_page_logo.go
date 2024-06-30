package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/prov/usecase"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

const (
	companyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
	// CompanyPageLogo - comment struct.
	CompanyPageLogo struct {
		parser  mrserver.RequestParserImage
		sender  mrserver.ResponseSender
		useCase usecase.CompanyPageLogoUseCase
	}
)

// NewCompanyPageLogo - создаёт объект CompanyPageLogo.
func NewCompanyPageLogo(parser mrserver.RequestParserImage, sender mrserver.ResponseSender, useCase usecase.CompanyPageLogoUseCase) *CompanyPageLogo {
	return &CompanyPageLogo{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - comment method.
func (ht *CompanyPageLogo) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPatch, URL: companyPageItemLogoURL, Func: ht.UploadLogo},
		{Method: http.MethodDelete, URL: companyPageItemLogoURL, Func: ht.RemoveLogo},
	}
}

// UploadLogo - comment method.
func (ht *CompanyPageLogo) UploadLogo(w http.ResponseWriter, r *http.Request) error {
	file, err := ht.parser.FormImage(r, module.ParamNameFileCompanyLogo)
	if err != nil {
		return mrparser.WrapImageError(err, module.ParamNameFileCompanyLogo)
	}

	defer file.Body.Close()

	if err = ht.useCase.StoreFile(r.Context(), uuid.MustParse(tmpAccountID), file); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// RemoveLogo - comment method.
func (ht *CompanyPageLogo) RemoveLogo(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.RemoveFile(r.Context(), uuid.MustParse(tmpAccountID)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPageLogo) wrapError(err error, _ *http.Request) error {
	return err
}
