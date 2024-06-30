package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/usecase"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	paperFactureListURL             = "/v1/dictionaries/paper-factures"
	paperFactureItemURL             = "/v1/dictionaries/paper-factures/{id}"
	paperFactureItemChangeStatusURL = "/v1/dictionaries/paper-factures/{id}/status"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    usecase.PaperFactureUseCase
		listSorter mrview.ListSorter
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperFactureUseCase,
	listSorter mrview.ListSorter,
) *PaperFacture {
	return &PaperFacture{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - comment method.
func (ht *PaperFacture) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: paperFactureListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: paperFactureListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: paperFactureItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: paperFactureItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: paperFactureItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: paperFactureItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *PaperFacture) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PaperFactureListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *PaperFacture) listParams(r *http.Request) entity.PaperFactureParams {
	return entity.PaperFactureParams{
		Filter: entity.PaperFactureListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *PaperFacture) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *PaperFacture) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePaperFactureRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperFacture{
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
func (ht *PaperFacture) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePaperFactureRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperFacture{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *PaperFacture) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperFacture{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// Remove - comment method.
func (ht *PaperFacture) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperFacture) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *PaperFacture) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *PaperFacture) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrPaperFactureNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
