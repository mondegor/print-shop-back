package http_v1

import (
    "calc-user-data-back-adm/internal/controller/dto"
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/internal/usecase"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrentity"
    "fmt"
    "net/http"
)

const (
    catalogPrintFormatGetListURL = "/v1/catalog-print-formats"
    catalogPrintFormatGetItemURL = "/v1/catalog-print-formats/:id"
    catalogPrintFormatCreateURL = "/v1/catalog-print-formats"
    catalogPrintFormatStoreURL = "/v1/catalog-print-formats/:id"
    catalogPrintFormatChangeStatusURL = "/v1/catalog-print-formats/:id/status"
    catalogPrintFormatRemove = "/v1/catalog-print-formats/:id"
)

type CatalogPrintFormat struct {
    service usecase.CatalogPrintFormatService
}

func NewCatalogPrintFormat(service usecase.CatalogPrintFormatService) *CatalogPrintFormat {
    return &CatalogPrintFormat{
        service: service,
    }
}

func (f *CatalogPrintFormat) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogPrintFormatGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogPrintFormatCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogPrintFormatChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogPrintFormatRemove, f.Remove())
}

func (f *CatalogPrintFormat) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogPrintFormat) newListFilter(c mrapp.ClientData) *entity.CatalogPrintFormatListFilter {
    var listFilter entity.CatalogPrintFormatListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogPrintFormat) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogPrintFormat) Create() mrapp.HttpHandlerFunc {
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

        err := f.service.Create(c.Context(), &item)

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

func (f *CatalogPrintFormat) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPrintFormat{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPrintFormat{
            Id: f.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPrintFormat) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPrintFormat{
            Id: f.getItemId(c),
            Version: request.Version,
            Status: request.Status,
        }

        err := f.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPrintFormat) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPrintFormat) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
