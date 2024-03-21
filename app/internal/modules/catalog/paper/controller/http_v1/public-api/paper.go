package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/catalog/paper/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/paper/entity/public-api"
	usecase "print-shop-back/internal/modules/catalog/paper/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	paperListURL = "/v1/catalog/papers"
)

type (
	Paper struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.PaperUseCase
	}
)

func NewPaper(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperUseCase,
) *Paper {
	return &Paper{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *Paper) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperListURL, "", ht.GetList},
	}
}

func (ht *Paper) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *Paper) listParams(r *http.Request) entity.PaperParams {
	return entity.PaperParams{}
}
