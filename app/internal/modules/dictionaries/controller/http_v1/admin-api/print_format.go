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
	printFormatURL             = "/v1/dictionaries/print-formats"
	printFormatItemURL         = "/v1/dictionaries/print-formats/:id"
	printFormatChangeStatusURL = "/v1/dictionaries/print-formats/:id/status"
)

type (
	PrintFormat struct {
		section    mrcore.ClientSection
		service    usecase.PrintFormatService
		listSorter mrview.ListSorter
	}
)

func NewPrintFormat(
	section mrcore.ClientSection,
	service usecase.PrintFormatService,
	listSorter mrview.ListSorter,
) *PrintFormat {
	return &PrintFormat{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *PrintFormat) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitPrintFormatPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(printFormatURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(printFormatURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(printFormatItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(printFormatItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(printFormatItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(printFormatChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *PrintFormat) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.PrintFormatListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *PrintFormat) listParams(c mrcore.ClientContext) entity.PrintFormatParams {
	return entity.PrintFormatParams{
		Filter: entity.PrintFormatListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Length:     view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterLengthRange),
			Width:      view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterWidthRange),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *PrintFormat) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *PrintFormat) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreatePrintFormatRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.PrintFormat{
			Caption: request.Caption,
			Length:  request.Length,
			Width:   request.Width,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgDictionariesPrintFormatSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *PrintFormat) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StorePrintFormatRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.PrintFormat{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Caption:    request.Caption,
			Length:     request.Length,
			Width:      request.Width,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *PrintFormat) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.PrintFormat{
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

func (ht *PrintFormat) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *PrintFormat) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *PrintFormat) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *PrintFormat) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return dictionaries.FactoryErrPrintFormatNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
