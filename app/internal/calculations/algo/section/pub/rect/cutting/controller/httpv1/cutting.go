package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/entity"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/usecase"
)

const (
	rectCuttingQuantityURL = "/v1/calculations/algo/rect/cutting-quantity"
)

type (
	// RectCutting - comment struct.
	RectCutting struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase usecase.RectCuttingUseCase
	}
)

// NewRectCutting - создаёт контроллер RectCutting.
func NewRectCutting(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase usecase.RectCuttingUseCase) *RectCutting {
	return &RectCutting{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера RectCutting.
func (ht *RectCutting) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: rectCuttingQuantityURL, Func: ht.CalcQuantity},
	}
}

// CalcQuantity - comment method.
func (ht *RectCutting) CalcQuantity(w http.ResponseWriter, r *http.Request) error {
	request := CalcCuttingQuantityRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.RawData{
		Fragments:      request.Fragments,
		DistanceFormat: request.DistanceFormat,
	}

	calcResponse, err := ht.useCase.CalcQuantity(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}
