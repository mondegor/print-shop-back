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
    catalogPaperColorGetListURL = "/v1/catalog-paper-colors"
    catalogPaperColorGetItemURL = "/v1/catalog-paper-colors/:id"
    catalogPaperColorCreateURL = "/v1/catalog-paper-colors"
    catalogPaperColorStoreURL = "/v1/catalog-paper-colors/:id"
    catalogPaperColorChangeStatusURL = "/v1/catalog-paper-colors/:id/status"
    catalogPaperColorRemove = "/v1/catalog-paper-colors/:id"
)

type CatalogPaperColor struct {
    service usecase.CatalogPaperColorService
}

func NewCatalogPaperColor(service usecase.CatalogPaperColorService) *CatalogPaperColor {
    return &CatalogPaperColor{
        service: service,
    }
}

func (f *CatalogPaperColor) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperColorGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogPaperColorGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperColorCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperColorStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperColorChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperColorRemove, f.Remove())
}

func (f *CatalogPaperColor) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogPaperColor) newListFilter(c mrapp.ClientData) *entity.CatalogPaperColorListFilter {
    var listFilter entity.CatalogPaperColorListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogPaperColor) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogPaperColor) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogPaperColor{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
            Caption: request.Caption,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogPaperColorSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogPaperColor) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPaperColor{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
            Id: f.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPaperColor) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
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

func (f *CatalogPaperColor) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPaperColor) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
