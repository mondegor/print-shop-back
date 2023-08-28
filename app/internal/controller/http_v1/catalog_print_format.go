package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
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

func (ht *CatalogPrintFormat) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogPrintFormatListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogPrintFormatItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogPrintFormat) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPrintFormat) newListFilter(c mrapp.ClientData) *entity.CatalogPrintFormatListFilter {
    var listFilter entity.CatalogPrintFormatListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPrintFormat) Get() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPrintFormat) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogPrintFormat{}

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

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogPrintFormatSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogPrintFormat) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPrintFormat{}

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

func (ht *CatalogPrintFormat) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

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

func (ht *CatalogPrintFormat) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPrintFormat) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
