package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/public-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	submitFormListURL = "/v1/controls/submit-forms"
	submitFormItemURL = "/v1/controls/submit-forms/:rewriteName"
)

type (
	SubmitForm struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.SubmitFormUseCase
	}
)

func NewSubmitForm(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.SubmitFormUseCase,
) *SubmitForm {
	return &SubmitForm{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *SubmitForm) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, submitFormListURL, "", ht.GetList},
		{http.MethodGet, submitFormItemURL, "", ht.Get},
	}
}

func (ht *SubmitForm) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *SubmitForm) listParams(r *http.Request) entity.SubmitFormParams {
	return entity.SubmitFormParams{}
}

func (ht *SubmitForm) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItemByRewriteName(r.Context(), ht.parser.PathParamString(r, "rewriteName"))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}
