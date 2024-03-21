package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/dictionaries/paper-facture/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/public-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-facture/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	paperFactureListURL = "/v1/dictionaries/paper-factures"
)

type (
	PaperFacture struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.PaperFactureUseCase
	}
)

func NewPaperFacture(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperFactureUseCase,
) *PaperFacture {
	return &PaperFacture{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *PaperFacture) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperFactureListURL, "", ht.GetList},
	}
}

func (ht *PaperFacture) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *PaperFacture) listParams(r *http.Request) entity.PaperFactureParams {
	return entity.PaperFactureParams{}
}
