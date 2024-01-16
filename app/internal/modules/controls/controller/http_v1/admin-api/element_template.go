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
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
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
		section    mrcore.ClientSection
		service    usecase.ElementTemplateService
		listSorter mrview.ListSorter
	}
)

func NewElementTemplate(
	section mrcore.ClientSection,
	service usecase.ElementTemplateService,
	listSorter mrview.ListSorter,
) *ElementTemplate {
	return &ElementTemplate{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *ElementTemplate) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitElementTemplatePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(elementTemplateURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(elementTemplateURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(elementTemplateItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(elementTemplateItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(elementTemplateItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(elementTemplateChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *ElementTemplate) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.ElementTemplateListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *ElementTemplate) listParams(c mrcore.ClientContext) entity.ElementTemplateParams {
	return entity.ElementTemplateParams{
		Filter: entity.ElementTemplateListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Detailing:  view_shared.ParseFilterElementDetailingList(c, module.ParamNameFilterElementDetailing),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *ElementTemplate) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *ElementTemplate) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateElementTemplateRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.ElementTemplate{
			ParamName: request.ParamName,
			Caption:   request.Caption,
			Type:      request.Type,
			Detailing: request.Detailing,
			Body:      request.Body,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgControlsElementTemplateSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *ElementTemplate) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreElementTemplateRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.ElementTemplate{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			ParamName:  request.ParamName,
			Caption:    request.Caption,
			Type:       request.Type,
			Detailing:  request.Detailing,
			Body:       request.Body,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *ElementTemplate) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.ElementTemplate{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			Status:     request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *ElementTemplate) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *ElementTemplate) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *ElementTemplate) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *ElementTemplate) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrElementTemplateNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
