package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	paperFactureListURL = "/v1/dictionaries/paper-factures"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.PaperFactureUseCase
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.PaperFactureUseCase) *PaperFacture {
	return &PaperFacture{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - comment method.
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
