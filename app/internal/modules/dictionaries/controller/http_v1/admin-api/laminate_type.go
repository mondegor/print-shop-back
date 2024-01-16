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
		section    mrcore.ClientSection
		service    usecase.LaminateTypeService
		listSorter mrview.ListSorter
	}
)

func NewLaminateType(
	section mrcore.ClientSection,
	service usecase.LaminateTypeService,
	listSorter mrview.ListSorter,
) *LaminateType {
	return &LaminateType{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *LaminateType) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitLaminateTypePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(laminateTypeURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(laminateTypeURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(laminateTypeItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(laminateTypeItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(laminateTypeItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(laminateTypeChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *LaminateType) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.LaminateTypeListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *LaminateType) listParams(c mrcore.ClientContext) entity.LaminateTypeParams {
	return entity.LaminateTypeParams{
		Filter: entity.LaminateTypeListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *LaminateType) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *LaminateType) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateLaminateTypeRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.LaminateType{
			Caption: request.Caption,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgDictionariesLaminateTypeSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *LaminateType) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreLaminateTypeRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.LaminateType{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Caption:    request.Caption,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *LaminateType) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.LaminateType{
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

func (ht *LaminateType) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *LaminateType) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *LaminateType) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *LaminateType) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return dictionaries.FactoryErrLaminateTypeNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
