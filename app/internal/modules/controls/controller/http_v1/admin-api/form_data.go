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
	formDataURL             = "/v1/controls/forms"
	formDataItemURL         = "/v1/controls/forms/:id"
	formDataChangeStatusURL = "/v1/controls/forms/:id/status"
	formDataCompileURL      = "/v1/controls/forms/:id/compile"
)

type (
	FormData struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		service usecase.FormDataService
		// serviceUIFormData usecase.UIFormDataService
		listSorter mrview.ListSorter
	}
)

func NewFormData(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.FormDataService,
	// serviceUIFormData usecase.UIFormDataService,
	listSorter mrview.ListSorter,
) *FormData {
	return &FormData{
		parser:  parser,
		sender:  sender,
		service: service,
		// serviceUIFormData: serviceUIFormData,
		listSorter: listSorter,
	}
}

func (ht *FormData) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, formDataURL, "", ht.GetList},
		{http.MethodPost, formDataURL, "", ht.Create},

		{http.MethodGet, formDataItemURL, "", ht.Get},
		{http.MethodPut, formDataItemURL, "", ht.Store},
		{http.MethodDelete, formDataItemURL, "", ht.Remove},

		{http.MethodPut, formDataChangeStatusURL, "", ht.ChangeStatus},

		{http.MethodPatch, formDataCompileURL, "", ht.Compile},
	}
}

func (ht *FormData) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.FormDataListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *FormData) listParams(r *http.Request) entity.FormDataParams {
	return entity.FormDataParams{
		Filter: entity.FormDataListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Detailing:  ht.parser.FilterElementDetailingList(r, module.ParamNameFilterElementDetailing),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *FormData) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *FormData) Create(w http.ResponseWriter, r *http.Request) error {
	request := view.CreateFormDataRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormData{
		ParamName: request.ParamName,
		Caption:   request.Caption,
		Detailing: request.Detailing,
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
				"msgControlsFormDataSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *FormData) Store(w http.ResponseWriter, r *http.Request) error {
	request := view.StoreFormDataRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormData{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Detailing:  request.Detailing,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormData) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormData{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormData) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormData) Compile(w http.ResponseWriter, r *http.Request) error {
	// item, err := ht.serviceUIFormData.CompileForm(r.Context(), ht.getItemID(r))
	err := fmt.Errorf("of")
	item := ""

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *FormData) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *FormData) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *FormData) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormDataNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
