package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/catalog/laminate/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/public-api"
	usecase "print-shop-back/internal/modules/catalog/laminate/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	laminateListURL = "/v1/catalog/laminates"
)

type (
	Laminate struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.LaminateUseCase
	}
)

func NewLaminate(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.LaminateUseCase,
) *Laminate {
	return &Laminate{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *Laminate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, laminateListURL, "", ht.GetList},
	}
}

func (ht *Laminate) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *Laminate) listParams(r *http.Request) entity.LaminateParams {
	return entity.LaminateParams{}
}
