package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/usecase"
)

const (
	circulationPackInBoxURL = "/v1/calculations/algo/circulation/pack-in-box"
)

type (
	// PackInBox - comment struct.
	PackInBox struct {
		parser  mrserver.RequestParserValidate
		sender  mrserver.ResponseSender
		useCase usecase.CirculationPackInBoxUseCase
	}
)

// NewPackInBox - создаёт контроллер PackInBox.
func NewPackInBox(parser mrserver.RequestParserValidate, sender mrserver.ResponseSender, useCase usecase.CirculationPackInBoxUseCase) *PackInBox {
	return &PackInBox{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера PackInBox.
func (ht *PackInBox) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: circulationPackInBoxURL, Func: ht.Calc},
	}
}

// Calc - comment method.
func (ht *PackInBox) Calc(w http.ResponseWriter, r *http.Request) error {
	request := CalcCirculationPackInBoxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.RawData{
		Format: request.Format,
	}

	calcResponse, err := ht.useCase.CalcQuantity(r.Context(), item)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, calcResponse)
}
