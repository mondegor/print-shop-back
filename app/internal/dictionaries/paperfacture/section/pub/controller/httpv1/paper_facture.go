package httpv1

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	paperFactureListURL = "/v1/dictionaries/paper-factures"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase pub.PaperFactureUseCase
	}
)

// NewPaperFacture - создаёт контроллер PaperFacture.
func NewPaperFacture(parser validate.RequestParser, sender mrserver.ResponseSender, useCase pub.PaperFactureUseCase) *PaperFacture {
	return &PaperFacture{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера PaperFacture.
func (ht *PaperFacture) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: paperFactureListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
func (ht *PaperFacture) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *PaperFacture) listParams(_ *http.Request) entity.PaperFactureParams {
	return entity.PaperFactureParams{}
}
