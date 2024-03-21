package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/catalog/box/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/box/entity/public-api"
	usecase "print-shop-back/internal/modules/catalog/box/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	boxListURL = "/v1/catalog/boxes"
)

type (
	Box struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.BoxUseCase
	}
)

func NewBox(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.BoxUseCase,
) *Box {
	return &Box{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *Box) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, boxListURL, "", ht.GetList},
	}
}

func (ht *Box) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *Box) listParams(r *http.Request) entity.BoxParams {
	return entity.BoxParams{}
}
