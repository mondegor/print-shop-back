package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1/model"
)

const (
	sheetQuantityInsideOnOutsideURL = "/v1/calculations/algo/sheet/inside-on-outside-quantity"
	sheetMaxInsideOnOutsideURL      = "/v1/calculations/algo/sheet/inside-on-outside-max"
)

type (
	// SheetInsideOutside - comment struct.
	SheetInsideOutside struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.SheetInsideOutsideUseCase
	}
)

// NewSheetInsideOutside - создаёт контроллер SheetInsideOutside.
func NewSheetInsideOutside(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.SheetInsideOutsideUseCase) *SheetInsideOutside {
	return &SheetInsideOutside{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера SheetInsideOutside.
func (ht *SheetInsideOutside) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: sheetQuantityInsideOnOutsideURL, Func: ht.CalcQuantity},
		{Method: http.MethodPost, URL: sheetMaxInsideOnOutsideURL, Func: ht.CalcMax},
	}
}

// CalcQuantity - comment method.
func (ht *SheetInsideOutside) CalcQuantity(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetInsideOutsideQuantityRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item, err := ht.parseRequest(request)
	if err != nil {
		return err
	}

	calcResponse, err := ht.useCase.CalcQuantity(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}

// CalcMax - comment method.
func (ht *SheetInsideOutside) CalcMax(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetInsideOutsideMaxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item, err := ht.parseRequest(model.SheetInsideOutsideQuantityRequest(request))
	if err != nil {
		return err
	}

	calcResponse, err := ht.useCase.CalcMax(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}
