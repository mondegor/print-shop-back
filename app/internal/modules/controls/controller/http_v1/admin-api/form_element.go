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

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	formElementURL     = "/v1/controls/form-elements"
	formElementItemURL = "/v1/controls/form-elements/:id"
	formElementMoveURL = "/v1/controls/form-elements/:id/move"
)

type (
	FormElement struct {
		section    mrcore.ClientSection
		service    usecase.FormElementService
		listSorter mrview.ListSorter
	}
)

func NewFormElement(
	section mrcore.ClientSection,
	service usecase.FormElementService,
	listSorter mrview.ListSorter,
) *FormElement {
	return &FormElement{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *FormElement) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitFormElementPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(formElementURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(formElementURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(formElementItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(formElementItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(formElementItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPatch, ht.section.Path(formElementMoveURL), moduleAccessFunc(ht.Move()))
}

func (ht *FormElement) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		formID, err := view_shared.ParseFormDataID(c)

		if err != nil {
			return err
		}

		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c, formID))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.FormElementListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *FormElement) listParams(c mrcore.ClientContext, formID mrtype.KeyInt32) entity.FormElementParams {
	return entity.FormElementParams{
		FormID: formID,
		Filter: entity.FormElementListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Detailing:  view_shared.ParseFilterElementDetailingList(c, module.ParamNameFilterElementDetailing),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *FormElement) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *FormElement) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateFormElementRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.FormElement{
			FormID:     request.FormID,
			TemplateID: request.TemplateID,
			ParamName:  request.ParamName,
			Caption:    request.Caption,
			Required:   request.Required,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgControlsFormElementSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *FormElement) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreFormElementRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.FormElement{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			ParamName:  request.ParamName,
			Caption:    request.Caption,
			Required:   request.Required,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *FormElement) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *FormElement) Move() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.MoveFormElementRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		if err := ht.service.MoveAfterID(c.Context(), ht.getItemID(c), request.AfterNodeID); err != nil {
			return ht.wrapErrorNode(c, err)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *FormElement) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *FormElement) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *FormElement) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if usecase_shared.FactoryErrFormDataNotFound.Is(err) {
		return mrerr.NewFieldError("formId", err)
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if usecase_shared.FactoryErrFormElementParamNameAlreadyExists.Is(err) {
		return mrerr.NewFieldError("paramName", err)
	}

	if usecase_shared.FactoryErrElementTemplateNotFound.Is(err) {
		return mrerr.NewFieldError("templateId", err)
	}

	return err
}

func (ht *FormElement) wrapErrorNode(c mrcore.ClientContext, err error) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrorderer.FactoryErrAfterNodeNotFound.Is(err) {
		return mrerr.NewFieldError("afterNodeId", err)
	}

	return err
}
