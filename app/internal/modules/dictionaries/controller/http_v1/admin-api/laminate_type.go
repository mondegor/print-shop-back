package http_v1

import (
	"fmt"
	"net/http"
	module "print-shop-back/internal/modules/dictionaries"
	"print-shop-back/internal/modules/dictionaries/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/dictionaries/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/usecase/admin-api"
	"print-shop-back/pkg/modules/dictionaries"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	laminateTypeURL             = "/v1/dictionaries/laminate-types"
	laminateTypeItemURL         = "/v1/dictionaries/laminate-types/:id"
	laminateTypeChangeStatusURL = "/v1/dictionaries/laminate-types/:id/status"
)

type (
	LaminateType struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.LaminateTypeService
		listSorter mrview.ListSorter
	}
)

func NewLaminateType(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.LaminateTypeService,
	listSorter mrview.ListSorter,
) *LaminateType {
	return &LaminateType{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *LaminateType) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, laminateTypeURL, "", ht.GetList},
		{http.MethodPost, laminateTypeURL, "", ht.Create},

		{http.MethodGet, laminateTypeItemURL, "", ht.Get},
		{http.MethodPut, laminateTypeItemURL, "", ht.Store},
		{http.MethodDelete, laminateTypeItemURL, "", ht.Remove},

		{http.MethodPut, laminateTypeChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *LaminateType) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.LaminateTypeListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *LaminateType) listParams(r *http.Request) entity.LaminateTypeParams {
	return entity.LaminateTypeParams{
		Filter: entity.LaminateTypeListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *LaminateType) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *LaminateType) Create(w http.ResponseWriter, r *http.Request) error {
	request := view.CreateLaminateTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
		Caption: request.Caption,
	}

	if err := ht.service.Create(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemResponse{
			ItemID: fmt.Sprintf("%d", item.ID),
			Message: mrctx.Locale(r.Context()).TranslateMessage(
				"msgDictionariesLaminateTypeSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *LaminateType) Store(w http.ResponseWriter, r *http.Request) error {
	request := view.StoreLaminateTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
		ID:         ht.getItemID(r),
		TagVersion: request.Version,
		Caption:    request.Caption,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *LaminateType) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *LaminateType) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return dictionaries.FactoryErrLaminateTypeNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
