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
    catalogPaperColorRemoveURL = "/v1/catalog-paper-colors/:id"
)

type CatalogPaperColor struct {
    service usecase.CatalogPaperColorService
}

func NewCatalogPaperColor(service usecase.CatalogPaperColorService) *CatalogPaperColor {
    return &CatalogPaperColor{
        service: service,
    }
}

func (ht *CatalogPaperColor) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperColorGetListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogPaperColorGetItemURL, ht.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperColorCreateURL, ht.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperColorStoreURL, ht.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperColorChangeStatusURL, ht.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperColorRemoveURL, ht.Remove())
}

func (ht *CatalogPaperColor) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPaperColor) newListFilter(c mrapp.ClientData) *entity.CatalogPaperColorListFilter {
    var listFilter entity.CatalogPaperColorListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPaperColor) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPaperColor) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogPaperColor{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
            Caption: request.Caption,
        }

        err := ht.service.Create(c.Context(), &item)

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

func (ht *CatalogPaperColor) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPaperColor{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaperColor) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperColor{
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

func (ht *CatalogPaperColor) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaperColor) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
