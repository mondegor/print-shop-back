package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogPaperListURL = "/v1/catalog-papers"
    catalogPaperItemURL = "/v1/catalog-papers/:id"
    catalogPaperChangeStatusURL = "/v1/catalog-papers/:id/status"
)

type CatalogPaper struct {
    service usecase.CatalogPaperService
}

func NewCatalogPaper(service usecase.CatalogPaperService) *CatalogPaper {
    return &CatalogPaper{
        service: service,
    }
}

func (ht *CatalogPaper) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogPaperItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogPaperChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogPaper) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPaper) newListFilter(c mrcore.ClientData) *entity.CatalogPaperListFilter {
    var listFilter entity.CatalogPaperListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPaper) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPaper) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogPaper{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Density: request.Density,
            ColorId: request.ColorId,
            FactureId: request.FactureId,
            Thickness: request.Thickness,
            Sides: request.Sides,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogPaperSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogPaper) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogPaper{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Id:        ht.getItemId(c),
            Version:   request.Version,
            Article:   request.Article,
            Caption:   request.Caption,
            Length:    request.Length,
            Width:     request.Width,
            Density:   request.Density,
            ColorId:   request.ColorId,
            FactureId: request.FactureId,
            Thickness: request.Thickness,
            Sides:     request.Sides,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaper) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaper{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Status:  request.Status,
        }

        err := ht.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaper) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaper) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
