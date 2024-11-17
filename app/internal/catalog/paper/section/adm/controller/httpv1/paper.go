package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/dictionaries/api"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/validate"
	"github.com/mondegor/print-shop-back/pkg/view"
)

const (
	paperListURL             = "/v1/catalog/papers"
	paperItemURL             = "/v1/catalog/papers/{id}"
	paperItemChangeStatusURL = "/v1/catalog/papers/{id}/status"
)

type (
	// Paper - comment struct.
	Paper struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    adm.PaperUseCase
		listSorter mrview.ListSorter
	}
)

// NewPaper - создаёт контроллер Paper.
func NewPaper(parser validate.RequestExtendParser, sender mrserver.ResponseSender, useCase adm.PaperUseCase, listSorter mrview.ListSorter) *Paper {
	return &Paper{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера Paper.
func (ht *Paper) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: paperListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: paperListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: paperItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: paperItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: paperItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: paperItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *Paper) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PaperListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Paper) listParams(r *http.Request) entity.PaperParams {
	return entity.PaperParams{
		Filter: entity.PaperListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			TypeIDs:    ht.parser.FilterUint64List(r, module.ParamNameFilterCatalogPaperTypeIDs),
			ColorIDs:   ht.parser.FilterUint64List(r, module.ParamNameFilterCatalogPaperColorIDs),
			FactureIDs: ht.parser.FilterUint64List(r, module.ParamNameFilterCatalogPaperFactureIDs),
			Width:      measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange).Transform(measure.OneThousandth)), // mm -> m
			Height:     measure.RangeMeter(ht.parser.FilterRangeInt64(r, module.ParamNameFilterHeightRange).Transform(measure.OneThousandth)), // mm -> m
			Density: measure.RangeKilogramPerMeter2(
				ht.parser.FilterRangeInt64(r, module.ParamNameFilterDensityRange).Transform(measure.OneThousandth), // g/m2 -> kg/m2
			),
			Statuses: ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *Paper) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Paper) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePaperRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
		Article:   request.Article,
		Caption:   request.Caption,
		TypeID:    request.TypeID,
		ColorID:   request.ColorID,
		FactureID: request.FactureID,
		Width:     measure.Meter(request.Length * measure.OneThousandth),
		Height:    measure.Meter(request.Height * measure.OneThousandth),
		Thickness: measure.Meter(request.Thickness * measure.OneMillionth),
		Density:   measure.KilogramPerMeter2(request.Density * measure.OneThousandth),
		Sides:     request.Sides,
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
func (ht *Paper) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePaperRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		TypeID:     request.TypeID,
		ColorID:    request.ColorID,
		FactureID:  request.FactureID,
		Width:      measure.Meter(request.Length * measure.OneThousandth),
		Height:     measure.Meter(request.Height * measure.OneThousandth),
		Thickness:  measure.Meter(request.Thickness * measure.OneMillionth),
		Density:    measure.KilogramPerMeter2(request.Density * measure.OneThousandth),
		Sides:      request.Sides,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *Paper) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
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
func (ht *Paper) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Paper) getItemID(r *http.Request) uint64 {
	return ht.parser.PathParamUint64(r, "id")
}

func (ht *Paper) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Paper) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrPaperNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if module.ErrPaperArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if api.ErrMaterialTypeRequired.Is(err) ||
		api.ErrMaterialTypeNotFound.Is(err) {
		return mrerr.NewCustomError("typeId", err)
	}

	if api.ErrPaperColorRequired.Is(err) ||
		api.ErrPaperColorNotFound.Is(err) {
		return mrerr.NewCustomError("colorId", err)
	}

	if api.ErrPaperFactureRequired.Is(err) ||
		api.ErrPaperFactureNotFound.Is(err) {
		return mrerr.NewCustomError("factureId", err)
	}

	return err
}
