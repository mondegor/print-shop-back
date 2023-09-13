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
    catalogPaperFactureListURL = "/v1/catalog-paper-factures"
    catalogPaperFactureItemURL = "/v1/catalog-paper-factures/:id"
    catalogPaperFactureChangeStatusURL = "/v1/catalog-paper-factures/:id/status"
)

type CatalogPaperFacture struct {
    service usecase.CatalogPaperFactureService
}

func NewCatalogPaperFacture(service usecase.CatalogPaperFactureService) *CatalogPaperFacture {
    return &CatalogPaperFacture{
        service: service,
    }
}

func (ht *CatalogPaperFacture) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogPaperFactureListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogPaperFactureListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogPaperFactureItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogPaperFactureItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogPaperFactureItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogPaperFactureChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogPaperFacture) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogPaperFacture) newListFilter(c mrcore.ClientData) *entity.CatalogPaperFactureListFilter {
    var listFilter entity.CatalogPaperFactureListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogPaperFacture) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogPaperFacture) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogPaperFacture{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
            Caption: request.Caption,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := view.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogPaperFactureSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogPaperFacture) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogPaperFacture{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
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

func (ht *CatalogPaperFacture) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogPaperFacture{
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

func (ht *CatalogPaperFacture) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogPaperFacture) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
