package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/entity"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	boxListURL = "/v1/catalog/boxes"
)

type (
	// Box - comment struct.
	// Box - comment struct.
	Box struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.BoxUseCase
	}
)

// NewBox - создаёт контроллер Box.
func NewBox(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.BoxUseCase) *Box {
	return &Box{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера Box.
func (ht *Box) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: boxListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
func (ht *Box) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *Box) listParams(_ *http.Request) entity.BoxParams {
	return entity.BoxParams{}
}
