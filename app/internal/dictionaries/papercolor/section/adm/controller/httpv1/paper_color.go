package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/usecase"
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
	paperColorListURL             = "/v1/dictionaries/paper-colors"
	paperColorItemURL             = "/v1/dictionaries/paper-colors/{id}"
	paperColorItemChangeStatusURL = "/v1/dictionaries/paper-colors/{id}/status"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    usecase.PaperColorUseCase
		listSorter mrview.ListSorter
	}
)

// NewPaperColor - создаёт контроллер PaperColor.
func NewPaperColor(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperColorUseCase,
	listSorter mrview.ListSorter,
) *PaperColor {
	return &PaperColor{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера PaperColor.
func (ht *PaperColor) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: paperColorListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: paperColorListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: paperColorItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: paperColorItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: paperColorItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: paperColorItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *PaperColor) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PaperColorListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *PaperColor) listParams(r *http.Request) entity.PaperColorParams {
	return entity.PaperColorParams{
		Filter: entity.PaperColorListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *PaperColor) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *PaperColor) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
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
func (ht *PaperColor) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
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
func (ht *PaperColor) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
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
func (ht *PaperColor) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperColor) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *PaperColor) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *PaperColor) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrPaperColorNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
