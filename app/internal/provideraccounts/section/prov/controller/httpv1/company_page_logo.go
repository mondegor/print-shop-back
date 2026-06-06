package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"

	"print-shop-back/internal/provideraccounts/module"
	"print-shop-back/internal/provideraccounts/section/prov"
)

const (
	companyPageItemLogoURL = "/v1/account/company-page/logo"
)

type (
	// CompanyPageLogo - comment struct.
	CompanyPageLogo struct {
		parser            requestParser
		sender            mrserver.ResponseSender
		useCase           prov.CompanyPageLogoUseCase
		imageErrorWrapper errors.CustomWrapper
	}

	requestParser interface {
		request.ParserImage
		request.ParserUser
	}
)

// NewCompanyPageLogo - создаёт контроллер CompanyPageLogo.
func NewCompanyPageLogo(
	parser requestParser,
	sender mrserver.ResponseSender,
	useCase prov.CompanyPageLogoUseCase,
) *CompanyPageLogo {
	return &CompanyPageLogo{
		parser:            parser,
		sender:            sender,
		useCase:           useCase,
		imageErrorWrapper: errors.NewDownloadImageWrapper(module.ParamNameFileCompanyLogo),
	}
}

// Handlers - возвращает обработчики контроллера CompanyPageLogo.
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
		return ht.imageErrorWrapper.Wrap(err)
	}

	defer func() {
		_ = file.Body.Close()
	}()

	if err = ht.useCase.StoreFile(r.Context(), ht.parser.UserID(r), file); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// RemoveLogo - comment method.
func (ht *CompanyPageLogo) RemoveLogo(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.RemoveFile(r.Context(), ht.parser.UserID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CompanyPageLogo) wrapError(err error, _ *http.Request) error {
	return err
}
