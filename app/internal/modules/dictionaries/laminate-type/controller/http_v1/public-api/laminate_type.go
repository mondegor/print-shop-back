package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/dictionaries/laminate-type/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/public-api"
	usecase "print-shop-back/internal/modules/dictionaries/laminate-type/usecase/public-api"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	laminateTypeListURL = "/v1/dictionaries/laminate-types"
)

type (
	LaminateType struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.LaminateTypeUseCase
	}
)

func NewLaminateType(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.LaminateTypeUseCase,
) *LaminateType {
	return &LaminateType{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *LaminateType) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, laminateTypeListURL, "", ht.GetList},
	}
}

func (ht *LaminateType) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *LaminateType) listParams(r *http.Request) entity.LaminateTypeParams {
	return entity.LaminateTypeParams{}
}
