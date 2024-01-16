package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/catalog"
	"print-shop-back/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/catalog/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	usecase "print-shop-back/internal/modules/catalog/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/usecase/shared"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	boxURL             = "/v1/catalog/boxes"
	boxItemURL         = "/v1/catalog/boxes/:id"
	boxChangeStatusURL = "/v1/catalog/boxes/:id/status"
)

type (
	Box struct {
		section    mrcore.ClientSection
		service    usecase.BoxService
		listSorter mrview.ListSorter
	}
)

func NewBox(
	section mrcore.ClientSection,
	service usecase.BoxService,
	listSorter mrview.ListSorter,
) *Box {
	return &Box{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Box) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitBoxPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(boxURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(boxURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(boxItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(boxItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(boxItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(boxChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *Box) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.BoxListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Box) listParams(c mrcore.ClientContext) entity.BoxParams {
	return entity.BoxParams{
		Filter: entity.BoxListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Length:     view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterLengthRange),
			Width:      view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterWidthRange),
			Depth:      view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterDepthRange),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Box) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Box) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateBoxRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Box{
			Article: request.Article,
			Caption: request.Caption,
			Length:  request.Length,
			Width:   request.Width,
			Depth:   request.Depth,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(item.ID)),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgCatalogBoxSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Box) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreBoxRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Box{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Article:    request.Article,
			Caption:    request.Caption,
			Length:     request.Length,
			Width:      request.Width,
			Depth:      request.Depth,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Box) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Box{
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

func (ht *Box) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Box) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *Box) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Box) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrBoxNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	if usecase_shared.FactoryErrBoxArticleAlreadyExists.Is(err) {
		return mrerr.NewFieldError("article", err)
	}

	return err
}
