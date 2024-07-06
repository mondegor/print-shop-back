package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/usecase"
)

const (
	rectQuantityInsideOnOutsideURL = "/v1/calculations/algo/rect/inside-on-outside-quantity"
	rectMaxInsideOnOutsideURL      = "/v1/calculations/algo/rect/inside-on-outside-max"
)

type (
	// RectInsideOutside - comment struct.
	RectInsideOutside struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase usecase.RectInsideOutsideUseCase
	}
)

// NewRectInsideOutside - создаёт контроллер RectInsideOutside.
func NewRectInsideOutside(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase usecase.RectInsideOutsideUseCase) *RectInsideOutside {
	return &RectInsideOutside{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера RectInsideOutside.
func (ht *RectInsideOutside) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: rectQuantityInsideOnOutsideURL, Func: ht.CalcQuantity},
		{Method: http.MethodPost, URL: rectMaxInsideOnOutsideURL, Func: ht.CalcMax},
	}
}

// CalcQuantity - comment method.
func (ht *RectInsideOutside) CalcQuantity(w http.ResponseWriter, r *http.Request) error {
	request := CalcRectInsideOutsideRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.RawData{
		InFormat:  request.InFormat,
		OutFormat: request.OutFormat,
	}

	calcResponse, err := ht.useCase.CalcQuantity(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}

// CalcMax - comment method.
func (ht *RectInsideOutside) CalcMax(w http.ResponseWriter, r *http.Request) error {
	request := CalcRectInsideOutsideRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.RawData{
		InFormat:  request.InFormat,
		OutFormat: request.OutFormat,
	}

	calcResponse, err := ht.useCase.CalcMax(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}
