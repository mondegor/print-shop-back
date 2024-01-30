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
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	paperColorURL             = "/v1/dictionaries/paper-colors"
	paperColorItemURL         = "/v1/dictionaries/paper-colors/:id"
	paperColorChangeStatusURL = "/v1/dictionaries/paper-colors/:id/status"
)

type (
	PaperColor struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.PaperColorService
		listSorter mrview.ListSorter
	}
)

func NewPaperColor(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.PaperColorService,
	listSorter mrview.ListSorter,
) *PaperColor {
	return &PaperColor{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *PaperColor) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperColorURL, "", ht.GetList},
		{http.MethodPost, paperColorURL, "", ht.Create},

		{http.MethodGet, paperColorItemURL, "", ht.Get},
		{http.MethodPut, paperColorItemURL, "", ht.Store},
		{http.MethodDelete, paperColorItemURL, "", ht.Remove},

		{http.MethodPut, paperColorChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *PaperColor) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.PaperColorListResponse{
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

func (ht *PaperColor) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *PaperColor) Create(w http.ResponseWriter, r *http.Request) error {
	request := view.CreatePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
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
			Message: mrlang.Ctx(r.Context()).TranslateMessage(
				"msgDictionariesPaperColorSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *PaperColor) Store(w http.ResponseWriter, r *http.Request) error {
	request := view.StorePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
		ID:         ht.getItemID(r),
		TagVersion: request.Version,
		Caption:    request.Caption,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

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

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperColor) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
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
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return dictionaries.FactoryErrPaperColorNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
