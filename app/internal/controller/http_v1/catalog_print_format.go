package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

const (
    catalogPrintFormatListURL = "/v1/catalog-print-formats"
    catalogPrintFormatItemURL = "/v1/catalog-print-formats/:id"
    catalogPrintFormatChangeStatusURL = "/v1/catalog-print-formats/:id/status"
)

type CatalogPrintFormat struct {
    service usecase.CatalogPrintFormatService
}

func NewCatalogPrintFormat(service usecase.CatalogPrintFormatService) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        service: service,
    }
}

func (ht *CatalogPrintFormat) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogPrintFormatListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogPrintFormatItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogPrintFormat) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPrintFormat) newListFilter(c mrcore.ClientData) *entity.CatalogPrintFormatListFilter {
    var listFilter entity.CatalogPrintFormatListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPrintFormat) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPrintFormat) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogPrintFormat{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPrintFormat{
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := view.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogPrintFormatSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogPrintFormat) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogPrintFormat{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPrintFormat{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
            Length:  request.Length,
            Width:   request.Width,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPrintFormat) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPrintFormat{
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

func (ht *CatalogPrintFormat) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPrintFormat) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
