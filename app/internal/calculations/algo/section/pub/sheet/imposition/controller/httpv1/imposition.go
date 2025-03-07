package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1/model"
)

const (
	sheetImpositionCalcURL         = "/v1/calculations/algo/sheet/imposition"
	sheetImpositionCalcVariantsURL = "/v1/calculations/algo/sheet/imposition/variants"
)

type (
	// SheetImposition - comment struct.
	SheetImposition struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.SheetImpositionUseCase
	}
)

// NewSheetImposition - создаёт контроллер SheetImposition.
func NewSheetImposition(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.SheetImpositionUseCase) *SheetImposition {
	return &SheetImposition{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера SheetImposition.
func (ht *SheetImposition) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: sheetImpositionCalcURL, Func: ht.Calc},
		{Method: http.MethodPost, URL: sheetImpositionCalcVariantsURL, Func: ht.CalcVariants},
	}
}

// Calc - comment method.
func (ht *SheetImposition) Calc(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetImpositionRequest{}

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
func (ht *SheetImposition) CalcVariants(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetImpositionRequest{}

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
