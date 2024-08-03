package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
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
	laminateListURL             = "/v1/catalog/laminates"
	laminateItemURL             = "/v1/catalog/laminates/{id}"
	laminateItemChangeStatusURL = "/v1/catalog/laminates/{id}/status"
)

type (
	// Laminate - comment struct.
	Laminate struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    adm.LaminateUseCase
		listSorter mrview.ListSorter
	}
)

// NewLaminate - создаёт контроллер Laminate.
func NewLaminate(parser validate.RequestExtendParser, sender mrserver.ResponseSender, useCase adm.LaminateUseCase, listSorter mrview.ListSorter) *Laminate {
	return &Laminate{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера Laminate.
func (ht *Laminate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: laminateListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: laminateListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: laminateItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: laminateItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: laminateItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: laminateItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *Laminate) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		LaminateListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Laminate) listParams(r *http.Request) entity.LaminateParams {
	return entity.LaminateParams{
		Filter: entity.LaminateListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			TypeIDs:    ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogMaterialTypeIDs),
			Length:     measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange).Transform(measure.OneThousandth)), // mm -> m
			Width:      measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange).Transform(measure.OneThousandth)),  // mm -> m
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *Laminate) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Laminate) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateLaminateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
		Article:   request.Article,
		Caption:   request.Caption,
		TypeID:    request.TypeID,
		Length:    request.Length,
		Width:     measure.Meter(request.Width * measure.OneThousandth),
		Thickness: measure.Meter(request.Thickness * measure.OneMillionth),
		WeightM2:  measure.KilogramPerMeter2(request.WeightM2 * measure.OneThousandth),
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
func (ht *Laminate) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreLaminateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		TypeID:     request.TypeID,
		Length:     request.Length,
		Width:      measure.Meter(request.Width * measure.OneThousandth),
		Thickness:  measure.Meter(request.Thickness * measure.OneMillionth),
		WeightM2:   measure.KilogramPerMeter2(request.WeightM2 * measure.OneThousandth),
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *Laminate) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
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
func (ht *Laminate) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Laminate) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Laminate) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Laminate) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrLaminateNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if module.ErrLaminateArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if api.ErrMaterialTypeRequired.Is(err) ||
		api.ErrMaterialTypeNotFound.Is(err) {
		return mrerr.NewCustomError("typeId", err)
	}

	return err
}
