package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	paperColorListURL = "/v1/dictionaries/paper-colors"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.PaperColorUseCase
	}
)

// NewPaperColor - создаёт контроллер PaperColor.
func NewPaperColor(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.PaperColorUseCase) *PaperColor {
	return &PaperColor{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера PaperColor.
func (ht *PaperColor) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: paperColorListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
func (ht *PaperColor) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *PaperColor) listParams(_ *http.Request) entity.PaperColorParams {
	return entity.PaperColorParams{}
}
