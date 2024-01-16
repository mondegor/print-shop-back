package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/catalog"
	"print-shop-back/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "print-shop-back/internal/modules/catalog/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"
	usecase "print-shop-back/internal/modules/catalog/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/usecase/shared"
	"print-shop-back/pkg/modules/dictionaries"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	laminateURL             = "/v1/catalog/laminates"
	laminateItemURL         = "/v1/catalog/laminates/:id"
	laminateChangeStatusURL = "/v1/catalog/laminates/:id/status"
)

type (
	Laminate struct {
		section    mrcore.ClientSection
		service    usecase.LaminateService
		listSorter mrview.ListSorter
	}
)

func NewLaminate(
	section mrcore.ClientSection,
	service usecase.LaminateService,
	listSorter mrview.ListSorter,
) *Laminate {
	return &Laminate{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Laminate) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitLaminatePermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(laminateURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(laminateURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(laminateItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(laminateItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(laminateItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(laminateChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *Laminate) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.LaminateListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Laminate) listParams(c mrcore.ClientContext) entity.LaminateParams {
	return entity.LaminateParams{
		Filter: entity.LaminateListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			TypeIDs:    view_shared.ParseFilterKeyInt32List(c, module.ParamNameFilterCatalogLaminateTypeIDs),
			Length:     view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterLengthRange),
			Weight:     view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterWeightRange),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Laminate) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Laminate) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateLaminateRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Laminate{
			Article:   request.Article,
			Caption:   request.Caption,
			TypeID:    request.TypeID,
			Length:    request.Length,
			Weight:    request.Weight,
			Thickness: request.Thickness,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(item.ID)),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgCatalogLaminateSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Laminate) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreLaminateRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Laminate{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Article:    request.Article,
			Caption:    request.Caption,
			TypeID:     request.TypeID,
			Length:     request.Length,
			Weight:     request.Weight,
			Thickness:  request.Thickness,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Laminate) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Laminate{
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

func (ht *Laminate) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Laminate) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *Laminate) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Laminate) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrLaminateNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	if usecase_shared.FactoryErrLaminateArticleAlreadyExists.Is(err) {
		return mrerr.NewFieldError("article", err)
	}

	if dictionaries.FactoryErrLaminateTypeNotFound.Is(err) {
		return mrerr.NewFieldError("typeId", err)
	}

	return err
}
