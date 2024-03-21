package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/provider-accounts"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/providers-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

const (
	companyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
	CompanyPageLogo struct {
		parser  mrserver.RequestParserImage
		sender  mrserver.ResponseSender
		useCase usecase.CompanyPageLogoUseCase
	}
)

func NewCompanyPageLogo(
	parser mrserver.RequestParserImage,
	sender mrserver.ResponseSender,
	useCase usecase.CompanyPageLogoUseCase,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *CompanyPageLogo) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodPatch, companyPageItemLogoURL, "", ht.UploadLogo},
		{http.MethodDelete, companyPageItemLogoURL, "", ht.RemoveLogo},
	}
}

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

func (ht *CompanyPageLogo) RemoveLogo(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.RemoveFile(r.Context(), uuid.MustParse(tmpAccountID)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPageLogo) wrapError(err error, r *http.Request) error {
	return err
}
