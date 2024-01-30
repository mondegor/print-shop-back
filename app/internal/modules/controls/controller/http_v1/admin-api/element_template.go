package http_v1

import (
	"fmt"
	"net/http"
	module "print-shop-back/internal/modules/controls"
	"print-shop-back/internal/modules/controls/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	elementTemplateURL             = "/v1/controls/element-templates"
	elementTemplateItemURL         = "/v1/controls/element-templates/:id"
	elementTemplateChangeStatusURL = "/v1/controls/element-templates/:id/status"
)

type (
	ElementTemplate struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.ElementTemplateService
		listSorter mrview.ListSorter
	}
)

func NewElementTemplate(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.ElementTemplateService,
	listSorter mrview.ListSorter,
) *ElementTemplate {
	return &ElementTemplate{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *ElementTemplate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, elementTemplateURL, "", ht.GetList},
		{http.MethodPost, elementTemplateURL, "", ht.Create},

		{http.MethodGet, elementTemplateItemURL, "", ht.Get},
		{http.MethodPut, elementTemplateItemURL, "", ht.Store},
		{http.MethodDelete, elementTemplateItemURL, "", ht.Remove},

		{http.MethodPut, elementTemplateChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *ElementTemplate) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.ElementTemplateListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *ElementTemplate) listParams(r *http.Request) entity.ElementTemplateParams {
	return entity.ElementTemplateParams{
		Filter: entity.ElementTemplateListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Detailing:  ht.parser.FilterElementDetailingList(r, module.ParamNameFilterElementDetailing),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *ElementTemplate) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *ElementTemplate) Create(w http.ResponseWriter, r *http.Request) error {
	request := view.CreateElementTemplateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.ElementTemplate{
		ParamName: request.ParamName,
		Caption:   request.Caption,
		Type:      request.Type,
		Detailing: request.Detailing,
		Body:      request.Body,
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
				"msgControlsElementTemplateSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *ElementTemplate) Store(w http.ResponseWriter, r *http.Request) error {
	request := view.StoreElementTemplateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.ElementTemplate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Type:       request.Type,
		Detailing:  request.Detailing,
		Body:       request.Body,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *ElementTemplate) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.ElementTemplate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *ElementTemplate) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *ElementTemplate) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *ElementTemplate) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *ElementTemplate) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrElementTemplateNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
