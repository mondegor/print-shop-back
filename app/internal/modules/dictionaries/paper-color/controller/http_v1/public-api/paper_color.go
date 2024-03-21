package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/dictionaries/paper-color/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/public-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-color/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	paperColorListURL = "/v1/dictionaries/paper-colors"
)

type (
	PaperColor struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.PaperColorUseCase
	}
)

func NewPaperColor(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperColorUseCase,
) *PaperColor {
	return &PaperColor{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *PaperColor) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperColorListURL, "", ht.GetList},
	}
}

func (ht *PaperColor) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *PaperColor) listParams(r *http.Request) entity.PaperColorParams {
	return entity.PaperColorParams{}
}
