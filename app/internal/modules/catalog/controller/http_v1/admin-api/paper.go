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
	paperURL             = "/v1/catalog/papers"
	paperItemURL         = "/v1/catalog/papers/:id"
	paperChangeStatusURL = "/v1/catalog/papers/:id/status"
)

type (
	Paper struct {
		section    mrcore.ClientSection
		service    usecase.PaperService
		listSorter mrview.ListSorter
	}
)

func NewPaper(
	section mrcore.ClientSection,
	service usecase.PaperService,
	listSorter mrview.ListSorter,
) *Paper {
	return &Paper{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Paper) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitPaperPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(paperURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(paperURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(paperItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(paperItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(paperItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(paperChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *Paper) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.PaperListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Paper) listParams(c mrcore.ClientContext) entity.PaperParams {
	return entity.PaperParams{
		Filter: entity.PaperListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			ColorIDs:   view_shared.ParseFilterKeyInt32List(c, module.ParamNameFilterCatalogPaperColorIDs),
			FactureIDs: view_shared.ParseFilterKeyInt32List(c, module.ParamNameFilterCatalogPaperFactureIDs),
			Length:     view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterLengthRange),
			Width:      view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterWidthRange),
			Density:    view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterDensityRange),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Paper) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Paper) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreatePaperRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Paper{
			Article:   request.Article,
			Caption:   request.Caption,
			ColorID:   request.ColorID,
			FactureID: request.FactureID,
			Length:    request.Length,
			Width:     request.Width,
			Density:   request.Density,
			Thickness: request.Thickness,
			Sides:     request.Sides,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(item.ID)),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgCatalogPaperSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Paper) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StorePaperRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Paper{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Article:    request.Article,
			Caption:    request.Caption,
			ColorID:    request.ColorID,
			FactureID:  request.FactureID,
			Length:     request.Length,
			Width:      request.Width,
			Density:    request.Density,
			Thickness:  request.Thickness,
			Sides:      request.Sides,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Paper) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Paper{
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

func (ht *Paper) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Paper) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *Paper) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Paper) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrPaperNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	if usecase_shared.FactoryErrPaperArticleAlreadyExists.Is(err) {
		return mrerr.NewFieldError("article", err)
	}

	if dictionaries.FactoryErrPaperColorNotFound.Is(err) {
		return mrerr.NewFieldError("colorId", err)
	}

	if dictionaries.FactoryErrPaperFactureNotFound.Is(err) {
		return mrerr.NewFieldError("factureId", err)
	}

	return err
}
