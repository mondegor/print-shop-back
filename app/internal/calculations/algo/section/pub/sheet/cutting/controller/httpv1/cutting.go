package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1/model"
)

const (
	sheetCuttingQuantityURL = "/v1/calculations/algo/sheet/cutting-quantity"
)

type (
	// SheetCutting - comment struct.
	SheetCutting struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.SheetCuttingUseCase
	}
)

// NewSheetCutting - создаёт контроллер SheetCutting.
func NewSheetCutting(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.SheetCuttingUseCase) *SheetCutting {
	return &SheetCutting{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера SheetCutting.
func (ht *SheetCutting) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: sheetCuttingQuantityURL, Func: ht.CalcQuantity},
	}
}

// CalcQuantity - comment method.
func (ht *SheetCutting) CalcQuantity(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetCuttingQuantityRequest{}

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
