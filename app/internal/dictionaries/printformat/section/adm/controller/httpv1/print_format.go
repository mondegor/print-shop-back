package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	printFormatListURL             = "/v1/dictionaries/print-formats"
	printFormatItemURL             = "/v1/dictionaries/print-formats/{id}"
	printFormatItemChangeStatusURL = "/v1/dictionaries/print-formats/{id}/status"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    adm.PrintFormatUseCase
		listSorter mrview.ListSorter
	}
)

// NewPrintFormat - создаёт контроллер PrintFormat.
func NewPrintFormat(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	useCase adm.PrintFormatUseCase,
	listSorter mrview.ListSorter,
) *PrintFormat {
	return &PrintFormat{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера PrintFormat.
func (ht *PrintFormat) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: printFormatListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: printFormatListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: printFormatItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: printFormatItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: printFormatItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: printFormatItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *PrintFormat) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PrintFormatListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *PrintFormat) listParams(r *http.Request) entity.PrintFormatParams {
	return entity.PrintFormatParams{
		Filter: entity.PrintFormatListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Width:      measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange).Transform(measure.OneThousandth)),  // mm -> m
			Height:     measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterHeightRange).Transform(measure.OneThousandth)), // mm -> m
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *PrintFormat) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *PrintFormat) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePrintFormatRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
		Caption: request.Caption,
		Width:   measure.Meter(request.Width * measure.OneThousandth),
		Height:  measure.Meter(request.Height * measure.OneThousandth),
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
func (ht *PrintFormat) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePrintFormatRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
		Width:      measure.Meter(request.Width * measure.OneThousandth),
		Height:     measure.Meter(request.Height * measure.OneThousandth),
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *PrintFormat) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
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
func (ht *PrintFormat) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PrintFormat) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *PrintFormat) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *PrintFormat) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrPrintFormatNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
