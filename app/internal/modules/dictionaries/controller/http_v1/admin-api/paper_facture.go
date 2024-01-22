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
	paperFactureURL             = "/v1/dictionaries/paper-factures"
	paperFactureItemURL         = "/v1/dictionaries/paper-factures/:id"
	paperFactureChangeStatusURL = "/v1/dictionaries/paper-factures/:id/status"
)

type (
	PaperFacture struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.PaperFactureService
		listSorter mrview.ListSorter
	}
)

func NewPaperFacture(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.PaperFactureService,
	listSorter mrview.ListSorter,
) *PaperFacture {
	return &PaperFacture{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *PaperFacture) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperFactureURL, "", ht.GetList},
		{http.MethodPost, paperFactureURL, "", ht.Create},

		{http.MethodGet, paperFactureItemURL, "", ht.Get},
		{http.MethodPut, paperFactureItemURL, "", ht.Store},
		{http.MethodDelete, paperFactureItemURL, "", ht.Remove},

		{http.MethodPut, paperFactureChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *PaperFacture) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.PaperFactureListResponse{
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

func (ht *PaperFacture) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *PaperFacture) Create(w http.ResponseWriter, r *http.Request) error {
	request := view.CreatePaperFactureRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperFacture{
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
				"msgDictionariesPaperFactureSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *PaperFacture) Store(w http.ResponseWriter, r *http.Request) error {
	request := view.StorePaperFactureRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperFacture{
		ID:         ht.getItemID(r),
		TagVersion: request.Version,
		Caption:    request.Caption,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

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

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperFacture) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
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
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return dictionaries.FactoryErrPaperFactureNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
