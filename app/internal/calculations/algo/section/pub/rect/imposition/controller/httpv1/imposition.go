package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
)

const (
	rectImpositionCalcURL         = "/v1/calculations/algo/rect/imposition"
	rectImpositionCalcVariantsURL = "/v1/calculations/algo/rect/imposition/variants"
)

type (
	// RectImposition - comment struct.
	RectImposition struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.RectImpositionUseCase
	}
)

// NewRectImposition - создаёт контроллер RectImposition.
func NewRectImposition(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.RectImpositionUseCase) *RectImposition {
	return &RectImposition{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера RectImposition.
func (ht *RectImposition) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: rectImpositionCalcURL, Func: ht.Calc},
		{Method: http.MethodPost, URL: rectImpositionCalcVariantsURL, Func: ht.CalcVariants},
	}
}

// Calc - comment method.
func (ht *RectImposition) Calc(w http.ResponseWriter, r *http.Request) error {
	request := CalcRectImpositionRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item, err := ht.parseRequest(request)
	if err != nil {
		return err
	}

	calcResponse, err := ht.useCase.Calc(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}

// CalcVariants - comment method.
func (ht *RectImposition) CalcVariants(w http.ResponseWriter, r *http.Request) error {
	request := CalcRectImpositionRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item, err := ht.parseRequest(request)
	if err != nil {
		return err
	}

	calcResponse, err := ht.useCase.CalcVariants(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}
