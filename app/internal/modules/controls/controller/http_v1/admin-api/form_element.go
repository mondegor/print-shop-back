package http_v1

import (
	"fmt"
	"net/http"
	module "print-shop-back/internal/modules/controls"
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	formElementListURL     = "/v1/controls/form-elements"
	formElementItemURL     = "/v1/controls/form-elements/:id"
	formElementItemMoveURL = "/v1/controls/form-elements/:id/move"
)

type (
	FormElement struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.FormElementService
		listSorter mrview.ListSorter
	}
)

func NewFormElement(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.FormElementService,
	listSorter mrview.ListSorter,
) *FormElement {
	return &FormElement{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *FormElement) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, formElementListURL, "", ht.GetList},
		{http.MethodPost, formElementListURL, "", ht.Create},

		{http.MethodGet, formElementItemURL, "", ht.Get},
		{http.MethodPut, formElementItemURL, "", ht.Store},
		{http.MethodDelete, formElementItemURL, "", ht.Remove},

		{http.MethodPatch, formElementItemMoveURL, "", ht.Move},
	}
}

func (ht *FormElement) GetList(w http.ResponseWriter, r *http.Request) error {
	formID, err := ht.parser.FormDataID(r)

	if err != nil {
		return err
	}

	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r, formID))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		FormElementListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *FormElement) listParams(r *http.Request, formID mrtype.KeyInt32) entity.FormElementParams {
	return entity.FormElementParams{
		FormID: formID,
		Filter: entity.FormElementListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Detailing:  ht.parser.FilterElementDetailingList(r, module.ParamNameFilterElementDetailing),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *FormElement) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *FormElement) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateFormElementRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormElement{
		FormID:     request.FormID,
		TemplateID: request.TemplateID,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Required:   request.Required,
	}

	if err := ht.service.Create(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessCreatedItemResponse{
			ItemID: fmt.Sprintf("%d", item.ID),
			Message: mrlang.Ctx(r.Context()).TranslateMessage(
				"msgControlsFormElementSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *FormElement) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreFormElementRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormElement{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Required:   request.Required,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) Move(w http.ResponseWriter, r *http.Request) error {
	request := MoveFormElementRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.service.MoveAfterID(r.Context(), ht.getItemID(r), request.AfterNodeID); err != nil {
		return ht.wrapErrorNode(r, err)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *FormElement) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *FormElement) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if usecase_shared.FactoryErrFormDataNotFound.Is(err) {
		return mrerr.NewCustomError("formId", err)
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if usecase_shared.FactoryErrFormElementParamNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("paramName", err)
	}

	if usecase_shared.FactoryErrElementTemplateNotFound.Is(err) {
		return mrerr.NewCustomError("templateId", err)
	}

	return err
}

func (ht *FormElement) wrapErrorNode(r *http.Request, err error) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrorderer.FactoryErrAfterNodeNotFound.Is(err) {
		return mrerr.NewCustomError("afterNodeId", err)
	}

	return err
}
