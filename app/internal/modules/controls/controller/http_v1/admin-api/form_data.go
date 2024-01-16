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
	formDataURL             = "/v1/controls/forms"
	formDataItemURL         = "/v1/controls/forms/:id"
	formDataChangeStatusURL = "/v1/controls/forms/:id/status"
	formDataCompileURL      = "/v1/controls/forms/:id/compile"
)

type (
	FormData struct {
		section mrcore.ClientSection
		service usecase.FormDataService
		// serviceUIFormData usecase.UIFormDataService
		listSorter mrview.ListSorter
	}
)

func NewFormData(
	section mrcore.ClientSection,
	service usecase.FormDataService,
	// serviceUIFormData usecase.UIFormDataService,
	listSorter mrview.ListSorter,
) *FormData {
	return &FormData{
		section: section,
		service: service,
		// serviceUIFormData: serviceUIFormData,
		listSorter: listSorter,
	}
}

func (ht *FormData) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitFormDataPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(formDataURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(formDataURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(formDataItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(formDataItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(formDataItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(formDataChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))

	router.HttpHandlerFunc(http.MethodPatch, ht.section.Path(formDataCompileURL), moduleAccessFunc(ht.Compile()))
}

func (ht *FormData) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.FormDataListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *FormData) listParams(c mrcore.ClientContext) entity.FormDataParams {
	return entity.FormDataParams{
		Filter: entity.FormDataListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Detailing:  view_shared.ParseFilterElementDetailingList(c, module.ParamNameFilterElementDetailing),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *FormData) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *FormData) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateFormDataRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.FormData{
			ParamName: request.ParamName,
			Caption:   request.Caption,
			Detailing: request.Detailing,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgControlsFormDataSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *FormData) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreFormDataRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.FormData{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			ParamName:  request.ParamName,
			Caption:    request.Caption,
			Detailing:  request.Detailing,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *FormData) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.FormData{
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

func (ht *FormData) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *FormData) Compile() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		// item, err := ht.serviceUIFormData.CompileForm(c.Context(), ht.getItemID(c))
		err := fmt.Errorf("of")
		item := ""

		if err != nil {
			return err
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *FormData) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *FormData) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *FormData) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormDataNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
