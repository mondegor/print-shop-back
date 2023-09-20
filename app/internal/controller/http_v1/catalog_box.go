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
    catalogBoxListURL = "/v1/catalog-boxes"
    catalogBoxItemURL = "/v1/catalog-boxes/:id"
    catalogBoxChangeStatusURL = "/v1/catalog-boxes/:id/status"
)

type CatalogBox struct {
    service usecase.CatalogBoxService
}

func NewCatalogBox(service usecase.CatalogBoxService) *CatalogBox {
    return &CatalogBox{
        service: service,
    }
}

func (ht *CatalogBox) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogBoxListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogBoxListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogBoxItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogBoxItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogBoxItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogBoxChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogBox) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogBox) newListFilter(c mrcore.ClientData) *entity.CatalogBoxListFilter {
    var listFilter entity.CatalogBoxListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogBox) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogBox) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogBoxRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
            Article: request.Article,
            Caption: request.Caption,
            Length: request.Length,
            Width: request.Width,
            Depth: request.Depth,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogBoxSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogBox) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogBoxRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Article: request.Article,
            Caption: request.Caption,
            Length:  request.Length,
            Width:   request.Width,
            Depth:   request.Depth,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogBox) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogBox{
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

func (ht *CatalogBox) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogBox) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
