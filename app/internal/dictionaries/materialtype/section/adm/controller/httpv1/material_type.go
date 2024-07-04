package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/usecase"
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
	materialTypeListURL             = "/v1/dictionaries/material-types"
	materialTypeItemURL             = "/v1/dictionaries/material-types/{id}"
	materialTypeItemChangeStatusURL = "/v1/dictionaries/material-types/{id}/status"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    usecase.MaterialTypeUseCase
		listSorter mrview.ListSorter
	}
)

// NewMaterialType - создаёт контроллер MaterialType.
func NewMaterialType(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	useCase usecase.MaterialTypeUseCase,
	listSorter mrview.ListSorter,
) *MaterialType {
	return &MaterialType{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера MaterialType.
func (ht *MaterialType) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: materialTypeListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: materialTypeListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: materialTypeItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: materialTypeItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: materialTypeItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: materialTypeItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *MaterialType) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		MaterialTypeListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *MaterialType) listParams(r *http.Request) entity.MaterialTypeParams {
	return entity.MaterialTypeParams{
		Filter: entity.MaterialTypeListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *MaterialType) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *MaterialType) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateMaterialTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.MaterialType{
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
func (ht *MaterialType) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreMaterialTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.MaterialType{
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
func (ht *MaterialType) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.MaterialType{
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
func (ht *MaterialType) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *MaterialType) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *MaterialType) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *MaterialType) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrMaterialTypeNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
