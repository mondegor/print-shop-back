package http_v1

import (
	"net/http"
	usecase "print-shop-back/internal/modules/provider-accounts/usecase/public-api"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	companyPageItemURL = "/v1/company/:rewriteName"
)

type (
	CompanyPage struct {
		parser     mrserver.RequestParserString
		sender     mrserver.ResponseSender
		useCase    usecase.CompanyPageUseCase
		imgBaseURL mrlib.BuilderPath
	}
)

func NewCompanyPage(
	parser mrserver.RequestParserString,
	sender mrserver.ResponseSender,
	useCase usecase.CompanyPageUseCase,
	imgBaseURL mrlib.BuilderPath,
) *CompanyPage {
	return &CompanyPage{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		imgBaseURL: imgBaseURL,
	}
}

func (ht *CompanyPage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, companyPageItemURL, "", ht.Get},
	}
}

func (ht *CompanyPage) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItemByRewriteName(r.Context(), ht.parser.PathParamString(r, "rewriteName"))

	if err != nil {
		return err
	}

	item.LogoURL = ht.imgBaseURL.FullPath(item.LogoURL)

	return ht.sender.Send(w, http.StatusOK, item)
}
