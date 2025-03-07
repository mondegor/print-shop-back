package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1/model"
)

const (
	sheetPackInStackURL = "/v1/calculations/algo/sheet/pack-in-stack"
)

type (
	// PackInStack - comment struct.
	PackInStack struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.SheetPackInStackUseCase
	}
)

// NewSheetPackInStack - создаёт контроллер PackInStack.
func NewSheetPackInStack(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.SheetPackInStackUseCase) *PackInStack {
	return &PackInStack{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера PackInStack.
func (ht *PackInStack) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: sheetPackInStackURL, Func: ht.Calc},
	}
}

// Calc - comment method.
func (ht *PackInStack) Calc(w http.ResponseWriter, r *http.Request) error {
	request := model.SheetPackInStackRequest{}

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
