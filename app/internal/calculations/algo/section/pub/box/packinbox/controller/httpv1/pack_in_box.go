package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1/model"
)

const (
	boxPackInBoxURL = "/v1/calculations/algo/box/pack-in-box"
)

type (
	// BoxPackInBox - comment struct.
	BoxPackInBox struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase pub.BoxPackInBoxUseCase
	}
)

// NewBoxPackInBox - создаёт контроллер BoxPackInBox.
func NewBoxPackInBox(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase pub.BoxPackInBoxUseCase) *BoxPackInBox {
	return &BoxPackInBox{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера BoxPackInBox.
func (ht *BoxPackInBox) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: boxPackInBoxURL, Func: ht.Calc},
	}
}

// Calc - comment method.
func (ht *BoxPackInBox) Calc(w http.ResponseWriter, r *http.Request) error {
	request := model.CalcBoxPackInBoxRequest{}

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
