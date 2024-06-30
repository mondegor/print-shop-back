package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/calculations/box/module"

	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/entity"
	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	calcResultListURL = "/v1/calculations"
	calcResultItemURL = "/v1/calculations/{id}"
)

type (
	// Box - comment struct.
	Box struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.CalcResultUseCase
	}
)

// NewBox - создаёт объект Box.
func NewBox(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.CalcResultUseCase) *Box {
	return &Box{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - comment method.
func (ht *Box) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: calcResultListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: calcResultItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: calcResultItemURL, Func: ht.Store},
	}
}

// Get - comment method.
func (ht *Box) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Box) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateCalcResultRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CalcResult{
		Caption: request.Caption,
	}

	itemID, err := ht.useCase.Create(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemInt32Response{
			ItemID: itemID,
		},
	)
}

// Store - comment method.
func (ht *Box) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreCalcResultRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.CalcResult{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Box) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Box) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Box) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrCalcResultNotFound.Wrap(err, ht.getRawItemID(r))
	}

	return err
}
