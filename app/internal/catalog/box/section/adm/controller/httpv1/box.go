package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/catalog/box/module"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
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
	boxListURL             = "/v1/catalog/boxes"
	boxItemURL             = "/v1/catalog/boxes/{id}"
	boxItemChangeStatusURL = "/v1/catalog/boxes/{id}/status"
)

type (
	// Box - comment struct.
	Box struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    adm.BoxUseCase
		listSorter mrview.ListSorter
	}
)

// NewBox - создаёт контроллер Box.
func NewBox(parser validate.RequestExtendParser, sender mrserver.ResponseSender, useCase adm.BoxUseCase, listSorter mrview.ListSorter) *Box {
	return &Box{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера Box.
func (ht *Box) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: boxListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: boxListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: boxItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: boxItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: boxItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: boxItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *Box) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		BoxListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Box) listParams(r *http.Request) entity.BoxParams {
	return entity.BoxParams{
		Filter: entity.BoxListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Length:     measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange).Transform(measure.OneThousandth)),    // mm -> m
			Width:      measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange).Transform(measure.OneThousandth)),     // mm -> m
			Height:     measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterHeightRange).Transform(measure.OneThousandth)),    // mm -> m
			Weight:     measure.RangeKilogram(ht.parser.FilterRangeInt64(r, module.ParamNameFilterWeightRange).Transform(measure.OneThousandth)), // g -> kg
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
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
	request := CreateBoxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
		Article: request.Article,
		Caption: request.Caption,
		Length:  measure.Meter(request.Length * measure.OneThousandth),
		Width:   measure.Meter(request.Width * measure.OneThousandth),
		Height:  measure.Meter(request.Height * measure.OneThousandth),
		Weight:  measure.Kilogram(request.Weight * measure.OneThousandth),
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
	request := StoreBoxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		Length:     measure.Meter(request.Length * measure.OneThousandth),
		Width:      measure.Meter(request.Width * measure.OneThousandth),
		Height:     measure.Meter(request.Height * measure.OneThousandth),
		Weight:     measure.Kilogram(request.Weight * measure.OneThousandth),
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *Box) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
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
func (ht *Box) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
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
		return module.ErrBoxNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if module.ErrBoxArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	return err
}
