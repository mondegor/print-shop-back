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
    catalogPaperFactureGetListURL = "/v1/catalog-paper-factures"
    catalogPaperFactureGetItemURL = "/v1/catalog-paper-factures/:id"
    catalogPaperFactureCreateURL = "/v1/catalog-paper-factures"
    catalogPaperFactureStoreURL = "/v1/catalog-paper-factures/:id"
    catalogPaperFactureChangeStatusURL = "/v1/catalog-paper-factures/:id/status"
    catalogPaperFactureRemove = "/v1/catalog-paper-factures/:id"
)

type CatalogPaperFacture struct {
    service usecase.CatalogPaperFactureService
}

func NewCatalogPaperFacture(service usecase.CatalogPaperFactureService) *CatalogPaperFacture {
    return &CatalogPaperFacture{
        service: service,
    }
}

func (f *CatalogPaperFacture) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperFactureGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogPaperFactureGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperFactureCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperFactureStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperFactureChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperFactureRemove, f.Remove())
}

func (f *CatalogPaperFacture) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogPaperFacture) newListFilter(c mrapp.ClientData) *entity.CatalogPaperFactureListFilter {
    var listFilter entity.CatalogPaperFactureListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogPaperFacture) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogPaperFacture) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogPaperFacture{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
            Caption: request.Caption,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogPaperFactureSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogPaperFacture) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogPaperFacture{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
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

func (f *CatalogPaperFacture) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
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

func (f *CatalogPaperFacture) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogPaperFacture) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
