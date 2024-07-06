package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/validate"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	materialTypeListURL = "/v1/dictionaries/material-types"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase pub.MaterialTypeUseCase
	}
)

// NewMaterialType - создаёт контроллер MaterialType.
func NewMaterialType(parser validate.RequestParser, sender mrserver.ResponseSender, useCase pub.MaterialTypeUseCase) *MaterialType {
	return &MaterialType{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера MaterialType.
func (ht *MaterialType) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: materialTypeListURL, Func: ht.GetList},
	}
}

// GetList - comment method.
func (ht *MaterialType) GetList(w http.ResponseWriter, r *http.Request) error {
	items, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, items)
}

func (ht *MaterialType) listParams(_ *http.Request) entity.MaterialTypeParams {
	return entity.MaterialTypeParams{}
}
