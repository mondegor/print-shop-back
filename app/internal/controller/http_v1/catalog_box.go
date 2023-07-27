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
    catalogBoxGetListURL = "/v1/catalog-boxes"
    catalogBoxGetItemURL = "/v1/catalog-boxes/:id"
    catalogBoxCreateURL = "/v1/catalog-boxes"
    catalogBoxStoreURL = "/v1/catalog-boxes/:id"
    catalogBoxChangeStatusURL = "/v1/catalog-boxes/:id/status"
    catalogBoxRemoveURL = "/v1/catalog-boxes/:id"
)

type CatalogBox struct {
    service usecase.CatalogBoxService
}

func NewCatalogBox(service usecase.CatalogBoxService) *CatalogBox {
    return &CatalogBox{
        service: service,
    }
}

func (ht *CatalogBox) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogBoxGetListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogBoxGetItemURL, ht.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogBoxCreateURL, ht.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogBoxStoreURL, ht.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogBoxChangeStatusURL, ht.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogBoxRemoveURL, ht.Remove())
}

func (ht *CatalogBox) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogBox) newListFilter(c mrapp.ClientData) *entity.CatalogBoxListFilter {
    var listFilter entity.CatalogBoxListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogBox) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogBox) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogBox{}

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

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogBoxSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogBox) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogBox{}

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

func (ht *CatalogBox) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

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

func (ht *CatalogBox) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogBox) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
