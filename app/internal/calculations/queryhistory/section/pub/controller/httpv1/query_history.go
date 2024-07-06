package httpv1

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/module"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/entity"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/usecase"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	queryHistoryURL     = "/v1/calculations/query-history"
	queryHistoryItemURL = "/v1/calculations/query-history/{id}"
)

type (
	// QueryHistory - comment struct.
	QueryHistory struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.QueryResultUseCase
	}
)

// NewQueryHistory - создаёт контроллер QueryHistory.
func NewQueryHistory(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.QueryResultUseCase) *QueryHistory {
	return &QueryHistory{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера QueryHistory.
func (ht *QueryHistory) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: queryHistoryItemURL, Func: ht.Get},
		{Method: http.MethodPost, URL: queryHistoryURL, Func: ht.Create},
	}
}

// Get - comment method.
func (ht *QueryHistory) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *QueryHistory) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateQueryHistoryRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.QueryHistoryItem{
		Caption: request.Caption,
		Params:  request.Params,
		Result:  request.Result,
	}

	itemID, err := ht.useCase.Create(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemResponse{
			ItemID: itemID.String(),
		},
	)
}

func (ht *QueryHistory) getItemID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}

func (ht *QueryHistory) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *QueryHistory) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrQueryHistoryNotFound.Wrap(err, ht.getRawItemID(r))
	}

	return err
}
